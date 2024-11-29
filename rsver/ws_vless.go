package rsver

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mojocn/felix/util"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func startDstConnection(vd *util.SchemaVLESS, timeout time.Duration) (net.Conn, []byte, error) {
	conn, err := net.DialTimeout(vd.DstProtocol, vd.HostPort(), timeout)
	if err != nil {
		return nil, nil, fmt.Errorf("connecting to destination: %w", err)
	}
	return conn, []byte{vd.Version, 0x00}, nil
}

func (a *App) wsVless(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	earlyDataHeader := r.Header.Get("sec-websocket-protocol")
	earlyData, err := base64.RawURLEncoding.DecodeString(earlyDataHeader)
	if err != nil {
		log.Println("Error decoding early data:", err)
	}

	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	defer ws.Close()

	if len(earlyData) == 0 {
		mt, p, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		if mt == websocket.BinaryMessage {
			earlyData = p
		}
	}

	vData, err := util.VlessParse(earlyData)
	if err != nil {
		log.Println("Error parsing vless data:", err)
		return
	}
	if a.IsUserNotAllowed(vData.UUID()) {
		return
	}
	if vData.DstProtocol == "udp" {
		vlessUDP(ctx, vData, ws)
	} else if vData.DstProtocol == "tcp" {
		vlessTCP(ctx, vData, ws)
	}
}

func vlessTCP(_ context.Context, sv *util.SchemaVLESS, ws *websocket.Conn) {
	logger := sv.Logger()
	conn, headerVLESS, err := startDstConnection(sv, time.Millisecond*1000)
	if err != nil {
		logger.Error("Error starting session:", "err", err)
		return
	}
	defer conn.Close()
	logger.Info("Session started tcp")

	//write early data
	_, err = conn.Write(sv.DataTcp())
	if err != nil {
		logger.Error("Error writing early data to TCP connection:", "err", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				logger.Error("Error reading message:", "err", err)
				return
			}
			if mt != websocket.BinaryMessage {
				continue
			}
			_, err = conn.Write(message)
			if err != nil {
				logger.Error("Error writing to TCP connection:", "err", err)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		hasNotSentHeader := true
		for {
			buf := make([]byte, buffSize)
			n, err := conn.Read(buf)
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				logger.Error("Error reading from TCP connection:", "err", err)
				return
			}
			// send header data only for the first time
			if hasNotSentHeader {
				hasNotSentHeader = false
				buf = append(headerVLESS, buf[:n]...)
			} else {
				buf = buf[:n]
			}
			err = ws.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				logger.Error("Error writing to websocket:", "err", err)
				return
			}
		}
	}()
	wg.Wait()
}

func vlessUDP(_ context.Context, sv *util.SchemaVLESS, ws *websocket.Conn) {
	logger := sv.Logger()
	conn, headerVLESS, err := startDstConnection(sv, time.Millisecond*1000)
	if err != nil {
		logger.Error("Error starting session:", "err", err)
		return
	}
	defer conn.Close()

	//write early data
	_, err = conn.Write(sv.DataUdp())
	if err != nil {
		logger.Error("Error writing early data to TCP connection:", "err", err)
		return
	}

	buf := make([]byte, buffSize)
	n, err := conn.Read(buf)
	if err != nil {
		logger.Error("Error reading from TCP connection:", "err", err)
		return
	}
	udpDataLen1 := (n >> 8) & 0xff
	udpDataLen2 := n & 0xff
	headerVLESS = append(headerVLESS, byte(udpDataLen1), byte(udpDataLen2))
	headerVLESS = append(headerVLESS, buf[:n]...)

	err = ws.WriteMessage(websocket.BinaryMessage, headerVLESS)
	if err != nil {
		logger.Error("Error writing to websocket:", "err", err)
		return
	}
}
