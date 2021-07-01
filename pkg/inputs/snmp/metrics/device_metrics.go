package metrics

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
)

type DeviceMetrics struct {
	log     logger.ContextL
	conf    *kt.SnmpDeviceConfig
	gconf   *kt.SnmpGlobalConfig
	metrics *kt.SnmpDeviceMetric
}

func NewDeviceMetrics(gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, profileMetrics map[string]*kt.Mib, log logger.ContextL) *DeviceMetrics {
	if conf.DeviceOids == nil && len(profileMetrics) > 0 {
		conf.DeviceOids = profileMetrics
	} else if len(profileMetrics) > 0 {
		for oid, m := range profileMetrics {
			noid := oid
			if !strings.HasPrefix(noid, ".") {
				noid = "." + noid
			}
			oidName := m.Name
			if m.Tag != "" {
				oidName = m.Tag
			}
			log.Infof("Adding device metric %s -> %s", noid, oidName)
			conf.DeviceOids[noid] = m
		}
	}

	// These are defined per device in the yaml conf.
	for k, v := range conf.UserTags {
		log.Infof("Adding user tag %s -> %s", k, v)
	}

	return &DeviceMetrics{
		gconf:   gconf,
		log:     log,
		conf:    conf,
		metrics: metrics,
	}
}

type deviceMetricRow struct {
	Error  string
	Uptime int64

	// Custom Specific
	customStr    map[string]string
	customInt    map[string]int32
	customBigInt map[string]int64
}

var (
	sysUpTime = "1.3.6.1.2.1.1.3.0"
)

func (dm *DeviceMetrics) Poll(server *gosnmp.GoSNMP) ([]*kt.JCHF, error) {
	return dm.pollFromConfig(server)
}

func (dm *DeviceMetrics) pollFromConfig(server *gosnmp.GoSNMP) ([]*kt.JCHF, error) {
	var results []gosnmp.SnmpPDU
	m := map[string]*deviceMetricRow{}

	for oid, mib := range dm.conf.DeviceOids {
		oidResults, err := snmp_util.WalkOID(oid, server, dm.log, "CustomDeviceMetrics")
		if err != nil {
			m[fmt.Sprintf("err-%s", mib.Name)] = &deviceMetricRow{
				Error:        fmt.Sprintf("Walking %s: %v", oid, err),
				customStr:    map[string]string{},
				customInt:    map[string]int32{},
				customBigInt: map[string]int64{},
			}
			dm.metrics.Errors.Mark(1)
			continue
		}

		results = append(results, oidResults...)
	}

	// Get uptime manually here.
	var uptime int64
	uptimeResults, err := snmp_util.WalkOID(sysUpTime, server, dm.log, "CustomDeviceMetrics")
	if err == nil {
		// You might think that if err == nil then you definitely got back some
		// results.  Not exactly.  The result might be "No such object", which
		// is not an error, but also not what you're looking for.
		if len(uptimeResults) > 0 {
			uptime = snmp_util.ToInt64(uptimeResults[0].Value)
		}
	} else {
		m["uptime"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", sysUpTime, err),
			customStr:    map[string]string{},
			customInt:    map[string]int32{},
			customBigInt: map[string]int64{},
		}
	}

	// Map back into types we know about.
	metricsFound := map[string]string{"Uptime": sysUpTime}
	for _, variable := range results {
		if variable.Value == nil { // You can get nil w/out getting an error, though.
			continue
		}

		var mib *kt.Mib = nil
		idx := ""
		for oid, m := range dm.conf.DeviceOids {
			if strings.HasPrefix(variable.Name, oid) {
				idx = snmp_util.GetIndex(variable.Name, oid)
				mib = m
				break
			}
		}

		if mib == nil {
			dm.log.Infof("Missing Custom oid: %+v, Value: %T %+v", variable, variable.Value, variable.Value)
			continue
		}
		oidName := mib.Name
		if mib.Tag != "" {
			oidName = mib.Tag
		}

		dmr := assureDeviceMetrics(m, idx)
		metricsFound[oidName] = mib.Oid
		switch variable.Type {
		case gosnmp.OctetString, gosnmp.BitString:
			value := string(variable.Value.([]byte))
			if mib.Enum != nil {
				dmr.customStr[kt.StringPrefix+oidName] = value // Save the string valued field as an attribute.
				if val, ok := mib.Enum[strings.ToLower(value)]; ok {
					dmr.customBigInt[oidName] = val
				} else {
					dm.log.Warnf("Missing enum value for device metric %s %s", oidName, value)
					dmr.customBigInt[oidName] = 0
				}
			} else {
				// Try to parse this as a number. If its not though, just store as a string.
				if s, err := strconv.ParseInt(value, 10, 64); err == nil {
					dmr.customBigInt[oidName] = s
				} else {
					dm.log.Debugf("unable to set string valued metric as numeric: %s %s", oidName, value)
					dmr.customStr[kt.StringPrefix+oidName] = value // Still save this as a string valued field.
					dmr.customBigInt[oidName] = 0
				}
			}
		default:
			dmr.customBigInt[oidName] = snmp_util.ToInt64(variable.Value)
		}
	}

	// Convert to JCFH and pass on.
	flows := make([]*kt.JCHF, 0, len(m))
	for idx, dmr := range m {
		dst := kt.NewJCHF()
		dst.CustomStr = dmr.customStr
		dst.CustomInt = dmr.customInt
		dst.CustomBigInt = dmr.customBigInt
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		dst.Provider = dm.conf.Provider
		dst.CustomBigInt["Uptime"] = uptime
		dst.CustomStr["Error"] = dmr.Error
		dst.CustomStr[kt.IndexVar] = idx
		dst.DeviceName = dm.conf.DeviceName
		dst.SrcAddr = dm.conf.DeviceIP
		dst.Timestamp = time.Now().Unix()
		dst.CustomMetrics = metricsFound // Add this in so that we know what metrics to pull out down the road.
		for k, v := range dm.conf.UserTags {
			dst.CustomStr[k] = v
		}
		flows = append(flows, dst)
	}

	// In this case, we need to send just a blank line with the uptime.
	if len(m) == 0 {
		dst := kt.NewJCHF()
		dst.CustomStr = map[string]string{}
		dst.CustomInt = map[string]int32{}
		dst.CustomBigInt = map[string]int64{}
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		dst.Provider = dm.conf.Provider
		dst.CustomBigInt["Uptime"] = uptime
		dst.DeviceName = dm.conf.DeviceName
		dst.SrcAddr = dm.conf.DeviceIP
		dst.Timestamp = time.Now().Unix()
		dst.CustomMetrics = metricsFound // Add this in so that we know what metrics to pull out down the road.
		for k, v := range dm.conf.UserTags {
			dst.CustomStr[k] = v
		}
		flows = append(flows, dst)
	}

	dm.metrics.DeviceMetrics.Mark(int64(len(flows)))
	return flows, nil
}

func assureDeviceMetrics(m map[string]*deviceMetricRow, index string) *deviceMetricRow {
	dm, ok := m[index]
	if !ok {
		dm = &deviceMetricRow{
			customStr:    map[string]string{},
			customInt:    map[string]int32{},
			customBigInt: map[string]int64{},
		}
		m[index] = dm
	}
	return dm
}
