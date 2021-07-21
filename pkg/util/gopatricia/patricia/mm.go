package patricia

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/kentik/golog/logger"
	"github.com/oschwald/geoip2-golang"
)

type MMMap struct {
	mmdb *geoip2.Reader
	mapr map[uint32]string
}

func NewMapFromMM(file string, log *logger.Logger) (*MMMap, error) {
	t := &MMMap{
		mapr: map[uint32]string{},
	}
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
	return uint32(entry.AutonomousSystemNumber), t.mapr[uint32(entry.AutonomousSystemNumber)], nil
}

func (t *MMMap) LoadNames(fileName string) (int, error) {
	found := 0
	file, err := os.Open(fileName)
	if err != nil {
		return found, err
	}
	defer file.Close()

	mapr := map[uint32]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		pts := strings.Fields(line)
		if len(pts) > 0 {
			num, err := strconv.Atoi(pts[0][2:])
			if err != nil {
				return found, fmt.Errorf("Error adding name with line %s -> %v", line, err)
			}
			more := strings.Split(pts[1], ",")
			mapr[uint32(num)] = strings.TrimSpace(more[0])
			found++
		}
	}

	if err := scanner.Err(); err != nil {
		return found, err
	}

	t.mapr = mapr
	return found, nil
}

func (t *MMMap) FromMap(val uint32) string {
	return t.mapr[val]
}
