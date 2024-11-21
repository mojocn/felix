package shadowos

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mojocn/felix/util"
	"io"
	"log"
	"net"
	"os"
	"time"
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
		go ss.handleConnection(conn)
	}
}

func socks5packet(conn net.Conn, uuidS string) (connData []byte, err error) {
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

	connBytes, err := socks5packet(conn, ss.UUID)
	if err != nil {
		log.Printf("failed to parse SOCKS5 request: %v", err)
		return
	}
	// Read the version and number of authentication methods

	// Connect to the target server
	ws, _, err := websocket.DefaultDialer.Dial(ss.AddrWs, nil)
	if err != nil {
		log.Printf("failed to connect to target: %v", err)
		conn.Write([]byte{SOCKS5VERSION, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	} else { // Send success response
		conn.Write([]byte{SOCKS5VERSION, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	}
	defer func() {
		err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		log.Println("send close ws msg", err)
		ws.Close()
		log.Println("closed ws")
	}()
	// Relay data between client and target server
	pipeWebsocketSocks5(context.Background(), ws, conn, connBytes)
}

func pipeWebsocketSocks5(ctx context.Context, ws *websocket.Conn, s5 net.Conn, firstData []byte) {
	ws.SetCloseHandler(func(code int, text string) error {
		log.Println("ws closed", code, text)
		return nil
	})
	done := make(chan struct{})

	go func() {
		isFirstData := true
		defer func() {
			log.Println("ddd ws -> s5")
			done <- struct{}{}
		}()
		for {
			select {
			case <-done:
				log.Println("wss5 done")
				return
			case <-ctx.Done():
				log.Println("done: ws -> s5")
				return
			default:
				//ws.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
				_, data, err := ws.ReadMessage()
				n := len(data)
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
					log.Println("EOF from ws")
					return
				}
				var e *net.OpError
				if errors.As(err, &e) {
					log.Println("net.OpError", e)
					return
				}
				if err != nil {
					log.Printf("other ws -> socks5 error %T", err)
					log.Println("other ws -> socks5 error", err)
					continue
				}
				fromByteIndex := 0
				// skip the first data
				if isFirstData && n >= 2 {
					extraN := data[1]
					isFirstData = false
					fromByteIndex = 2 + int(extraN)
				}
				log.Println("write back socks5", n)
				s5.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
				_, err = s5.Write(data[fromByteIndex:n])
				if err != nil {
					log.Println(" ws -> socks5 error", err)
					return
				}
			}
		}
	}()
	go func() { // s5 -> ws
		defer func() {
			log.Println("ddd s5 -> ws")
			done <- struct{}{}
		}()
		for {
			select {
			case <-done:
				log.Println("s5ws done")
				return
			case <-ctx.Done():
				log.Println("done: s5 -> ws")
				return
			default:
				buf := make([]byte, 8<<10)
				//s5.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
				n, err := s5.Read(buf)
				if errors.Is(err, os.ErrDeadlineExceeded) {
					continue
				}
				if errors.Is(err, io.EOF) {
					log.Println("EOF from socks5")
					return
				}
				if err != nil {
					log.Printf("read from socks5 error %T", err)
					log.Println("read from socks5 error", err)
					continue
				}
				log.Println("read from socks5", n)
				data := buf[:n]
				if len(firstData) > 0 {
					log.Println("send version header only once")
					data = append(firstData, buf[:n]...)
					firstData = nil
				}
				err = ws.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					log.Println("write error", err)
					return
				}
			}
		}
	}()
	<-done
	log.Println("2 done")

}
