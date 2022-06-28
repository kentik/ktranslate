package api

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	synthetics "github.com/kentik/api-schema-public/gen/go/kentik/synthetics/v202202"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	CACHE_TIME_DEVICE             = 1 * time.Hour
	API_TIMEOUT                   = 60 * time.Second
	USER_AGENT_BASE               = "KentikFirehose"
	HTTP_USER_AGENT               = "User-Agent"
	API_EMAIL_HEADER              = "X-CH-Auth-Email"
	API_PASSWORD_HEADER           = "X-CH-Auth-API-Token"
	MIN_TIME_BETWEEN_SYNTH_CHECKS = 60 * time.Second
)

var (
	deviceFile = flag.String("api_device_file", "", "File to sideload devices without hitting API")
)

func (api *KentikApi) getDeviceInfo(ctx context.Context, apiUrl string) ([]byte, error) {
	if *deviceFile != "" {
		api.Infof("Reading devices from local file: %s", *deviceFile)
		return os.ReadFile(*deviceFile)
	}

	// Try to make a request, parse the result as json.
	req, err := http.NewRequestWithContext(ctx, "GET", apiUrl, nil)
	if err != nil {
		api.Errorf("Error with Launcher: %v", err)
		return nil, err
	}

	userAgentString := USER_AGENT_BASE

	req.Header.Add(API_EMAIL_HEADER, api.conf.ApiEmail)
	req.Header.Add(API_PASSWORD_HEADER, api.conf.ApiToken)
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

func (api *KentikApi) UpdateTests(ctx context.Context) {
	if api == nil {
		return
	}

	if time.Now().Before(api.lastSynth.Add(MIN_TIME_BETWEEN_SYNTH_CHECKS)) { // Only check this often.
		return
	}

	go func() {
		err := api.getSynthInfo(ctx)
		if err != nil {
			api.Errorf("Cannot get API synth on demand: %v", err)
		}
	}()
	return
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

func (api *KentikApi) GetAgentByIP(ip string) *synthetics.Agent {
	if api == nil {
		return nil
	}
	return api.synAgentsByIP[ip]
}

func (api *KentikApi) GetDevicesAsMap(cid kt.Cid) map[string]*kt.Device {
	if api == nil {
		return nil
	}
	res := map[string]*kt.Device{}
	if cid == 0 {
		for _, cd := range api.devices {
			for _, d := range cd {
				for _, ip := range d.SendingIps {
					res[ip.String()] = d
				}
			}
		}
	} else {
		for _, d := range api.devices[cid] {
			for _, ip := range d.SendingIps {
				res[ip.String()] = d
			}
		}
	}
	return res
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

func (api *KentikApi) getDevices(ctx context.Context) error {
	stime := time.Now()
	res, err := api.getDeviceInfo(ctx, api.conf.ApiRoot+"/api/v5/devices")
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
		myd := device
		if _, ok := resDev[device.CompanyID]; !ok {
			resDev[device.CompanyID] = map[kt.DeviceID]*kt.Device{}
		}
		device.Interfaces = map[kt.IfaceID]kt.Interface{}
		for _, intf := range device.AllInterfaces {
			device.Interfaces[intf.SnmpID] = intf
		}
		resDev[device.CompanyID][device.ID] = &myd
		num++
	}

	api.setTime = time.Now()
	api.Infof("Loaded %d Kentik devices via API in %v", num, api.setTime.Sub(stime))
	api.devices = resDev
	return nil
}

type KentikApi struct {
	logger.ContextL
	tr            *http.Transport
	client        *http.Client
	devices       map[kt.Cid]kt.Devices
	synAgents     map[kt.AgentId]*synthetics.Agent
	synAgentsByIP map[string]*synthetics.Agent
	synTests      map[kt.TestId]*synthetics.Test
	setTime       time.Time
	apiTimeout    time.Duration
	synClient     synthetics.SyntheticsAdminServiceClient
	conf          *kt.KentikConfig
	mux           sync.RWMutex
	lastSynth     time.Time
}

func NewKentikApi(ctx context.Context, conf *kt.KentikConfig, log logger.ContextL) (*KentikApi, error) {
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
		conf:       conf,
		tr:         tr,
		client:     client,
		apiTimeout: apiTimeout,
	}

	// Now, check to see if synthetics API works.
	err := kapi.connectSynth(ctx)
	if err != nil {
		return nil, err
	}

	// check api works and also pre-populate cache.
	err = kapi.getDevices(ctx)

	// Drop cache every hour to keep up to date
	go kapi.manageCache(ctx)

	return kapi, err
}

func NewKentikApiFromLocalDevices(localDevices map[string]*kt.Device, log logger.ContextL) *KentikApi {
	if localDevices == nil {
		return nil
	}

	api := &KentikApi{
		ContextL: log,
	}

	resDev := map[kt.Cid]kt.Devices{}
	num := 0
	for _, devicel := range localDevices {
		device := devicel
		if _, ok := resDev[device.CompanyID]; !ok {
			resDev[device.CompanyID] = map[kt.DeviceID]*kt.Device{}
		}
		device.Interfaces = map[kt.IfaceID]kt.Interface{}
		for _, intf := range device.AllInterfaces {
			device.Interfaces[intf.SnmpID] = intf
		}
		resDev[device.CompanyID][device.ID] = device
		num++
	}

	api.setTime = time.Now()
	api.Infof("Loaded %d Kentik devices via local file", num)
	api.devices = resDev

	return api
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

	address, err := getAddressFromApiRoot(api.conf.ApiRoot)
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
	api.mux.Lock() // Guard this function totally.
	defer api.mux.Unlock()
	api.lastSynth = time.Now()

	md := metadata.New(map[string]string{
		"X-CH-Auth-Email":     api.conf.ApiEmail,
		"X-CH-Auth-API-Token": api.conf.ApiToken,
	})
	ctxo := metadata.NewOutgoingContext(ctx, md)

	lt := &synthetics.ListTestsRequest{}
	r, err := api.synClient.ListTests(ctxo, lt)
	if err != nil {
		return err
	}

	synTests := map[kt.TestId]*synthetics.Test{}
	for _, test := range r.GetTests() {
		localt := test
		synTests[kt.NewTestId(test.GetId())] = localt
	}

	la := &synthetics.ListAgentsRequest{}
	ra, err := api.synClient.ListAgents(ctxo, la)
	if err != nil {
		return err
	}

	synAgents := map[kt.AgentId]*synthetics.Agent{}
	synAgentsByIP := map[string]*synthetics.Agent{}
	for _, agent := range ra.GetAgents() {
		locala := agent
		synAgents[kt.NewAgentId(agent.GetId())] = locala
		lip := locala.GetLocalIp() // Store local ip seperately from public one, if a local is set.
		if lip != "" {
			synAgentsByIP[lip] = locala
		}
		synAgentsByIP[locala.GetIp()] = locala
	}

	api.synAgents = synAgents
	api.synAgentsByIP = synAgentsByIP
	api.synTests = synTests
	api.Infof("Loaded %d Kentik Tests and %d Agents via API", len(api.synTests), len(api.synAgents))

	return nil
}

func (api *KentikApi) EnsureDevice(ctx context.Context, conf *kt.SnmpDeviceConfig) error {
	if api == nil || api.conf == nil {
		return nil
	}

	// If there's no plan id to create devices on, just silently return here.
	if api.conf.ApiPlan == 0 {
		return nil
	}

	// If the device isn't in the list of devices we have, return it too
	dname := strings.ToLower(strings.Replace(conf.DeviceName, ".", "_", -1)) // Normalize to make the API happy.
	for _, c := range api.devices {
		for _, d := range c {
			if d.Name == dname {
				api.Infof("The %s device already exists in Kentik.")
				return nil // No need to keep going.
			}
		}
	}

	desc := conf.Description
	if len(desc) > 128 {
		desc = desc[0:127]
	}
	api.Infof("Adding the %s device in Kentik.", dname)
	dev := &deviceCreate{
		Name:        dname,
		Type:        "router",
		Description: desc,
		SampleRate:  1,
		BgpType:     "none",
		PlanID:      api.conf.ApiPlan,
		IPs:         []net.IP{net.ParseIP(conf.DeviceIP)},
		Subtype:     "router",
		MinSnmp:     false,
	}

	err := api.createDevice(ctx, dev, api.conf.ApiRoot+"/api/v5/device")
	if err != nil {
		return err
	}

	return nil
}

func (api *KentikApi) createDevice(ctx context.Context, create *deviceCreate, url string) error {
	payload, err := json.Marshal(map[string]*deviceCreate{
		"device": create,
	})
	if err != nil {
		return err
	}

	// Try to make a request, parse the result as json.
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	userAgentString := USER_AGENT_BASE
	req.Header.Add(API_EMAIL_HEADER, api.conf.ApiEmail)
	req.Header.Add(API_PASSWORD_HEADER, api.conf.ApiToken)
	req.Header.Add(HTTP_USER_AGENT, userAgentString+" AGENT")
	req.Header.Add("Content-Type", "application/json")

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if err != nil {
		api.Errorf("Error with device create: %v", err)
		api.client = &http.Client{Transport: api.tr, Timeout: api.apiTimeout}
		return err
	}
	if resp.StatusCode != 201 {
		bodymap := map[string]string{}
		json.NewDecoder(resp.Body).Decode(&bodymap)
		return fmt.Errorf("Invalid status code %d -> %v", resp.StatusCode, bodymap["error"])
	}
	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

type deviceCreate struct {
	Name        string   `json:"device_name"`
	Type        string   `json:"device_type"`
	Description string   `json:"device_description"`
	SampleRate  int      `json:"device_sample_rate,string"`
	BgpType     string   `json:"device_bgp_type"`
	PlanID      int      `json:"plan_id,omitempty"`
	SiteID      int      `json:"site_id,omitempty"`
	IPs         []net.IP `json:"sending_ips"`
	CdnAttr     string   `json:"-"`
	ExportId    int      `json:"cloud_export_id,omitempty"`
	Subtype     string   `json:"device_subtype"`
	Region      string   `json:"cloud_region,omitempty"`
	Zone        string   `json:"cloud_zone,omitempty"`
	MinSnmp     bool     `json:"minimize_snmp"`
}
