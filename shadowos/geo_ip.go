package shadowos

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
)

type GeoDns struct {
	db *geoip2.Reader
}

func NewGeoDns(dbFile string) *GeoDns {
	db, err := geoip2.Open("GeoLite2-Country.mmdb") //https://github.com/P3TERX/GeoLite.mmdb?tab=readme-ov-file
	if err != nil {
		panic(fmt.Sprintf("Failed to open GeoIP database: %v", err))
	}
	return &GeoDns{db: db}
}

func (g *GeoDns) Close() error {
	if g.db != nil {
		return g.db.Close()
	}
	return nil
}

func (g *GeoDns) country(domainOrIp string) (string, error) {
	db, err := geoip2.Open("GeoLite2-Country.mmdb") //https://github.com/P3TERX/GeoLite.mmdb?tab=readme-ov-file
	if err != nil {
		return "", err
	}
	defer db.Close()
	ip := net.ParseIP(domainOrIp)
	if ip == nil {
		ips, err := net.LookupIP(domainOrIp)
		if err != nil {
			return "", err
		}
		ip = ips[0]
	}
	record, err := db.Country(ip)
	if err != nil {
		return "", err
	}
	return record.Country.IsoCode, nil
}
