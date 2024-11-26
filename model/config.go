package model

import "log/slog"

type Config struct {
	PortSocks5 int    `json:"port_socks5"`
	PortHttp   int    `json:"port_http"`
	AuthUser   string `json:"auth_user"`
	AuthPass   string `json:"auth_pass"`
}

var (
	cfg    *Config
	defCfg = Config{
		PortSocks5: 1080,
		PortHttp:   1080 + 5,
		AuthUser:   "admin",
		AuthPass:   "admin",
	}
)

func Cfg() *Config {
	if cfg == nil {
		row := new(Meta)
		row.Config = defCfg
		err := db.FirstOrCreate(&row).Error
		if err != nil {
			slog.Error("get config error", "err", err)
		} else {
			cfg = &defCfg
		}
	}
	return cfg
}
