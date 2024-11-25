package shadowos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"
)

type App struct {
	AddrSocks5 string
	geo        *GeoIP
	Timeout    time.Duration
}

func NewApp(addr, geoIpPath string) (*App, error) {
	geo, err := NewGeoIP(geoIpPath)
	if err != nil {
		return nil, err
	}
	return &App{
		AddrSocks5: addr,
		geo:        geo,
		Timeout:    5 * time.Minute,
	}, nil

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
	return parseSocks5Request(data)
}

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
		relayTcp, err := ss.dispatchRelayTcpServer(ctx, cfg, req)
		if checkSocks5Request(conn, err) {
			return
		}
		defer relayTcp.Close()
		ss.pipTcp(ctx, conn, relayTcp)
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
	checkSocks5Request(conn, err)
}

func (ss *App) shouldGoDirect(req *Socks5Request) (goDirect bool) {
	countryCode, err := ss.geo.country(req.host())
	if err != nil {
		slog.Error("geoip failed", "err", err.Error())
		return true
	}
	slog.Info("countryCode", "code", countryCode, "host", req.host())
	if countryCode == "CN" || countryCode == "" {
		//empty means geoip failed or local address
		return true
	}
	return false
}

func checkSocks5Request(socks5conn net.Conn, err error) (hasError bool) {
	hasError = err != nil
	if hasError {
		slog.Error("failed reason:", "err", err.Error())
		_, err = socks5conn.Write(socks5ReplyBytesFailed)
	} else {
		_, err = socks5conn.Write(socks5ReplyBytesSuccess)
	}
	if err != nil {
		slog.Error("socks5 request rely failed to write", "err", err.Error())
	}
	return hasError
}

func (ss *App) dispatchRelayTcpServer(ctx context.Context, cfg *ProxyCfg, req *Socks5Request) (io.ReadWriteCloser, error) {
	if ss.shouldGoDirect(req) {
		return NewRelayTcpDirect(req)
	}
	return NewRelayTcpSocks5e(ctx, cfg, req)
}

func (ss *App) pipTcp(ctx context.Context, s5 net.Conn, relayRw io.ReadWriter) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		span := slog.With("fn", "ws -> s5")
		defer func() {
			span.Debug("wg done")
			wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				span.Info("ctx.Done exit")
				return
			default:
				//ws.SetReadDeadline(time.Now().Add(1 * time.Second))
				buf := make([]byte, 8<<10)
				n, err := relayRw.Read(buf)
				if err != nil {
					span.Debug("other ws -> socks5 error", "err", err.Error())
					return
				}
				_, err = s5.Write(buf[:n])
				if err != nil {
					span.Error(" ws -> socks5", "err", err.Error())
					return
				}
			}
		}
	}()
	go func() { // s5 -> ws
		span := slog.With("fn", "s5 -> ws")
		defer func() {
			span.Debug("wg done")
			wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				span.Debug("ctx.Done exit")
				return
			default:
				buf := make([]byte, 8<<10)
				//s5.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
				n, err := s5.Read(buf)
				if errors.Is(err, io.EOF) {
					return
				}
				if err != nil {
					et := fmt.Sprintf("%T", err)
					span.With("errType", et).Error("s5 read", "err", err.Error())
					return
				}
				//ws.SetWriteDeadline(time.Now().Add(1 * time.Second))
				n, err = relayRw.Write(buf[:n])
				if err != nil {
					span.Error("write error", "err", err.Error())
					return
				}
			}
		}
	}()
	wg.Wait()
	slog.Debug("2 goroutines is Done")
}
