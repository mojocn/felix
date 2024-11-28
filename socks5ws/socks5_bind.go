package socks5ws

import (
	"io"
	"log/slog"
	"net"
	"sync"
)

func relayBind(s5 net.Conn, _ *Socks5Request) {
	bindListener, err := net.Listen("tcp4", ":0")
	if err != nil {
		slog.Error("bind tcp failed", "err", err)
		socks5Response(s5, net.IPv4zero, 0, socks5ReplyFail)
		return
	}
	defer bindListener.Close()
	//first reply
	localAddr := bindListener.Addr().(*net.TCPAddr)
	socks5Response(s5, localAddr.IP, localAddr.Port, socks5ReplyOkay)

	targetConn, err := bindListener.Accept()
	if err != nil {
		slog.Error("bind tcp failed", "err", err)
		return
	}
	defer targetConn.Close()
	//sec reply
	targetAddr := targetConn.RemoteAddr().(*net.TCPAddr)
	socks5Response(s5, targetAddr.IP, targetAddr.Port, socks5ReplyOkay)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(targetConn, s5)
		if err != nil {
			slog.Error("bind tcp failed", "err", err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(s5, targetConn)
		if err != nil {
			slog.Error("bind tcp failed", "err", err)
		}
	}()
	wg.Wait()
}
