package socks5ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mojocn/felix/model"
	"log"
	"net/http"
)

const (
	browserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

func webSocketConn(ctx context.Context, proxy *model.Proxy, req *Socks5Request) (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Set("x-req-id", req.id)
	headers.Set("Authorization", proxy.UserID)
	headers.Set("User-Agent", browserAgent)
	ws, resp, err := websocket.DefaultDialer.DialContext(ctx, proxy.RelayURL(), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote proxy server: %s ,error:%v", proxy.RelayURL(), err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, fmt.Errorf("failed to connect to remote proxy server: %s ,error:%v", proxy.RelayURL(), err)
	}
	err = ws.WriteMessage(websocket.BinaryMessage, toDstConn(proxy, req))
	if err != nil {
		return nil, fmt.Errorf("failed to send dst conn info to remote proxy server %w", err)
	}
	return ws, nil
}

type DstConn struct {
	Network  string `json:"network"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"` //vless,ws_socks5,
}

func toDstConn(proxy *model.Proxy, req *Socks5Request) []byte {
	info := &DstConn{
		Network:  req.Network(),
		Host:     req.host(),
		Port:     req.port(),
		Protocol: proxy.Protocol,
	}
	data, err := json.Marshal(info)
	if err != nil {
		log.Println("failed to marshal dst conn info", err)
		return nil
	}
	return data
}
