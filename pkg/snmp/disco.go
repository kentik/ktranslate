package snmp

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/liamg/furious/scan" // Discovery
	"gopkg.in/yaml.v2"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/snmp/metadata"
	"github.com/kentik/ktranslate/pkg/snmp/mibs"
	snmp_util "github.com/kentik/ktranslate/pkg/snmp/util"
)

const (
	permissions = 0644
)

func Discover(ctx context.Context, snmpFile string, log logger.ContextL) error {
	// First, parse the config file and see what we're doing.
	log.Infof("SNMP Discovery, loading config from %s", snmpFile)
	conf, err := parseConfig(snmpFile)
	if err != nil {
		return err
	}

	if conf.Disco == nil {
		return fmt.Errorf("Missing the discovery config %+v", conf)
	}

	if conf.Global == nil || conf.Global.MibProfileDir == "" {
		return fmt.Errorf("Add a global section and mib profile directory %+v", conf)
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
	mdb, err := mibs.NewMibDB(conf.Global.MibDB, conf.Global.MibProfileDir, conf.Global.PyMibProfileDir, log)
	if err != nil {
		return fmt.Errorf("Cannot set up mibDB -- db: %s, profiles: %s -> %v", conf.Global.MibDB, conf.Global.MibProfileDir, err)
	}
	defer mdb.Close()

	foundDevices := map[string]*kt.SnmpDeviceConfig{}
	for _, ipr := range conf.Disco.Cidrs {
		log.Infof("Discovering snmp speaking devices on %s", ipr)
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
		log.Infof("Starting to check %d ips in %s, checkall=%v", len(results), ipr, conf.Disco.CheckAll)
		for _, result := range results {
			if result.IsHostUp() && (conf.Disco.CheckAll || result.Manufacturer != "") {
				wg.Add(1)
				go doubleCheckHost(result, timeout, ctl, &mux, &wg, foundDevices, mdb, conf, log)
			}
		}
		wg.Wait()
		log.Infof("Checked %d ips in %v", len(results), time.Now().Sub(st))
	}

	if conf.Disco.AddDevices {
		err := addDevices(foundDevices, snmpFile, conf, log)
		if err != nil {
			return err
		}
	}

	return nil
}

func doubleCheckHost(result scan.Result, timeout time.Duration, ctl chan bool, mux *sync.RWMutex, wg *sync.WaitGroup,
	foundDevices map[string]*kt.SnmpDeviceConfig, mdb *mibs.MibDB, conf *kt.SnmpConfig, log logger.ContextL) {

	// Get the token to allow us to run.
	_ = <-ctl
	defer func() {
		wg.Done()
		ctl <- true
	}()

	log.Infof("Host found at %s, Manufacturer: %s, Name: %s -- now attepting checking snmp connectivity", result.Host.String(), result.Manufacturer, result.Name)
	var device kt.SnmpDeviceConfig
	var md *kt.DeviceMetricsMetadata
	var err error
	if conf.Disco.DefaultV3 != nil {
		device = kt.SnmpDeviceConfig{
			DeviceName: result.Name,
			DeviceIP:   result.Host.String(),
			Community:  "", // Run using v3 here.
			V3:         conf.Disco.DefaultV3,
			Debug:      conf.Disco.Debug,
			Port:       uint16(conf.Disco.Ports[0]),
			Checked:    time.Now(),
		}
		serv, err := snmp_util.InitSNMP(&device, timeout, conf.Global.Retries, log)
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
				V3:         conf.Disco.DefaultV3,
				Debug:      conf.Disco.Debug,
				Port:       uint16(conf.Disco.Ports[0]),
				Checked:    time.Now(),
			}
			serv, err := snmp_util.InitSNMP(&device, timeout, conf.Global.Retries, log)
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
	log.Infof("Success connecting to %s -- %v", result.Host.String(), md)

	// Stick in the profile too for future use.
	mibProfile := mdb.FindProfile(md.SysObjectID)
	if mibProfile != nil {
		log.Infof("Found profile for %s: %v", md.SysObjectID, mibProfile)
		device.MibProfile = mibProfile.From
	}

	// Now, see what mibs this sucker can use.
	// TODO, actually store this mibs.
	mibs, provider, first, err := mdb.GetForOidRecur(md.SysObjectID, device.MibProfile, device.Description)
	if err != nil {
		log.Warnf("Issue loading mibs: %v", err)
	} else {
		if first {
			device.DeviceOids = mibs
		}
		device.Provider = provider
	}

	mux.Lock()
	defer mux.Unlock()
	foundDevices[result.Host.String()] = &device
}

func addDevices(foundDevices map[string]*kt.SnmpDeviceConfig, snmpFile string, conf *kt.SnmpConfig, log logger.ContextL) error {
	// Now add the new.
	added := 0
	replaced := 0
	if conf.Devices == nil {
		conf.Devices = map[string]*kt.SnmpDeviceConfig{}
	}

	for _, d := range foundDevices {
		if conf.Devices[d.DeviceName] == nil {
			conf.Devices[d.DeviceName] = d
			added++
		} else {
			if conf.Disco.ReplaceDevices {
				conf.Devices[d.DeviceName] = d
				replaced++
			} else {
				conf.Devices[d.DeviceName].Checked = time.Now()
			}
		}
	}
	log.Infof("Adding %d new snmp devices to the config, %d replaced from %d", added, replaced, len(foundDevices))

	// Save out the config file.
	t, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(snmpFile, t, permissions)
}
