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
	metrics  *kt.SnmpDeviceMetric
	client   *apiclient.DashboardAPIGolang
	auth     runtime.ClientAuthInfoWriter
	orgs     []orgDesc
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

func NewMerakiClient(jchfChan chan []*kt.JCHF, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, log logger.ContextL) (*MerakiClient, error) {
	c := MerakiClient{
		log:      log,
		jchfChan: jchfChan,
		conf:     conf,
		metrics:  metrics,
		orgs:     []orgDesc{},
		auth:     httptransport.APIKeyAuth("X-Cisco-Meraki-API-Key", "header", conf.Ext.MerakiConfig.ApiKey),
	}

	host := conf.Ext.MerakiConfig.Host
	if host == "" {
		host = apiclient.DefaultHost
	}

	trans := apiclient.DefaultTransportConfig().WithHost(host)
	client := apiclient.NewHTTPClientWithConfig(nil, trans)
	c.log.Infof("Verifying %s connectivity", c.GetName())

	// First, list out all of the organizations present.
	params := organizations.NewGetOrganizationsParams()
	prod, err := client.Organizations.GetOrganizations(params, c.auth)
	if err != nil {
		return nil, err
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
			return nil, err
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
	params := organizations.NewGetOrganizationNetworksParams()
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
	doDevices := c.conf.Ext.MerakiConfig.MonitorDevices
	doOrgChanges := c.conf.Ext.MerakiConfig.MonitorOrgChanges
	doNetworkClients := c.conf.Ext.MerakiConfig.MonitorNetworkClients
	doVpnStatus := c.conf.Ext.MerakiConfig.MonitorVpnStatus
	if !doUplinks && !doDevices && !doOrgChanges && !doNetworkClients && !doVpnStatus {
		doUplinks = true
	}
	c.log.Infof("Running Every %v with uplinks=%v, devices=%v, orgs=%v, networks=%v, vpn status=%v", dur, doUplinks, doDevices, doOrgChanges, doNetworkClients, doVpnStatus)

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

			if doDevices {
				if res, err := c.getDeviceClients(dur); err != nil {
					c.log.Infof("Meraki cannot get Device Client Info: %v", err)
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
		for _, network := range org.networks {
			params := organizations.NewGetOrganizationConfigurationChangesParams()
			params.SetOrganizationID(org.ID)
			params.SetNetworkID(&(network.ID))
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
				res = append(res, c.parseOrgLog(lg, network, org))
			}
		}
	}

	return res, nil
}

func (c *MerakiClient) parseOrgLog(l *orgLog, network networkDesc, org orgDesc) *kt.JCHF {
	dst := kt.NewJCHF()
	dst.Timestamp = l.TimeStamp.Unix()

	dst.CustomStr = map[string]string{
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
			params := networks.NewGetNetworkClientsParams()
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
			params := networks.NewGetNetworkDevicesParams()
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
		params := networks.NewGetNetworkClientsApplicationUsageParams()
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
		return nil, err
	}

	c.log.Infof("Got devices for %d networks", len(networkDevs))
	clientSet := []*client{}
	durs := float32(dur.Seconds())
	for network, deviceSet := range networkDevs {
		c.log.Infof("Looking at %d devices for network %s", len(deviceSet), network)
		for _, device := range deviceSet {
			params := devices.NewGetDeviceClientsParams()
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
	network        *networkDesc
}

/**
1) Get all the uplinks with status
2) Get the latency for these uplinks for any which have them
3) Get usage for each network.

*/

func (c *MerakiClient) getUplinks(dur time.Duration) ([]*kt.JCHF, error) {

	uplinkSet := map[string]deviceUplink{}
	for _, org := range c.orgs {
		params := organizations.NewGetOrganizationUplinksStatusesParams()
		params.SetOrganizationID(org.ID)

		prod, err := c.client.Organizations.GetOrganizationUplinksStatuses(params, c.auth)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return nil, err
		}

		var uplinks []deviceUplink
		err = json.Unmarshal(b, &uplinks)
		if err != nil {
			return nil, err
		}

		for _, uplink := range uplinks {
			if network, ok := org.networks[uplink.NetworkID]; ok {
				lnet := network
				uplink.network = &lnet
			}

			if uplink.network == nil {
				// Skip this uplink because its from a network we don't care about.
				// c.log.Errorf("Missing Network for Uplink %s -- %s", uplink.NetworkID, uplink.Serial)
			} else {
				if _, ok := uplinkSet[uplink.Serial]; ok {
					c.log.Errorf("Duplicate Uplink %s", uplink.Serial)
				} else {
					uplinkSet[uplink.Serial] = uplink
				}
			}
		}
	}

	err := c.getUplinkUsage(dur, uplinkSet)
	if err != nil {
		return nil, err
	}

	// Now, load latency for any of these which have them:
	err = c.getUplinkLatencyLoss(dur, uplinkSet)
	if err != nil {
		return nil, err
	}

	return c.parseUplinks(uplinkSet)
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
		params := organizations.NewGetOrganizationDevicesUplinksLossAndLatencyParams()
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
			} else {
				// Not on a network we care about.
				// c.log.Errorf("Missing Uplink %s In LatencyLoss", uplink.Serial)
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
				c.log.Infof("Uplink Usage: %s 429, sleeping", org.Name)
				time.Sleep(3 * time.Second) // For right now guess on this, need to add 429 to spec.
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
		params := appliance.NewGetOrganizationApplianceUplinksUsageByNetworkParams()
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

func (c *MerakiClient) parseUplinks(uplinkMap map[string]deviceUplink) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)
	for _, device := range uplinkMap {
		for _, uplink := range device.Uplinks {
			dst := kt.NewJCHF()
			dst.SrcAddr = uplink.IP
			dst.DeviceName = strings.Join([]string{device.network.Name, uplink.Interface}, ".")

			dst.CustomStr = map[string]string{
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

			c.conf.SetUserTags(dst.CustomStr)
			res = append(res, dst)
		}
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

	var getVpnStatus func(nextToken string, org orgDesc, vpns *[]*vpnStatus) error
	getVpnStatus = func(nextToken string, org orgDesc, vpns *[]*vpnStatus) error {
		params := appliance.NewGetOrganizationApplianceVpnStatusesParams()
		params.SetOrganizationID(org.ID)
		if nextToken != "" {
			params.SetStartingAfter(&nextToken)
		}

		prod, err := c.client.Appliance.GetOrganizationApplianceVpnStatuses(params, c.auth)
		if err != nil {
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
			return getVpnStatus(nextLink, org, vpns)
		} else {
			return nil
		}
	}

	vpns := make([]*vpnStatus, 0)
	for _, org := range c.orgs {
		err := getVpnStatus("", org, &vpns)
		if err != nil {
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
			"network":    vpn.NetworkName,
			"network_id": vpn.NetworkID,
			"serial":     vpn.DeviceSerial,
			"status":     vpn.DeviceStatus,
			"vpn_mode":   vpn.VpnMode,
			"org_name":   vpn.org.Name,
			"org_id":     vpn.org.ID,
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
				dst.CustomStr["peer_reachablity"] = peer.Reachability
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
				dst.CustomStr["peer_reachablity"] = peer.Reachability
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
