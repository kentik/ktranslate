package metadata

import (
	"context"
	"net"
	"sort"
	"unicode/utf8"

	"github.com/elliotchance/orderedmap"

	"github.com/gosnmp/gosnmp"
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
	SNMP_engineID    = "sysEngineID"

	SNMP_device_metadata_oids = func() *orderedmap.OrderedMap {
		m := orderedmap.NewOrderedMap()
		m.Set(".1.3.6.1.2.1.1.1.0", SNMP_sysDescr)
		m.Set(".1.3.6.1.2.1.1.2.0", SNMP_sysObjectID)
		m.Set(".1.3.6.1.2.1.1.4.0", SNMP_sysContact)
		m.Set(".1.3.6.1.2.1.1.5.0", SNMP_sysName)
		m.Set(".1.3.6.1.2.1.1.6.0", SNMP_sysLocation)
		m.Set(".1.3.6.1.2.1.1.7.0", SNMP_sysServices)
		m.Set(".1.3.6.1.6.3.10.2.1.1.0", SNMP_engineID)
		return m
	}()
)

type DeviceMetadata struct {
	log     logger.ContextL
	mibs    map[string]*kt.Mib
	conf    *kt.SnmpDeviceConfig
	metrics *kt.SnmpDeviceMetric
	missing map[string]bool
	basic   map[string]string
}

func NewDeviceMetadata(deviceMetadataMibs map[string]*kt.Mib, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, log logger.ContextL) *DeviceMetadata {
	mibs := map[string]*kt.Mib{}
	basic := map[string]string{}
	for el := SNMP_device_metadata_oids.Front(); el != nil; el = el.Next() {
		basic[el.Key.(string)] = el.Value.(string)
	}

	if len(deviceMetadataMibs) > 0 {
		oids := getFromCustomMap(deviceMetadataMibs)
		for _, oid := range oids {
			mib := deviceMetadataMibs[oid]
			mibs[oid] = mib
			log.Infof("Adding custom device metadata oid: %s -> %s %v", oid, mib.GetName(), mib.OtherTables)
		}
	} else {
		// Copy the global values into this map which is per device.
		for el := SNMP_device_metadata_oids.Front(); el != nil; el = el.Next() {
			mibs[el.Key.(string)] = nil
		}
	}

	return &DeviceMetadata{
		log:     log,
		mibs:    mibs,
		conf:    conf,
		basic:   basic,
		metrics: metrics,
		missing: map[string]bool{},
	}
}

// Poll device-level metadata.
func (dm *DeviceMetadata) Poll(ctx context.Context, server *gosnmp.GoSNMP) (*kt.DeviceMetricsMetadata, error) {
	return dm.poll(ctx, server)
}

type wrapper struct {
	variable gosnmp.SnmpPDU
	mib      *kt.Mib
	oid      string
}

func (dm *DeviceMetadata) poll(ctx context.Context, server *gosnmp.GoSNMP) (*kt.DeviceMetricsMetadata, error) {
	var results []wrapper
	md := kt.DeviceMetricsMetadata{
		Customs:    map[string]string{},
		CustomInts: map[string]int64{},
		Tables:     map[string]kt.DeviceTableMetadata{},
	}

	missing := int64(0)
	for oid, mib := range dm.mibs {
		if !mib.IsPollReady() { // Skip this mib because its time to poll hasn't elapsed yet.
			continue
		}
		oidResults, err := snmp_util.WalkOID(ctx, dm.conf, oid, server, dm.log, "CustomDeviceMetadata")
		if err != nil {
			dm.metrics.Errors.Mark(1)
			continue
		}

		if len(oidResults) == 0 {
			missing++
			if _, ok := dm.missing[oid]; ok {
				dm.log.Debugf("OID %s failed to return results, Metric Name: %s", oid, mib.GetName())
			} else {
				dm.missing[oid] = true
				dm.log.Warnf("OID %s failed to return results, Metric Name: %s", oid, mib.GetName())
			}
		}
		for _, result := range oidResults {
			results = append(results, wrapper{variable: result, mib: mib, oid: oid})
		}
	}

	// Update the number of missing metrics metric here.
	dm.metrics.MissingMeta.Update(missing)

	// Map back into types we know about.
	for _, wrapper := range results {
		if wrapper.variable.Value == nil { // You can get nil w/out getting an error, though.
			continue
		}

		// Calculate the index first out here.
		idx := snmp_util.GetIndex(wrapper.variable.Name[1:], wrapper.oid)

		oidName := ""
		if name, ok := dm.basic[wrapper.variable.Name]; ok {
			oidName = name
		} else if wrapper.mib != nil {
			oidName = wrapper.mib.GetName()
		}
		if oidName == "" {
			dm.log.Warnf("Missing metadata name: %s", wrapper.oid)
			continue
		}

		value := wrapper.variable.Value
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
		case SNMP_engineID:
			_, md.EngineID = snmp_util.GetFromConv(wrapper.variable, snmp_util.CONV_ENGINE_ID, dm.log)
		default:
			// Now we're actually in the range of custom fields.
			if wrapper.mib == nil { // This should never happen here.
				dm.log.Warnf("Missing Custom metadata oid: %+v, Value: %T %+v", wrapper.variable, wrapper.variable.Value, wrapper.variable.Value)
				continue
			}

			if idx != "" {
				dm.handleTable(idx, wrapper, oidName, &md)
			} else {
				switch vt := value.(type) {
				case string:
					md.Customs[oidName] = vt
				case []byte:
					if wrapper.mib.Conversion != "" { // Adjust for any hard coded values here.
						ival, sval := snmp_util.GetFromConv(wrapper.variable, wrapper.mib.Conversion, dm.log)
						if ival > 0 {
							md.CustomInts[oidName] = ival
							md.Customs[kt.StringPrefix+oidName] = sval
						} else {
							md.Customs[oidName] = sval
						}
					} else {
						md.Customs[oidName] = string(vt)
					}
				case net.IP:
					md.Customs[oidName] = vt.String()
				default:
					val := snmp_util.ToInt64(value)
					if wrapper.mib.EnumRev != nil {
						if ev, ok := wrapper.mib.EnumRev[val]; ok {
							md.Customs[oidName] = ev
						} else {
							md.CustomInts[oidName] = val
						}
					} else {
						md.CustomInts[oidName] = val
					}
				}
			}
		}
	}

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

func (dm *DeviceMetadata) handleTable(idx string, value wrapper, oidName string, md *kt.DeviceMetricsMetadata) {
	if idx[0:1] == "." {
		idx = idx[1:]
	}
	if _, ok := md.Tables[idx]; !ok {
		md.Tables[idx] = kt.NewDeviceTableMetadata()
	}
	switch value.variable.Type {
	case gosnmp.OctetString:
		val := string(value.variable.Value.([]byte))
		if value.mib.Conversion != "" { // Adjust for any hard coded values here.
			_, val = snmp_util.GetFromConv(value.variable, value.mib.Conversion, dm.log)
		}
		if utf8.ValidString(val) {
			md.Tables[idx].Customs[oidName] = kt.NewMetaValue(value.mib, val, 0)
		}
	case gosnmp.IPAddress: // Does this work?
		switch val := value.variable.Value.(type) {
		case string:
			md.Tables[idx].Customs[oidName] = kt.NewMetaValue(value.mib, val, 0)
		case []byte:
			md.Tables[idx].Customs[oidName] = kt.NewMetaValue(value.mib, string(val), 0)
		case net.IP:
			md.Tables[idx].Customs[oidName] = kt.NewMetaValue(value.mib, val.String(), 0)
		}
	case gosnmp.ObjectIdentifier:
		val := string(value.variable.Value.(string))
		if value.mib.Conversion != "" { // Adjust for any hard coded values here.
			_, val = snmp_util.GetFromConv(value.variable, value.mib.Conversion, dm.log)
		}
		if utf8.ValidString(val) {
			md.Tables[idx].Customs[oidName] = kt.NewMetaValue(value.mib, val, 0)
		}
	default:
		// Try to just use as a number, either via an enum or directly.
		val := gosnmp.ToBigInt(value.variable.Value).Int64()
		if value.mib != nil && value.mib.EnumRev != nil {
			if ev, ok := value.mib.EnumRev[val]; ok {
				md.Tables[idx].Customs[oidName] = kt.NewMetaValue(value.mib, ev, 0)
			} else {
				md.Tables[idx].Customs[oidName] = kt.NewMetaValue(value.mib, "", val)
			}
		} else {
			md.Tables[idx].Customs[oidName] = kt.NewMetaValue(value.mib, "", val)
		}
	}
}

// Super basic loop to get info for discovery.
func GetBasicDeviceMetadata(log logger.ContextL, server *gosnmp.GoSNMP) (*kt.DeviceMetricsMetadata, error) {
	md := kt.DeviceMetricsMetadata{}

	var oids []string
	for el := SNMP_device_metadata_oids.Front(); el != nil; el = el.Next() {
		oids = append(oids, el.Key.(string))
	}

	result, err := server.Get(oids)
	if err != nil {
		return nil, err
	}

	for _, pdu := range result.Variables {
		log.Debugf("pdu: %+v", pdu)
		oidVal, value := pdu.Name, pdu.Value

		// You can get a nil value w/out getting an error.
		if value == nil || pdu.Type == gosnmp.NoSuchObject {
			continue
		}

		var oidName string
		thing, ok := SNMP_device_metadata_oids.Get(oidVal)
		if !ok {
			if oidVal == ".1.3.6.1.6.3.15.1.1.3.0" { // This is a bad v3 config.
				log.Errorf("User found who is not known to the SNMP engine. Likely this is an invalid v3 config.")
			} else {
				log.Errorf("SNMP Device Metadata: Unknown oid retrieved: %v %v", oidVal, value)
			}
			continue
		}
		oidName = thing.(string)

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
		case SNMP_engineID:
			_, md.EngineID = snmp_util.GetFromConv(pdu, snmp_util.CONV_ENGINE_ID, log)
		}
	}

	return &md, nil
}
