package metadata

import (
	"net"
	"sort"

	"github.com/elliotchance/orderedmap"

	"github.com/kentik/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	SNMP_sysDescr    = "sysDescr"
	SNMP_sysObjectID = "sysObjectID"
	SNMP_sysContact  = "sysContact"
	SNMP_sysName     = "sysName"
	SNMP_sysLocation = "sysLocation"
	SNMP_sysServices = "sysServices"

	SNMP_device_metadata_oids = func() *orderedmap.OrderedMap {
		m := orderedmap.NewOrderedMap()
		m.Set(".1.3.6.1.2.1.1.1.0", SNMP_sysDescr)
		m.Set(".1.3.6.1.2.1.1.2.0", SNMP_sysObjectID)
		m.Set(".1.3.6.1.2.1.1.4.0", SNMP_sysContact)
		m.Set(".1.3.6.1.2.1.1.5.0", SNMP_sysName)
		m.Set(".1.3.6.1.2.1.1.6.0", SNMP_sysLocation)
		m.Set(".1.3.6.1.2.1.1.7.0", SNMP_sysServices)
		return m
	}()
)

// Poll device-level metadata, which we only need once(?).  Works for (at least) Juniper, Cisco, and Arista.
func GetDeviceMetadata(log logger.ContextL, server *gosnmp.GoSNMP, deviceMetadataMibs map[string]*kt.Mib) (*kt.DeviceMetricsMetadata, error) {
	md := kt.DeviceMetricsMetadata{
		Customs:    map[string]string{},
		CustomInts: map[string]int64{},
		Tables:     map[string]kt.DeviceTableMetadata{},
	}

	var oids []string
	if len(deviceMetadataMibs) == 0 {
		for el := SNMP_device_metadata_oids.Front(); el != nil; el = el.Next() {
			oids = append(oids, el.Key.(string))
		}
	} else {
		log.Infof("Getting device metadata from custom map: %v", deviceMetadataMibs)
		oids = getFromCustomMap(deviceMetadataMibs)
	}

	result, err := server.Get(oids)
	if err != nil {
		return nil, err
	}

	hasData := false
	for _, pdu := range result.Variables {
		log.Debugf("pdu: %+v", pdu)

		oidVal, value := pdu.Name, pdu.Value

		// You can get a nil value w/out getting an error.
		if value == nil || pdu.Type == gosnmp.NoSuchObject {
			if oidInfo, ok := deviceMetadataMibs[oidVal[1:]]; ok {
				log.Infof("Trying to walk %s -> %s as a table", oidInfo.Name, oidVal)
				err := getTable(log, server, oidVal, oidInfo, &md)
				if err != nil {
					log.Warnf("Dropping %s because of nil value or missing object: %+v", oidVal, pdu)
				}
			}
			continue
		}

		var oidName string
		oid, ok := deviceMetadataMibs[oidVal[1:]]
		if !ok {
			thing, ok := SNMP_device_metadata_oids.Get(oidVal)
			if !ok {
				log.Errorf("SNMP Device Metadata: Unknown oid retrieved: %v", oidVal)
				continue
			}
			oidName = thing.(string)
		} else {
			oidName = oid.Name
		}

		hasData = true
		switch oidName {
		case SNMP_sysDescr:
			md.SysDescr = string(value.([]byte))
		case SNMP_sysObjectID:
			md.SysObjectID = value.(string)
		case SNMP_sysContact:
			md.SysContact = string(value.([]byte))
		case SNMP_sysName:
			md.SysName = string(value.([]byte))
		case SNMP_sysLocation:
			md.SysLocation = string(value.([]byte))
		case SNMP_sysServices:
			md.SysServices = int(snmp_util.ToInt64(value))
		default:
			if oid.Tag != "" {
				oidName = oid.Tag
			}
			switch vt := value.(type) {
			case string:
				md.Customs[oidName] = vt
			case []byte:
				md.Customs[oidName] = string(vt)
			case net.IP:
				md.Customs[oidName] = vt.String()
			default:
				md.CustomInts[oidName] = snmp_util.ToInt64(value)
			}
		}
	}

	// If no fields in md were set, return nil.  (Trust me on the (). :)
	if !hasData {
		log.Infof("SNMP Device Metadata: No data received")
		return nil, nil
	}
	log.Infof("SNMP Device Metadata: Data received: %+v", md)

	return &md, nil
}

func getFromCustomMap(mibs map[string]*kt.Mib) []string {
	keys := []string{}
	for k, _ := range mibs {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func getTable(log logger.ContextL, g *gosnmp.GoSNMP, oid string, mib *kt.Mib, md *kt.DeviceMetricsMetadata) error {
	results, err := g.WalkAll(oid)
	if err != nil {
		return err
	}

	oidName := mib.Name
	if mib.Tag != "" {
		oidName = mib.Tag
	}

	log.Infof("TableWalk Results: %s: %s -> %d", oidName, oid, len(results))
	// Save as index -> oid -> value
	for _, variable := range results {
		if len(variable.Name) <= len(oid)+1 {
			log.Warnf("Skipping invalid table, could not get index: %s -> %s", oid, variable.Name)
			continue
		}
		idx := variable.Name[len(oid)+1:]
		if _, ok := md.Tables[idx]; !ok {
			md.Tables[idx] = kt.NewDeviceTableMetadata()
		}

		switch variable.Type {
		case gosnmp.OctetString:
			value := string(variable.Value.([]byte))
			if mib.Conversion != "" { // Adjust for any hard coded values here.
				value = snmp_util.GetFromConv(variable, mib.Conversion, log)
			}
			md.Tables[idx].Customs[oidName] = value
		case gosnmp.IPAddress: // Does this work?
			switch val := variable.Value.(type) {
			case string:
				md.Tables[idx].Customs[oidName] = val
			case []byte:
				md.Tables[idx].Customs[oidName] = string(val)
			case net.IP:
				md.Tables[idx].Customs[oidName] = val.String()
			}
		default:
			// Try to just use as a number
			md.Tables[idx].CustomInts[oidName] = gosnmp.ToBigInt(variable.Value).Int64()
		}
	}

	return nil
}
