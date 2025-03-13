package snmp

/**
Use netbox API to pull down list of devices to target. Get all devices matching a given tag.
*/
import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/liamg/furious/scan"
	"github.com/netbox-community/go-netbox/v4"
)

func getDevicesFromNetbox(ctx context.Context, ctl chan bool, foundDevices map[string]*kt.SnmpDeviceConfig,
	mdb *mibs.MibDB, conf *kt.SnmpConfig, kentikDevices map[string]string, log logger.ContextL, ignoreMap map[string]bool) error {
	c := netbox.NewAPIClientFor(conf.Disco.NetboxAPIHost, conf.Disco.NetboxAPIToken)

	limit := int32(0) // This implies we see all. Do we need to implement paging?
	res, _, err := c.DcimAPI.
		DcimDevicesList(ctx).
		Status([]string{"active"}). // Only look at active devices.
		Limit(limit).
		Execute()

	if err != nil {
		return err
	}

	timeout := time.Millisecond * time.Duration(conf.Global.TimeoutMS)
	stb := time.Now()
	results := make([]scan.Result, 0)
	for _, res := range res.Results {
		hasTag := false
		for _, tag := range res.Tags {
			if tag.Display == conf.Disco.NetboxTag {
				hasTag = true
				continue
			}
		}

		if hasTag && res.PrimaryIp.IsSet() {
			addr := res.PrimaryIp.Get().GetAddress()
			ipv, err := netip.ParsePrefix(addr)
			if err != nil {
				log.Infof("Skipping %s with bad PrimaryIP: %v", res.Display, addr)
			} else {
				results = append(results, scan.Result{Name: res.Display, Manufacturer: res.DeviceType.GetDisplay(), Host: net.ParseIP(ipv.Addr().String())})
			}
		}

	}

	var wg sync.WaitGroup
	var mux sync.RWMutex
	st := time.Now()
	log.Infof("Starting to check %d ips from netbox", len(results))
	for i, result := range results {
		if ignoreMap[result.Host.String()] { // If we have marked this ip as to be ignored, don't do anything more with it.
			continue
		}
		wg.Add(1)
		posit := fmt.Sprintf("%d/%d)", i+1, len(results))
		go doubleCheckHost(result, timeout, ctl, &mux, &wg, foundDevices, mdb, conf, posit, kentikDevices, log)
	}
	wg.Wait()
	log.Infof("Checked %d ips in %v (from start: %v)", len(results), time.Now().Sub(st), time.Now().Sub(stb))

	return nil
}
