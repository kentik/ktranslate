package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	CACHE_TIME_DEVICE   = -1 * 3 * time.Hour
	API_TIMEOUT         = 60 * time.Second
	USER_AGENT_BASE     = "KentikFirehose"
	HTTP_USER_AGENT     = "User-Agent"
	API_EMAIL_HEADER    = "X-CH-Auth-Email"
	API_PASSWORD_HEADER = "X-CH-Auth-API-Token"
)

func (api *KentikApi) getDeviceInfo(ctx context.Context, apiUrl string) ([]byte, error) {

	// Try to make a request, parse the result as json.
	req, err := http.NewRequestWithContext(ctx, "GET", apiUrl, nil)
	if err != nil {
		api.Errorf("Error with Launcher: %v", err)
		return nil, err
	}

	userAgentString := USER_AGENT_BASE

	req.Header.Add(API_EMAIL_HEADER, api.email)
	req.Header.Add(API_PASSWORD_HEADER, api.token)
	req.Header.Add(HTTP_USER_AGENT, userAgentString+" AGENT")

	resp, err := api.client.Do(req)
	if err != nil {
		api.Errorf("Error with Launcher: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.Errorf("Error with Launcher: %v", err)
		api.client = &http.Client{Transport: api.tr, Timeout: api.apiTimeout}
		return nil, err
	} else if resp.StatusCode == 404 {
		// This should mean that the device is no longer valid -- it's been deleted,
		// or the company has expired, or some similar issue.  We'll stop the flow at
		// chfproxy-relay, anyway, but it's important to bring down the client agent
		// here so that if the customer deletes a device and then adds a new version
		// of it with the same sending-ip, we'll re-do the lookup and start treating
		// the flow as the new device.  So log the error we received, but don't return
		// an error -- return an empty json block.
		api.Warnf("Received a %d error from API: %s", resp.StatusCode, body)
		return []byte("{}"), nil
	} else if resp.StatusCode == 403 {
		// This should mean that the credentials that the agent was started with have become
		// invalid while the agent was running.  If this is true, we've decided to stop sending
		// flow from all device, in order to attract attention to the issue; otherwise, we
		// could keep sending flow indefinitely (as long as we keep sending flow, proxy-relay
		// doesn't have to re-authenticate), but couldn't add new devices, and all flow could be
		// cut off by a later proxy-relay restart.  Better to fail quickly.
		api.Warnf("Received a %d error from API: %s", resp.StatusCode, body)
		return []byte("{}"), nil
	} else if resp.StatusCode >= 300 {
		err = fmt.Errorf("HTTP error: status code %d, body %v", resp.StatusCode, string(body))
		api.Errorf("Lookup failed: %v", err)
		return nil, err
	}

	return body, nil
}

func (api *KentikApi) GetDevices(ctx context.Context) (map[kt.Cid]kt.Devices, error) {
	if api.devices != nil {
		if api.setTime.Before(time.Now().Add(CACHE_TIME_DEVICE)) {
			api.Infof("Re-generating cached devices")
		} else {
			return api.devices, nil // Cache here.
		}
	}

	stime := time.Now()
	res, err := api.getDeviceInfo(ctx, api.apiRoot+"/devices")
	if err != nil {
		return nil, err
	}
	var devices kt.DeviceList
	err = json.Unmarshal(res, &devices)
	if err != nil {
		return nil, err
	}

	resDev := map[kt.Cid]kt.Devices{}
	num := 0
	for _, device := range devices.Devices {
		if _, ok := resDev[device.CompanyID]; !ok {
			resDev[device.CompanyID] = map[kt.DeviceID]kt.Device{}
		}
		device.Interfaces = map[kt.IfaceID]kt.Interface{}
		for _, intf := range device.AllInterfaces {
			device.Interfaces[intf.SnmpID] = intf
		}
		resDev[device.CompanyID][device.ID] = device
		num++
	}

	api.setTime = time.Now()
	api.Infof("Loaded %d Kentik devices via API in %v", num, api.setTime.Sub(stime))
	api.devices = resDev

	return resDev, nil
}

type KentikApi struct {
	logger.ContextL
	email      string
	token      string
	apiRoot    string
	tr         *http.Transport
	client     *http.Client
	devices    map[kt.Cid]kt.Devices
	setTime    time.Time
	apiTimeout time.Duration
}

func NewKentikApi(ctx context.Context, email string, token string, apiRoot string, log logger.ContextL) (*KentikApi, error) {

	apiTimeoutStr := os.Getenv(kt.KentikAPITimeout)
	apiTimeout := API_TIMEOUT
	if apiTimeoutStr != "" {
		intv, _ := strconv.Atoi(apiTimeoutStr)
		apiTimeout = time.Duration(intv) * time.Second
	}
	log.Infof("Setting API timeout to %v", apiTimeout)

	tr := &http.Transport{
		DisableCompression: false,
		DisableKeepAlives:  false,
		Dial: (&net.Dialer{
			Timeout: apiTimeout,
		}).Dial,
		TLSHandshakeTimeout: apiTimeout,
	}
	client := &http.Client{Transport: tr, Timeout: apiTimeout}

	kapi := &KentikApi{
		ContextL:   log,
		email:      email,
		token:      token,
		apiRoot:    apiRoot,
		tr:         tr,
		client:     client,
		apiTimeout: apiTimeout,
	}

	_, err := kapi.GetDevices(ctx) // check api works and also pre-populate cache.
	return kapi, err
}
