package model

import (
	"fmt"
	"strings"
)

type Proxy struct {
	ModelBase
	Name string `json:"name" gorm:"varchar(255)"`

	Protocol string `json:"protocol" gorm:"varchar(16)"` //ws,wss,http2,tls,http3
	Host     string `json:"host" gorm:"varchar(255)"`
	Uri      string `json:"uri" gorm:"varchar(255)"`
	Sni      string `json:"sni" gorm:"varchar(255)"`
	Version  string `json:"version" gorm:"varchar(16)"` // one socks5

	UserID    string `json:"user_id"`
	Password  string `json:"password"`
	TrafficKb int64  `json:"traffic_kb" gorm:"default:0"`
	SpeedMs   int64  `json:"speed_ms" gorm:"default:0"`
	Status    string `json:"status" gorm:"varchar(16);default:''"` //active, inactive
}

func (p *Proxy) IsActive() bool {
	return p.Status == "active"
}
func (p *Proxy) RelayURL() string {
	switch p.Protocol {
	case "ws":
		return fmt.Sprintf("ws://%s/%s", p.Host, strings.TrimPrefix(p.Uri, "/"))
	case "wss":
		return fmt.Sprintf("wss://%s/%s", p.Host, strings.TrimPrefix(p.Uri, "/"))
	case "tcp+tls":
		return fmt.Sprintf("tcp-tls://%s/%s", p.Host, strings.TrimPrefix(p.Uri, "/"))
	default:
		return ""
	}
}
