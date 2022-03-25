package snmp

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
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

type SnmpDiscoDeviceStat struct {
	added    int
	replaced int
	delta    int
}

func Discover(ctx context.Context, snmpFile string, log logger.ContextL) (*SnmpDiscoDeviceStat, error) {
	// First, parse the config file and see what we're doing.
	log.Infof("SNMP Discovery, loading config from %s", snmpFile)
	conf, err := parseConfig(ctx, snmpFile, log)
	if err != nil {
		return nil, err
	}

	if conf.Disco == nil {
		return nil, fmt.Errorf("The discovery configuration is not set: %+v.", conf)
	}

	if conf.Global == nil || conf.Global.MibProfileDir == "" {
		return nil, fmt.Errorf("You need to specify a global section and mib profile directory: %v.", conf)
	}

	if *snmpOutFile != "" { // If we want to write somewhere else, swap the output file in here.
		snmpFile = *snmpOutFile
		log.Infof("Writing snmp config file to %s.", snmpFile)
	}

	if conf.Disco.AddDevices { // Verify that the output is writeable before diving into discoing.
		if _, err := addDevices(ctx, nil, snmpFile, conf, true, log); err != nil {
			return nil, fmt.Errorf("There was an error when writing the %s SNMP configuration file: %v.", snmpFile, err)
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
	mdb, err := mibs.NewMibDB(conf.Global.MibDB, conf.Global.MibProfileDir, false, log)
	if err != nil {
		return nil, fmt.Errorf("There was an error when setting up the %s mibDB database and the %s profiles: %v.", conf.Global.MibDB, conf.Global.MibProfileDir, err)
	}
	defer mdb.Close()

	ignoreMap := map[string]bool{}
	for _, ip := range conf.Disco.IgnoreList {
		ignoreMap[ip] = true
	}

	foundDevices := map[string]*kt.SnmpDeviceConfig{}
	for _, ipr := range conf.Disco.Cidrs {
		_, _, err := net.ParseCIDR(ipr)
		if err != nil {
			// Try defaulting this to a /32.
			ipr = ipr + "/32"
			_, _, err := net.ParseCIDR(ipr)
			if err != nil {
				log.Errorf("Invalid cidr, skipping: %s", ipr)
				continue
			} else {
				log.Infof("Defaulting to a /32 range for %s", ipr)
			}
		}

		log.Infof("Discovering SNMP devices on %s.", ipr)
		stb := time.Now()
		targetIterator := scan.NewTargetIterator(ipr)
		timeout := time.Millisecond * time.Duration(conf.Global.TimeoutMS)
		scanner := scan.NewDeviceScanner(targetIterator, timeout)
		if err := scanner.Start(); err != nil {
			return nil, err
		}
		results, err := scanner.Scan(ctx, conf.Disco.Ports)
		if err != nil {
			return nil, err
		}
		var wg sync.WaitGroup
		var mux sync.RWMutex
		st := time.Now()
		log.Infof("Starting to check %d ips in %s", len(results), ipr)
		for i, result := range results {
			if strings.HasSuffix(ipr, "/32") || result.IsHostUp() {
				if ignoreMap[result.Host.String()] { // If we have marked this ip as to be ignored, don't do anything more with it.
					continue
				}
				wg.Add(1)
				posit := fmt.Sprintf("%d/%d)", i+1, len(results))
				go doubleCheckHost(result, timeout, ctl, &mux, &wg, foundDevices, mdb, conf, posit, log)
			}
		}
		wg.Wait()
		log.Infof("Checked %d ips in %v (from start: %v)", len(results), time.Now().Sub(st), time.Now().Sub(stb))
	}

	var stats *SnmpDiscoDeviceStat
	if conf.Disco.AddDevices {
		stats, err = addDevices(ctx, foundDevices, snmpFile, conf, false, log)
		if err != nil {
			return nil, err
		}
	}

	time.Sleep(2 * time.Second) // Give logs time to get sent back.

	return stats, nil
}

func RunDiscoOnTimer(ctx context.Context, c chan os.Signal, snmpFile string, log logger.ContextL, pollTimeMin int) {
	pt := time.Duration(pollTimeMin) * time.Minute
	log.Infof("Running SNMP Discovery Loop every %v", pt)
	discoCheck := time.NewTicker(pt)
	defer discoCheck.Stop()
	for {
		select {
		case _ = <-discoCheck.C:
			stats, err := Discover(ctx, snmpFile, log)
			if err != nil {
				log.Errorf("Discovery SNMP Error: %v", err)
			} else {
				if stats.delta != 0 || stats.added > 0 { // Only restart if there's a different configuration.
					log.Infof("Discovery SNMP reloading: added: %d replaced: %d delta: %d", stats.added, stats.replaced, stats.delta)
					c <- kt.SIGUSR2 // Restart the main loop with a new config.
				} else {
					log.Infof("Discovery SNMP no change so not reloading: added: %d replaced: %d delta: %d", stats.added, stats.replaced, stats.delta)
				}
			}
		case <-ctx.Done():
			log.Infof("Discovery Loop Done")
			return
		}
	}
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
	if conf.Disco.DefaultV3 != nil || len(conf.Disco.OtherV3s) > 0 {
		v3configs := conf.Disco.OtherV3s // Need to keep default seperate for backwords compatibility.
		if v3configs == nil {
			v3configs = make([]*kt.V3SNMPConfig, 0)
		}
		if conf.Disco.DefaultV3 != nil {
			v3configs = append(v3configs, conf.Disco.DefaultV3)
		}

		for _, config := range v3configs { // Now loop over all the possible configs.
			testConfig := config // Capture this as a local var.
			device = kt.SnmpDeviceConfig{
				DeviceName: result.Name,
				DeviceIP:   result.Host.String(),
				Community:  "", // Run using v3 here.
				UseV1:      conf.Disco.UseV1,
				V3:         testConfig,
				Debug:      conf.Disco.Debug,
				Port:       uint16(conf.Disco.Ports[0]),
				Checked:    time.Now(),
			}
			serv, err := snmp_util.InitSNMP(&device, timeout, conf.Global.Retries, posit, log)
			if err != nil {
				log.Warnf("There was an error when starting SNMP interface component -- %v.", err)
				return
			}
			md, err = metadata.GetBasicDeviceMetadata(log, serv)
			if err != nil {
				log.Warnf("Cannot get device metadata on %s: %v. Check for correct snmp credentials.", result.Host.String(), err)
				continue
			}
		}
	} else { // Loop over all possibe v2c options here.
		for _, community := range conf.Disco.DefaultCommunities {
			device = kt.SnmpDeviceConfig{
				DeviceName: result.Name,
				DeviceIP:   result.Host.String(),
				Community:  community,
				UseV1:      conf.Disco.UseV1,
				Debug:      conf.Disco.Debug,
				Port:       uint16(conf.Disco.Ports[0]),
				Checked:    time.Now(),
			}
			serv, err := snmp_util.InitSNMP(&device, timeout, conf.Global.Retries, posit, log)
			if err != nil {
				log.Warnf("There was an error when starting SNMP interface component -- %v.", err)
				return
			}
			md, err = metadata.GetBasicDeviceMetadata(log, serv)
			if err != nil {
				log.Warnf("Cannot get device metadata on %s: %v. Check for correct snmp credentials.", result.Host.String(), err)
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
	if md.EngineID != "" {
		device.EngineID = md.EngineID
	}
	if device.DeviceName == "" { // If for whatever reason we can't get a name, fall back to IP.
		device.DeviceName = device.DeviceIP
	}
	log.Infof("%s Success connecting to %s -- %v", posit, result.Host.String(), md)

	// Stick in the profile too for future use.
	mibProfile := mdb.FindProfile(md.SysObjectID, md.SysDescr, "")
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
	_, provider, _, err := mdb.GetForOidRecur(md.SysObjectID, device.MibProfile, device.Description)
	if err != nil {
		log.Warnf("There was an error when loading the mibs: %v.", err)
	} else {
		// Use the profile's provider if it is set.
		if mibProfile != nil && mibProfile.Provider != "" {
			device.Provider = mibProfile.Provider
		} else {
			device.Provider = provider
		}
	}

	mux.Lock()
	defer mux.Unlock()
	foundDevices[result.Host.String()] = &device
}

func addDevices(ctx context.Context, foundDevices map[string]*kt.SnmpDeviceConfig, snmpFile string, conf *kt.SnmpConfig, isTest bool, log logger.ContextL) (*SnmpDiscoDeviceStat, error) {
	// Now add the new.
	stats := SnmpDiscoDeviceStat{}
	if conf.Devices == nil {
		conf.Devices = map[string]*kt.SnmpDeviceConfig{}
	}
	byIP := map[string]*kt.SnmpDeviceConfig{}
	byEngineID := map[string]*kt.SnmpDeviceConfig{}
	for _, d := range conf.Devices {
		byIP[d.DeviceIP] = d
	}
	origCount := len(conf.Devices)

	for dip, d := range foundDevices {
		key := d.DeviceName
		keyAlt := d.DeviceName + "__" + dip
		if byIP[dip] == nil && conf.Devices[d.DeviceName] != nil {
			log.Warnf("Common device name found with different IPs. %s has %s and %s", d.DeviceName, dip, conf.Devices[d.DeviceName].DeviceIP)
			key = keyAlt
		}

		if conf.Devices[key] == nil && conf.Devices[keyAlt] == nil {
			// Start adding new devices based on deviceName__ip
			conf.Devices[keyAlt] = d
			stats.added++
		} else if conf.Devices[key] != nil {
			if conf.Disco.ReplaceDevices { // But keep backwards compatible with existing devices and don't change their entries.
				d.UpdateFrom(conf.Devices[key], conf)
				conf.Devices[key] = d
				stats.replaced++
			} else {
				conf.Devices[key].Checked = time.Now()
			}
		} else {
			if conf.Disco.ReplaceDevices { // Else, new style keys all use keyAlt.
				d.UpdateFrom(conf.Devices[keyAlt], conf)
				conf.Devices[keyAlt] = d
				stats.replaced++
			} else {
				conf.Devices[keyAlt].Checked = time.Now()
			}
		}
	}

	// Remove any duplicate devices based on Engine ID here.
	for dip, d := range conf.Devices {
		if !conf.Disco.NoDedup && d.EngineID != "" {
			if _, ok := byEngineID[d.EngineID]; ok {
				// Someone else has this engine ID. Delete this device.
				log.Warnf("Removing device %s because of duplicate EngineID %s.", d.DeviceName, d.EngineID)
				delete(conf.Devices, dip)
				stats.added--
			} else {
				byEngineID[d.EngineID] = d
			}
		}
	}

	// Calculate total number of new devices.
	stats.delta = origCount - len(conf.Devices)
	if !isTest {
		log.Infof("Adding %d new SNMP devices to the configuration. %d replaced from %d. Delta: %d", stats.added, stats.replaced, len(foundDevices), stats.delta)
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
			return nil, err
		}
		err = ioutil.WriteFile(conf.Disco.CidrOrig, t, permissions)
		if err != nil {
			return nil, err
		}
		if !isTest {
			conf.Disco.Cidrs = nil
		}
	}

	if conf.Disco.IgnoreOrig != "" {
		t, err := yaml.Marshal(conf.Disco.IgnoreList)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(conf.Disco.IgnoreOrig, t, permissions)
		if err != nil {
			return nil, err
		}
		if !isTest {
			conf.Disco.IgnoreList = nil
		}
	}

	if conf.DeviceOrig != "" {
		t, err := yaml.Marshal(conf.Devices)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(conf.DeviceOrig, t, permissions)
		if err != nil {
			return nil, err
		}
		if !isTest {
			conf.Devices = nil
		}
	}

	// Save out the config file.
	t, err := yaml.Marshal(conf)
	if err != nil {
		return nil, err
	}

	// Swap for our external sections.
	if conf.Disco.CidrOrig != "" {
		t = bytes.Replace(t, []byte("cidrs: []"), []byte(`cidrs: "@`+conf.Disco.CidrOrig+`"`), 1)
	}
	if conf.Disco.IgnoreOrig != "" {
		t = bytes.Replace(t, []byte("ignore_list: []"), []byte(`ignore_list: "@`+conf.Disco.IgnoreOrig+`"`), 1)
	}
	if conf.DeviceOrig != "" {
		t = bytes.Replace(t, []byte("devices: {}"), []byte(`devices: "@`+conf.DeviceOrig+`"`), 1)
	}

	return &stats, snmp_util.WriteFile(ctx, snmpFile, t, permissions)
}
