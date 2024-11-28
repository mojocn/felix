package util

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
)

type GeoIP struct {
	db *geoip2.Reader
}

func NewGeoIP(geoIpFilePath string) (*GeoIP, error) {
	if geoIpFilePath == "" {
		geoIpFilePath = "GeoLite2-Country.mmdb" //https://github.com/P3TERX/GeoLite.mmdb?tab=readme-ov-file
	}
	db, err := geoip2.Open(geoIpFilePath)
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

func (g *GeoIP) Country(host string) (isoCountryCode string, err error) {
	if g == nil {
		return "", fmt.Errorf("geo databse is nil")
	}
	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return "", fmt.Errorf("failed to lookup IP: %w", err)
		}
		if len(ips) == 0 {
			return "", fmt.Errorf("no IP found for %s", host)
		}
		ip = ips[0]
	}
	record, err := g.db.Country(ip)
	if err != nil {
		return "", fmt.Errorf("failed to get Country: %w", err)
	}
	return record.Country.IsoCode, nil
}
