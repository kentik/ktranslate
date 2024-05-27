package util

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/gosnmp/gosnmp"
)

// Load a snmpwalk response from the given file and use it for debugging.
func LoadFromWalk(ctx context.Context, file string, log logger.ContextL) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	res := map[string]gosnmp.SnmpPDU{}
	for _, line := range strings.Split(string(data), "\n") {
		pts := strings.SplitN(line, " = ", 2)
		if len(pts) != 2 {
			continue
		}
		oid := strings.TrimSpace(pts[0])
		if oid[0:1] == "." { // Strip off leading dot.
			oid = oid[1:]
		}
		pp := strings.SplitN(pts[1], ": ", 2)
		if len(pp) != 2 {
			continue
		}
		val := pp[1]

		switch strings.ToLower(pp[0]) {
		case "string":
			res[oid] = gosnmp.SnmpPDU{Value: []byte(val), Type: gosnmp.OctetString, Name: oid}
		case "integer":
			res[oid] = gosnmp.SnmpPDU{Value: val, Type: gosnmp.Integer, Name: oid}
		case "counter64":
			res[oid] = gosnmp.SnmpPDU{Value: val, Type: gosnmp.Counter64, Name: oid}
		case "counter32":
			res[oid] = gosnmp.SnmpPDU{Value: val, Type: gosnmp.Counter32, Name: oid}
		default:
			log.Errorf("Skipping unknown walk type: %s", pp[0])
		}
	}

	walkCacheMap = res
	log.Infof("Loaded up %d entries to the walk cache map", len(res))

	return nil
}

func useCachedMap(oid string) ([]gosnmp.SnmpPDU, error) {
	res := []gosnmp.SnmpPDU{}

	if oid[0:1] == "." { // Strip off leading dot.
		oid = oid[1:]
	}
	for k, v := range walkCacheMap {
		if strings.HasPrefix(k, oid) {
			res = append(res, v)
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("Nothing to see")
	}
	return res, nil
}
