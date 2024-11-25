package shadowos

import (
	"context"
	"github.com/gorilla/websocket"
	"log/slog"
	"time"
)

var _ RelayTcp = (*RelayTcpSocks5e)(nil)

type RelayTcpSocks5e struct {
	cfg  *ProxyCfg
	req  *Socks5Request
	conn *websocket.Conn
}

func NewRelayTcpSocks5e(ctx context.Context, cfg *ProxyCfg, req *Socks5Request) (*RelayTcpSocks5e, error) {
	ws, err := webSocketConn(ctx, cfg, req)
	if err != nil {
		return nil, err
	}
	ws.SetCloseHandler(func(code int, text string) error {
		slog.Debug("ws has closed", "code", code, "text", text)
		return nil
	})
	return &RelayTcpSocks5e{cfg: cfg, req: req, conn: ws}, nil
}

func (r RelayTcpSocks5e) Read(data []byte) (n int, err error) {
	if r.conn != nil {
		_, p, err := r.conn.ReadMessage()
		if err != nil {
			slog.Error("failed to read ws", "err", err.Error())
		}
		return copy(data, p), err
	}
	return 0, nil
}

func (r RelayTcpSocks5e) Write(data []byte) (n int, err error) {
	if r.conn != nil {
		err = r.conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			slog.Error("failed to write ws", "err", err.Error())
		}
		return len(data), err
	}
	return 0, nil
}

func (r RelayTcpSocks5e) Close() error {
	if r.conn != nil {
		err := r.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Millisecond*20))
		if err != nil {
			slog.Error("failed to close ws", "err", err.Error())
		}
		return r.conn.Close()
	}
	return nil
}
