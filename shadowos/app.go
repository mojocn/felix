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

var socks5ReplyFailBytes = []byte{socksVersion, socks5ReplyFail, 0x00, 0x01, 0, 0, 0, 0, 0, 0}

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
		session := SessionUdp{
			SessionTcp: SessionTcp{
				req:      req,
				s5:       conn,
				proxyCfg: &ProxyCfg{AddrWs: ss.AddrWs}, //todo config dynamic
				wsExitCh: make(chan struct{}, 1),
			},
			udpConn: nil,
		}
		defer session.Close()
		err = session.connectProxyServer()
		if err == nil {
			session.pipe()
		}

		return

	} else if req.socks5Cmd == socks5CmdConnect { //tcp
		session := &SessionTcp{
			req:      req,
			s5:       conn,
			proxyCfg: &ProxyCfg{AddrWs: ss.AddrWs}, //todo config dynamic
			wsExitCh: make(chan struct{}, 1),
		}
		defer session.Close()
		err = session.connectProxyServer()
		if err == nil {
			session.pipe(ctx)
		}
	} else {
		err = fmt.Errorf("unsupported command: %d", req.socks5Cmd)
	}
	//handle all error
	if err != nil {
		log.Println(err)
		conn.Write(socks5ReplyFailBytes)
	}
}
