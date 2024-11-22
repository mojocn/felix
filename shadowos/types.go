package shadowos

import "net/http"

type ProxyCfg struct {
	WsUrl    string
	WsHeader http.Header
	UUID     [16]byte
}
