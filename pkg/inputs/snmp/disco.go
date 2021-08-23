package snmp

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/liamg/furious/scan" // Discovery
	"gopkg.in/yaml.v2"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/metadata"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	permissions = 0644
)

func Discover(ctx context.Context, snmpFile string, log logger.ContextL) error {
	// First, parse the config file and see what we're doing.
	log.Infof("SNMP Discovery, loading config from %s", snmpFile)
	conf, err := parseConfig(snmpFile, log)
	if err != nil {
		return err
	}

	if conf.Disco == nil {
		return fmt.Errorf("Missing the discovery config %+v", conf)
	}

	if conf.Global == nil || conf.Global.MibProfileDir == "" {
		return fmt.Errorf("Add a global section and mib profile directory %+v", conf)
	}

	if conf.Disco.AddDevices { // Verify that the output is writeable before diving into discoing.
		if err := addDevices(nil, snmpFile, conf, true, log); err != nil {
			return fmt.Errorf("Can not write snmp config file %s - %v", snmpFile, err)
		}
	}

	if len(conf.Disco.Ports) == 0 {
		conf.Disco.Ports = []int{int(snmp_util.SNMP_PORT)}
	}

	// Use this to limit how much parellelism is going on.
	ctl := make(chan bool, conf.Disco.Threads)
	for i := 0; i < conf.Disco.Threads; i++ {
		ctl <- true
	}

	// Use this for auto-discovering metrics to pull.
	mdb, err := mibs.NewMibDB(conf.Global.MibDB, conf.Global.MibProfileDir, log)
	if err != nil {
		return fmt.Errorf("Cannot set up mibDB -- db: %s, profiles: %s -> %v", conf.Global.MibDB, conf.Global.MibProfileDir, err)
	}
	defer mdb.Close()

	foundDevices := map[string]*kt.SnmpDeviceConfig{}
	for _, ipr := range conf.Disco.Cidrs {
		log.Infof("Discovering SNMP devices on %s.", ipr)
		stb := time.Now()
		targetIterator := scan.NewTargetIterator(ipr)
		timeout := time.Millisecond * time.Duration(conf.Global.TimeoutMS)
		scanner := scan.NewDeviceScanner(targetIterator, timeout)
		if err := scanner.Start(); err != nil {
			return err
		}
		results, err := scanner.Scan(ctx, conf.Disco.Ports)
		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		var mux sync.RWMutex
		st := time.Now()
		log.Infof("Starting to check %d ips in %s", len(results), ipr)
		for i, result := range results {
			if strings.HasSuffix(ipr, "/32") || result.IsHostUp() {
				wg.Add(1)
				posit := fmt.Sprintf("%d/%d)", i+1, len(results))
				go doubleCheckHost(result, timeout, ctl, &mux, &wg, foundDevices, mdb, conf, posit, log)
			}
		}
		wg.Wait()
		log.Infof("Checked %d ips in %v (from start: %v)", len(results), time.Now().Sub(st), time.Now().Sub(stb))
	}

	if conf.Disco.AddDevices {
		err := addDevices(foundDevices, snmpFile, conf, false, log)
		if err != nil {
			return err
		}
	}

	return nil
}

func doubleCheckHost(result scan.Result, timeout time.Duration, ctl chan bool, mux *sync.RWMutex, wg *sync.WaitGroup,
	foundDevices map[string]*kt.SnmpDeviceConfig, mdb *mibs.MibDB, conf *kt.SnmpConfig, posit string, log logger.ContextL) {

	// Get the token to allow us to run.
	_ = <-ctl
	defer func() {
		wg.Done()
		ctl <- true
	}()

	log.Infof("%s Host found at %s, Manufacturer: %s, Name: %s -- now attepting checking snmp connectivity", posit, result.Host.String(), result.Manufacturer, result.Name)
	var device kt.SnmpDeviceConfig
	var md *kt.DeviceMetricsMetadata
	var err error
	if conf.Disco.DefaultV3 != nil {
		device = kt.SnmpDeviceConfig{
			DeviceName: result.Name,
			DeviceIP:   result.Host.String(),
			Community:  "", // Run using v3 here.
			UseV1:      conf.Disco.UseV1,
			V3:         conf.Disco.DefaultV3,
			Debug:      conf.Disco.Debug,
			Port:       uint16(conf.Disco.Ports[0]),
			Checked:    time.Now(),
		}
		serv, err := snmp_util.InitSNMP(&device, timeout, conf.Global.Retries, posit, log)
		if err != nil {
			log.Warnf("Init Issue starting SNMP interface component -- %v", err)
			return
		}
		md, err = metadata.GetDeviceMetadata(log, serv, nil)
		if err != nil {
			log.Debugf("Cannot get device metadata on %s: %v", result.Host.String(), err)
			return
		}
	} else { // Loop over all possibe v2c options here.
		for _, community := range conf.Disco.DefaultCommunities {
			device = kt.SnmpDeviceConfig{
				DeviceName: result.Name,
				DeviceIP:   result.Host.String(),
				Community:  community,
				UseV1:      conf.Disco.UseV1,
				V3:         conf.Disco.DefaultV3,
				Debug:      conf.Disco.Debug,
				Port:       uint16(conf.Disco.Ports[0]),
				Checked:    time.Now(),
			}
			serv, err := snmp_util.InitSNMP(&device, timeout, conf.Global.Retries, posit, log)
			if err != nil {
				log.Warnf("Init Issue starting SNMP interface component -- %v", err)
				return
			}
			md, err = metadata.GetDeviceMetadata(log, serv, nil)
			if err != nil {
				log.Debugf("Cannot get device metadata on %s: %v", result.Host.String(), err)
				continue
			}
			break // We're good to go here.
		}
	}

	if md == nil { // No way to establish comminications
		return
	}

	// Map in any discovered values here.
	device.OID = md.SysObjectID
	device.Description = md.SysDescr
	if md.SysName != "" && device.DeviceName == "" { // Swap this in.
		device.DeviceName = md.SysName
	}
	log.Infof("%s Success connecting to %s -- %v", posit, result.Host.String(), md)

	// Stick in the profile too for future use.
	mibProfile := mdb.FindProfile(md.SysObjectID)
	if mibProfile != nil {
		log.Infof("Found profile for %s: %v", md.SysObjectID, mibProfile)
		device.MibProfile = mibProfile.From
		mibs := mibProfile.GetMibs()
		if len(mibs) > 0 {
			device.DiscoveredMibs = make([]string, len(mibs))
			i := 0
			for m, _ := range mibs {
				device.DiscoveredMibs[i] = m
				i++
			}
			sort.Strings(device.DiscoveredMibs) // Put them in a common ordering.
		}
	}

	// Now, see what mibs this sucker can use.
	// TODO, actually store this mibs.
	mibs, provider, first, err := mdb.GetForOidRecur(md.SysObjectID, device.MibProfile, device.Description)
	if err != nil {
		log.Warnf("Issue loading mibs: %v", err)
	} else {
		if first && conf.Disco.AddFromMibDB {
			device.DeviceOids = mibs
		}
		device.Provider = provider
	}

	mux.Lock()
	defer mux.Unlock()
	foundDevices[result.Host.String()] = &device
}

func addDevices(foundDevices map[string]*kt.SnmpDeviceConfig, snmpFile string, conf *kt.SnmpConfig, isTest bool, log logger.ContextL) error {
	// Now add the new.
	added := 0
	replaced := 0
	if conf.Devices == nil {
		conf.Devices = map[string]*kt.SnmpDeviceConfig{}
	}
	byIP := map[string]*kt.SnmpDeviceConfig{}
	for _, d := range conf.Devices {
		byIP[d.DeviceIP] = d
	}

	for dip, d := range foundDevices {
		key := d.DeviceName
		if byIP[dip] == nil && conf.Devices[d.DeviceName] != nil {
			log.Warnf("Common device name found with different IPs. %s has %s and %s", d.DeviceName, dip, conf.Devices[d.DeviceName].DeviceIP)
			key = d.DeviceName + "-" + dip
		}

		if conf.Devices[key] == nil {
			conf.Devices[key] = d
			added++
		} else {
			if conf.Disco.ReplaceDevices {
				conf.Devices[key] = d
				replaced++
			} else {
				conf.Devices[key].Checked = time.Now()
			}
		}
	}
	if !isTest {
		log.Infof("Adding the %s new SNMP device to the configurationAdding %d new snmp devices to the config, %d replaced from %d", added, replaced, len(foundDevices))
	}

	// Fill up list of mibs to run on here.
	if conf.Disco.AddAllMibs {
		fullMibSet := map[string]bool{}
		for _, device := range conf.Devices {
			for _, mib := range device.DiscoveredMibs {
				fullMibSet[mib] = true
			}
		}
		for _, mib := range conf.Global.MibsEnabled {
			fullMibSet[mib] = true
		}

		mibList := []string{}
		for mib, _ := range fullMibSet {
			mibList = append(mibList, mib)
		}
		conf.Global.MibsEnabled = mibList
	}

	sort.Strings(conf.Global.MibsEnabled) // Put them in a common ordering.

	// Write out to seperate files any sections which need this.
	if conf.Disco.CidrOrig != "" {
		t, err := yaml.Marshal(conf.Disco.Cidrs)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(conf.Disco.CidrOrig, t, permissions)
		if err != nil {
			return err
		}
		conf.Disco.Cidrs = nil
	}

	if conf.DeviceOrig != "" {
		t, err := yaml.Marshal(conf.Devices)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(conf.DeviceOrig, t, permissions)
		if err != nil {
			return err
		}
		conf.Devices = nil
	}

	// Save out the config file.
	t, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	// Swap for our external sections.
	if conf.Disco.CidrOrig != "" {
		t = bytes.Replace(t, []byte("cidrs: []"), []byte(`cidrs: "@`+conf.Disco.CidrOrig+`"`), 1)
	}
	if conf.DeviceOrig != "" {
		t = bytes.Replace(t, []byte("devices: {}"), []byte(`devices: "@`+conf.DeviceOrig+`"`), 1)
	}

	return ioutil.WriteFile(snmpFile, t, permissions)
}
