package patricia

import (
	"net"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/oschwald/geoip2-golang"
)

type MMMap struct {
	mmdb *geoip2.Reader
	log  logger.ContextL
}

func NewMapFromMM(file string, log logger.ContextL) (*MMMap, error) {
	t := &MMMap{log: log}
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

// Utility to turn Geo -> packed int
func PackGeo(cnty []byte) uint32 {
	if len(cnty) == 2 {
		return uint32(rune((cnty)[0])*256 + rune((cnty)[1]))
	}
	return 0
}
