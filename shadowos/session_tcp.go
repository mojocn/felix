package shadowos

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type SessionTcp struct {
	req      *Socks5Request
	s5       net.Conn
	ws       *websocket.Conn
	proxyCfg *ProxyCfg
}

func (st *SessionTcp) connectProxyServer() error {
	wsAddr := st.proxyCfg.AddrWs
	ws, _, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	if err != nil {
		_, err = st.s5.Write([]byte{socksVersion, socks5ReplyFail, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		if err != nil {
			log.Printf("failed to send response to client: %v", err)
		}
		return fmt.Errorf("failed to connect to remote proxy server: %s ,error:%v", wsAddr, err)
	} // Send success response
	_, err = st.s5.Write([]byte{socksVersion, socks5ReplySuccess, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		log.Printf("failed to send response to client: %v", err)
	}
	st.ws = ws
	return nil
}

func (st *SessionTcp) pipe(ctx context.Context) {
	exitWs := make(chan struct{}, 1)
	ws := st.ws
	s5 := st.s5
	firstData, err := st.req.vlessHeader(st.proxyCfg.UUID)
	if err != nil {
		log.Println("failed to generate vless header")
		return
	}
	ws.SetCloseHandler(func(code int, text string) error {
		log.Println("ws closed", code, text)
		exitWs <- struct{}{}
		return nil
	})
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		isFirstReceive := true
		defer func() {
			log.Println("[ddd] ws -> s5")
			wg.Done()
		}()
		for {
			select {
			case <-exitWs:
				log.Println("exitWs")
				return
			case <-ctx.Done():
				log.Println("doneCh: ws -> s5")
				return
			default:
				ws.SetReadDeadline(time.Now().Add(1 * time.Second))
				_, data, err := ws.ReadMessage()
				n := len(data)
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
					log.Println("EOF from ws")
					return
				}
				if err != nil {
					log.Printf("other ws -> socks5 error %T", err)
					log.Println("other ws -> socks5 error", err)
					return
				}
				fromByteIndex := 0
				// skip the first data
				if isFirstReceive && n >= 2 {
					extraN := data[1]
					isFirstReceive = false
					fromByteIndex = 2 + int(extraN)
				}
				log.Println("write back socks5", n)
				s5.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
				_, err = s5.Write(data[fromByteIndex:n])
				if err != nil {
					log.Println(" ws -> socks5 error", err)
					return
				}
			}
		}
	}()
	go func() { // s5 -> ws
		defer func() {
			log.Println("[ddd] s5 -> ws")
			wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				log.Println("doneCh: s5 -> ws")
				return
			default:
				buf := make([]byte, 8<<10)
				s5.SetReadDeadline(time.Now().Add(1 * time.Second))
				n, err := s5.Read(buf)
				if errors.Is(err, os.ErrDeadlineExceeded) {
					continue
				}
				if errors.Is(err, io.EOF) {
					log.Println("EOF from socks5")
					return
				}
				var opErr *net.OpError
				if err != nil && errors.As(err, &opErr) {
					log.Println("opErr", opErr)
					return
				}

				if err != nil {
					log.Printf("read from socks5 error %T", err)
					log.Println("read from socks5 error", err)
					continue
				}
				log.Println("read from socks5", n)
				data := buf[:n]
				if len(firstData) > 0 {
					log.Println("send version header only once")
					data = append(firstData, buf[:n]...)
					firstData = nil
				}
				err = ws.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					log.Println("write error", err)
					return
				}
			}
		}
	}()
	wg.Wait()
	log.Println("2 doneCh")

}

func (st *SessionTcp) Close() {
	//s5 has already been closed in outside
	ws := st.ws
	if ws == nil {
		return
	}
	err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("send websocket close message failed: ", err)
	}
	err = ws.Close()
	if err != nil {
		log.Println("close websocket conn failed: ", err)
	}
}
