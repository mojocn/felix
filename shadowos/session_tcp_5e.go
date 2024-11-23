package shadowos

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
)

type SessionTcp5e struct {
	req      *Socks5Request
	s5       net.Conn
	ws       *websocket.Conn
	wsExitCh chan struct{}
}

func (st *SessionTcp5e) Logger() *slog.Logger {
	return st.req.Logger()
}

func (st *SessionTcp5e) breakGfwSvr(ctx context.Context, proxy *ProxyCfg) error {
	ws, err := webSocketConn(ctx, proxy, st.req)
	if err != nil {
		return err
	}
	ws.SetCloseHandler(func(code int, text string) error {
		st.Logger().Debug("ws has closed", "code", code, "text", text)
		st.wsExitCh <- struct{}{}
		return nil
	})
	st.ws = ws
	_, err = st.s5.Write([]byte{socks5Version, socks5ReplySuccess, socks5ReplyReserved, 0x01, 0, 0, 0, 0, 0, 0})
	return nil
}

func (st *SessionTcp5e) Close() {
	//s5 has already been closed in outside
	if ws := st.ws; ws != nil {
		span := st.Logger()
		err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			span.Debug("send websocket close message failed: ", err)
		}
		err = ws.Close()
		if err != nil {
			span.Debug("close websocket conn failed: ", err)
		}
	}
}

func (st *SessionTcp5e) pipe(ctx context.Context, _ [16]byte) {
	ws := st.ws
	s5 := st.s5

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		span := st.Logger().With("fn", "ws -> s5")
		defer func() {
			span.Debug("wg done")
			wg.Done()
		}()
		for {
			select {
			case <-st.wsExitCh:
				span.Info("exitWs")
				return
			case <-ctx.Done():
				span.Info("ctx.Done exit")
				return
			default:
				//ws.SetReadDeadline(time.Now().Add(1 * time.Second))
				_, data, err := ws.ReadMessage()
				n := len(data)
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
					span.Info("EOF from ws")
					return
				}
				if err != nil {
					span.Error("other ws -> socks5 error %T", err)
					span.Debug("other ws -> socks5 error", err)
					return
				}
				fromByteIndex := 0
				//log.Println("write back socks5", n)
				//s5.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
				_, err = s5.Write(data[fromByteIndex:n])
				if err != nil {
					span.Error(" ws -> socks5 error", err)
					return
				}
			}
		}
	}()
	go func() { // s5 -> ws
		span := st.Logger().With("fn", "s5 -> ws")
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
				if errors.Is(err, os.ErrDeadlineExceeded) {
					continue
				}
				if errors.Is(err, io.EOF) {
					span.Debug("EOF from socks5")
					return
				}
				var opErr *net.OpError
				if err != nil && errors.As(err, &opErr) {
					span.Error("net.OpError", "err", opErr)
					return
				}

				if err != nil {
					et := fmt.Sprintf("%T", err)
					span.With("errType", et).Error("s5 read", err)
					continue
				}
				span.Debug("s5 read", "n", n)

				//ws.SetWriteDeadline(time.Now().Add(1 * time.Second))
				err = ws.WriteMessage(websocket.BinaryMessage, buf[:n])
				if err != nil {
					span.Debug("write error", err)
					return
				}
			}
		}
	}()
	wg.Wait()
	st.Logger().Debug("2 goroutines is Done")
}
