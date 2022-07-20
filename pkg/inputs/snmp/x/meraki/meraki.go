package meraki

import (
	"context"
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
	client   *apiclient.MerakiDashboard
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
	ID   string `json:"id"`
	Name string `json:"name"`
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

	orgs := map[string]bool{}
	nets := map[string]bool{}
	for _, org := range conf.Ext.MerakiConfig.Orgs {
		orgs[org] = true
	}
	for _, net := range conf.Ext.MerakiConfig.Networks {
		nets[net] = true
	}

	numNets := 0
	for _, org := range prod.GetPayload() {
		lorg := org

		if len(orgs) > 0 && !orgs[org.Name] {
			continue // This organization isn't opted in.
		}

		// Now list the networks for this org.
		params := organizations.NewGetOrganizationNetworksParams()
		params.SetOrganizationID(org.ID)
		prod, err := client.Organizations.GetOrganizationNetworks(params, c.auth)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return nil, err
		}
		var networks []networkDesc
		err = json.Unmarshal(b, &networks)
		if err != nil {
			return nil, err
		}

		netSet := map[string]networkDesc{}
		for _, network := range networks {
			if len(nets) > 0 && !nets[network.Name] {
				continue // This network isn't opted in.
			}
			network.org = lorg
			c.log.Infof("Adding network %s %s to list to track", network.Name, network.ID)
			netSet[network.ID] = network
			numNets++
		}

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
	c.client = client

	return &c, nil
}

func (c *MerakiClient) GetName() string {
	return "Meraki API"
}

func (c *MerakiClient) Run(ctx context.Context, dur time.Duration) {
	poll := time.NewTicker(dur)
	defer poll.Stop()

	doUplinks := c.conf.Ext.MerakiConfig.MonitorUplinks
	doDevices := c.conf.Ext.MerakiConfig.MonitorDevices
	if !doUplinks && !doDevices {
		doUplinks = true
	}
	c.log.Infof("Running Every %v with uplinks=%v, devices=%v", dur, doUplinks, doDevices)

	for {
		select {

		// Track the counters here, to convert from raw counters to differences
		case _ = <-poll.C:
			if doDevices {
				if res, err := c.getDeviceClients(dur); err != nil {
					c.log.Infof("Meraki cannot get Device Client Info: %v", err)
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

		case <-ctx.Done():
			c.log.Infof("Meraki Poll Done")
			return
		}
	}
}

type client struct {
	Usage            map[string]float64 `json:"usage"`
	ID               string             `json:"id"`
	Description      string             `json:"description"`
	Mac              string             `json:"mac"`
	IP               string             `json:"ip"`
	User             string             `json:"user"`
	Vlan             int                `json:"vlan"`
	NamedVlan        string             `json:"namedVlan"`
	IPv6             string             `json:"ip6"`
	Manufacturer     string             `json:"manufacturer"`
	DeviceType       string             `json:"deviceTypePrediction"`
	RecentDeviceName string             `json:"recentDeviceName"`
	Status           string             `json:"status"`
	MdnsName         string             `json:"mdnsName"`
	DhcpHostname     string             `json:"dhcpHostname"`
	network          string
	device           networkDevice
}

func (c *MerakiClient) getNetworkClients() ([]*kt.JCHF, error) {
	clientSet := []client{}
	for _, org := range c.orgs {
		for _, network := range org.networks {
			params := networks.NewGetNetworkClientsParams()
			params.SetNetworkID(network.ID)

			prod, err := c.client.Networks.GetNetworkClients(params, c.auth)
			if err != nil {
				return nil, err
			}

			b, err := json.Marshal(prod.GetPayload())
			if err != nil {
				return nil, err
			}

			var clients []client
			err = json.Unmarshal(b, &clients)
			if err != nil {
				return nil, err
			}
			for _, client := range clients {
				client.network = network.Name
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

func (c *MerakiClient) getDeviceClients(dur time.Duration) ([]*kt.JCHF, error) {
	networkDevs, err := c.getNetworkDevices()
	if err != nil {
		return nil, err
	}

	c.log.Infof("Got devices for %d networks", len(networkDevs))
	clientSet := []client{}
	durs := float32(3600) // Look back 1 hour in seconds, get all devices using APs in this range.
	for network, deviceSet := range networkDevs {
		c.log.Infof("Looking at %d devices for network %s", len(deviceSet), network)
		for _, device := range deviceSet {
			params := devices.NewGetDeviceClientsParams()
			params.SetSerial(device.Serial)
			params.SetTimespan(&durs)

			prod, err := c.client.Devices.GetDeviceClients(params, c.auth)
			if err != nil {
				return nil, err
			}

			b, err := json.Marshal(prod.GetPayload())
			if err != nil {
				return nil, err
			}

			var clients []client
			err = json.Unmarshal(b, &clients)
			if err != nil {
				return nil, err
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

func (c *MerakiClient) parseClients(cs []client) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)
	for _, client := range cs {
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
			"dhcp_hostname":      client.DhcpHostname,
			"mdsn_name":          client.MdnsName,
		}
		dst.CustomInt = map[string]int32{
			"vlan": int32(client.Vlan),
		}
		dst.CustomBigInt = map[string]int64{}
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		//dst.Provider = c.conf.Provider // @TODO, pick a provider for this one.

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

		dst.CustomBigInt["Sent"] = int64(client.Usage["sent"] * 1000) // Unit is kilobytes, convert to bytes
		dst.CustomMetrics["Sent"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

		dst.CustomBigInt["Recv"] = int64(client.Usage["recv"] * 1000) // Same, convert to bytes.
		dst.CustomMetrics["Recv"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

		c.conf.SetUserTags(dst.CustomStr)
		res = append(res, dst)
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
	Usage          []uplinkInterfaceUsage
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

func (du *deviceUplink) SetUsage(uplinkHistories []uplinkUsage) {
	for _, ts := range uplinkHistories {
		for _, u := range ts.ByInterface {
			found := false
			for i, ul := range du.Uplinks {
				if ul.Interface == u.Interface {
					if du.Uplinks[i].Usage == nil {
						du.Uplinks[i].Usage = make([]uplinkInterfaceUsage, 0)
					}
					du.Uplinks[i].Usage = append(du.Uplinks[i].Usage, u)
					found = true
				}
			}
			if !found {
				du.Uplinks = append(du.Uplinks, uplink{
					Interface: u.Interface,
					Usage: []uplinkInterfaceUsage{
						u,
					},
				})
			}
		}
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

	uplinkSet := map[string]*deviceUplink{}
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
				// Skip this uplink.
				continue
				//c.log.Errorf("Missing Network for Uplink %s -- %s", uplink.NetworkID, uplink.Serial)
			}
			if _, ok := uplinkSet[uplink.Serial]; ok {
				c.log.Errorf("Duplicate Uplink %s", uplink.Serial)
			} else {
				uplinkSet[uplink.Serial] = &uplink
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

func (c *MerakiClient) getUplinkLatencyLoss(dur time.Duration, uplinkMap map[string]*deviceUplink) error {

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
				c.log.Errorf("Missing Uplink %s In LatencyLoss", uplink.Serial)
			}
		}
	}

	return nil
}

type uplinkInterfaceUsage struct {
	Interface string  `json:"interface"`
	Sent      float64 `json:"sent"`
	Received  float64 `json:"received"`
}

type uplinkUsage struct {
	StartTime   time.Time              `json:"startTime"`
	EndTime     time.Time              `json:"endTime"`
	ByInterface []uplinkInterfaceUsage `json:"byInterface"`
}

func (c *MerakiClient) getUplinkUsage(dur time.Duration, uplinkMap map[string]*deviceUplink) error {

	for _, org := range c.orgs {
		for _, network := range org.networks {
			params := appliance.NewGetNetworkApplianceUplinksUsageHistoryParams()
			params.SetNetworkID(network.ID)

			prod, err := c.client.Appliance.GetNetworkApplianceUplinksUsageHistory(params, c.auth)
			if err != nil {
				c.log.Warnf("Cannot get Uplink Usage: %s %v", network.Name, err)
				continue
			}

			b, err := json.Marshal(prod.GetPayload())
			if err != nil {
				return err
			}

			var uplinkHistories []uplinkUsage
			err = json.Unmarshal(b, &uplinkHistories)
			if err != nil {
				return err
			}

			if len(uplinkHistories) > 0 {
				for _, du := range uplinkMap {
					if du.network.ID == network.ID {
						du.SetUsage(uplinkHistories)
					}
				}

			}
		}
	}

	return nil
}

func (c *MerakiClient) parseUplinks(uplinkMap map[string]*deviceUplink) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)
	for _, device := range uplinkMap {
		for _, uplink := range device.Uplinks {
			dst := kt.NewJCHF()
			dst.SrcAddr = uplink.IP
			dst.DeviceName = device.network.Name + uplink.Interface

			dst.CustomStr = map[string]string{
				"network":         device.network.Name,
				"status":          uplink.Status,
				"connection_type": uplink.ConnectionType,
				"interface":       uplink.Interface,
				"provider":        uplink.Provider,
				"signal_type":     uplink.SignalType,
				"signal_rsrp":     uplink.SignalStat.Rsrp,
				"signal_rsrq":     uplink.SignalStat.Rsrq,
			}
			dst.CustomInt = map[string]int32{}
			dst.CustomBigInt = map[string]int64{}
			dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
			//dst.Provider = c.conf.Provider // @TODO, pick a provider for this one.

			dst.Timestamp = time.Now().Unix()
			dst.CustomMetrics = map[string]kt.MetricInfo{}

			sent, recv := getUsage(uplink)
			dst.CustomBigInt["Sent"] = sent
			dst.CustomMetrics["Sent"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.uplinks", Type: "meraki.uplinks"}

			dst.CustomBigInt["Recv"] = recv
			dst.CustomMetrics["Recv"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.uplinks", Type: "meraki.uplinks"}

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

func getUsage(u uplink) (int64, int64) {
	sent := int64(0)
	rcv := int64(0)
	for _, t := range u.Usage {
		sent += int64(t.Sent * 1000)    // Kilobytes -> Bytes
		rcv += int64(t.Received * 1000) // Same.
	}

	return sent, rcv
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
