package metadata

import (
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
	}

	var oids []string
	if len(deviceMetadataMibs) == 0 {
		for el := SNMP_device_metadata_oids.Front(); el != nil; el = el.Next() {
			oids = append(oids, el.Key.(string))
		}
	} else {
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
			switch vt := value.(type) {
			case string:
				md.Customs[oid.Name] = vt
			case []byte:
				md.Customs[oid.Name] = string(vt)
			default:
				md.CustomInts[oid.Name] = snmp_util.ToInt64(value)
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
