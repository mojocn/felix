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
	"sync"
	"time"
)

type VlessCmd byte
type VlessAddrType byte

const (
	socksVersion    = 0x05
	cmdUDPAssociate = 0x03
	replySuccess    = 0x00

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

	addressTypeIPv4   = 0x01
	addressTypeDomain = 0x03
	addressTypeIPv6   = 0x04
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
func handleUDPPacket(conn *net.UDPConn, clientAddr *net.UDPAddr, packet []byte) {
	// Parse UDP packet
	headerLen := 3 + int(packet[4]) // Assuming SOCKS5 UDP header
	data := packet[headerLen:]

	targetAddr := packet[3:]
	log.Printf("Received UDP packet for %v from %v", targetAddr, clientAddr)

	// Forward data (implement logic for forwarding here)

	// Example: Echo back to client
	if _, err := conn.WriteToUDP(data, clientAddr); err != nil {
		log.Printf("Failed to send response to %v: %v", clientAddr, err)
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
	case addressTypeIPv4: // IPv4
		vlessAddrType = VlessAddrTypeIPv4
		addrLen = net.IPv4len
	case addressTypeDomain: // Domain name
		vlessAddrType = VlessAddrTypeDomain
		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			return nil, fmt.Errorf("failed to read domain length: %w", err)
		}
		addrLen = int(buf[0])
	case addressTypeIPv6: // IPv6
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

	if len(connBytes) >= 18 && connBytes[18] == byte(VlessCmdUdp) {
		udpAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
		udpConn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			log.Printf("failed to bind UDP socket: %v", err)
			return
		}
		defer udpConn.Close()
		boundAddr := udpConn.LocalAddr().(*net.UDPAddr)
		response := []byte{
			socksVersion, replySuccess, 0x00, addressTypeIPv4,
			boundAddr.IP[0], boundAddr.IP[1], boundAddr.IP[2], boundAddr.IP[3],
			byte(boundAddr.Port >> 8), byte(boundAddr.Port & 0xFF),
		}
		conn.Write(response)
		go func() {
			for {
				packet := make([]byte, 65535)
				n, clientAddr, err := udpConn.ReadFromUDP(packet)
				if err != nil {
					log.Printf("UDP read error: %v", err)
					continue
				}

				go handleUDPPacket(udpConn, clientAddr, packet[:n])
			}
		}()
		buf := make([]byte, 1)
		conn.Read(buf) // Block until client closes the connection
		return
	}

	// Connect to the target server
	ws, _, err := websocket.DefaultDialer.Dial(ss.AddrWs, nil)
	if err != nil {
		log.Printf("failed to connect to target: %v", err)
		conn.Write([]byte{socksVersion, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	} else { // Send success response
		conn.Write([]byte{socksVersion, replySuccess, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	}
	defer func() {
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		log.Println("send close ws msg")
		ws.Close()
		log.Println("closed ws")
	}()
	// Relay data between client and target server
	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()
	pipeWebsocketSocks5(ctx, ws, conn, connBytes)
}

func pipeWebsocketSocks5(ctx context.Context, ws *websocket.Conn, s5 net.Conn, firstData []byte) {
	exitWs := make(chan struct{}, 1)
	ws.SetCloseHandler(func(code int, text string) error {
		log.Println("ws closed", code, text)
		exitWs <- struct{}{}
		return nil
	})
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		isFirstData := true
		defer func() {
			log.Println("[ddd] ws -> s5")
			wg.Done()
		}()
		for {
			select {
			case <-exitWs:
				log.Println("exitWs")
				return
			case <-ctx.Done():
				log.Println("doneCh: ws -> s5")
				return
			default:
				ws.SetReadDeadline(time.Now().Add(1 * time.Second))
				_, data, err := ws.ReadMessage()
				n := len(data)
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
					log.Println("EOF from ws")
					return
				}
				if err != nil {
					log.Printf("other ws -> socks5 error %T", err)
					log.Println("other ws -> socks5 error", err)
					return
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
			log.Println("[ddd] s5 -> ws")
			wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				log.Println("doneCh: s5 -> ws")
				return
			default:
				buf := make([]byte, 8<<10)
				s5.SetReadDeadline(time.Now().Add(1 * time.Second))
				n, err := s5.Read(buf)
				if errors.Is(err, os.ErrDeadlineExceeded) {
					continue
				}
				if errors.Is(err, io.EOF) {
					log.Println("EOF from socks5")
					return
				}
				var opErr *net.OpError
				if err != nil && errors.As(err, &opErr) {
					log.Println("opErr", opErr)
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
	wg.Wait()
	log.Println("2 doneCh")

}
