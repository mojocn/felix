package cmd

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mojocn/felix/shadowos"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	app = &shadowos.App{
		AddrSocks5: "127.0.0.1:1080",
		Timeout:    time.Second * 60,
	}
	userUUID = "53881505-c10c-464a-8949-e57184a576a9"
	url      = "ws://demo.libragen.cn/5sdfasdf"
)

var socks5Cmd = &cobra.Command{
	Use:   "socks5",
	Short: "socks5 over websocket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
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
		}
		app.Run(ctx, cfg)
	},
}

func init() {
	rootCmd.AddCommand(socks5Cmd)
}
