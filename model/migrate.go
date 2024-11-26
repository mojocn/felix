package model

import "log"

func migrate() {
	if db == nil {
		log.Print("db is nil")
		return
	}
	for _, m := range []interface{}{&CfIp{}, &Meta{}, &Proxy{}} {
		if err := db.AutoMigrate(m); err != nil {
			log.Print(err)
		}
	}
}
