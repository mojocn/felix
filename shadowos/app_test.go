package shadowos

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"testing"
)

var (
	url = "ws://127.0.0.1:8787/53881505-c10c-464a-8949-e57184a576a9"
	app = &ShadowosApp{
		AddrWs:     url,
		AddrSocks5: "127.0.0.1:1080",
		UUID:       "53881505-c10c-464a-8949-e57184a576a9",
	}
)

func TestShadowosApp_Run(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logrus.SetReportCaller(true)
	app.Run()
}

func TestWsReadMessage(t *testing.T) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal("dial:", err)
	}

	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		t.Log("recv: ", messageType, string(message))
		t.Logf("recv: %s", message)
		log.Printf("recv: %s", message)
	}
	t.Log("test done")
}
