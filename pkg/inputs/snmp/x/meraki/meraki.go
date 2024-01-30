package meraki

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	apiclient "github.com/kentik/dashboard-api-golang/client"
	"github.com/kentik/dashboard-api-golang/client/appliance"
	"github.com/kentik/dashboard-api-golang/client/devices"
	"github.com/kentik/dashboard-api-golang/client/networks"
	"github.com/kentik/dashboard-api-golang/client/organizations"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

type MerakiClient struct {
	log      logger.ContextL
	jchfChan chan []*kt.JCHF
	conf     *kt.SnmpDeviceConfig
	gconf    *kt.SnmpGlobalConfig
	metrics  *kt.SnmpDeviceMetric
	client   *apiclient.DashboardAPIGolang
	auth     runtime.ClientAuthInfoWriter
	orgs     []orgDesc
	timeout  time.Duration
	cache    *clientCache
	maxRetry int
}

type orgDesc struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	networks map[string]networkDesc
	org      *organizations.GetOrganizationsOKBodyItems0
}

type networkDesc struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	org  *organizations.GetOrganizationsOKBodyItems0
}

const (
	ControllerKey         = "meraki_controller_name"
	MerakiApiKey          = "KENTIK_MERAKI_API_KEY"
	DeviceCacheDuration   = time.Duration(24) * time.Hour
	MAX_TIMEOUT_RETRY     = 10 // Don't retry a call more than this many times.
	MAX_TIMEOUT_SEC       = 5  // Sleep this many sec each 429.
	DEFAULT_TIMEOUT_RETRY = 2
)

func NewMerakiClient(jchfChan chan []*kt.JCHF, gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, log logger.ContextL) (*MerakiClient, error) {
	c := MerakiClient{
		log:      log,
		jchfChan: jchfChan,
		conf:     conf,
		gconf:    gconf,
		metrics:  metrics,
		orgs:     []orgDesc{},
		auth:     httptransport.APIKeyAuth("X-Cisco-Meraki-API-Key", "header", kt.LookupEnvString(MerakiApiKey, conf.Ext.MerakiConfig.ApiKey)),
		timeout:  30 * time.Second,
		cache:    newClientCache(log),
		maxRetry: conf.Ext.MerakiConfig.MaxAPIRetry,
	}

	host := conf.Ext.MerakiConfig.Host
	if host == "" {
		host = apiclient.DefaultHost
	}

	// Figure out max retries.
	if c.maxRetry == 0 || c.maxRetry > MAX_TIMEOUT_RETRY {
		c.maxRetry = DEFAULT_TIMEOUT_RETRY
	}

	// Figure out global or local timeout here.
	if conf.TimeoutMS > 0 {
		c.timeout = time.Duration(conf.TimeoutMS) * time.Millisecond
	} else if gconf.TimeoutMS > 0 {
		c.timeout = time.Duration(gconf.TimeoutMS) * time.Millisecond
	}
	c.log.Infof("Using a timeout of %v", c.timeout)

	trans := apiclient.DefaultTransportConfig().WithHost(host)
	client := apiclient.NewHTTPClientWithConfig(nil, trans)
	c.log.Infof("Verifying %s connectivity", c.GetName())

	// First, list out all of the organizations present.
	params := organizations.NewGetOrganizationsParamsWithTimeout(c.timeout)
	prod, err := client.Organizations.GetOrganizations(params, c.auth)
	if err != nil {
		return nil, fmt.Errorf("Invalid API Key. Cannot get organizations. Check config and try again")
	}

	// There's some options which are disabled now, we need to check and error.
	if c.conf.Ext.MerakiConfig.MonitorDevices &&
		!c.conf.Ext.MerakiConfig.Prefs["device_status_only"] {
		return nil, fmt.Errorf("monitor_devices option is not supported for Meraki in this version.")
	}

	if c.conf.Ext.MerakiConfig.MonitorNetworkClients {
		return nil, fmt.Errorf("monitor_clients option is not supported for Meraki in this version.")
	}

	orgs := []*regexp.Regexp{}
	nets := []*regexp.Regexp{}
	c.client = client
	for _, org := range conf.Ext.MerakiConfig.Orgs {
		re := regexp.MustCompile(org)
		orgs = append(orgs, re)
	}
	for _, net := range conf.Ext.MerakiConfig.Networks {
		re := regexp.MustCompile(net)
		nets = append(nets, re)
	}

	numNets := 0
	for _, org := range prod.GetPayload() {
		lorg := org

		if len(orgs) > 0 {
			foundOrg := false
			for _, orgR := range orgs {
				if orgR.MatchString(org.Name) {
					foundOrg = true
					break
				}
			}
			if !foundOrg {
				continue // This organization isn't opted in.
			}
		}

		// Now list the networks for this org.
		netSet := map[string]networkDesc{}
		numAdded, err := c.getOrgNetworks(netSet, "", lorg, nets, 0)
		if err != nil {
			c.log.Warnf("Skipping organization %s because it does not have permission to list networks.", org.Name)
			continue
		}
		numNets += numAdded

		if len(netSet) > 0 { // Only add this org in to track if it has some networks.
			c.log.Infof("Adding organization %s to list to track", org.Name)
			c.orgs = append(c.orgs, orgDesc{
				ID:       org.ID,
				Name:     org.Name,
				networks: netSet,
				org:      lorg,
			})
		}
	}

	// If we get this far, we have a list of things to look at.
	c.log.Infof("%s connected to API for %s with %d organization(s) and %d network(s).", c.GetName(), conf.DeviceName, len(c.orgs), numNets)

	return &c, nil
}

var reLink = regexp.MustCompile(`startingAfter=(.*)>;`)

func getNextLink(linkSet string) string {
	for _, link := range strings.Split(linkSet, ",") {
		if strings.Contains(link, "rel=next") {
			links := reLink.FindStringSubmatch(link)
			if len(links) > 1 {
				return links[1]
			}
		}
	}
	return ""
}

func (c *MerakiClient) getOrgNetworks(netSet map[string]networkDesc, nextToken string,
	org *organizations.GetOrganizationsOKBodyItems0, nets []*regexp.Regexp, numNets int) (int, error) {

	perPageLimit := int64(100) // Seems like a good default.
	params := organizations.NewGetOrganizationNetworksParamsWithTimeout(c.timeout)
	params.SetOrganizationID(org.ID)
	params.SetPerPage(&perPageLimit)
	if nextToken != "" {
		params.SetStartingAfter(&nextToken)
	}
	prod, err := c.client.Organizations.GetOrganizationNetworks(params, c.auth)
	if err != nil {
		return 0, err
	}

	b, err := json.Marshal(prod.GetPayload())
	if err != nil {
		return 0, err
	}
	var networks []networkDesc
	err = json.Unmarshal(b, &networks)
	if err != nil {
		return 0, err
	}

	for _, network := range networks {
		if len(nets) > 0 {
			foundNet := false
			for _, net := range nets {
				if net.MatchString(network.Name) {
					foundNet = true
					break
				}
			}
			if !foundNet {
				continue // This network isn't opted in.
			}
		}
		network.org = org
		c.log.Infof("Adding network %s %s to list to track", network.Name, network.ID)
		netSet[network.ID] = network
		numNets++
	}

	// Recursion!
	nextLink := getNextLink(prod.Link)
	if nextLink != "" {
		return c.getOrgNetworks(netSet, nextLink, org, nets, numNets)
	} else {
		return numNets, nil
	}
}

func (c *MerakiClient) GetName() string {
	return "Meraki API"
}

func (c *MerakiClient) Run(ctx context.Context, dur time.Duration) {
	poll := time.NewTicker(dur)
	defer poll.Stop()

	doUplinks := c.conf.Ext.MerakiConfig.MonitorUplinks
	doDeviceClients := c.conf.Ext.MerakiConfig.MonitorDevices &&
		!c.conf.Ext.MerakiConfig.Prefs["device_status_only"]
	doDeviceStatus := c.conf.Ext.MerakiConfig.MonitorDevices &&
		(c.conf.Ext.MerakiConfig.Prefs["device_status_only"] || c.conf.Ext.MerakiConfig.Prefs["device_status"])
	doOrgChanges := c.conf.Ext.MerakiConfig.MonitorOrgChanges
	doNetworkClients := c.conf.Ext.MerakiConfig.MonitorNetworkClients
	doVpnStatus := c.conf.Ext.MerakiConfig.MonitorVpnStatus
	doNetworkAttr := c.conf.Ext.MerakiConfig.Prefs["show_network_attr"]
	if !doUplinks && !doDeviceClients && !doDeviceStatus && !doOrgChanges && !doNetworkClients && !doVpnStatus && !doNetworkAttr {
		doUplinks = true
	}
	c.log.Infof("Running Every %v with uplinks=%v, device_clients=%v, device_status=%v, orgs=%v, networks=%v, vpn_status=%v, network_attr=%v",
		dur, doUplinks, doDeviceClients, doDeviceStatus, doOrgChanges, doNetworkClients, doVpnStatus, doNetworkAttr)

	for {
		select {

		// Track the counters here, to convert from raw counters to differences
		case _ = <-poll.C:
			if doOrgChanges {
				if res, err := c.getOrgChanges(dur); err != nil {
					c.log.Infof("Meraki cannot get Organization Changes: %v", err)
				} else if len(res) > 0 {
					c.jchfChan <- res
				}
			}

			if doDeviceClients {
				if res, err := c.getDeviceClients(dur); err != nil {
					c.log.Infof("Meraki cannot get Device Client Info: %v", err)
				} else if len(res) > 0 {
					c.jchfChan <- res
				}
			}

			if doDeviceStatus {
				if res, err := c.getDeviceStatus(dur); err != nil {
					c.log.Infof("Meraki cannot get Device Status Info: %v", err)
				} else if len(res) > 0 {
					c.jchfChan <- res
				}
			}

			if doNetworkClients {
				if res, err := c.getNetworkClients(dur); err != nil {
					c.log.Infof("Meraki cannot get Network Client Info: %v", err)
				} else if len(res) > 0 {
					c.jchfChan <- res
				}
			}

			if doUplinks {
				if res, err := c.getUplinks(dur); err != nil {
					c.log.Infof("Meraki cannot get Uplink Info: %v", err)
				} else if len(res) > 0 {
					c.jchfChan <- res
				}
			}

			if doVpnStatus {
				if res, err := c.getVpnStatus(dur); err != nil {
					c.log.Infof("Meraki cannot get vpn status Info: %v", err)
				} else if len(res) > 0 {
					c.jchfChan <- res
				}
			}

			if doNetworkAttr {
				if res, err := c.getNetworkAttr(dur); err != nil {
					c.log.Infof("Meraki cannot get network attr: %v", err)
				} else if len(res) > 0 {
					c.jchfChan <- res
				}
			}

		case <-ctx.Done():
			c.log.Infof("Meraki Poll Done")
			return
		}
	}
}

type orgLog struct {
	TimeStamp   time.Time `json:"ts"`
	AdminName   string    `json:"adminName"`
	NetworkName string    `json:"networkName"`
	Label       string    `json:"label"`
	NewValue    string    `json:"newValue"`
	AdminEmail  string    `json:"adminEmail"`
	AdminId     string    `json:"adminId"`
	NetworkId   string    `json:"networkId"`
	Page        string    `json:"page"`
	OldValue    string    `json:"oldValue"`
}

func (c *MerakiClient) getOrgChanges(dur time.Duration) ([]*kt.JCHF, error) {
	startTime := time.Now().Add(-1 * dur)
	startTimeStr := fmt.Sprintf("%v", startTime.Unix())

	res := []*kt.JCHF{}
	for _, org := range c.orgs {
		params := organizations.NewGetOrganizationConfigurationChangesParamsWithTimeout(c.timeout)
		params.SetOrganizationID(org.ID)
		params.SetT0(&startTimeStr)

		prod, err := c.client.Organizations.GetOrganizationConfigurationChanges(params, c.auth)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return nil, err
		}

		var logs []*orgLog
		err = json.Unmarshal(b, &logs)
		if err != nil {
			return nil, err
		}

		for _, lg := range logs {
			// Filter for networks here.
			if _, ok := org.networks[lg.NetworkId]; !ok {
				continue
			}
			res = append(res, c.parseOrgLog(lg, org))
		}
	}

	return res, nil
}

func (c *MerakiClient) parseOrgLog(l *orgLog, org orgDesc) *kt.JCHF {
	dst := kt.NewJCHF()
	dst.Timestamp = l.TimeStamp.Unix()

	dst.CustomStr = map[string]string{
		ControllerKey:  c.conf.DeviceName,
		"admin_name":   l.AdminName,
		"network_name": l.NetworkName,
		"label":        l.Label,
		"new_value":    l.NewValue,
		"admin_email":  l.AdminEmail,
		"admin_id":     l.AdminId,
		"network_id":   l.NetworkId,
		"page":         l.Page,
		"old_value":    l.OldValue,
		"org_name":     org.Name,
		"org_id":       org.ID,
	}
	dst.EventType = kt.KENTIK_EVENT_EXT // This gets sent as event, not metric.
	dst.Provider = kt.ProviderMerakiCloud

	c.conf.SetUserTags(dst.CustomStr)
	return dst
}

type client struct {
	Usage              map[string]float64 `json:"usage"`
	ID                 string             `json:"id"`
	Description        string             `json:"description"`
	Mac                string             `json:"mac"`
	IP                 string             `json:"ip"`
	User               string             `json:"user"`
	Vlan               string             `json:"vlan"`
	NamedVlan          string             `json:"namedVlan"`
	IPv6               string             `json:"ip6"`
	Manufacturer       string             `json:"manufacturer"`
	DeviceType         string             `json:"deviceTypePrediction"`
	RecentDeviceName   string             `json:"recentDeviceName"`
	RecentDeviceSerial string             `json:"recentDeviceSerial"`
	RecentDeviceMac    string             `json:"recentDeviceMac"`
	SSID               string             `json:"ssid"`
	Status             string             `json:"status"`
	MdnsName           string             `json:"mdnsName"`
	DhcpHostname       string             `json:"dhcpHostname"`
	network            string
	orgName            string
	orgId              string
	device             networkDevice
	appUsage           []appUsage
}

func (c *MerakiClient) getNetworkClients(dur time.Duration) ([]*kt.JCHF, error) {
	clientSet := []*client{}
	durs := float32(dur.Seconds())
	for _, org := range c.orgs {
		for _, network := range org.networks {
			params := networks.NewGetNetworkClientsParamsWithTimeout(c.timeout)
			params.SetNetworkID(network.ID)
			params.SetTimespan(&durs)

			prod, err := c.client.Networks.GetNetworkClients(params, c.auth)
			if err != nil {
				c.log.Warnf("Cannot get network clients for %s: %v", network.Name, err)
				continue
			}

			b, err := json.Marshal(prod.GetPayload())
			if err != nil {
				return nil, err
			}

			var clients []*client
			err = json.Unmarshal(b, &clients)
			if err != nil {
				return nil, err
			}
			for _, client := range clients {
				client.network = network.Name
				client.orgName = org.Name
				client.orgId = org.ID
				clientSet = append(clientSet, client) // Toss these all in together
			}
		}
	}

	return c.parseClients(clientSet)
}

type networkDevice struct {
	Name      string   `json:"name"`
	Serial    string   `json:"serial"`
	Mac       string   `json:"mac"`
	Model     string   `json:"model"`
	Notes     string   `json:"notes"`
	LanIP     string   `json:"lanIp"`
	Tags      []string `json:"tags"`
	NetworkID string   `json:"networkId"`
	Firmware  string   `json:"firmware"`
	network   networkDesc
}

func (c *MerakiClient) getNetworkDevices() (map[string][]networkDevice, error) {
	deviceSet := map[string][]networkDevice{}
	for _, org := range c.orgs {
		for _, network := range org.networks {
			params := networks.NewGetNetworkDevicesParamsWithTimeout(c.timeout)
			params.SetNetworkID(network.ID)

			prod, err := c.client.Networks.GetNetworkDevices(params, c.auth)
			if err != nil {
				return nil, err
			}

			b, err := json.Marshal(prod.GetPayload())
			if err != nil {
				return nil, err
			}

			var devices []networkDevice
			err = json.Unmarshal(b, &devices)
			if err != nil {
				return nil, err
			}

			for i, _ := range devices {
				if devices[i].Name == "" {
					devices[i].Name = devices[i].Serial
				}
				devices[i].network = network
			}

			if len(devices) > 0 {
				deviceSet[network.Name] = devices
			}
		}
	}

	return deviceSet, nil
}

type appUsage struct {
	Application string `json:"application"`
	Received    int64  `json:"received"`
	Sent        int64  `json:"sent"`
}

type applicationUsage struct {
	AppUsage  []appUsage `json:"applicationUsage"`
	ClientId  string     `json:"clientId"`
	ClientIp  string     `json:"clientIp"`
	ClientMac string     `json:"clientMac"`
}

func (c *MerakiClient) getDeviceClientApplications(dur time.Duration, clients []*client, network networkDesc) error {
	// Max of 10 clients to check application usage per call.
	getApps := func(client *client) error {
		params := networks.NewGetNetworkClientsApplicationUsageParamsWithTimeout(c.timeout)
		params.SetNetworkID(network.ID)
		params.SetClients(client.ID)

		prod, err := c.client.Networks.GetNetworkClientsApplicationUsage(params, c.auth)
		if err != nil {
			return err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return err
		}

		var appUsage []applicationUsage
		err = json.Unmarshal(b, &appUsage)
		if err != nil {
			return err
		}

		for _, appU := range appUsage {
			if appU.ClientId == client.ID {
				client.appUsage = appU.AppUsage
				break
			}
		}

		return nil
	}

	// Have to go 1 by 1 here because some clients arn't valid in the app call for some reason.
	for _, client := range clients {
		err := getApps(client)
		if err != nil {
			c.log.Debugf("Cannot get application usage for %s %v: %v", network.ID, client.ID, err)
		}
	}

	return nil
}

func (c *MerakiClient) getDeviceClients(dur time.Duration) ([]*kt.JCHF, error) {
	networkDevs, err := c.getNetworkDevices()
	if err != nil {
		if strings.Contains(err.Error(), "(status 400)") { // There are no valid network devices to worry about here.
			return nil, nil
		}
		return nil, err
	}

	c.log.Infof("Got devices for %d networks", len(networkDevs))
	clientSet := []*client{}
	durs := float32(dur.Seconds())
	for network, deviceSet := range networkDevs {
		c.log.Infof("Looking at %d devices for network %s", len(deviceSet), network)
		for _, device := range deviceSet {
			params := devices.NewGetDeviceClientsParamsWithTimeout(c.timeout)
			params.SetSerial(device.Serial)
			params.SetTimespan(&durs)

			prod, err := c.client.Devices.GetDeviceClients(params, c.auth)
			if err != nil {
				c.log.Warnf("Cannot get device clients for %s: %v", device.Serial, err)
				continue
			}

			b, err := json.Marshal(prod.GetPayload())
			if err != nil {
				return nil, err
			}

			var clients []*client
			err = json.Unmarshal(b, &clients)
			if err != nil {
				return nil, err
			}

			// Now, look up application usage for this client set.
			err = c.getDeviceClientApplications(dur, clients, device.network)
			if err != nil {
				c.log.Warnf("Cannot get application usage for %s: %v", device.network.ID, err)
			}

			for _, client := range clients {
				client.network = network
				client.device = device
				clientSet = append(clientSet, client) // Toss these all in together
			}
		}
	}

	return c.parseClients(clientSet)
}

func (c *MerakiClient) parseClients(cs []*client) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)

	makeJCHF := func(client *client) *kt.JCHF {
		dst := kt.NewJCHF()
		if client.IPv6 != "" {
			dst.DstAddr = client.IPv6
		} else {
			dst.DstAddr = client.IP
		}
		dst.CustomStr = map[string]string{
			ControllerKey:        c.conf.DeviceName,
			"network":            client.network,
			"client_id":          client.ID,
			"description":        client.Description,
			"status":             client.Status,
			"vlan_name":          client.NamedVlan,
			"client_mac_addr":    client.Mac,
			"user":               client.User,
			"manufacturer":       client.Manufacturer,
			"device_type":        client.DeviceType,
			"recent_device_name": client.RecentDeviceName,
			"device_mac_addr":    client.RecentDeviceMac,
			"device_serial":      client.RecentDeviceSerial,
			"ssid":               client.SSID,
			"dhcp_hostname":      client.DhcpHostname,
			"mdns_name":          client.MdnsName,
			"vlan":               client.Vlan,
			"org_name":           client.orgName,
			"org_id":             client.orgId,
		}

		dst.CustomBigInt = map[string]int64{}
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		dst.Provider = kt.ProviderMerakiCloud

		if client.device.Serial != "" {
			dst.DeviceName = client.device.Name // Here, device is this device's name.
			dst.SrcAddr = client.device.LanIP
			dst.CustomStr["device_serial"] = client.device.Serial
			dst.CustomStr["device_firmware"] = client.device.Firmware
			dst.CustomStr["device_mac_addr"] = client.device.Mac
			dst.CustomStr["device_tags"] = strings.Join(client.device.Tags, ",")
			dst.CustomStr["device_notes"] = client.device.Notes
			dst.CustomStr["device_model"] = client.device.Model
			dst.CustomStr["src_ip"] = client.device.LanIP
			if client.device.network.org != nil {
				dst.CustomStr["org_name"] = client.device.network.org.Name
				dst.CustomStr["org_id"] = client.device.network.org.ID
			}
		} else {
			dst.DeviceName = client.network // Here, device is this network's name.
			dst.SrcAddr = c.conf.DeviceIP
		}

		dst.Timestamp = time.Now().Unix()
		dst.CustomMetrics = map[string]kt.MetricInfo{}

		dst.CustomBigInt["SentTotal"] = int64(client.Usage["sent"] * 1000) // Unit is kilobytes, convert to bytes
		dst.CustomMetrics["SentTotal"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

		dst.CustomBigInt["RecvTotal"] = int64(client.Usage["recv"] * 1000) // Same, convert to bytes.
		dst.CustomMetrics["RecvTotal"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

		c.conf.SetUserTags(dst.CustomStr)

		return dst
	}

	for _, client := range cs {
		if len(client.appUsage) > 0 {
			// If there is app usage, record per app.
			for _, appU := range client.appUsage {
				dst := makeJCHF(client)

				dst.CustomStr["application"] = appU.Application
				dst.CustomBigInt["Sent"] = int64(appU.Sent * 1000) // Unit is kilobytes, convert to bytes
				dst.CustomMetrics["Sent"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

				dst.CustomBigInt["Recv"] = int64(appU.Received * 1000) // Same, convert to bytes.
				dst.CustomMetrics["Recv"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

				res = append(res, dst)
			}
		} else {
			// Just record totals
			dst := makeJCHF(client)
			res = append(res, dst)
		}
	}

	return res, nil
}

type signalStat struct {
	Rsrp string `json:"rsrp"`
	Rsrq string `json:"rsrq"`
}

// Needed because sometimes meraki makes this a struct and other times an empty array, just to be f-ing annoying.
type signalStatAlias signalStat

func (a *signalStat) UnmarshalJSON(data []byte) error {
	var stat = signalStatAlias{}
	err := json.Unmarshal(data, &stat)
	if err != nil {
		set := []signalStatAlias{}
		err := json.Unmarshal(data, &set)
		if err != nil {
			return err
		}

		if len(set) > 0 {
			*a = signalStat(set[0])
		}
	} else {
		*a = signalStat(stat)
	}

	return nil
}

type uplink struct {
	PrimaryDNS     string     `json:"primaryDns"`
	SecondaryDNS   string     `json:"secondaryDns"`
	IpAssignedBy   string     `json:"ipAssignedBy"`
	Interface      string     `json:"interface"`
	Status         string     `json:"status"`
	IP             string     `json:"ip"`
	Gateway        string     `json:"gateway"`
	PublicIP       string     `json:"publicIp"`
	Provider       string     `json:"provider"`
	ICCID          string     `json:"iccid"`
	ConnectionType string     `json:"connectionType"`
	Model          string     `json:"model"`
	APN            string     `json:"apn"`
	SignalStat     signalStat `json:"signalStat"`
	DNS1           string     `json:"dns1"`
	DNS2           string     `json:"dns2"`
	SignalType     string     `json:"signalType"`
	Usage          *appliance.GetOrganizationApplianceUplinksUsageByNetworkOKBodyItems0ByUplinkItems0
	LatencyLoss    deviceUplinkLatency
}

func (du *deviceUplink) SetLatencyLoss(u deviceUplinkLatency) {
	found := false
	for i, ul := range du.Uplinks {
		if ul.Interface == u.Uplink {
			du.Uplinks[i].LatencyLoss = u
			found = true
		}
	}
	if !found {
		du.Uplinks = append(du.Uplinks, uplink{
			Interface:   u.Uplink,
			LatencyLoss: u,
		})
	}
}

type deviceUplink struct {
	Serial         string    `json:"serial"`
	Model          string    `json:"model"`
	LastReportedAt time.Time `json:"lastReportedAt"`
	Uplinks        []uplink  `json:"uplinks"`
	NetworkID      string    `json:"networkId"`
	network        networkDesc
}

/**
1) Get all the uplinks with status
2) Get the latency for these uplinks for any which have them
3) Get usage for each network.

*/

func (c *MerakiClient) getUplinks(dur time.Duration) ([]*kt.JCHF, error) {

	var getUplinkStatus func(nextToken string, org orgDesc, uplinks *map[string]deviceUplink, timeouts int) error
	getUplinkStatus = func(nextToken string, org orgDesc, uplinks *map[string]deviceUplink, timeouts int) error {
		params := organizations.NewGetOrganizationUplinksStatusesParamsWithTimeout(c.timeout)
		params.SetOrganizationID(org.ID)
		if nextToken != "" {
			params.SetStartingAfter(&nextToken)
		}

		prod, err := c.client.Organizations.GetOrganizationUplinksStatuses(params, c.auth)
		if err != nil {
			if strings.Contains(err.Error(), "(status 429)") && timeouts < c.maxRetry {
				sleepDur := time.Duration(MAX_TIMEOUT_SEC) * time.Second
				c.log.Warnf("Uplink Status: %s 429, sleeping %v", org.Name, sleepDur)
				time.Sleep(sleepDur) // For right now guess on this, need to add 429 to spec.
				timeouts++
				return getUplinkStatus(nextToken, org, uplinks, timeouts)
			}
			return err
		}

		// Store these for some tail recursion.
		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return err
		}

		var raw []deviceUplink
		err = json.Unmarshal(b, &raw)
		if err != nil {
			return err
		}

		for _, uplink := range raw {
			// Filter for networks here.
			if _, ok := org.networks[uplink.NetworkID]; !ok {
				continue
			}
			uplink.network = org.networks[uplink.NetworkID]
			(*uplinks)[uplink.Serial] = uplink
		}

		// Recursion!
		nextLink := getNextLink(prod.Link)
		if nextLink != "" {
			return getUplinkStatus(nextLink, org, uplinks, timeouts)
		} else {
			return nil
		}
	}

	uplinks := map[string]deviceUplink{}
	for _, org := range c.orgs {
		err := getUplinkStatus("", org, &uplinks, 0)
		if err != nil {
			if strings.Contains(err.Error(), "(status 400)") { // There are no valid uplinks to worry about here.
				return nil, nil
			}
			return nil, err
		}
	}

	if !c.conf.Ext.MerakiConfig.Prefs["hide_uplink_usage"] {
		err := c.getUplinkUsage(dur, uplinks)
		if err != nil {
			return nil, err
		}

		// Now, load latency for any of these which have them:
		err = c.getUplinkLatencyLoss(dur, uplinks)
		if err != nil {
			return nil, err
		}
	}

	return c.parseUplinks(uplinks)
}

type uplinkTS struct {
	TS          time.Time `json:"ts"`
	LossPercent float64   `json:"lossPercent"`
	LatencyMS   float64   `json:"latencyMs"`
}

type deviceUplinkLatency struct {
	Serial     string     `json:"serial"`
	NetworkID  string     `json:"networkId"`
	Uplink     string     `json:"uplink"`
	IP         string     `json:"ip"`
	TimeSeries []uplinkTS `json:"timeSeries"`
}

func (c *MerakiClient) getUplinkLatencyLoss(dur time.Duration, uplinkMap map[string]deviceUplink) error {

	for _, org := range c.orgs {
		params := organizations.NewGetOrganizationDevicesUplinksLossAndLatencyParamsWithTimeout(c.timeout)
		params.SetOrganizationID(org.ID)

		prod, err := c.client.Organizations.GetOrganizationDevicesUplinksLossAndLatency(params, c.auth)
		if err != nil {
			return err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return err
		}

		var uplinks []deviceUplinkLatency
		err = json.Unmarshal(b, &uplinks)
		if err != nil {
			return err
		}

		for _, uplink := range uplinks {
			if uu, ok := uplinkMap[uplink.Serial]; ok { // We've matched the serial, now add this loss.
				uu.SetLatencyLoss(uplink)
			}
		}
	}

	return nil
}

func (c *MerakiClient) getUplinkUsage(dur time.Duration, uplinkMap map[string]deviceUplink) error {

	var getUsage func(params *appliance.GetOrganizationApplianceUplinksUsageByNetworkParams, org orgDesc) ([]*appliance.GetOrganizationApplianceUplinksUsageByNetworkOKBodyItems0, error)
	getUsage = func(params *appliance.GetOrganizationApplianceUplinksUsageByNetworkParams, org orgDesc) ([]*appliance.GetOrganizationApplianceUplinksUsageByNetworkOKBodyItems0, error) {
		prod, err := c.client.Appliance.GetOrganizationApplianceUplinksUsageByNetwork(params, c.auth)
		if err != nil {
			if strings.Contains(err.Error(), "status 429") {
				sleepDur := time.Duration(MAX_TIMEOUT_SEC) * time.Second
				c.log.Infof("Uplink Usage: %s 429, sleeping %v", org.Name, sleepDur)
				time.Sleep(sleepDur) // For right now guess on this, need to add 429 to spec.
				return getUsage(params, org)
			} else {
				c.log.Warnf("Cannot get Uplink Usage: %s %v", org.Name, err)
				return nil, err
			}
		}

		return prod.GetPayload(), nil
	}

	ts := float32(dur.Seconds())
	for _, org := range c.orgs {
		params := appliance.NewGetOrganizationApplianceUplinksUsageByNetworkParamsWithTimeout(c.timeout)
		params.SetOrganizationID(org.ID)
		params.SetTimespan(&ts)

		uplinkUsage, err := getUsage(params, org)
		if err != nil {
			continue
		}

		if len(uplinkUsage) > 0 {
			for _, network := range uplinkUsage {
				for _, uplink := range network.ByUplink {
					if _, ok := uplinkMap[uplink.Serial]; ok {
						for i, _ := range uplinkMap[uplink.Serial].Uplinks {
							if uplinkMap[uplink.Serial].Uplinks[i].Interface == uplink.Interface {
								uplinkMap[uplink.Serial].Uplinks[i].Usage = uplink
							}
						}
					}
				}
			}
		}
	}

	return nil
}

var uplinkStatus = map[string]int64{
	"active":        1,
	"failed":        2,
	"not connected": 3,
	"ready":         4,
}

func (c *MerakiClient) parseUplinks(uplinkMap map[string]deviceUplink) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)
	for _, device := range uplinkMap {
		for _, uplink := range device.Uplinks {
			dst := kt.NewJCHF()
			dst.SrcAddr = uplink.IP
			dst.DeviceName = strings.Join([]string{device.network.Name, uplink.Interface}, ".")

			dst.CustomStr = map[string]string{
				ControllerKey:       c.conf.DeviceName,
				"network":           device.network.Name,
				"network_id":        device.NetworkID,
				"serial":            device.Serial,
				"model":             device.Model,
				"status":            uplink.Status,
				"connection_type":   uplink.ConnectionType,
				"interface":         uplink.Interface,
				"cellular_provider": uplink.Provider,
				"signal_type":       uplink.SignalType,
				"signal_rsrp":       uplink.SignalStat.Rsrp,
				"signal_rsrq":       uplink.SignalStat.Rsrq,
				"org_name":          device.network.org.Name,
				"org_id":            device.network.org.ID,
			}
			dst.CustomInt = map[string]int32{}
			dst.CustomBigInt = map[string]int64{}
			dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
			dst.Provider = kt.ProviderMerakiCloud

			dst.Timestamp = time.Now().Unix()
			dst.CustomMetrics = map[string]kt.MetricInfo{}

			dst.CustomBigInt["Status"] = uplinkStatus[uplink.Status]
			dst.CustomMetrics["Status"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.uplink_status", Type: "meraki.uplink_status"}

			if !c.conf.Ext.MerakiConfig.Prefs["hide_uplink_usage"] {
				if uplink.Usage != nil {
					dst.CustomBigInt["Sent"] = uplink.Usage.Sent
					dst.CustomMetrics["Sent"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.uplinks", Type: "meraki.uplinks"}

					dst.CustomBigInt["Recv"] = uplink.Usage.Received
					dst.CustomMetrics["Recv"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.uplinks", Type: "meraki.uplinks"}
				}

				latency, loss := getLatencyLoss(uplink)
				dst.CustomBigInt["LatencyMS"] = int64(latency * 1000.)
				dst.CustomMetrics["LatencyMS"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.uplinks", Type: "meraki.uplinks", Format: kt.FloatMS}

				dst.CustomBigInt["LossPct"] = int64(loss * 1000.)
				dst.CustomMetrics["LossPct"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.uplinks", Type: "meraki.uplinks", Format: kt.FloatMS}
			}

			c.conf.SetUserTags(dst.CustomStr)
			res = append(res, dst)
		}
	}

	return res, nil
}

func (c *MerakiClient) getNetworkAttr(dur time.Duration) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)
	for _, org := range c.orgs {
		for _, network := range org.networks {
			dst := kt.NewJCHF()
			dst.DeviceName = network.Name

			dst.CustomStr = map[string]string{
				ControllerKey: c.conf.DeviceName,
				"network":     network.Name,
				"network_id":  network.ID,
				"org_name":    network.org.Name,
				"org_id":      network.org.ID,
			}
			for i, tag := range network.Tags {
				dst.CustomStr[fmt.Sprintf("tag_%d", i+1)] = tag
			}

			dst.CustomInt = map[string]int32{}
			dst.CustomBigInt = map[string]int64{}
			dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
			dst.Provider = kt.ProviderMerakiCloud

			dst.Timestamp = time.Now().Unix()
			dst.CustomMetrics = map[string]kt.MetricInfo{}

			dst.CustomBigInt["Count"] = 1
			dst.CustomMetrics["Count"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.network", Type: "meraki.network"}

			c.conf.SetUserTags(dst.CustomStr)
			res = append(res, dst)
		}

		// Now add in 1 count per organization.
		dst := kt.NewJCHF()
		dst.DeviceName = org.Name

		dst.CustomStr = map[string]string{
			ControllerKey: c.conf.DeviceName,
			"org_name":    org.Name,
			"org_id":      org.ID,
		}

		dst.CustomInt = map[string]int32{}
		dst.CustomBigInt = map[string]int64{}
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		dst.Provider = kt.ProviderMerakiCloud

		dst.Timestamp = time.Now().Unix()
		dst.CustomMetrics = map[string]kt.MetricInfo{}

		dst.CustomBigInt["Count"] = 1
		dst.CustomMetrics["Count"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.organization", Type: "meraki.organization"}

		c.conf.SetUserTags(dst.CustomStr)
		res = append(res, dst)
	}

	return res, nil
}

func getLatencyLoss(u uplink) (float64, float64) {
	latency := float64(0)
	loss := float64(0)

	if len(u.LatencyLoss.TimeSeries) == 0 { // Kinda a noop but 0 is best here.
		return 0, 0
	}

	for _, t := range u.LatencyLoss.TimeSeries {
		latency += t.LatencyMS
		loss += t.LossPercent
	}

	return latency / float64(len(u.LatencyLoss.TimeSeries)), loss / float64(len(u.LatencyLoss.TimeSeries))
}

/**
1) Get all the vpns with status.
*/

type subnet struct {
	Subnet string `json:"subnet"`
	Name   string `json:"name"`
}

type vpnPeer struct {
	NetworkID    string `json:"networkId"`
	NetworkName  string `json:"networkName"`
	Reachability string `json:"reachability"`
	Name         string `json:"name"`
	PublicIp     string `json:"publicIp"`
}

type vpnStatus struct {
	NetworkID          string    `json:"networkId"`
	NetworkName        string    `json:"networkName"`
	DeviceSerial       string    `json:"deviceSerial"`
	DeviceStatus       string    `json:"deviceStatus"`
	Uplinks            []uplink  `json:"uplinks"`
	VpnMode            string    `json:"vpnMode"`
	ExportedSubnets    []subnet  `json:"exportedSubnets"`
	MerakiVpnPeers     []vpnPeer `json:"merakiVpnPeers"`
	ThirdPartyVpnPeers []vpnPeer `json:"thirdPartyVpnPeers"`
	org                orgDesc
}

func (c *MerakiClient) getVpnStatus(dur time.Duration) ([]*kt.JCHF, error) {

	var getVpnStatus func(nextToken string, org orgDesc, vpns *[]*vpnStatus, timeouts int) error
	getVpnStatus = func(nextToken string, org orgDesc, vpns *[]*vpnStatus, timeouts int) error {
		params := appliance.NewGetOrganizationApplianceVpnStatusesParamsWithTimeout(c.timeout)
		params.SetOrganizationID(org.ID)
		if nextToken != "" {
			params.SetStartingAfter(&nextToken)
		}

		prod, err := c.client.Appliance.GetOrganizationApplianceVpnStatuses(params, c.auth)
		if err != nil {
			if strings.Contains(err.Error(), "(status 429)") && timeouts < c.maxRetry {
				sleepDur := time.Duration(MAX_TIMEOUT_SEC) * time.Second
				c.log.Warnf("Vpn Status: %s 429, sleeping %v", org.Name, sleepDur)
				time.Sleep(sleepDur) // For right now guess on this, need to add 429 to spec.
				timeouts++
				return getVpnStatus(nextToken, org, vpns, timeouts)
			}
			return err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return err
		}

		var vpnSet []*vpnStatus
		err = json.Unmarshal(b, &vpnSet)
		if err != nil {
			return err
		}

		// Store these for some tail recursion.
		filtered := make([]*vpnStatus, 0, len(vpnSet))
		for _, vpn := range vpnSet {
			if _, ok := org.networks[vpn.NetworkID]; !ok {
				continue
			}
			vpn.org = org
			filtered = append(filtered, vpn)
		}

		*vpns = append(*vpns, filtered...)

		// Recursion!
		nextLink := getNextLink(prod.Link)
		if nextLink != "" {
			return getVpnStatus(nextLink, org, vpns, timeouts)
		} else {
			return nil
		}
	}

	vpns := make([]*vpnStatus, 0)
	for _, org := range c.orgs {
		err := getVpnStatus("", org, &vpns, 0)
		if err != nil {
			if strings.Contains(err.Error(), "(status 400)") { // There are no valid vpns to worry about here.
				return nil, nil
			}
			return nil, err
		}
	}

	if len(vpns) == 0 {
		return nil, nil
	}

	return c.parseVpnStatus(vpns)
}

func (c *MerakiClient) parseVpnStatus(vpns []*vpnStatus) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)

	makeChf := func(vpn *vpnStatus) *kt.JCHF {
		dst := kt.NewJCHF()
		//dst.SrcAddr = uplink.PublicIP
		dst.DeviceName = vpn.DeviceSerial

		dst.CustomStr = map[string]string{
			ControllerKey: c.conf.DeviceName,
			"network":     vpn.NetworkName,
			"network_id":  vpn.NetworkID,
			"serial":      vpn.DeviceSerial,
			"status":      vpn.DeviceStatus,
			"vpn_mode":    vpn.VpnMode,
			"org_name":    vpn.org.Name,
			"org_id":      vpn.org.ID,
		}

		for _, uplink := range vpn.Uplinks {
			dst.CustomStr[uplink.Interface] = uplink.PublicIP
		}

		//for _, subnet := range vpn.ExportedSubnets {
		//	dst.CustomStr["subnets"+subnet.Name] = subnet.Subnet
		//}

		dst.CustomInt = map[string]int32{}
		dst.CustomBigInt = map[string]int64{}
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		dst.Provider = kt.ProviderMerakiCloud

		dst.Timestamp = time.Now().Unix()
		dst.CustomMetrics = map[string]kt.MetricInfo{}

		c.conf.SetUserTags(dst.CustomStr)

		return dst
	}

	for _, vpn := range vpns {

		// Basic status here.
		dst := makeChf(vpn)
		status := int64(0)
		if vpn.DeviceStatus == "online" { // Online is 1, others are 0.
			status = 1
		}
		dst.CustomBigInt["Status"] = status
		dst.CustomMetrics["Status"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.vpn_status", Type: "meraki.vpn_status"}

		res = append(res, dst)

		// Now add in metrics for peer reachibility
		if c.conf.Ext.MerakiConfig.Prefs["show_vpn_peers"] {
			for _, peer := range vpn.MerakiVpnPeers {
				dst := makeChf(vpn)
				dst.CustomStr["peer_name"] = peer.NetworkName
				dst.CustomStr["peer_network_id"] = peer.NetworkID
				dst.CustomStr["peer_reachability"] = peer.Reachability
				dst.CustomStr["peer_type"] = "Meraki"

				status := int64(0)
				if peer.Reachability == "reachable" { // Reachable is 1, others are 0.
					status = 1
				}
				dst.CustomBigInt["PeerStatus"] = status
				dst.CustomMetrics["PeerStatus"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.vpn_status", Type: "meraki.vpn_status"}

				res = append(res, dst)
			}

			for _, peer := range vpn.ThirdPartyVpnPeers {
				dst := makeChf(vpn)
				dst.CustomStr["peer_name"] = peer.Name
				dst.CustomStr["peer_public_ip"] = peer.PublicIp
				dst.CustomStr["peer_reachability"] = peer.Reachability
				dst.CustomStr["peer_type"] = "ThirdParty"

				status := int64(0)
				if peer.Reachability == "reachable" { // Reachable is 1, others are 0.
					status = 1
				}
				dst.CustomBigInt["PeerStatus"] = status
				dst.CustomMetrics["PeerStatus"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.vpn_status", Type: "meraki.vpn_status"}

				res = append(res, dst)
			}
		}
	}

	return res, nil
}

type deviceStatusWrapper struct {
	org         orgDesc
	NetworkName string `json:"networkName"`
	device      *organizations.GetOrganizationDevicesStatusesOKBodyItems0
	info        *organizations.GetOrganizationDevicesOKBodyItems0
}

func (c *MerakiClient) getDeviceStatus(dur time.Duration) ([]*kt.JCHF, error) {

	// Build up a map of product types to filter on.
	productTypes := map[string]bool{}
	for _, pt := range c.conf.Ext.MerakiConfig.ProductTypes {
		productTypes[pt] = true
	}

	var getDeviceStatus func(nextToken string, org orgDesc, devices *[]*deviceStatusWrapper, timeouts int) error
	getDeviceStatus = func(nextToken string, org orgDesc, devices *[]*deviceStatusWrapper, timeouts int) error {
		params := organizations.NewGetOrganizationDevicesStatusesParamsWithTimeout(c.timeout)
		params.SetOrganizationID(org.ID)
		if nextToken != "" {
			params.SetStartingAfter(&nextToken)
		}

		prod, err := c.client.Organizations.GetOrganizationDevicesStatuses(params, c.auth)
		if err != nil {
			if strings.Contains(err.Error(), "(status 429)") && timeouts < c.maxRetry {
				sleepDur := time.Duration(MAX_TIMEOUT_SEC) * time.Second
				c.log.Warnf("Device Status: %s 429, sleeping %v", org.Name, sleepDur)
				time.Sleep(sleepDur) // For right now guess on this, need to add 429 to spec.
				timeouts++
				return getDeviceStatus(nextToken, org, devices, timeouts)
			}
			return err
		}

		// Store these for some tail recursion.
		raw := prod.GetPayload()
		filtered := make([]*deviceStatusWrapper, 0, len(raw))
		for _, device := range raw {
			// Filter for networks here.
			if _, ok := org.networks[device.NetworkID]; !ok {
				continue
			}
			// Also filter on product types.
			if len(productTypes) > 0 && !productTypes[device.ProductType] {
				continue
			}
			nd := deviceStatusWrapper{
				device:      device,
				org:         org,
				NetworkName: org.networks[device.NetworkID].Name,
			}
			filtered = append(filtered, &nd)
		}

		*devices = append(*devices, filtered...)

		// Recursion!
		nextLink := getNextLink(prod.Link)
		if nextLink != "" {
			return getDeviceStatus(nextLink, org, devices, timeouts)
		} else {
			return nil
		}
	}

	devices := make([]*deviceStatusWrapper, 0)
	for _, org := range c.orgs {
		err := getDeviceStatus("", org, &devices, 0)
		if err != nil {
			if strings.Contains(err.Error(), "(status 400)") { // There are no valid devices to worry about here.
				return nil, nil
			}
			return nil, err
		}
	}

	if len(devices) == 0 {
		return nil, nil
	}

	// Next, get device info to add more info about these devices.
	err := c.getDeviceInfo(devices)
	if err != nil {
		return nil, err
	}

	return c.parseDeviceStatus(devices)
}

type clientCache struct {
	log            logger.ContextL
	deviceInfoTime time.Time
	deviceInfo     []*organizations.GetOrganizationDevicesOKBodyItems0
}

func newClientCache(log logger.ContextL) *clientCache {
	return &clientCache{
		log: log,
	}
}

func (c *clientCache) getDeviceInfo() ([]*organizations.GetOrganizationDevicesOKBodyItems0, bool) {
	if c.deviceInfoTime.Add(DeviceCacheDuration).Before(time.Now()) { // No information, cache invalid or old.
		return nil, false
	}

	return c.deviceInfo, true
}

func (c *clientCache) setDeviceInfo(infos []*organizations.GetOrganizationDevicesOKBodyItems0) {
	c.deviceInfo = infos
	c.deviceInfoTime = time.Now()
}

func (c *MerakiClient) getDeviceInfo(devices []*deviceStatusWrapper) error {
	// Build up a map of product types to filter on.
	productTypes := map[string]bool{}
	for _, pt := range c.conf.Ext.MerakiConfig.ProductTypes {
		productTypes[pt] = true
	}

	var getDeviceInfo func(nextToken string, org orgDesc, devices *[]*deviceStatusWrapper, cache *[]*organizations.GetOrganizationDevicesOKBodyItems0) error
	getDeviceInfo = func(nextToken string, org orgDesc, devices *[]*deviceStatusWrapper, cache *[]*organizations.GetOrganizationDevicesOKBodyItems0) error {
		params := organizations.NewGetOrganizationDevicesParamsWithTimeout(c.timeout)
		params.SetOrganizationID(org.ID)
		if nextToken != "" {
			params.SetStartingAfter(&nextToken)
		}

		prod, err := c.client.Organizations.GetOrganizationDevices(params, c.auth)
		if err != nil {
			return err
		}

		// Store these for some tail recursion.
		raw := prod.GetPayload()

		for _, device := range raw {
			// Filter for networks here.
			if _, ok := org.networks[device.NetworkID]; !ok {
				continue
			}
			// Also filter on product types.
			if len(productTypes) > 0 && !productTypes[device.ProductType] {
				continue
			}

			ld := device
			*cache = append(*cache, ld)
			for _, wrap := range *devices {
				if wrap.device.Serial == device.Serial {
					wrap.info = ld
				}
			}
		}

		// Recursion!
		nextLink := getNextLink(prod.Link)
		if nextLink != "" {
			return getDeviceInfo(nextLink, org, devices, cache)
		} else {
			return nil
		}
	}

	// First check the cache.
	if cc, ok := c.cache.getDeviceInfo(); ok {
		// We have a valid cache, don't use rest of call.
		for _, device := range cc {
			for _, wrap := range devices {
				if wrap.device.Serial == device.Serial {
					ld := device
					wrap.info = ld
				}
			}
		}
		return nil
	}

	// Need to build up cache case here.
	cache := []*organizations.GetOrganizationDevicesOKBodyItems0{}
	for _, org := range c.orgs {
		err := getDeviceInfo("", org, &devices, &cache)
		if err != nil {
			return err
		}
	}
	c.cache.setDeviceInfo(cache)

	return nil
}

func (c *MerakiClient) parseDeviceStatus(devices []*deviceStatusWrapper) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)

	makeChf := func(wrap *deviceStatusWrapper) *kt.JCHF {
		dst := kt.NewJCHF()
		dst.SrcAddr = wrap.device.PublicIP
		dst.DeviceName = wrap.device.Name

		dst.CustomStr = map[string]string{
			ControllerKey:  c.conf.DeviceName,
			"network":      wrap.NetworkName,
			"network_id":   wrap.device.NetworkID,
			"serial":       wrap.device.Serial,
			"status":       wrap.device.Status,
			"tags":         strings.Join(wrap.device.Tags, ","),
			"org_name":     wrap.org.Name,
			"org_id":       wrap.org.ID,
			"mac":          wrap.device.Mac,
			"model":        wrap.device.Model,
			"product_type": wrap.device.ProductType,
			"lat":          fmt.Sprintf("%f", wrap.info.Lat),
			"lng":          fmt.Sprintf("%f", wrap.info.Lng),
			"address":      wrap.info.Address,
			"notes":        wrap.info.Notes,
		}

		dst.CustomInt = map[string]int32{}
		dst.CustomBigInt = map[string]int64{}
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		dst.Provider = kt.ProviderMerakiCloud

		dst.Timestamp = time.Now().Unix()
		dst.CustomMetrics = map[string]kt.MetricInfo{}

		c.conf.SetUserTags(dst.CustomStr)

		return dst
	}

	for _, wrap := range devices {

		// Basic status here.
		dst := makeChf(wrap)
		status := int64(0)
		if wrap.device.Status == "online" { // Online is 1, others are 0.
			status = 1
		}
		dst.CustomBigInt["Status"] = status
		dst.CustomMetrics["Status"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.device_status", Type: "meraki.device_status"}

		res = append(res, dst)
	}

	return res, nil
}
