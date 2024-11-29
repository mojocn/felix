package rsver

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

const buffSize = 8 << 10

var upGrader = websocket.Upgrader{
	ReadBufferSize:  buffSize,
	WriteBufferSize: buffSize,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		return true
	},
}

func (a *App) wsS5(w http.ResponseWriter, r *http.Request) {
	//authorization := r.Header.Get("authorization")
	//reqID := r.Header.Get("x-request-id")
	dstNetwork := r.Header.Get("x-dst-network")
	if dstNetwork == "" {
		dstNetwork = "tcp"
	}
	dstAddr := r.Header.Get("x-dst-addr")
	dstPort := r.Header.Get("x-dst-port")
	if dstPort == "" {
		dstPort = "0"
	}
	//dstVersion := r.Header.Get("x-dst-vlessResponseHeader")
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	tcpConn, err := net.Dial(dstNetwork, net.JoinHostPort(dstAddr, dstPort))
	if err != nil {
		log.Println("Error connecting to destination:", err)
		return
	}
	defer tcpConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			if mt != websocket.BinaryMessage {
				continue
			}
			_, err = tcpConn.Write(message)
			if err != nil {
				log.Println("Error writing to TCP connection:", err)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			buf := make([]byte, buffSize)
			n, err := tcpConn.Read(buf)
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				log.Println("Error reading from TCP connection:", err)
				return
			}
			err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				log.Println("Error writing to websocket:", err)
				return
			}
		}
	}()
	wg.Wait()
}
