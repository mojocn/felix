package shadowos

import (
	"encoding/hex"
	"net/http"
)

type ProxyCfg struct {
	WsUrl    string
	WsHeader http.Header
	UUID     [16]byte
}

func (c ProxyCfg) uuidHex() string {
	return hex.EncodeToString(c.UUID[:])
}
