package main

import (
	"log"

	"github.com/mojocn/felix/shadowos"
	"github.com/sirupsen/logrus"
)

var (
	url = "wss://demo.libragen.cn/53881505-c10c-464a-8949-e57184a576a9"
	app = &shadowos.ShadowosApp{
		AddrWs:     url,
		AddrSocks5: "127.0.0.1:1080",
		UUID:       "53881505-c10c-464a-8949-e57184a576a9",
	}
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logrus.SetReportCaller(true)
	app.Run()
}
