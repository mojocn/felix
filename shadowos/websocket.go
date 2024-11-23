package shadowos

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func webSocketConn(proxy *ProxyCfg, req *Socks5Request) (*websocket.Conn, error) {
	headers := http.Header{}
	for k, v := range proxy.WsHeader {
		headers[k] = v
	}
	headers.Set("x-req-id", req.id)
	headers.Set("Authorization", proxy.uuidHex())
	headers.Set("x-felix-network", "tcp")
	headers.Set("x-felix-addr", req.addr())
	headers.Set("x-felix-port", req.port())

	ws, resp, err := websocket.DefaultDialer.Dial(proxy.WsUrl, proxy.WsHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote proxy server: %s ,error:%v", proxy.WsUrl, err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		log.Println("ws connected failed", resp.Status)
	}
	return ws, nil
}
