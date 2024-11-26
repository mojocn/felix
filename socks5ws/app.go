package socks5ws

import (
	"context"
	"errors"
	"fmt"
	"github.com/mojocn/felix/model"
	"io"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"
)

type ClientLocalSocks5Server struct {
	AddrSocks5 string
	geo        *GeoIP
	Timeout    time.Duration
	proxy      *model.Proxy
}

func NewClientLocalSocks5Server(addr, geoIpPath string) (*ClientLocalSocks5Server, error) {
	geo, err := NewGeoIP(geoIpPath)
	if err != nil {
		return nil, err
	}
	return &ClientLocalSocks5Server{
		AddrSocks5: addr,
		geo:        geo,
		Timeout:    5 * time.Minute,
	}, nil

}

func (ss *ClientLocalSocks5Server) fetchActiveProxy() {
	var proxies []model.Proxy
	err := model.DB().Find(&proxies).Error
	if err != nil {
		slog.Error("failed to get proxy setting", "err", err.Error())
		return
	}
	if len(proxies) == 0 {
		slog.Error("no proxy setting found")
		return
	}
	ss.proxy = &proxies[0]
	for _, proxy := range proxies {
		if proxy.IsActive() {
			ss.proxy = &proxy
			break
		}
	}
}

func (ss *ClientLocalSocks5Server) Run(ctx context.Context) {
	ss.fetchActiveProxy()

	listener, err := net.Listen("tcp", ss.AddrSocks5)
	if err != nil {
		listener, err = net.Listen("tcp4", "127.0.0.1:0")
	}
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", ss.AddrSocks5, err)
	}
	ss.AddrSocks5 = listener.Addr().String()
	slog.Info("socks5 server listening on", "addr", ss.AddrSocks5)

	defer listener.Close()
	log.Println("SOCKS5 server listening on: " + ss.AddrSocks5)
	//proxySettingOn(ss.AddrSocks5)
	//defer proxySettingOff()
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
			go ss.handleConnection(ctx, conn)
		}
	}
}

func (ss *ClientLocalSocks5Server) socks5HandShake(conn net.Conn) error {
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

func (ss *ClientLocalSocks5Server) socks5Request(conn net.Conn) (*Socks5Request, error) {
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

func (ss *ClientLocalSocks5Server) handleConnection(outerCtx context.Context, conn net.Conn) {
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
	req.Logger().Info("remote target")
	if req.socks5Cmd == socks5CmdConnect { //tcp
		relayTcpSvr, err := ss.dispatchRelayTcpServer(ctx, req)
		if checkSocks5Request(conn, err) {
			return
		}
		defer relayTcpSvr.Close()
		ss.pipeTcp(ctx, conn, relayTcpSvr)
		return
	} else if req.socks5Cmd == socks5CmdUdpAssoc {
		session := SessionUdp{
			req:      req,
			s5:       conn,
			wsExitCh: make(chan struct{}, 1),
			udpConn:  nil,
		}
		defer session.Close()
		err = session.breakGfwSvr(ss.proxy)
		if err == nil {
			session.pipe(ctx, [16]byte{})
		}
		return
	} else if req.socks5Cmd == socks5CmdBind {
		err = fmt.Errorf("unsupported command: BIND")
	} else {
		err = fmt.Errorf("unknown command: %d", req.socks5Cmd)
	}
	//handle all error
	checkSocks5Request(conn, err)
}

func (ss *ClientLocalSocks5Server) shouldGoDirect(req *Socks5Request) (goDirect bool) {
	if req.CountryCode == "CN" || req.CountryCode == "" {
		//empty means geo ip failed or local address
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

func (ss *ClientLocalSocks5Server) dispatchRelayTcpServer(ctx context.Context, req *Socks5Request) (io.ReadWriteCloser, error) {
	if ss.shouldGoDirect(req) {
		return NewRelayTcpDirect(req)
	}
	cfg := ss.proxy
	return NewRelayTcpSocks5e(ctx, cfg, req)
}

func (ss *ClientLocalSocks5Server) pipeTcp(ctx context.Context, s5 net.Conn, relayRw io.ReadWriter) {
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
					span.Debug("relay read", "err", err.Error())
					return
				}
				_, err = s5.Write(buf[:n])
				if err != nil {
					span.Error("s5 write", "err", err.Error())
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
					span.Error("relay write", "err", err.Error())
					return
				}
			}
		}
	}()
	wg.Wait()
	slog.Debug("2 goroutines is Done")
}
