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
	AddrSocks5 string
	geo        *GeoIP
	Timeout    time.Duration
}

func (ss *App) Run(ctx context.Context, cfg *ProxyCfg) {
	listener, err := net.Listen("tcp", ss.AddrSocks5)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", ss.AddrSocks5, err)
	}
	EnableInternetSetting(ss.AddrSocks5)
	defer DisableInternetSetting()
	defer listener.Close()
	log.Println("SOCKS5 server listening on: " + ss.AddrSocks5)
	for {
		select {
		case <-ctx.Done():
			log.Println("socks5 server exit")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
			go ss.handleConnection(ctx, conn, cfg)
		}
	}
}

func (ss *App) socks5HandShake(conn net.Conn) error {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return fmt.Errorf("failed to read version and nmethods: %w", err)
	}
	if buf[0] != socks5Version {
		return fmt.Errorf("socks5 only. unsupported SOCKS version: %d", buf[0])
	}

	// Read the supported authentication methods
	nMethods := int(buf[1])
	nMethodsData := make([]byte, nMethods)
	if _, err := io.ReadFull(conn, nMethodsData); err != nil {
		return fmt.Errorf("failed to read methods: %w", err)
	}

	// Select no authentication (0x00)
	if _, err := conn.Write([]byte{socks5Version, 0x00}); err != nil {
		return fmt.Errorf("failed to write method selection: %w", err)
	}
	return nil
}

func (ss *App) socks5Request(conn net.Conn) (*Socks5Request, error) {
	buf := make([]byte, 8<<10)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %w", err)
	}
	data := buf[:n]
	if len(data) < 4 {
		return nil, fmt.Errorf("request too short")
	}
	return parseSocks5Request(data, ss.geo)
}

var socks5ReplyFailBytes = []byte{socks5Version, socks5ReplyFail, socks5ReplyReserved, socks5AtypeIPv4, 0, 0, 0, 0, 0, 0}

func (ss *App) handleConnection(outerCtx context.Context, conn net.Conn, cfg *ProxyCfg) {
	defer conn.Close() // the outer for loop is not suitable for defer, so defer close here
	ctx, cf := context.WithTimeout(outerCtx, ss.Timeout)
	defer cf()

	err := ss.socks5HandShake(conn)
	if err != nil {
		log.Printf("failed to shake hand: %v", err)
		return
	}
	req, err := ss.socks5Request(conn)
	if err != nil {
		log.Printf("failed to parse SOCKS5 request: %v", err)
		return
	}
	req.Logger().Info("connect to->")
	if req.socks5Cmd == socks5CmdConnect { //tcp
		session := &SessionTcp5e{
			req:      req,
			s5:       conn,
			wsExitCh: make(chan struct{}, 1),
		}
		defer session.Close()
		err = session.breakGfwSvr(ctx, cfg)
		if err == nil {
			session.pipe(ctx, cfg.UUID)
		}
	} else if req.socks5Cmd == socks5CmdUdpAssoc {
		session := SessionUdp{
			SessionTcp: SessionTcp{
				req:      req,
				s5:       conn,
				wsExitCh: make(chan struct{}, 1),
			},
			udpConn: nil,
		}
		defer session.Close()
		err = session.breakGfwSvr(cfg)
		if err == nil {
			session.pipe(ctx, cfg.UUID)
		}
	} else if req.socks5Cmd == socks5CmdBind {
		err = fmt.Errorf("unsupported command: BIND")
	} else {
		err = fmt.Errorf("unknown command: %d", req.socks5Cmd)
	}
	//handle all error
	if err != nil {
		log.Println(err)
		conn.Write(socks5ReplyFailBytes)
	}
}
