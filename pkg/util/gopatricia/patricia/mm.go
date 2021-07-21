package patricia

import (
	"net"

	"github.com/kentik/golog/logger"
	"github.com/oschwald/geoip2-golang"
)

type MMMap struct {
	mmdb *geoip2.Reader
}

func NewMapFromMM(file string, log *logger.Logger) (*MMMap, error) {
	t := &MMMap{}
	db, err := geoip2.Open(file)
	if err != nil {
		return nil, err
	}
	t.mmdb = db
	return t, nil
}

func (t *MMMap) Close() {
	if t.mmdb != nil {
		t.mmdb.Close()
	}
}

func (t *MMMap) SearchBestFromHostGeo(ip net.IP) (string, error) {
	entry, err := t.mmdb.Country(ip)
	if err != nil {
		return "", err
	}
	return entry.Country.IsoCode, nil
}

func (t *MMMap) SearchBestFromHostAsn(ip net.IP) (uint32, string, error) {
	entry, err := t.mmdb.ASN(ip)
	if err != nil {
		return 0, "", err
	}
	return uint32(entry.AutonomousSystemNumber), entry.AutonomousSystemOrganization, nil
}
