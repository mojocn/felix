package shadowos

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
)

type GeoIP struct {
	db *geoip2.Reader
}

func NewGeoDns(dbFile string) (*GeoIP, error) {
	if dbFile == "" {
		dbFile = "GeoLite2-Country.mmdb" //https://github.com/P3TERX/GeoLite.mmdb?tab=readme-ov-file
	}
	db, err := geoip2.Open(dbFile)
	if err != nil {
		return nil, err
	}
	return &GeoIP{db: db}, nil
}

func (g *GeoIP) Close() error {
	if g.db != nil {
		return g.db.Close()
	}
	return nil
}

func (g *GeoIP) country(ip, domain []byte) (isoCountryCode string, err error) {
	if g == nil {
		return "", fmt.Errorf("geo databse is nil")
	}
	if len(domain) > 0 {
		ips, err := net.LookupIP(string(domain))
		if err != nil {
			return "", err
		}
		if len(ips) == 0 {
			return "", fmt.Errorf("no IP found for %s", domain)
		}
		record, err := g.db.Country(ips[0])
		if err != nil {
			return "", err
		}
		return record.Country.IsoCode, nil
	}
	if len(ip) > 0 {
		record, err := g.db.Country(ip)
		if err != nil {
			return "", err
		}
		return record.Country.IsoCode, nil
	}
	return "", fmt.Errorf("no ip or domain")
}
