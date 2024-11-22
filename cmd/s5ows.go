package cmd

import (
	"github.com/mojocn/felix/shadowos"
	"github.com/sirupsen/logrus"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var (
	url = "ws://127.0.0.1:8787/afasdfasdf"
	app = &shadowos.App{
		AddrWs:     url,
		AddrSocks5: "127.0.0.1:1080",
		UUID:       "53881505-c10c-464a-8949-e57184a576a9",
		Timeout:    time.Second * 60,
	}
)

var vlessClient = &cobra.Command{
	Use:   "vless",
	Short: "socks5 over websocket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		logrus.SetReportCaller(true)
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(vlessClient)
}
