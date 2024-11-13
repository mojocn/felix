package cmd

import (
	"github.com/mojocn/felix/shadowos"
	"github.com/sirupsen/logrus"
	"log"

	"github.com/spf13/cobra"
)

var (
	url = "wss://demo.libragen.cn/53881505-c10c-464a-8949-e57184a576a9"
	app = &shadowos.ShadowosApp{
		AddrWs:     url,
		AddrSocks5: "127.0.0.1:1080",
		UUID:       "53881505-c10c-464a-8949-e57184a576a9",
	}
)

var s5OverWsCmd = &cobra.Command{
	Use:   "s5ows",
	Short: "socks5 over websocket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		logrus.SetReportCaller(true)
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(s5OverWsCmd)
}
