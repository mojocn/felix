package shadowos

import (
	"encoding/hex"
	"net/http"
)

type ProxyCfg struct {
	WsUrl    string
	WsHeader http.Header
	UUID     [16]byte
	Protocol string // vless or socks5e
}

func (c ProxyCfg) uuidHex() string {
	return hex.EncodeToString(c.UUID[:])
}
