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
	"strings"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/liamg/furious/scan"
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

type NBRespOK struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []NBResult `json:"results"`
}

type NBResult struct {
	ID           int                    `json:"id"`
	Url          *string                `json:"url"`
	Name         *string                `json:"name"`
	DeviceType   *NBDeviceType          `json:"device_type"`
	PrimaryIp    *NBIP                  `json:"primary_ip"`
	PrimaryIpv4  *NBIP                  `json:"primary_ip4"`
	PrimaryIpv6  *NBIP                  `json:"primary_ip6"`
	OobIp        *NBIP                  `json:"oob_ip"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type NBDeviceType struct {
	Name  *string `json:"name"`
	Model *string `json:"model"`
}

type NBIP struct {
	Address *string `json:"address"`
}

func (i *NBIP) GetVal() string {
	if i == nil || i.Address == nil {
		return "<nil>"
	}
	return *i.Address
}

func getToken(oauthTokenUrl string) (string, error) {
	if time.Now().Before(tokenExp) {
		return accessToken, nil
	}

	client_id := os.Getenv("KTRANS_OAUTH_CLIENT_ID")
	if client_id == "" {
		return "", fmt.Errorf("Set envroment variable KTRANS_OAUTH_CLIENT_ID")
	}
	client_secret := os.Getenv("KTRANS_OAUTH_CLIENT_SECRET")
	if client_secret == "" {
		return "", fmt.Errorf("Set envroment variable KTRANS_OAUTH_CLIENT_SECRET")
	}
	client_scope := os.Getenv("KTRANS_OAUTH_SCOPE")
	if client_scope == "" {
		return "", fmt.Errorf("Set envroment variable KTRANS_OAUTH_SCOPE")
	}

	// Load up the url encoded payload here.
	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")
	payload.Set("client_id", client_id)
	payload.Set("client_secret", client_secret)
	payload.Set("scope", client_scope)

	resp, err := http.Post(oauthTokenUrl, "application/x-www-form-urlencoded", bytes.NewBufferString(payload.Encode()))
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		var res OauthResp
		err = json.Unmarshal(body, &res)
		if err != nil {
			return "", err
		}

		accessToken = res.AccessToken
		tokenExp = time.Now().Add(time.Duration(res.ExpiresIn-60) * time.Second)
		return accessToken, nil
	} else {
		return "", fmt.Errorf("Invalid response from oauth server: %v %v", resp.StatusCode, string(body))
	}
}

func setupDcimFilter(conf *kt.SnmpConfig, log logger.ContextL, offset int32, limit int32) url.Values {
	// Set up some url filters here.
	v := url.Values{}
	v.Add("limit", fmt.Sprintf("%d", limit))
	v.Add("offset", fmt.Sprintf("%d", offset))
	v.Add("interface_count__gt", "0")

	for _, t := range conf.Disco.Netbox.Tag {
		log.Infof("Adding netbox filter for tag %s", t)
		v.Add("tag", t)
	}
	for _, s := range conf.Disco.Netbox.Site {
		log.Infof("Adding netbox filter for site %s", s)
		v.Add("site", s)
	}
	for _, t := range conf.Disco.Netbox.Tenant {
		log.Infof("Adding netbox filter for tenant %s", t)
		v.Add("tenant", t)
	}
	if conf.Disco.Netbox.Status != "" {
		log.Infof("Adding netbox filter for status %s", conf.Disco.Netbox.Status)
		v.Add("status", conf.Disco.Netbox.Status)
	} else { // Default to active.
		log.Infof("Adding netbox filter for status active")
		v.Add("status", "active")
	}
	for _, r := range conf.Disco.Netbox.Role {
		log.Infof("Adding netbox filter for role %s", r)
		v.Add("role", r)
	}
	for _, l := range conf.Disco.Netbox.Location {
		log.Infof("Adding netbox filter for location %s", l)
		v.Add("location", l)
	}

	return v
}

func getDcimDevicesApi(ctx context.Context, conf *kt.SnmpConfig, log logger.ContextL, offset int32, limit int32) (*NBRespOK, error) {

	var url *url.URL
	if conf.Disco.Netbox.NetboxAPIUrl != "" {
		u, err := url.Parse(conf.Disco.Netbox.NetboxAPIUrl)
		if err != nil {
			return nil, err
		}
		url = u
	} else {
		return nil, fmt.Errorf("Missing url for netbox.")
	}
	url.RawQuery = setupDcimFilter(conf, log, offset, limit).Encode()
	log.Infof("Calling netbox at %s", url.String())

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if conf.Disco.Netbox.NetboxAPIToken.String() == "" {
		oauthTokenUrl := os.Getenv("KTRANS_OAUTH_TOKEN_URL")
		if oauthTokenUrl != "" {
			log.Infof("Trying to get bearer token from %s", oauthTokenUrl)
			t, err := getToken(oauthTokenUrl)
			if err != nil {
				return nil, err
			}
			req.Header.Set(authHeaderName, fmt.Sprintf(bearerHeaderFormat, t))
		} else {
			log.Infof("Skipping authentication")
		}
	} else {
		req.Header.Set(authHeaderName, fmt.Sprintf(authHeaderFormat, conf.Disco.Netbox.NetboxAPIToken.String()))
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		nbRes := NBRespOK{}
		err = json.Unmarshal(body, &nbRes)
		if err != nil {
			return nil, err
		}
		return &nbRes, nil
	} else {
		log.Warnf("Invalid response from netbox server: %v", res.Status)
		myErr := map[string]interface{}{}
		err := json.Unmarshal(body, &myErr)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%v", myErr)
	}
}

// See if any the fields in the conf are in the devices custom_fields map.
func checkCustomFields(conf *kt.SnmpConfig, res NBResult) bool {
	for k, v := range conf.Disco.Netbox.CustomFields {
		if val, ok := res.CustomFields[k]; ok {
			switch tv := val.(type) {
			case string:
				if v == tv {
					return true
				}
			}
		}
	}

	return false
}

type netboxDeviceToCheck struct {
	Name    string
	Results []scan.Result
}

func getDevicesFromNetbox(ctx context.Context, ctl chan bool, foundDevices map[string]*kt.SnmpDeviceConfig,
	mdb *mibs.MibDB, conf *kt.SnmpConfig, kentikDevices map[string]string, log logger.ContextL, ignoreMap map[string]bool, ignoreList []netip.Prefix) error {

	log.Infof("Discovering devices from Netbox.")

	var getDeviceList func(offset int32, results *[]netboxDeviceToCheck) error
	getDeviceList = func(offset int32, results *[]netboxDeviceToCheck) error {
		limit := int32(500) // @TODO, what's a good limit here?
		res, err := getDcimDevicesApi(ctx, conf, log, offset, limit)
		if err != nil {
			return err
		}
		for _, res := range res.Results {
			if len(conf.Disco.Netbox.CustomFields) > 0 && len(res.CustomFields) > 0 {
				if !checkCustomFields(conf, res) {
					continue
				}
			}

			ipvs, err := getIPs(res, conf.Disco.Netbox, log)
			if err != nil {
				if res.Name != nil {
					log.Warnf("Skipping %v with bad IP: %v", *res.Name, err)
				} else {
					log.Warnf("Skipping null device with bad IP: %v", err)
				}
			} else {
				if res.Name != nil {
					model := "unknown"
					if res.DeviceType != nil && res.DeviceType.Model != nil {
						model = *res.DeviceType.Model
					}
					rr := make([]scan.Result, len(ipvs))
					for i, ipv := range ipvs {
						rr[i] = scan.Result{Name: *res.Name, Manufacturer: model, Host: net.ParseIP(ipv.Addr().String())}
					}
					*results = append(*results, netboxDeviceToCheck{Name: *res.Name, Results: rr})
				} else {
					log.Warnf("Skipping device with IP %v because of null Name value.", ipvs)
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

	results := make([]netboxDeviceToCheck, 0)
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
		for _, rr := range result.Results {
			if checkIfIgnored(rr.Host.String(), ignoreMap, ignoreList) {
				continue
			}
			wg.Add(1)
			posit := fmt.Sprintf("%d/%d)", i+1, len(results))
			go doubleCheckHost(rr, timeout, ctl, &mux, &wg, foundDevices, mdb, conf, posit, kentikDevices, log)
		}
	}
	wg.Wait()
	log.Infof("Checked %d ips in %v (from start: %v)", len(results), time.Now().Sub(st), time.Now().Sub(stb))

	// Dedupe any extra devices which got added to the found devices set.
	for _, result := range results {
		for ii, rr := range result.Results {
			if _, ok := foundDevices[rr.Host.String()]; ok {
				// Since this current ip is valid, remove any subsequent ips which might have been added.
				for i := ii + 1; i < len(result.Results); i++ {
					delete(foundDevices, result.Results[i].Host.String())
				}
				break // We're good here, move on to next top level result.
			}
		}
	}

	return nil
}

// Pull out the offset value from a url.
func getNextOffset(next *string, log logger.ContextL) int32 {
	if next == nil {
		return 0
	}
	u, err := url.Parse(*next)
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

func getIPs(res NBResult, conf *kt.NetboxConfig, log logger.ContextL) ([]netip.Prefix, error) {

	getMyIP := func(target string) (netip.Prefix, error) {
		switch target {
		case "primary":
			log.Infof("Looking at primary_ip %s", res.PrimaryIp.GetVal())
			if res.PrimaryIp != nil && res.PrimaryIp.Address != nil {
				addr := *res.PrimaryIp.Address
				if addr != "" {
					ipv, err := netip.ParsePrefix(addr)
					return ipv, err
				}
			}
		case "primary_ip4":
			log.Infof("Looking at primary_ip4 %s", res.PrimaryIpv4.GetVal())
			if res.PrimaryIpv4 != nil && res.PrimaryIpv4.Address != nil {
				addr := *res.PrimaryIpv4.Address
				if addr != "" {
					ipv, err := netip.ParsePrefix(addr)
					return ipv, err
				}
			}
		case "primary_ip6":
			log.Infof("Looking at primary_ip6 %s", res.PrimaryIpv6.GetVal())
			if res.PrimaryIpv6 != nil && res.PrimaryIpv6.Address != nil {
				addr := *res.PrimaryIpv6.Address
				if addr != "" {
					ipv, err := netip.ParsePrefix(addr)
					return ipv, err
				}
			}
		case "oob":
			log.Infof("Looking at oob %v", res.OobIp.GetVal())
			if res.OobIp != nil && res.OobIp.Address != nil {
				addr := *res.OobIp.Address
				if addr != "" {
					ipv, err := netip.ParsePrefix(addr)
					return ipv, err
				}
			}
		default:
			log.Infof("Looking at primary_ip %v", res.PrimaryIp.GetVal())
			if res.PrimaryIp != nil && res.PrimaryIp.Address != nil {
				addr := *res.PrimaryIp.Address
				if addr != "" {
					ipv, err := netip.ParsePrefix(addr)
					return ipv, err
				}
			}
		}
		return netip.Prefix{}, fmt.Errorf("IP %s not set", conf.NetboxIP)
	}

	pts := strings.Split(conf.NetboxIP, ",")
	ires := []netip.Prefix{}
	for _, target := range pts {
		ipv, err := getMyIP(target)
		if err != nil {
			log.Warnf("Cannot get ip from target %s: %v", target, err)
		} else {
			ires = append(ires, ipv)
		}
	}

	if len(ires) > 0 {
		return ires, nil
	}

	return nil, fmt.Errorf("IP %s not set", conf.NetboxIP)
}
