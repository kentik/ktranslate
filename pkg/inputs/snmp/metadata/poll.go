package metadata

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gosnmp/gosnmp"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/tick"
)

type Poller struct {
	log                   logger.ContextL
	server                *gosnmp.GoSNMP
	interval              time.Duration
	interfaceMetadata     *InterfaceMetadata
	deviceMetadata        *DeviceMetadata
	jchfChan              chan []*kt.JCHF
	conf                  *kt.SnmpDeviceConfig
	metrics               *kt.SnmpDeviceMetric
	gconf                 *kt.SnmpGlobalConfig
	deviceMetadataMibs    map[string]*kt.Mib
	interfaceMetadataMibs map[string]*kt.Mib
	matchAttr             map[string]*regexp.Regexp
}

const (
	DEFAULT_INTERVAL = 30 * 60 * time.Second // Run every 30 min.
	vendorPrefix     = ".1.3.6.1.4.1."
)

func NewPoller(server *gosnmp.GoSNMP, gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL) *Poller {

	attrMap := map[string]*regexp.Regexp{} // List of attributes to pass though, if empty all are passed.

	// If there's a profile passed in, look at the mibs set for this.
	var deviceMetadataMibs, interfaceMetadataMibs map[string]*kt.Mib
	if profile != nil {
		deviceMetadataMibs, interfaceMetadataMibs = profile.GetMetadata(gconf.MibsEnabled)
		if len(deviceMetadataMibs) > 0 {
			log.Infof("Custom device metadata")
			for n, d := range deviceMetadataMibs {
				log.Infof("   -> : %s -> %s", n, d.Name)
				for k, v := range d.MatchAttr {
					attrMap[k] = v
				}
			}
		}
		if len(interfaceMetadataMibs) > 0 {
			log.Infof("Custom interface metadata")
			for n, d := range interfaceMetadataMibs {
				log.Infof("   -> : %s -> %s", n, d.Name)
				for k, v := range d.MatchAttr {
					attrMap[k] = v
				}
			}
		}

		// See if there's a device comment for this profile.
		dev := profile.GetDeviceSysComment(conf.OID)
		if dev != "" {
			log.Debugf("Adding model tag %s", dev)
			conf.AddUserTag("kentik.model", dev)
		}
	}

	// Load any attribute level whiltelist info here.
	for attr, restr := range conf.MatchAttr {
		re, err := regexp.Compile(restr)
		if err != nil {
			log.Errorf("Ignoring Match Attribute %s: %s -- invalid regex %v", attr, restr, err)
		}
		attrMap[attr] = re
	}
	if !conf.MonitorAdminShut { // This one is common and so is set explicitly.
		attrMap[kt.AdminStatus] = regexp.MustCompile("up")
	}
	if len(attrMap) > 0 {
		log.Infof("Added %d Match Attribute(s)", len(attrMap))
	}

	return &Poller{
		gconf:                 gconf,
		conf:                  conf,
		log:                   log,
		server:                server,
		interval:              DEFAULT_INTERVAL,
		interfaceMetadata:     NewInterfaceMetadata(interfaceMetadataMibs, log),
		deviceMetadata:        NewDeviceMetadata(deviceMetadataMibs, conf, metrics, log),
		jchfChan:              jchfChan,
		metrics:               metrics,
		deviceMetadataMibs:    deviceMetadataMibs,
		interfaceMetadataMibs: interfaceMetadataMibs,
		matchAttr:             attrMap,
	}
}

func (p *Poller) StartLoop(ctx context.Context) {
	interfaceCheck := tick.NewJitterTicker(p.interval, 25, 100)
	runPoll := func() {
		p.log.Infof("Start: Polling SNMP Interface")
		deviceDataNew, err := p.PollSNMPMetadata(ctx)
		if err != nil {
			p.log.Warnf("There was an error when polling for SNMP devices: %v.", err)
			p.metrics.Errors.Mark(1)
			return
		}

		// Do something with this data.
		flows, err := p.toFlows(deviceDataNew)
		if err != nil {
			p.metrics.Errors.Mark(1)
			p.log.Warnf("There was an error when converting the metadata: %v.", err)
			return
		}
		p.metrics.Metadata.Mark(1)
		p.jchfChan <- flows
	}

	runPoll() // Get first set of metadata
	go func() {
		for {
			select {
			case _ = <-interfaceCheck.C:
				runPoll()

			case <-ctx.Done():
				p.log.Infof("Metadata Poll Done")
				interfaceCheck.Stop()
				return
			}
		}
	}()
}

// PollSNMPMetadata checks for relatively static metadata about devices and interfaces
func (p *Poller) PollSNMPMetadata(ctx context.Context) (*kt.DeviceData, error) {
	intLine, deviceManufacturer, err := p.interfaceMetadata.Poll(ctx, p.conf, p.server)
	if err != nil {
		return nil, err
	}

	// If there's no interfaces, note this this might be an issue for some devices but keep on going.
	if len(intLine) == 0 {
		p.log.Warnf("No SNMP devices were found when polling the metadata.")
	}

	deviceData := &kt.DeviceData{
		Manufacturer:  deviceManufacturer,
		InterfaceData: intLine,
	}

	// TODO: I'm not going to move this right now, but it seems almost impossible to me that this is correctly placed.
	//       Surely both the kafka topic and filterDeviceData need to see the mapped interface IDs to do the correct thing.
	//       Also, since UpdateForHuawei operates on deviceData, not deviceDataNew, and sets the InterfaceData map in that
	//       struct, I'm pretty sure it's been broken for years.  Fortunately, I don't think we have any Huawei devices
	//       that this is necessary for at the moment.
	//
	// If the device is Huawei, update the snmp ids
	if isBrokenHuawei(deviceManufacturer) {
		err := p.interfaceMetadata.UpdateForHuawei(p.server, deviceData)
		if err != nil {
			return nil, err
		}
	}

	deviceMetadata, err := p.deviceMetadata.Poll(ctx, p.server)
	if err != nil {
		return nil, err
	}
	if deviceMetadata != nil {
		deviceData.DeviceMetricsMetadata = deviceMetadata
	}

	return deviceData, nil
}

func (p *Poller) toFlows(dd *kt.DeviceData) ([]*kt.JCHF, error) {
	dst := kt.NewJCHF()
	dst.CustomStr = make(map[string]string)
	dst.CustomInt = make(map[string]int32)
	dst.CustomBigInt = make(map[string]int64)
	dst.CustomMetrics = make(map[string]kt.MetricInfo)
	dst.EventType = kt.KENTIK_EVENT_SNMP_METADATA
	dst.Provider = p.conf.Provider

	dst.CustomStr["Manufacturer"] = dd.Manufacturer
	dst.DeviceName = p.conf.DeviceName
	dst.SrcAddr = p.conf.DeviceIP
	dst.Timestamp = time.Now().Unix()
	if dd.DeviceMetricsMetadata != nil {
		dst.CustomStr["SysName"] = dd.DeviceMetricsMetadata.SysName
		dst.CustomStr["SysObjectID"] = dd.DeviceMetricsMetadata.SysObjectID
		dst.CustomStr["SysDescr"] = dd.DeviceMetricsMetadata.SysDescr
		dst.CustomStr["SysLocation"] = dd.DeviceMetricsMetadata.SysLocation
		dst.CustomStr["SysContact"] = dd.DeviceMetricsMetadata.SysContact
		dst.CustomInt["SysServices"] = int32(dd.DeviceMetricsMetadata.SysServices)
		if dst.DeviceName == "" {
			dst.DeviceName = dd.DeviceMetricsMetadata.SysName
		}
		for k, v := range dd.DeviceMetricsMetadata.Customs {
			if utf8.ValidString(v) {
				dst.CustomStr[k] = v
			}
		}
		for k, v := range dd.DeviceMetricsMetadata.CustomInts {
			dst.CustomInt[k] = int32(v)
		}
		if len(dd.DeviceMetricsMetadata.Tables) > 0 {
			dst.CustomTables = dd.DeviceMetricsMetadata.Tables

			for _, table := range dst.CustomTables {
				for k, v := range table.Customs {
					dst.CustomMetrics[k] = kt.MetricInfo{Table: v.TableName, Tables: v.TableNames}
				}
			}
		}

		// Compute vendor int here.
		if dst.Provider == kt.ProviderDefault { // Add this to trigger a UI element.
			if strings.HasPrefix(dst.CustomStr["SysObjectID"], vendorPrefix) {
				pts := strings.SplitN(dst.CustomStr["SysObjectID"][len(vendorPrefix):], ".", 2)
				if vendorId, err := strconv.Atoi(pts[0]); err == nil {
					dst.CustomInt["sysoid_vendor"] = int32(vendorId)
				}
			}
			dst.CustomStr["sysoid_profile"] = p.conf.MibProfile
		}
	}

	// Now, pass along lookup info if we find any.
	for k, _ := range dst.CustomStr {
		if mib, ok := p.lookupMib(k, false); ok {
			dst.CustomMetrics[k] = kt.MetricInfo{Oid: mib.Oid, Mib: mib.Mib, Table: mib.Table, Tables: mib.OtherTables}
		}
	}
	for k, _ := range dst.CustomInt {
		if mib, ok := p.lookupMib(k, false); ok {
			dst.CustomMetrics[k] = kt.MetricInfo{Oid: mib.Oid, Mib: mib.Mib, Table: mib.Table, Tables: mib.OtherTables}
		}
	}

	// Also any device level tags
	cs := map[string]string{}
	p.conf.SetUserTags(cs)
	for k, v := range cs {
		dst.CustomStr[k] = v
		dst.CustomMetrics[k] = kt.MetricInfo{Table: kt.DeviceTagTable}
	}

	for intr, id := range dd.InterfaceData {
		dst.CustomStr["if."+intr+".Address"] = id.Address
		dst.CustomStr["if."+intr+".Netmask"] = id.Netmask
		dst.CustomStr["if."+intr+".Index"] = id.Index
		dst.CustomInt["if."+intr+".Speed"] = int32(id.Speed)
		dst.CustomStr["if."+intr+".Description"] = id.Description
		dst.CustomStr["if."+intr+".Alias"] = id.Alias
		dst.CustomStr["if."+intr+".Type"] = id.Type

		// And in anything extra which came out here.
		for k, v := range id.ExtraInfo {
			if utf8.ValidString(v) {
				dst.CustomStr["if."+intr+"."+k] = v
				if mib, ok := p.lookupMib(k, true); ok {
					dst.CustomMetrics["if_"+k] = kt.MetricInfo{Oid: mib.Oid, Mib: mib.Mib, Table: mib.Table, Tables: mib.OtherTables}
				}
			}
		}
	}

	if len(p.matchAttr) > 0 {
		dst.MatchAttr = p.matchAttr
	}

	return []*kt.JCHF{dst}, nil
}

func (p *Poller) lookupMib(key string, isInterface bool) (*kt.Mib, bool) {
	if !isInterface {
		for _, mib := range p.deviceMetadataMibs {
			if mib.Name == key {
				return mib, true
			}
			if mib.Tag == key {
				return mib, true
			}
		}
	} else {
		for _, mib := range p.interfaceMetadataMibs {
			if mib.Name == key || mib.Name == "if_"+key {
				return mib, true
			}
			if mib.Tag == key || mib.Tag == "if_"+key {
				return mib, true
			}
		}
	}

	return nil, false
}
