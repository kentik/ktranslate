package snmp

/**
Use netbox API to pull down list of devices to target. Get all devices matching a given tag.
*/
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/liamg/furious/scan"
	"github.com/netbox-community/go-netbox/v4"
)

var (
	tokenExp    time.Time
	accessToken string
)

const (
	authHeaderName      = "Authorization"
	authHeaderFormat    = "Token %v"
	bearerHeaderFormat  = "Bearer %v"
	languageHeaderName  = "Accept-Language"
	languageHeaderValue = "en-US"
)

type OauthResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func getToken() (string, error) {
	if time.Now().Before(tokenExp) {
		return accessToken, nil
	}

	payload := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     os.Getenv("KTRANS_OAUTH_CLIENT_ID"),
		"client_secret": os.Getenv("KTRANS_OAUTH_CLIENT_SECRET"),
		"scope":         os.Getenv("KTRANS_OAUTH_SCOPE"),
	}
	jsonValue, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(os.Getenv("KTRANS_OAUTH_TOKEN_URL"), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res OauthResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	accessToken = res.AccessToken
	tokenExp = time.Now().Add(time.Duration(res.ExpiresIn-60) * time.Second)
	return accessToken, nil
}

func newAPIClientFor(host string, token string) (*netbox.APIClient, error) {
	cfg := netbox.NewConfiguration()

	// If needed, get a token via oauth
	if token == "" {
		t, err := getToken()
		if err != nil {
			return nil, err
		}
		cfg.AddDefaultHeader(
			authHeaderName,
			fmt.Sprintf(bearerHeaderFormat, t),
		)
	} else {
		cfg.AddDefaultHeader(
			authHeaderName,
			fmt.Sprintf(authHeaderFormat, token),
		)
	}

	cfg.Servers[0].URL = host

	cfg.AddDefaultHeader(
		languageHeaderName,
		languageHeaderValue,
	)

	return netbox.NewAPIClient(cfg), nil
}

func getDcimDevicesApi(ctx context.Context, conf *kt.SnmpConfig, log logger.ContextL) (string, error) {

	v := url.Values{}
	v.Add("limit", "500")
	v.Add("status", "active")
	v.Add("offset", "0")
	v.Add("tag", conf.Disco.Netbox.NetboxTag)
	v.Add("site", conf.Disco.Netbox.NetboxSite)
	v.Add("interface_count__gt", "0")

	u, err := url.Parse(conf.Disco.Netbox.NetboxAPIHost + "/api/dcim/devices/")
	if err != nil {
		return "", err
	}
	u.RawQuery = v.Encode()
	log.Infof("XXX %v", u.String())

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set(authHeaderName, fmt.Sprintf(authHeaderFormat, conf.Disco.Netbox.NetboxAPIToken.String()))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getDevicesFromNetbox(ctx context.Context, ctl chan bool, foundDevices map[string]*kt.SnmpDeviceConfig,
	mdb *mibs.MibDB, conf *kt.SnmpConfig, kentikDevices map[string]string, log logger.ContextL, ignoreMap map[string]bool) error {
	c, err := newAPIClientFor(conf.Disco.Netbox.NetboxAPIHost, conf.Disco.Netbox.NetboxAPIToken.String())
	if err != nil {
		return err
	}

	res, err := getDcimDevicesApi(ctx, conf, log)
	if err != nil {
		return err
	}
	log.Infof("XX %s", res)

	var getDeviceList func(offset int32, results *[]scan.Result) error
	getDeviceList = func(offset int32, results *[]scan.Result) error {
		limit := int32(500) // @TODO, what's a good limit here?
		res, _, err := c.DcimAPI.
			DcimDevicesList(ctx).
			Status([]string{"active"}). // Only look at active devices.
			Limit(limit).
			Offset(offset).
			InterfaceCountGt([]int32{0}). // We want only devices with at least one iterface.
			Execute()
		if err != nil {
			return err
		}
		for _, res := range res.Results {
			keep := true // Default to add all.

			if conf.Disco.Netbox.NetboxTag != "" {
				keep = false // If we are filtering with tags, only allow those tags which match though.
				for _, tag := range res.Tags {
					if tag.Display == conf.Disco.Netbox.NetboxTag {
						keep = true
						continue
					}
				}
			}
			if conf.Disco.Netbox.NetboxSite != "" { // Furthermore, if we filter with site, only the selected site.
				if res.Site.GetDisplay() != conf.Disco.Netbox.NetboxSite {
					keep = false
				}
			}

			if keep {
				ipv, err := getIP(res, conf.Disco.Netbox)
				if err != nil {
					log.Infof("Skipping %s with bad IP: %v", res.Display, err)
				} else {
					*results = append(*results, scan.Result{Name: res.Display, Manufacturer: res.DeviceType.GetDisplay(), Host: net.ParseIP(ipv.Addr().String())})
				}
			}
		}

		// Now, do we need to recurse further to get more results?
		nextOffset := getNextOffset(res.Next, log)
		if nextOffset > offset {
			return getDeviceList(nextOffset, results)
		}

		return nil
	}

	results := make([]scan.Result, 0)
	timeout := time.Millisecond * time.Duration(conf.Global.TimeoutMS)
	stb := time.Now()
	err = getDeviceList(0, &results)
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
func getNextOffset(next netbox.NullableString, log logger.ContextL) int32 {
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

func getIP(res netbox.DeviceWithConfigContext, conf *kt.NetboxConfig) (netip.Prefix, error) {
	switch conf.NetboxIP {
	case "primary":
		if res.PrimaryIp.IsSet() {
			addr := res.PrimaryIp.Get().GetAddress()
			if addr != "" {
				ipv, err := netip.ParsePrefix(addr)
				return ipv, err
			}
		}
	case "oob":
		if res.OobIp.IsSet() {
			addr := res.OobIp.Get().GetAddress()
			if addr != "" {
				ipv, err := netip.ParsePrefix(addr)
				return ipv, err
			}
		}
	default:
		if res.PrimaryIp.IsSet() {
			addr := res.PrimaryIp.Get().GetAddress()
			if addr != "" {
				ipv, err := netip.ParsePrefix(addr)
				return ipv, err
			}
		}
	}
	return netip.Prefix{}, fmt.Errorf("IP %s not set", conf.NetboxIP)
}
