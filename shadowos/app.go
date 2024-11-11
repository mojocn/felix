package shadowos

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/gorilla/websocket"
	"github.com/mojocn/felix/util"
)

type VlessCmd byte
type VlessAddrType byte

const (
	VlessCmdTcp VlessCmd = 0x01
	VlessCmdUdp VlessCmd = 0x02
	VlessCmdMux VlessCmd = 0x03

	VlessAddrTypeIPv4   VlessAddrType = 0x01
	VlessAddrTypeDomain VlessAddrType = 0x02
	VlessAddrTypeIPv6   VlessAddrType = 0x03

	SOCKS5VERSION = 0x05
	CMD_CONNECT   = 0x01
	CMD_BIND      = 0x02
	CMD_UDP_ASSOC = 0x03
)

type ShadowosApp struct {
	AddrWs     string
	AddrSocks5 string
	UUID       string
}

func (ss *ShadowosApp) Run() {
	listener, err := net.Listen("tcp", ss.AddrSocks5)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", ss.AddrSocks5, err)
	}
	defer listener.Close()
	log.Println("SOCKS5 server listening on: " + ss.AddrSocks5)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		log.Println("new request:-->")
		go ss.handleConnection(conn)
	}
}

func handshake(conn net.Conn, uuidS string) (connData []byte, err error) {
	uuidBytes, err := util.UUID2bytes(uuidS)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UUID: %w", err)
	}
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("failed to read version and nmethods: %w", err)
	}
	if buf[0] != SOCKS5VERSION {
		return nil, fmt.Errorf("socks5 only. unsupported SOCKS version: %d", buf[0])
	}

	// Read the supported authentication methods
	nMethods := int(buf[1])
	nMethodsData := make([]byte, nMethods)
	if _, err := io.ReadFull(conn, nMethodsData); err != nil {
		return nil, fmt.Errorf("failed to read methods: %w", err)
	}

	// Select no authentication (0x00)
	if _, err := conn.Write([]byte{SOCKS5VERSION, 0x00}); err != nil {
		return nil, fmt.Errorf("failed to write method selection: %w", err)
	}

	// Read the request
	buf = make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("failed to read request: %w", err)
	}

	if buf[0] != SOCKS5VERSION {
		return nil, fmt.Errorf("unsupported SOCKS version: %d", buf[0])
	}
	var vlessCmd VlessCmd
	if buf[1] == CMD_CONNECT {
		vlessCmd = VlessCmdTcp
	} else if buf[1] == CMD_UDP_ASSOC {
		vlessCmd = VlessCmdUdp
	} else {
		return nil, fmt.Errorf("unsupported command: %d", buf[1])
	}
	var vlessAddrType VlessAddrType
	addrLen := 0
	// Read the address
	switch buf[3] {
	case 0x01: // IPv4
		vlessAddrType = VlessAddrTypeIPv4
		addrLen = net.IPv4len
	case 0x03: // Domain name
		vlessAddrType = VlessAddrTypeDomain
		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			return nil, fmt.Errorf("failed to read domain length: %w", err)
		}
		addrLen = int(buf[0])
	case 0x04: // IPv6
		vlessAddrType = VlessAddrTypeIPv6
		addrLen = net.IPv6len
	default:
		return nil, fmt.Errorf("unsupported address type: %d", buf[3])
	}
	addrBytes := make([]byte, addrLen)
	if _, err := io.ReadFull(conn, addrBytes); err != nil {
		return nil, fmt.Errorf("failed to read address: %w", err)
	}
	if vlessAddrType == VlessAddrTypeDomain {
		addrBytes = append([]byte{byte(addrLen)}, addrBytes...)
	}

	// Read the port
	remotePort := make([]byte, 2)

	if _, err = io.ReadFull(conn, remotePort); err != nil {
		log.Printf("Failed to read port: %v", err)
		return
	}
	// Construct vless packet
	connData = make([]byte, 0, 1+16+1+1+2+1+len(addrBytes))
	connData = append(connData, 0x01)                // version
	connData = append(connData, uuidBytes...)        //16 bytes of UUID
	connData = append(connData, 0x00)                // additional info length M
	connData = append(connData, byte(vlessCmd))      //1 byte of command
	connData = append(connData, remotePort...)       //2 bytes of port
	connData = append(connData, byte(vlessAddrType)) //1 byte of address type
	connData = append(connData, addrBytes...)        //n bytes of address
	return connData, nil
}

func (ss *ShadowosApp) handleConnection(conn net.Conn) {
	defer conn.Close()
	connBytes, err := handshake(conn, ss.UUID)
	if err != nil {
		log.Printf("failed to parse SOCKS5 request: %v", err)
		return
	}
	// Read the version and number of authentication methods

	// Connect to the target server
	session, err := NewProxySession(ss.AddrWs, connBytes)
	if err != nil {
		log.Printf("failed to connect to target: %v", err)
		conn.Write([]byte{SOCKS5VERSION, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	} else { // Send success response
		conn.Write([]byte{SOCKS5VERSION, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	}
	defer session.Close()

	session.doProxy(conn)
}

type ProxySession struct {
	ws          *websocket.Conn
	connData    []byte
	isFirstData bool
	ch          chan struct{}
	nextRead    chan struct{}
}

func NewProxySession(url string, initialData []byte) (*ProxySession, error) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket server: %w", err)
	}

	return &ProxySession{
		ws:          c,
		connData:    initialData,
		isFirstData: true,
		ch:          make(chan struct{}, 1),
		nextRead:    make(chan struct{}, 1),
	}, nil
}

func (ps ProxySession) Close() error {
	log.Println("websocket close message sent")

	err := ps.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("failed to send close message", err)
		return err
	}
	return ps.ws.Close()
	// return nil
}

func (ps *ProxySession) doProxy(socks net.Conn) {
	go func() {
		defer func() {
			ps.ch <- struct{}{}
		}()
		for {
			select {
			case <-ps.ch:
				return
			default:
				messageType, data, err := ps.ws.ReadMessage()
				log.Println("ws2socks ->", messageType, data, err)
				if err != nil {
					log.Println("failed to read from websocket to socks5", err)
				}
				if len(data) > 0 && messageType == websocket.BinaryMessage {
					if ps.isFirstData && len(data) > 1 {
						ps.isFirstData = false
						extraN := int(data[1]) + 2
						data = data[extraN:]
					}
					nn := 0
					nn, err = socks.Write(data)
					if err != nil {
						log.Println(err)
					}
					log.Println(nn)
				}
				if err != nil {
					log.Printf("%T", err)
					log.Println("messageType", messageType)
					log.Println("failed to read from websocket to socks5", err)
					return
				}
			}
		}
	}()

	go func() {
		defer func() {
			ps.ch <- struct{}{}
		}()
		for {
			buf := make([]byte, 1024)
			n, err := socks.Read(buf)
			if n > 0 {
				log.Println("socks read N:", n)
				if len(ps.connData) > 0 {
					buf = append(ps.connData, buf[:n]...)
					ps.connData = nil
				} else {
					buf = buf[:n]
				}
				err = ps.ws.WriteMessage(websocket.BinaryMessage, buf)
				if err != nil {
					log.Println("failed to write to websocket", err)
				}
			}
			//socks5 EOF
			if err != io.EOF {
				continue
			}
			if err != net.ErrClosed {
				log.Print("socks5 closed")
				return
			}
			if err != nil {
				log.Printf("%T", err)
				log.Println("failed to read from socks5 to websocket", err)
				return
			}
		}
	}()
	<-ps.ch
	log.Print("doProxy done")
}
