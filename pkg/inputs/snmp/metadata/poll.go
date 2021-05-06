package metadata

import (
	"context"
	"fmt"
	"time"

	"github.com/kentik/gosnmp"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/tick"
)

type Poller struct {
	log                logger.ContextL
	server             *gosnmp.GoSNMP
	interval           time.Duration
	interfaceMetadata  *InterfaceMetadata
	gotDeviceMetadata  bool
	jchfChan           chan []*kt.JCHF
	conf               *kt.SnmpDeviceConfig
	metrics            *kt.SnmpDeviceMetric
	gconf              *kt.SnmpGlobalConfig
	deviceMetadataMibs map[string]*kt.Mib
}

const (
	DEFUALT_INTERVAL = 30 * 60 * time.Second // Run every 30 min.
)

func NewPoller(server *gosnmp.GoSNMP, gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL) *Poller {

	// If there's a profile passed in, look at the mibs set for this.
	var deviceMetadataMibs, interfaceMetadataMibs map[string]*kt.Mib
	if profile != nil {
		deviceMetadataMibs, interfaceMetadataMibs = profile.GetMetadata(gconf.MibsEnabled)
		if len(deviceMetadataMibs) > 0 {
			log.Infof("Custom device metadata")
			for n, d := range deviceMetadataMibs {
				log.Infof("   -> : %s -> %s", n, d.Name)
			}
		}
		if len(interfaceMetadataMibs) > 0 {
			log.Infof("Custom interface metadata")
			for n, d := range interfaceMetadataMibs {
				log.Infof("   -> : %s -> %s", n, d.Name)
			}
		}
	}

	return &Poller{
		gconf:              gconf,
		conf:               conf,
		log:                log,
		server:             server,
		interval:           DEFUALT_INTERVAL,
		interfaceMetadata:  NewInterfaceMetadata(interfaceMetadataMibs, log),
		gotDeviceMetadata:  false,
		jchfChan:           jchfChan,
		metrics:            metrics,
		deviceMetadataMibs: deviceMetadataMibs,
	}
}

func (p *Poller) StartLoop(ctx context.Context) {
	interfaceCheck := tick.NewJitterTicker(p.interval, 25, 100)

	go func() {
		for {
			select {
			case _ = <-interfaceCheck.C:
				p.log.Infof("Start: Polling SNMP Interface")
				deviceDataNew, err := p.PollSNMPMetadata()
				if err != nil {
					p.log.Warnf("Issue polling SNMP Interface: %v", err)
					p.metrics.Errors.Mark(1)
					continue
				}

				// Do something with this data.
				flows, err := p.toFlows(deviceDataNew)
				if err != nil {
					p.metrics.Errors.Mark(1)
					p.log.Warnf("Issue converting metadata: %v", err)
					continue
				}
				p.metrics.Metadata.Mark(1)
				p.jchfChan <- flows

			case <-ctx.Done():
				p.log.Infof("Metadata Poll Done")
				return
			}
		}
	}()
}

// PollSNMPMetadata checks for relatively static metadata about devices and interfaces
func (p *Poller) PollSNMPMetadata() (*kt.DeviceData, error) {
	intLine, deviceManufacturer, err := p.interfaceMetadata.Poll(p.server)
	if err != nil {
		return nil, err
	}

	// extra check -- if we got no interfaces in that poll, something went badly wrong.
	// Don't delete the interfaces from previous polls from the database; report an error
	// and bail.
	if len(intLine) == 0 {
		err := fmt.Errorf("zero interfaces found")
		p.log.Errorf("SNMP: Error polling metadata: %v", err)
		return nil, err
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

	// Get device-level metadata -- sysDescr and the like, but only once.
	// (But retry on failure or blank data.)
	if !p.gotDeviceMetadata {
		deviceMetadata, err := GetDeviceMetadata(p.log, p.server, p.deviceMetadataMibs)
		if err != nil {
			return nil, err
		}
		if deviceMetadata != nil {
			p.gotDeviceMetadata = true
			deviceData.DeviceMetricsMetadata = deviceMetadata
		}
	}

	return deviceData, nil
}

func (p *Poller) toFlows(dd *kt.DeviceData) ([]*kt.JCHF, error) {
	dst := kt.NewJCHF()
	dst.CustomStr = make(map[string]string)
	dst.CustomInt = make(map[string]int32)
	dst.CustomBigInt = make(map[string]int64)
	dst.EventType = kt.KENTIK_EVENT_SNMP_METADATA
	dst.Provider = p.conf.Provider

	dst.CustomStr["Manufacturer"] = dd.Manufacturer
	dst.DeviceName = p.conf.DeviceName
	dst.SrcAddr = p.conf.DeviceIP
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
			dst.CustomStr[k] = v
		}
		for k, v := range dd.DeviceMetricsMetadata.CustomInts {
			dst.CustomInt[k] = int32(v)
		}
	}

	for intr, id := range dd.InterfaceData {
		dst.CustomStr["if."+intr+".Address"] = id.Address
		dst.CustomStr["if."+intr+".Netmask"] = id.Netmask
		dst.CustomStr["if."+intr+".Index"] = id.Index
		dst.CustomInt["if."+intr+".Speed"] = int32(id.Speed)
		dst.CustomStr["if."+intr+".Description"] = id.Description
		dst.CustomStr["if."+intr+".Alias"] = id.Alias
		dst.CustomStr["if."+intr+".VrfName"] = id.VrfName
		dst.CustomStr["if."+intr+".VrfDescr"] = id.VrfDescr
		dst.CustomStr["if."+intr+".VrfRD"] = id.VrfRD
	}

	return []*kt.JCHF{dst}, nil
}
