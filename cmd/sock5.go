package cmd

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mojocn/felix/shadowos"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	userUUID = "53881505-c10c-464a-8949-e57184a576a9"
	url      = "ws://demo.libragen.cn/5sdfasdf"
	//url = "ws://127.0.0.1:8787/5sdfasdf"
	protocol = "socks5e" // or vless
)

var socks5Cmd = &cobra.Command{
	Use:   "socks5",
	Short: "socks5 over websocket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := shadowos.NewApp("127.0.0.1:1080", "C:\\artwork\\felix\\GeoLite2-Country.mmdb")
		if err != nil {
			log.Fatal(err)
		}
		slog.With("socks5", app.AddrSocks5).Info("socks5 server listening on")
		log.SetFlags(log.Lmicroseconds | log.Lshortfile)
		slog.SetLogLoggerLevel(slog.LevelDebug)

		uid, err := uuid.Parse(userUUID)
		if err != nil {
			log.Fatal("invalid uuid", userUUID)
		}

		ctx, cancel := context.WithCancel(context.Background())
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
		go func() {
			sig := <-signalChan
			fmt.Printf("\nReceived signal: %s\n", sig)
			cancel() // Cancel the context
		}()

		cfg := &shadowos.ProxyCfg{
			WsUrl:    url,
			WsHeader: http.Header{},
			UUID:     uid,
			Protocol: protocol,
		}
		log.Println("using:", cfg.Protocol)
		app.Run(ctx, cfg)
	},
}

func init() {
	rootCmd.AddCommand(socks5Cmd)
}
