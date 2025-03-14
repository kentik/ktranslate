package snmp

/**
Use netbox API to pull down list of devices to target. Get all devices matching a given tag.
*/
import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"
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
	c := netbox.NewAPIClientFor(conf.Disco.NetboxAPIHost, conf.Disco.NetboxAPIToken.String())

	var getDeviceList func(offset int32, results *[]scan.Result) error
	getDeviceList = func(offset int32, results *[]scan.Result) error {
		limit := int32(500) // @TODO, what's a good limit here?
		res, _, err := c.DcimAPI.
			DcimDevicesList(ctx).
			Status([]string{"active"}). // Only look at active devices.
			Limit(limit).
			Offset(offset).
			Execute()
		if err != nil {
			return err
		}

		for _, res := range res.Results {
			keep := true // Default to add all.

			if conf.Disco.NetboxTag != "" {
				keep = false // If we are filtering with tags, only allow those tags which match though.
				for _, tag := range res.Tags {
					if tag.Display == conf.Disco.NetboxTag {
						keep = true
						continue
					}
				}
			}
			if conf.Disco.NetboxSite != "" { // Furthermore, if we filter with site, only the selected site.
				if res.Site.GetDisplay() != conf.Disco.NetboxSite {
					keep = false
				}
			}

			if keep && res.PrimaryIp.IsSet() {
				addr := res.PrimaryIp.Get().GetAddress()
				ipv, err := netip.ParsePrefix(addr)
				if err != nil {
					log.Infof("Skipping %s with bad PrimaryIP: %v", res.Display, addr)
				} else {
					*results = append(*results, scan.Result{Name: res.Display, Manufacturer: res.DeviceType.GetDisplay(), Host: net.ParseIP(ipv.Addr().String())})
				}
			}
		}

		// Now, do we need to recurse further to get more results?
		nextOffset := getNextOffset(res.Next)
		if nextOffset > offset {
			return getDeviceList(nextOffset, results)
		}

		return nil
	}

	results := make([]scan.Result, 0)
	timeout := time.Millisecond * time.Duration(conf.Global.TimeoutMS)
	stb := time.Now()
	err := getDeviceList(0, &results)
	if err != nil {
		return err
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

// Pull out the offset value from a url.
func getNextOffset(next netbox.NullableString) int32 {
	val := next.Get()
	if val == nil {
		return 0
	}
	u, err := url.Parse(*val)
	if err != nil {
		return 0
	}
	q := u.Query()
	no, err := strconv.Atoi(q.Get("offset"))
	if err != nil {
		return 0
	}
	return int32(no)
}
