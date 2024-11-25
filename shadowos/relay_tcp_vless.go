package shadowos

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
	"time"
)

var _ RelayTcp = (*RelayTcpVless)(nil)

type RelayTcpVless struct {
	cfg               *ProxyCfg
	req               *Socks5Request
	conn              *websocket.Conn
	hasSentHeader     bool
	hasReceivedHeader bool
}

func NewRelayTcpVless(ctx context.Context, cfg *ProxyCfg, req *Socks5Request) (*RelayTcpVless, error) {
	ws, err := webSocketConn(ctx, cfg, req)
	if err != nil {
		return nil, err
	}
	ws.SetCloseHandler(func(code int, text string) error {
		slog.Debug("ws has closed", "code", code, "text", text)
		return nil
	})
	return &RelayTcpVless{
		cfg:               cfg,
		req:               req,
		conn:              ws,
		hasSentHeader:     false,
		hasReceivedHeader: false,
	}, nil
}

func (r *RelayTcpVless) Close() error {
	if r.conn != nil {
		err := r.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Millisecond*20))
		if err != nil {
			slog.Error("failed to close ws", "err", err.Error())
		}
		return r.conn.Close()
	}
	return nil
}

func (r *RelayTcpVless) Read(data []byte) (n int, err error) {
	if r.conn == nil {
		return 0, nil
	}
	_, p, err := r.conn.ReadMessage()
	if err != nil {
		slog.Error("failed to read ws", "err", err.Error())
	}
	fromByteIndex := 0
	if !r.hasReceivedHeader && len(p) >= 2 {
		r.hasReceivedHeader = true
		extraN := p[1]
		fromByteIndex = 2 + int(extraN)
	}
	return copy(data, p[fromByteIndex:]), err
}

func (r *RelayTcpVless) Write(data []byte) (n int, err error) {
	if r.conn == nil {
		return 0, nil
	}
	//executed only once
	if !r.hasSentHeader {
		r.hasSentHeader = true
		header, err := r.req.vlessHeaderTcp(r.cfg.UUID)
		if err != nil {
			return 0, fmt.Errorf("failed to generate vless header:%w", err)
		}
		data = append(header, data...)
	}

	err = r.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		slog.Error("failed to write ws", "err", err.Error())
	}
	return len(data), err
}
