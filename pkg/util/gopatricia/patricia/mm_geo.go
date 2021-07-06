package patricia

import (
	"net"

	"github.com/kentik/golog/logger"
	"github.com/oschwald/geoip2-golang"
)

type MMGeo struct {
	mmdb *geoip2.Reader
}

func NewGeoFromMM(file string, log *logger.Logger) (*MMGeo, error) {
	t := &MMGeo{}
	dbCountry, err := geoip2.Open(file)
	if err != nil {
		return nil, err
	}
	t.mmdb = dbCountry
	return t, nil
}

func (t *MMGeo) Close() {
	if t.mmdb != nil {
		t.mmdb.Close()
	}
}

func (t *MMGeo) SearchBestFromHostGeo(ip net.IP) (string, error) {
	entry, err := t.mmdb.Country(ip)
	if err != nil {
		return "", err
	}
	return entry.Country.IsoCode, nil
}
