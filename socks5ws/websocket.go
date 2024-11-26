package socks5ws

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mojocn/felix/model"
	"log/slog"
	"net/http"
)

const (
	browserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

func webSocketConn(ctx context.Context, proxy *model.Proxy, req *Socks5Request) (*websocket.Conn, error) {
	wsDialer := websocket.DefaultDialer

	headers := http.Header{}
	headers.Set("Authorization", proxy.UserID)
	headers.Set("User-Agent", browserAgent)
	if proxy.Sni != "" {
		headers.Set("Host", proxy.Sni)
		wsDialer.TLSClientConfig = &tls.Config{
			ServerName: proxy.Sni, // Set the SNI to the hostname of the server
		}
	}
	headers.Set("x-req-id", req.id)
	headers.Set("x-felix-network", "tcp")
	headers.Set("x-felix-addr", req.host())
	headers.Set("x-felix-port", req.port())
	headers.Set("x-felix-protocol", proxy.Version)
	url := proxy.RelayURL()
	slog.Debug("connecting to remote proxy server", "url", url)
	ws, resp, err := websocket.DefaultDialer.DialContext(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote proxy server: %s ,error:%v", proxy.RelayURL(), err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, fmt.Errorf("failed to connect to remote proxy server: %s ,error:%v", proxy.RelayURL(), err)
	}
	return ws, nil
}
