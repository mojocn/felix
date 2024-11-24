package shadowos

import (
	"testing"
)

func TestGeoDns_country(t *testing.T) {
	g := NewGeoDns("GeoLite2-Country.mmdb")
	defer g.Close()
	country, err := g.country("baidu.com")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(country)
	}

}
