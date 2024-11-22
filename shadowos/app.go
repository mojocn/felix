package shadowos

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type App struct {
	AddrWs     string
	AddrSocks5 string
	UUID       string
	Timeout    time.Duration
}

func (ss *App) Run() {
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

func (ss *App) handshakeNoAuth(conn net.Conn) error {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return fmt.Errorf("failed to read version and nmethods: %w", err)
	}
	if buf[0] != SOCKS5VERSION {
		return fmt.Errorf("socks5 only. unsupported SOCKS version: %d", buf[0])
	}

	// Read the supported authentication methods
	nMethods := int(buf[1])
	nMethodsData := make([]byte, nMethods)
	if _, err := io.ReadFull(conn, nMethodsData); err != nil {
		return fmt.Errorf("failed to read methods: %w", err)
	}

	// Select no authentication (0x00)
	if _, err := conn.Write([]byte{SOCKS5VERSION, 0x00}); err != nil {
		return fmt.Errorf("failed to write method selection: %w", err)
	}
	return nil
}

type Socks5Request struct {
	socks5Cmd  byte
	socks5Atyp byte
	dstAddr    []byte
	dstPort    []byte
}

func (s Socks5Request) String() string {
	return fmt.Sprintf("socks5Cmd: %v, socks5Atyp: %v, dstAddr: %v, dstPort: %v", s.socks5Cmd, s.socks5Atyp, s.dstAddr, s.dstPort)
}
func (s Socks5Request) addressBytes() []byte {
	if s.socks5Atyp == socks5AtypeDomain {
		return append([]byte{byte(len(s.dstAddr))}, s.dstAddr...)
	}
	return s.dstAddr
}

func (s Socks5Request) vlessHeader(uuid [16]byte) ([]byte, error) {
	addrBytes := s.addressBytes()
	//https://xtls.github.io/development/protocols/vless.html
	headerBytes := make([]byte, 0, 1+16+1+1+2+1+len(addrBytes))

	headerBytes = append(headerBytes, 0x01)       // version
	headerBytes = append(headerBytes, uuid[:]...) //16 bytes of UUID
	headerBytes = append(headerBytes, 0x00)       // additional info length M

	//1 byte of command
	if s.socks5Cmd == socks5CmdUdpAssoc {
		headerBytes = append(headerBytes, byte(vlessCmdUdp))
	} else if s.socks5Cmd == socks5CmdConnect {
		headerBytes = append(headerBytes, byte(vlessCmdTcp))
	} else {
		return nil, fmt.Errorf("unsupported command: %d", s.socks5Cmd)
	}

	headerBytes = append(headerBytes, s.dstPort...) //2 bytes of port

	//1 byte of address type
	if s.socks5Atyp == socks5AtypeIPv4 {
		headerBytes = append(headerBytes, byte(vlessAtypeIPv4))
	} else if s.socks5Atyp == socks5AtypeIPv6 {
		headerBytes = append(headerBytes, byte(vlessAtypeIPv6))
	} else if s.socks5Atyp == socks5AtypeDomain {
		headerBytes = append(headerBytes, byte(vlessAtypeDomain))
	} else {
		return nil, fmt.Errorf("unsupported address type: %d", s.socks5Atyp)
	}

	headerBytes = append(headerBytes, addrBytes...) //n bytes of address
	return headerBytes, nil
}

func (*App) requestInfo(conn net.Conn) (*Socks5Request, error) {
	buf := make([]byte, 8<<10)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %w", err)
	}
	data := buf[:n]
	if len(data) < 4 {
		return nil, fmt.Errorf("request too short")
	}
	if data[0] != socks5Ver {
		return nil, fmt.Errorf("unsupported SOCKS version: %d", data[0])
	}
	info := new(Socks5Request)
	if data[1] == socks5CmdConnect {
		info.socks5Cmd = socks5CmdConnect
	} else if data[1] == socks5CmdUdpAssoc {
		info.socks5Cmd = socks5CmdUdpAssoc
	} else {
		//BIND is not supported
		return nil, fmt.Errorf("unsupported command: %d", data[1])
	}
	if data[2] != 0x00 {
		return nil, fmt.Errorf("RSV must be 0x00")
	}
	if data[3] == socks5AtypeIPv4 {
		if len(data) < 10 {
			return nil, fmt.Errorf("request too short for atyp IPv4")
		}
		info.socks5Atyp = socks5AtypeIPv4
		info.dstAddr = data[4:8]
		info.dstPort = data[8:10]
	} else if data[3] == socks5AtypeDomain {
		if len(data) < 5 {
			return nil, fmt.Errorf("request too short for atyp Domain")
		}
		addrLen := int(data[4])
		info.socks5Atyp = socks5AtypeDomain
		info.dstAddr = data[5 : 5+addrLen]
		info.dstPort = data[5+addrLen : 5+addrLen+2]
	} else if data[3] == socks5AtypeIPv6 {
		if len(data) < 22 {
			return nil, fmt.Errorf("request too short for atyp IPv6")
		}
		info.socks5Atyp = socks5AtypeIPv6
		info.dstAddr = data[4:20]
		info.dstPort = data[20:22]
	} else {
		return nil, fmt.Errorf("unsupported address type: %d", data[3])
	}
	return info, nil
}

type ProxyCfg struct {
	AddrWs string
	UUID   [16]byte
}

func (ss *App) handleConnection(conn net.Conn) {
	defer conn.Close()
	ctx, cf := context.WithTimeout(context.Background(), ss.Timeout)
	defer cf()

	err := ss.handshakeNoAuth(conn)
	if err != nil {
		log.Printf("failed to shake hand: %v", err)
		return
	}
	req, err := ss.requestInfo(conn)
	if err != nil {
		log.Printf("failed to parse SOCKS5 request: %v", err)
		return
	}
	if req.socks5Cmd == socks5CmdUdpAssoc {
		udpAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
		udpConn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			log.Printf("failed to bind UDP socket: %v", err)
			return
		}
		defer udpConn.Close()
		boundAddr := udpConn.LocalAddr().(*net.UDPAddr)
		response := []byte{
			socksVersion, socks5ReplySuccess, 0x00, socks5AtypeIPv4,
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

	} else if req.socks5Cmd == socks5CmdConnect { //tcp
		session := &SessionTcp{
			req:      req,
			s5:       conn,
			proxyCfg: &ProxyCfg{AddrWs: ss.AddrWs}, //todo config dynamic
		}
		err := session.connectProxyServer()
		if err != nil {
			log.Println(err)
			return
		}
		defer session.Close()
		session.pipe(ctx)

	} else {
		log.Printf("unsupported command: %d", req.socks5Cmd)
		return
	}

}
