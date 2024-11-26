package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mojocn/felix/api"
	"github.com/mojocn/felix/model"
	"github.com/mojocn/felix/socks5ws"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	buildTime, gitHash string
	userUUID           = "53881505-c10c-464a-8949-e57184a576a9"
	url                = "ws://demo.libragen.cn/5sdfasdf"
	protocol           = "socks5e" // or vless
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	model.DB()
	appCfg := model.Cfg()

	app, err := socks5ws.NewApp(fmt.Sprintf("127.0.0.1:%d", appCfg.PortSocks5), "GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	slog.With("socks5", app.AddrSocks5).Info("socks5 server listening on")

	uid, err := uuid.Parse(userUUID)
	if err != nil {
		log.Fatal("invalid uuid", userUUID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	httpS := api.ApiServer(fmt.Sprintf("127.0.0.1:%d", appCfg.PortHttp))
	go func() {
		if err := httpS.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		sig := <-signalChan
		fmt.Printf("\nReceived signal: %s\n", sig)
		cancel() // Cancel the context

		// Shutdown the server with a timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := httpS.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
	}()
	cfg := &socks5ws.ProxyCfg{
		WsUrl:    url,
		WsHeader: http.Header{},
		UUID:     uid,
		Protocol: protocol,
	}
	log.Println("using:", cfg.Protocol)
	app.Run(ctx, cfg)
}
