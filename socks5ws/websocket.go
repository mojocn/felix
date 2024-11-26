package socks5ws

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func webSocketConn(ctx context.Context, proxy *ProxyCfg, req *Socks5Request) (*websocket.Conn, error) {
	headers := http.Header{}
	for k, v := range proxy.WsHeader {
		headers[k] = v
	}
	headers.Set("x-req-id", req.id)
	headers.Set("Authorization", proxy.uuidHex())
	headers.Set("x-felix-network", "tcp")
	headers.Set("x-felix-addr", req.host())
	headers.Set("x-felix-port", req.port())
	headers.Set("x-felix-protocol", proxy.Protocol)

	ws, resp, err := websocket.DefaultDialer.DialContext(ctx, proxy.WsUrl, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote proxy server: %s ,error:%v", proxy.WsUrl, err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		log.Println("ws connected failed", resp.Status)
	}
	return ws, nil
}
