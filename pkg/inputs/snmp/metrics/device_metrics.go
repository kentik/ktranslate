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
	Error             string
	Component         string
	CPU               int64
	MemoryTotal       int64
	MemoryUsed        int64
	MemoryFree        int64
	MemoryUtilization int64
	Uptime            int64

	//juniper specific
	juniperOperatingDRAMSize int64
	juniperOperatingMemory   int64
	juniperOperatingBuffer   int64 //this looks like the utilization

	//hr aka arista specific
	hrStorageAllocationUnits int64
	hrStorageSize            int64
	hrStorageUsed            int64

	//cisco specific
	entPhysicalIndex string
	// "cpm" stands for "ciscoProcessMIB", I think
	// See http://oid-info.com/get/1.3.6.1.4.1.9.9.109
	cpmCPUTotal5min    int64
	cpmCPUTotal5minRev int64

	//cisco nexus specific
	cpmCPUMemoryUsed int64
	cpmCPUMemoryFree int64

	// Custom Specific
	customStr    map[string]string
	customInt    map[string]int32
	customBigInt map[string]int64
}

var (
	sysUpTime = "1.3.6.1.2.1.1.3.0"
)

func (dm *DeviceMetrics) Poll(server *gosnmp.GoSNMP) ([]*kt.JCHF, error) {

	// If there's mibs passed in, use these directly.
	if len(dm.conf.DeviceOids) > 0 {
		return dm.pollFromConfig(server)
	}

	// Otherwise, do it in a hard coded way.
	// Get the manufacturer
	deviceManufacturer := strings.ToLower(snmp_util.GetDeviceManufacturer(server, dm.log))
	// Query manufacturer-specific oids
	var dmrs []*deviceMetricRow
	switch {
	case snmp_util.ContainsAny(deviceManufacturer, "juniper", "junos"):
		dmrs = dm.getJuniperDeviceMetrics(dm.log, server)
	case snmp_util.ContainsAny(deviceManufacturer, "cisco"):
		dmrs = dm.getCiscoDeviceMetrics(dm.log, server)
	case snmp_util.ContainsAny(deviceManufacturer, "arista"):
		dmrs = dm.getAristaDeviceMetrics(dm.log, server)
	default:
		// Since we don't have any specific metrics here, go ahead and get the very basics with pollfromconfig
		return dm.pollFromConfig(server)
	}
	flows := dm.convertDMToCHF(dmrs)
	dm.metrics.DeviceMetrics.Mark(int64(len(flows)))

	return flows, nil
}

func (dm *DeviceMetrics) convertDMToCHF(dmrs []*deviceMetricRow) []*kt.JCHF {

	// Do something here.
	flows := make([]*kt.JCHF, 0, len(dmrs))
	for _, dmr := range dmrs {
		dst := kt.NewJCHF()
		dst.CustomStr = make(map[string]string)
		dst.CustomInt = make(map[string]int32)
		dst.CustomBigInt = make(map[string]int64)
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		dst.Provider = dm.conf.Provider
		dst.CustomStr["Error"] = dmr.Error
		dst.CustomStr["Component"] = dmr.Component
		dst.CustomBigInt["CPU"] = dmr.CPU
		dst.CustomBigInt["MemoryTotal"] = dmr.MemoryTotal
		dst.CustomBigInt["MemoryUsed"] = dmr.MemoryUsed
		dst.CustomBigInt["MemoryFree"] = dmr.MemoryFree
		dst.CustomBigInt["MemoryUtilization"] = dmr.MemoryUtilization
		dst.CustomBigInt["Uptime"] = dmr.Uptime
		dst.DeviceName = dm.conf.DeviceName
		dst.SrcAddr = dm.conf.DeviceIP
		dst.Timestamp = time.Now().Unix()
		metrics := map[string]string{"CPU": "", "MemoryUtilization": "", "Uptime": "sysUpTime", "MemoryFree": ""}
		for k, v := range dm.conf.UserTags {
			dst.CustomStr[k] = v
		}
		if dst.Provider == kt.ProviderDefault { // Add this to trigger a UI element.
			dst.CustomStr["profile_message"] = kt.DefaultProfileMessage
		}

		if dmr.juniperOperatingDRAMSize > 0 {
			dst.CustomBigInt["juniperOperatingDRAMSize"] = dmr.juniperOperatingDRAMSize
			dst.CustomBigInt["juniperOperatingMemory"] = dmr.juniperOperatingMemory
			dst.CustomBigInt["juniperOperatingBuffer"] = dmr.juniperOperatingBuffer
			metrics["juniperOperatingDRAMSize"] = jnxOperatingDRAMSize
			metrics["juniperOperatingMemory"] = jnxOperatingMemory
			metrics["juniperOperatingBuffer"] = jnxOperatingMemory
			metrics["CPU"] = jnxOperatingCPU
		}

		if dmr.hrStorageAllocationUnits > 0 {
			dst.CustomBigInt["hrStorageAllocationUnits"] = dmr.hrStorageAllocationUnits
			dst.CustomBigInt["hrStorageSize"] = dmr.hrStorageSize
			dst.CustomBigInt["hrStorageUsed"] = dmr.hrStorageUsed
			metrics["hrStorageAllocationUnits"] = hrStorageAllocationUnits
			metrics["hrStorageSize"] = hrStorageSize
			metrics["hrStorageUsed"] = hrStorageUsed
			metrics["CPU"] = hrProcessorTree
		}

		if dmr.cpmCPUTotal5min > 0 {
			dst.CustomStr["entPhysicalIndex"] = dmr.entPhysicalIndex
			dst.CustomBigInt["cpmCPUTotal5min"] = dmr.cpmCPUTotal5min
			dst.CustomBigInt["cpmCPUTotal5minRev"] = dmr.cpmCPUTotal5minRev
			dst.CustomBigInt["cpmCPUMemoryFree"] = dmr.cpmCPUMemoryFree
			metrics["cpmCPUTotal5min"] = cpmCPUTotal5min
			metrics["cpmCPUTotal5minRev"] = cpmCPUTotal5minRev
			metrics["cpmCPUMemoryFree"] = cpmCPUMemoryFree
			metrics["CPU"] = cpmCPUTree
			metrics["MemoryFree"] = ciscoMemoryPoolFree
		}

		dst.CustomMetrics = metrics // Add this in so that we know what metrics to pull out down the road.
		flows = append(flows, dst)
	}
	return flows
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
			if variable.Name[0:1] == "." { // Try again, this time not having a leading .
				for oid, m := range dm.conf.DeviceOids {
					if strings.HasPrefix(variable.Name[1:], oid) {
						idx = snmp_util.GetIndex(variable.Name[1:], oid)
						mib = m
						break
					}
				}
			}

			if mib == nil {
				dm.log.Warnf("Missing Custom oid: %+v, Value: %T %+v", variable, variable.Value, variable.Value)
				continue
			}
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
			if mib.Conversion != "" { // Adjust for any hard coded values here.
				value = snmp_util.GetFromConv(variable, mib.Conversion, dm.log)
			}
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
			if mib.EnumRev != nil {
				value := snmp_util.ToInt64(variable.Value)
				if val, ok := mib.EnumRev[value]; ok {
					dmr.customStr[kt.StringPrefix+oidName] = val // Save this string version as a attribute.
				} else {
					dm.log.Warnf("Missing enum value for device metric %s %d", oidName, value)
				}
			}
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
		if dst.Provider == kt.ProviderDefault { // Add this to trigger a UI element.
			dst.CustomStr["profile_message"] = kt.DefaultProfileMessage
		}

		// Memory can be compound value so need to do it here if present but not already set.
		if _, ok := dst.CustomBigInt["MemoryUtilization"]; !ok {
			memoryUsed, oku := dst.CustomBigInt["MemoryUsed"]
			memoryFree, okt := dst.CustomBigInt["MemoryFree"]
			if oku && okt {
				memoryTotal := memoryFree + memoryUsed
				if memoryTotal > 0 {
					dst.CustomBigInt["MemoryUtilization"] = int64(float32(memoryUsed) / float32(memoryTotal) * 100)
					dst.CustomMetrics["MemoryUtilization"] = "computed"
				}
			}
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
		if dst.Provider == kt.ProviderDefault { // Add this to trigger a UI element.
			dst.CustomStr["profile_message"] = kt.DefaultProfileMessage
		}
		flows = append(flows, dst)
	}

	dm.metrics.DeviceMetrics.Mark(int64(len(flows)))
	return flows, nil
}

const (
	jnxOperatingDescr    = "1.3.6.1.4.1.2636.3.1.13.1.5."
	jnxOperatingCPU      = "1.3.6.1.4.1.2636.3.1.13.1.8."
	jnxOperatingDRAMSize = "1.3.6.1.4.1.2636.3.1.13.1.10."
	jnxOperatingBuffer   = "1.3.6.1.4.1.2636.3.1.13.1.11."
	jnxOperatingMemory   = "1.3.6.1.4.1.2636.3.1.13.1.15."
)

// Polling each sub-tree individually is (we hope) a lot cheaper than polling
// their parent tree in its entirety and pulling out just the oids we actually
// want.  The parent tree has a lot of other stuff we don't care about (about
// 5x more than just the oids we want on one device) and there's no reason to
// get the whole thing.
var jnxOids = []string{
	// Do CPU poll at the very beginning, in hopes that it will not be impacted
	// by the CPU-load of polling the other OIDs.
	jnxOperatingCPU,
	jnxOperatingDescr,
	jnxOperatingDRAMSize,
	jnxOperatingBuffer,
	jnxOperatingMemory,
}

// getJuniperDeviceMetrics embeds any errors encountered in the deviceMetricRow.Error
// field and stores them as flow.
func (dm *DeviceMetrics) getJuniperDeviceMetrics(log logger.ContextL, server *gosnmp.GoSNMP) []*deviceMetricRow {
	m := map[string]*deviceMetricRow{}

	var results []gosnmp.SnmpPDU

	for i, oid := range jnxOids {
		// Strip trailing dot from the OID being polled.
		oidResults, err := snmp_util.WalkOID(oid[:len(oid)-1], server, log, "JuniperDeviceMetrics")
		if err != nil {
			m[fmt.Sprintf("err-%d", i)] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", oid, err)}
			dm.metrics.Errors.Mark(1)
			continue
		}

		results = append(results, oidResults...)
	}

	var uptime int64
	uptimeResults, err := snmp_util.WalkOID(sysUpTime, server, log, "JuniperDeviceMetrics")
	if err == nil {
		// You might think that if err == nil then you definitely got back some
		// results.  Not exactly.  The result might be "No such object", which
		// is not an error, but also not what you're looking for.
		if len(uptimeResults) > 0 {
			uptime = snmp_util.ToInt64(uptimeResults[0].Value)
		}
	} else {
		m["uptime"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", sysUpTime, err)}
	}

	// no-op if walkOIDs got errors.
	for _, variable := range results {
		if variable.Value == nil { // You can get nil w/out getting an error, though.
			continue
		}

		// log.Debugf(LOG_PREFIX, "Juniper oid: %+v, Value: %T %+v", variable, variable.Value, variable.Value)
		switch {
		case strings.Contains(variable.Name, jnxOperatingDescr):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, jnxOperatingDescr))
			dm.Component = string(variable.Value.([]byte))
		case strings.Contains(variable.Name, jnxOperatingCPU):
			// index := snmp_util.GetIndex(variable.Name, jnxOperatingCPU)
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, jnxOperatingCPU))
			// log.Infof(LOG_PREFIX, "jnxOperatingCPU: %T %v, index: %s", variable.Value, variable.Value, index)
			dm.CPU = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, jnxOperatingDRAMSize):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, jnxOperatingDRAMSize))
			// log.Debugf(LOG_PREFIX, "jnxOperatingDRAMSize: %T %v", variable.Value, variable.Value)
			dm.juniperOperatingDRAMSize = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, jnxOperatingBuffer):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, jnxOperatingBuffer))
			// log.Debugf(LOG_PREFIX, "jnxOperatingBuffer: %T %v", variable.Value, variable.Value)
			dm.juniperOperatingBuffer = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, jnxOperatingMemory):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, jnxOperatingMemory))
			// log.Debugf(LOG_PREFIX, "jnxOperatingMemory: %T %v", variable.Value, variable.Value)
			dm.juniperOperatingMemory = snmp_util.ToInt64(variable.Value)
		}
	}

	list := make([]*deviceMetricRow, 0, len(m))
	for _, dm := range m {
		dm.Uptime = uptime
		dm.calculateJuniperMemory()
		list = append(list, dm)
	}

	// I kind of wonder if I should sort the list by Component, just so the
	// output is consistent from poll to poll, but I imagine any query will do
	// an ORDER BY, and of it doesn't, then it obviously doesn't care, so
	// *shrug*.

	return list
}

// Adapted from kentik/topology-demo/devicemetrics/main.go
func (dm *deviceMetricRow) calculateJuniperMemory() {
	if dm.juniperOperatingDRAMSize == 0 {
		dm.MemoryTotal = dm.juniperOperatingMemory * 1024 * 1024
	} else {
		dm.MemoryTotal = dm.juniperOperatingDRAMSize
	}
	dm.MemoryUsed = (dm.MemoryTotal / 100.0) * dm.juniperOperatingBuffer
	dm.MemoryFree = dm.MemoryTotal - dm.MemoryUsed
	dm.MemoryUtilization = dm.juniperOperatingBuffer
}

const (
	entPhysicalName = "1.3.6.1.2.1.47.1.1.1.1.7"

	cpmCPUTree               = "1.3.6.1.4.1.9.9.109.1.1.1.1"
	cpmCPUTotalPhysicalIndex = "1.3.6.1.4.1.9.9.109.1.1.1.1.2."
	cpmCPUTotal5min          = "1.3.6.1.4.1.9.9.109.1.1.1.1.5." // values 1..100
	cpmCPUTotal5minRev       = "1.3.6.1.4.1.9.9.109.1.1.1.1.8." // values 0..100; deprecates the above.
	cpmCPUMemoryUsed         = "1.3.6.1.4.1.9.9.109.1.1.1.1.12."
	cpmCPUMemoryFree         = "1.3.6.1.4.1.9.9.109.1.1.1.1.13."
	cpmCPUMemoryHCUsed       = "1.3.6.1.4.1.9.9.109.1.1.1.1.17." // 64-bit version of MemoryUsed. "HC" = "high capacity"
	cpmCPUMemoryHCFree       = "1.3.6.1.4.1.9.9.109.1.1.1.1.19." // 64-bit version of MemoryFree

	ciscoMemoryTree     = "1.3.6.1.4.1.9.9.48.1.1.1"
	ciscoMemoryPoolName = "1.3.6.1.4.1.9.9.48.1.1.1.2."
	ciscoMemoryPoolUsed = "1.3.6.1.4.1.9.9.48.1.1.1.5."
	ciscoMemoryPoolFree = "1.3.6.1.4.1.9.9.48.1.1.1.6."
)

// getCiscoDeviceMetrics embeds any errors encountered in the deviceMetricRow.Error
// field and stores them as flow.
func (dm *DeviceMetrics) getCiscoDeviceMetrics(log logger.ContextL, server *gosnmp.GoSNMP) []*deviceMetricRow {
	m := map[string]*deviceMetricRow{}

	cpuResults, err := snmp_util.WalkOID(cpmCPUTree, server, log, "CiscoDeviceMetrics")
	if err != nil {
		m["err1"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", cpmCPUTree, err)}
	}

	physNameResults, err := snmp_util.WalkOID(entPhysicalName, server, log, "CiscoDeviceMetrics")
	if err != nil {
		m["err2"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", entPhysicalName, err)}
	}

	memResults, err := snmp_util.WalkOID(ciscoMemoryTree, server, log, "CiscoDeviceMetrics")
	if err != nil {
		m["err3"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", ciscoMemoryTree, err)}
	}

	var uptime int64
	uptimeResults, err := snmp_util.WalkOID(sysUpTime, server, log, "CiscoDeviceMetrics")
	if err == nil {
		// You might think that if err == nil then you definitely got back some
		// results.  Not exactly.  The result might be "No such object", which
		// is not an error, but also not what you're looking for.
		if len(uptimeResults) > 0 {
			uptime = snmp_util.ToInt64(uptimeResults[0].Value)
		}
	} else {
		m["uptime"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", sysUpTime, err)}
	}

	// Build this map up front.  Given e.g.
	// entPhysicalName.2 => 52690955
	// store entMap[2] => 52690955
	entMap := map[int64]string{}
	ePN := entPhysicalName + "."
	for _, variable := range physNameResults {
		if variable.Value == nil {
			continue
		}

		if strings.Contains(variable.Name, ePN) {
			index := snmp_util.GetIndex(variable.Name, ePN)
			indexI, err := strconv.ParseInt(index, 10, 64)
			if err != nil { // Shouldn't happen
				continue
			}
			entMap[indexI] = string(variable.Value.([]byte))
			// log.Infof(LOG_PREFIX, "entMap: %v => %v", index, entMap[indexI])
		}
	}

	for _, variable := range append(cpuResults, memResults...) {
		if variable.Value == nil {
			continue
		}
		if bytes, ok := variable.Value.([]byte); ok {
			variable.Value = string(bytes)
		}

		// log.Infof(LOG_PREFIX, "Cisco Name: %v, Type: %v, Value: %T %+v", variable.Name, variable.Type, variable.Value, variable.Value)
		switch {
		// Component set 1
		case strings.Contains(variable.Name, cpmCPUTotalPhysicalIndex):
			dm := assureDeviceMetrics(m, "cpu:"+snmp_util.GetIndex(variable.Name, cpmCPUTotalPhysicalIndex))
			dm.Component = entMap[snmp_util.ToInt64(variable.Value)]
			// log.Infof(LOG_PREFIX, "%v => %v", variable.Value, dm.Component)
		case strings.Contains(variable.Name, cpmCPUTotal5min):
			dm := assureDeviceMetrics(m, "cpu:"+snmp_util.GetIndex(variable.Name, cpmCPUTotal5min))
			dm.cpmCPUTotal5min = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, cpmCPUTotal5minRev):
			dm := assureDeviceMetrics(m, "cpu:"+snmp_util.GetIndex(variable.Name, cpmCPUTotal5minRev))
			dm.cpmCPUTotal5minRev = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, cpmCPUMemoryUsed):
			dm := assureDeviceMetrics(m, "cpu:"+snmp_util.GetIndex(variable.Name, cpmCPUMemoryUsed))
			if dm.cpmCPUMemoryUsed == 0 {
				dm.cpmCPUMemoryUsed = snmp_util.ToInt64(variable.Value)
			}
		case strings.Contains(variable.Name, cpmCPUMemoryHCUsed):
			dm := assureDeviceMetrics(m, "cpu:"+snmp_util.GetIndex(variable.Name, cpmCPUMemoryHCUsed))
			dm.cpmCPUMemoryUsed = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, cpmCPUMemoryFree):
			dm := assureDeviceMetrics(m, "cpu:"+snmp_util.GetIndex(variable.Name, cpmCPUMemoryFree))
			if dm.cpmCPUMemoryFree == 0 {
				dm.cpmCPUMemoryFree = snmp_util.ToInt64(variable.Value)
			}
		case strings.Contains(variable.Name, cpmCPUMemoryHCFree):
			dm := assureDeviceMetrics(m, "cpu:"+snmp_util.GetIndex(variable.Name, cpmCPUMemoryHCFree))
			dm.cpmCPUMemoryFree = snmp_util.ToInt64(variable.Value)

		// Component set 2, different from the above.  No CPU.
		case strings.Contains(variable.Name, ciscoMemoryPoolName):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, ciscoMemoryPoolName))
			dm.Component = variable.Value.(string)
		case strings.Contains(variable.Name, ciscoMemoryPoolUsed):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, ciscoMemoryPoolUsed))
			dm.MemoryUsed = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, ciscoMemoryPoolFree):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, ciscoMemoryPoolFree))
			dm.MemoryFree = snmp_util.ToInt64(variable.Value)
		}
	}

	list := make([]*deviceMetricRow, 0, len(m))
	for _, dm := range m {
		dm.Uptime = uptime
		dm.calculateCiscoCPU()
		dm.calculateCiscoMemory()
		list = append(list, dm)
	}

	// I kind of wonder if I should sort the list by Component, just so the
	// output is consistent from poll to poll, but I imagine any query will do
	// an ORDER BY, and of it doesn't, then it obviously doesn't care, so
	// *shrug*.

	return list
}

// Adapted from kentik/topology-demo/devicemetrics/main.go
func (dm *deviceMetricRow) calculateCiscoCPU() {
	if dm.cpmCPUTotal5minRev > 0 {
		dm.CPU = dm.cpmCPUTotal5minRev
	} else if dm.cpmCPUTotal5min > 0 {
		dm.CPU = dm.cpmCPUTotal5min
	}
}

// Adapted from kentik/topology-demo/devicemetrics/main.go
func (dm *deviceMetricRow) calculateCiscoMemory() {
	skippedCiscoComponents := map[string]bool{
		"image":    true,
		"reserved": true,
	}

	if dm.cpmCPUMemoryFree == 0 {
		dm.MemoryTotal = dm.MemoryFree + dm.MemoryUsed
		if dm.MemoryTotal > 0 {
			// This is a bit of a hack.  Or possibly a heuristic.  For the given
			// components, all the devices we've seen return MemoryFree == 0,
			// which makes MemoryUtilization == 100, which users don't like,
			// since it's alarming and inaccurate.  But maybe not all devices are
			// like that?  So set MemoryUtilization = 0 in the exact case we've
			// seen.
			if skippedCiscoComponents[dm.Component] && dm.MemoryFree == 0 {
				dm.MemoryUtilization = 0
			} else {
				dm.MemoryUtilization = int64(float32(dm.MemoryUsed) / float32(dm.MemoryTotal) * 100)
			}
		}
	} else {
		// in the future this should be based on the sysOID
		dm.MemoryUsed = dm.cpmCPUMemoryUsed * 1024
		dm.MemoryFree = dm.cpmCPUMemoryFree * 1024
		dm.MemoryTotal = dm.MemoryUsed + dm.MemoryFree
		dm.MemoryUtilization = int64(float32(dm.MemoryUsed) / float32(dm.MemoryTotal) * 100)
	}
}

const (
	hrDeviceTree             = "1.3.6.1.2.1.25.2.3.1"
	hrDeviceDescr            = "1.3.6.1.2.1.25.2.3.1.3."
	hrStorageAllocationUnits = "1.3.6.1.2.1.25.2.3.1.4."
	hrStorageSize            = "1.3.6.1.2.1.25.2.3.1.5."
	hrStorageUsed            = "1.3.6.1.2.1.25.2.3.1.6."

	hrProcessorTree = "1.3.6.1.2.1.25.3.3.1.2"
	hrProcessorLoad = "1.3.6.1.2.1.25.3.3.1.2."
)

// getAristaDeviceMetrics embeds any errors encountered in the deviceMetricRow.Error
// field and stores them as flow.
func (dm *DeviceMetrics) getAristaDeviceMetrics(log logger.ContextL, server *gosnmp.GoSNMP) []*deviceMetricRow {
	m := map[string]*deviceMetricRow{}

	// Walk the CPU tree first
	processorTreeResults, err := snmp_util.WalkOID(hrProcessorTree, server, log, "AristaDeviceMetrics")
	if err != nil {
		m["err2"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", hrProcessorTree, err)}
	}

	deviceTreeResults, err := snmp_util.WalkOID(hrDeviceTree, server, log, "AristaDeviceMetrics")
	if err != nil {
		m["err1"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", hrDeviceTree, err)}
	}

	var uptime int64
	uptimeResults, err := snmp_util.WalkOID(sysUpTime, server, log, "AristaDeviceMetrics")
	if err == nil {
		// You might think that if err == nil then you definitely got back some
		// results.  Not exactly.  The result might be "No such object", which
		// is not an error, but also not what you're looking for.
		if len(uptimeResults) > 0 {
			uptime = snmp_util.ToInt64(uptimeResults[0].Value)
		}
	} else {
		m["uptime"] = &deviceMetricRow{Error: fmt.Sprintf("Walking %s: %v", sysUpTime, err)}
	}

	for _, variable := range append(deviceTreeResults, processorTreeResults...) {
		if variable.Value == nil {
			continue
		}
		if bytes, ok := variable.Value.([]byte); ok {
			variable.Value = string(bytes)
		}

		// log.Infof(LOG_PREFIX, "Arista Name: %v, Type: %v, Value: %T %+v", variable.Name, variable.Type, variable.Value, variable.Value)
		switch {
		case strings.Contains(variable.Name, hrDeviceDescr):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, hrDeviceDescr))
			dm.Component = variable.Value.(string)
		case strings.Contains(variable.Name, hrStorageAllocationUnits):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, hrStorageAllocationUnits))
			dm.hrStorageAllocationUnits = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, hrStorageSize):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, hrStorageSize))
			dm.hrStorageSize = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, hrStorageUsed):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, hrStorageUsed))
			dm.hrStorageUsed = snmp_util.ToInt64(variable.Value)
		case strings.Contains(variable.Name, hrProcessorLoad):
			dm := assureDeviceMetrics(m, snmp_util.GetIndex(variable.Name, hrProcessorLoad))
			dm.CPU = snmp_util.ToInt64(variable.Value)
		}
	}

	list := make([]*deviceMetricRow, 0, len(m))
	for _, dm := range m {
		dm.Uptime = uptime
		dm.calculateAristaMemory()
		list = append(list, dm)
	}

	return list
}

// Adapted from kentik/topology-demo/devicemetrics/main.go
func (dm *deviceMetricRow) calculateAristaMemory() {
	dm.MemoryTotal = dm.hrStorageSize * dm.hrStorageAllocationUnits

	if dm.MemoryTotal != 0 {
		dm.MemoryUsed = dm.hrStorageUsed * dm.hrStorageAllocationUnits
		dm.MemoryFree = dm.MemoryTotal - dm.MemoryUsed
		dm.MemoryUtilization = int64(float32(dm.MemoryUsed) / float32(dm.MemoryTotal) * 100)
	}
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
