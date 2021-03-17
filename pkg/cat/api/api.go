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

	synthetics "github.com/kentik/api-schema/gen/go/kentik/synthetics/v202101beta1"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	CACHE_TIME_DEVICE   = 1 * time.Hour
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

func (api *KentikApi) GetTest(tid kt.TestId) *synthetics.Test {
	if api == nil {
		return nil
	}
	return api.synTests[tid]
}

func (api *KentikApi) GetAgent(aid kt.AgentId) *synthetics.Agent {
	if api == nil {
		return nil
	}
	return api.synAgents[aid]
}

func (api *KentikApi) GetDevice(cid kt.Cid, did kt.DeviceID) *kt.Device {
	if api == nil {
		return nil
	}
	if c, ok := api.devices[cid]; ok {
		return c[did]
	}
	return nil
}

func (api *KentikApi) CreateDeviceIfNotPresent(name string, ip string) error {
	//client, err = libkflow.NewSenderWithNewDevice(dconf, errors, config)
	return nil
}

func (api *KentikApi) getDevices(ctx context.Context) error {
	stime := time.Now()
	res, err := api.getDeviceInfo(ctx, api.apiRoot+"/api/v5/devices")
	if err != nil {
		return err
	}
	var devices kt.DeviceList
	err = json.Unmarshal(res, &devices)
	if err != nil {
		return err
	}

	resDev := map[kt.Cid]kt.Devices{}
	num := 0
	for _, device := range devices.Devices {
		if _, ok := resDev[device.CompanyID]; !ok {
			resDev[device.CompanyID] = map[kt.DeviceID]*kt.Device{}
		}
		device.Interfaces = map[kt.IfaceID]kt.Interface{}
		for _, intf := range device.AllInterfaces {
			device.Interfaces[intf.SnmpID] = intf
		}
		resDev[device.CompanyID][device.ID] = &device
		num++
	}

	api.setTime = time.Now()
	api.Infof("Loaded %d Kentik devices via API in %v", num, api.setTime.Sub(stime))
	api.devices = resDev
	return nil
}

type KentikApi struct {
	logger.ContextL
	email      string
	token      string
	apiRoot    string
	tr         *http.Transport
	client     *http.Client
	devices    map[kt.Cid]kt.Devices
	synAgents  map[kt.AgentId]*synthetics.Agent
	synTests   map[kt.TestId]*synthetics.Test
	setTime    time.Time
	apiTimeout time.Duration
	synClient  synthetics.SyntheticsAdminServiceClient
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

	err := kapi.connectSynth(ctx) // Now, check to see if synthetics API works.
	if err != nil {
		return nil, err
	}

	err = kapi.getDevices(ctx) // check api works and also pre-populate cache.

	// Drop cache every hour to keep up to date
	go kapi.manageCache(ctx)

	return kapi, err
}

func (api *KentikApi) manageCache(ctx context.Context) {
	checkTicker := time.NewTicker(CACHE_TIME_DEVICE)
	defer checkTicker.Stop()

	for {
		select {
		case _ = <-checkTicker.C:
			err := api.getDevices(ctx)
			if err != nil {
				api.Errorf("Cannot get API devices: %v", err)
			}
			err = api.getSynthInfo(ctx)
			if err != nil {
				api.Errorf("Cannot get API synth: %v", err)
			}

		case <-ctx.Done():
			api.Infof("manageApiCache Done")
			return
		}
	}
}

func (api *KentikApi) connectSynth(ctxIn context.Context) error {
	ctx, cancel := context.WithTimeout(ctxIn, api.apiTimeout)
	defer cancel()

	creds, err := clientTransportCredentials(true, "", "", "")
	if err != nil {
		return err
	}

	address, err := getAddressFromApiRoot(api.apiRoot)
	if err != nil {
		return err
	}

	api.Infof("Connecting to API server at %s", address)
	conn, err := grpc.DialContext(ctx, address, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}

	client := synthetics.NewSyntheticsAdminServiceClient(conn)
	api.Infof("Connected to Synth API server at %s", address)
	api.synClient = client

	return api.getSynthInfo(ctx)
}

func (api *KentikApi) getSynthInfo(ctx context.Context) error {
	md := metadata.New(map[string]string{
		"X-CH-Auth-Email":     api.email,
		"X-CH-Auth-API-Token": api.token,
	})
	ctxo := metadata.NewOutgoingContext(ctx, md)

	lt := &synthetics.ListTestsRequest{}
	r, err := api.synClient.ListTests(ctxo, lt)
	if err != nil {
		return err
	}

	synTests := map[kt.TestId]*synthetics.Test{}
	for _, test := range r.GetTests() {
		synTests[kt.NewTestId(test.GetId())] = test
	}

	la := &synthetics.ListAgentsRequest{}
	ra, err := api.synClient.ListAgents(ctxo, la)
	if err != nil {
		return err
	}

	synAgents := map[kt.AgentId]*synthetics.Agent{}
	for _, agent := range ra.GetAgents() {
		synAgents[kt.NewAgentId(agent.GetId())] = agent
	}

	api.synAgents = synAgents
	api.synTests = synTests
	api.Infof("Loaded %d Kentik Tests and %d Agents via API", len(api.synTests), len(api.synAgents))

	return nil
}
