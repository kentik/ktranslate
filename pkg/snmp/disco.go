package snmp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/liamg/furious/scan"

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

	if len(conf.Disco.Ports) == 0 {
		conf.Disco.Ports = []int{int(snmp_util.SNMP_PORT)}
	}

	// Use this to limit how much parellelism is going on.
	ctl := make(chan bool, conf.Disco.Threads)
	for i := 0; i < conf.Disco.Threads; i++ {
		ctl <- true
	}

	// Use this for auto-discovering metrics to pull.
	mdb, err := mibs.NewMibDB(conf.Disco.MibDB, log)
	if err != nil {
		return fmt.Errorf("Missing the mibs db config %s -> %v", conf.Disco.MibDB, err)
	}
	defer mdb.Close()

	foundDevices := map[string]*kt.SnmpDeviceConfig{}
	for _, ipr := range conf.Disco.Cidrs {
		log.Infof("Discovering snmp speaking devices on %s", ipr)
		targetIterator := scan.NewTargetIterator(ipr)
		timeout := time.Millisecond * time.Duration(conf.Disco.TimeoutMS)
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
		for _, result := range results {
			if result.IsHostUp() && result.Manufacturer != "" {
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
	device := kt.SnmpDeviceConfig{
		DeviceName: result.Name,
		DeviceIP:   result.Host.String(),
		Community:  conf.Disco.DefaultCommunity,
		V3:         conf.Disco.DefaultV3,
		Debug:      conf.Disco.Debug,
		Port:       uint16(conf.Disco.Ports[0]),
		Checked:    time.Now(),
	}
	serv, err := snmp_util.InitSNMP(&device, timeout, conf.Disco.Retries, log)
	if err != nil {
		log.Warnf("Init Issue starting SNMP interface component -- %v", err)
		return
	}
	md, err := metadata.GetDeviceMetadata(log, serv)
	if err != nil {
		log.Debugf("Cannot get device metadata on %s: %v", result.Host.String(), err)
		return
	}
	// Map in any discovered values here.
	device.OID = md.SysObjectID
	device.Description = md.SysDescr
	if md.SysName != "" && device.DeviceName == "" { // Swap this in.
		device.DeviceName = md.SysName
	}
	log.Infof("Success connecting to %s -- %v", result.Host.String(), md)

	// Now, see what mibs this sucker can use.
	mibs, err := mdb.GetForOidRecur(md.SysObjectID)
	if err != nil {
		log.Warnf("Issue loading mibs: %v", err)
	} else {
		for _, mib := range mibs {
			log.Infof("Mib: %v", mib)
		}
	}

	mux.Lock()
	defer mux.Unlock()
	foundDevices[result.Host.String()] = &device
}

func addDevices(foundDevices map[string]*kt.SnmpDeviceConfig, snmpFile string, conf *kt.SnmpConfig, log logger.ContextL) error {
	// List the old.
	oldDevices := map[string]*kt.SnmpDeviceConfig{}
	for _, d := range conf.Devices {
		oldDevices[d.DeviceIP] = d
	}

	// Now add the new.
	added := 0
	for ip, d := range foundDevices {
		if oldDevices[ip] == nil {
			conf.Devices = append(conf.Devices, d)
			added++
		} else {
			oldDevices[ip].Checked = time.Now()
		}
	}
	log.Infof("Adding %d new snmp devices to the config", added)

	// Save out the config file.
	t, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(snmpFile, t, permissions)
}
