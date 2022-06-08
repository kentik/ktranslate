package meraki

import (
	"context"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	apiclient "github.com/ddexterpark/dashboard-api-golang/client"
	"github.com/ddexterpark/dashboard-api-golang/client/networks"
	"github.com/ddexterpark/dashboard-api-golang/client/organizations"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

type MerakiClient struct {
	log      logger.ContextL
	jchfChan chan []*kt.JCHF
	conf     *kt.SnmpDeviceConfig
	metrics  *kt.SnmpDeviceMetric
	client   *apiclient.MerakiAPIGolang
	auth     runtime.ClientAuthInfoWriter
	orgs     []*organizations.GetOrganizationsOKBodyItems0
	networks []networkDesc
}

type networkDesc struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewMerakiClient(jchfChan chan []*kt.JCHF, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, log logger.ContextL) (*MerakiClient, error) {
	c := MerakiClient{
		log:      log,
		jchfChan: jchfChan,
		conf:     conf,
		metrics:  metrics,
		orgs:     []*organizations.GetOrganizationsOKBodyItems0{},
		networks: []networkDesc{},
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

	for _, org := range prod.GetPayload() {
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

		for _, network := range networks {
			c.log.Infof("Adding network %s to list to track", network.Name)
			c.networks = append(c.networks, network)
		}

		if len(networks) > 0 { // Only add this org in to track if it has some networks.
			c.log.Infof("Adding organization %s to list to track", org.Name)
			c.orgs = append(c.orgs, org)
		}
	}

	// If we get this far, we have a list of things to look at.
	c.log.Infof("%s connected to API for %s with %d organization(s) and %d network(s).", c.GetName(), conf.DeviceName, len(c.orgs), len(c.networks))
	c.client = client

	return &c, nil
}

func (c *MerakiClient) GetName() string {
	return "Meraki API"
}

func (c *MerakiClient) Run(ctx context.Context, dur time.Duration) {
	poll := time.NewTicker(dur)
	defer poll.Stop()

	for {
		select {

		// Track the counters here, to convert from raw counters to differences
		case _ = <-poll.C:
			if res, err := c.getClients(); err != nil {
				c.log.Infof("Meraki cannot get Client Info: %v", err)
			} else if len(res) > 0 {
				c.jchfChan <- res
			}

		case <-ctx.Done():
			c.log.Infof("Meraki Poll Done")
			return
		}
	}
}

type client struct {
	Usage            map[string]int `json:"usage"`
	ID               string         `json:"id"`
	Description      string         `json:"description"`
	Mac              string         `json:"mac"`
	IP               string         `json:"ip"`
	User             string         `json:"user"`
	Vlan             string         `json:"vlan"`
	NamedVlan        string         `json:"namedVlan"`
	IPv6             string         `json:"ip6"`
	Manufacturer     string         `json:"manufacturer"`
	DeviceType       string         `json:"deviceTypePrediction"`
	RecentDeviceName string         `json:"recentDeviceName"`
	Status           string         `json:"status"`
}

func (c *MerakiClient) getClients() ([]*kt.JCHF, error) {
	clientSet := map[string][]client{}
	for _, network := range c.networks {
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
		if len(clients) > 0 {
			clientSet[network.Name] = clients
		}
	}

	return c.parseClients(clientSet)
}

func (c *MerakiClient) parseClients(cs map[string][]client) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)
	for network, clients := range cs {
		for _, client := range clients {
			dst := kt.NewJCHF()
			if client.IPv6 != "" {
				dst.SrcAddr = client.IPv6
			} else {
				dst.SrcAddr = client.IP
			}
			dst.CustomStr = map[string]string{
				"network":            network,
				"client_id":          client.ID,
				"description":        client.Description,
				"status":             client.Status,
				"vlan_name":          client.NamedVlan,
				"mac_addr":           client.Mac,
				"ip":                 dst.SrcAddr,
				"user":               client.User,
				"vlan":               client.Vlan,
				"manufacturer":       client.Manufacturer,
				"device_type":        client.DeviceType,
				"recent_device_name": client.RecentDeviceName,
			}
			dst.CustomInt = map[string]int32{}
			dst.CustomBigInt = map[string]int64{}
			dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
			dst.Provider = c.conf.Provider
			dst.DeviceName = c.conf.DeviceName
			dst.SrcAddr = c.conf.DeviceIP
			dst.Timestamp = time.Now().Unix()
			dst.CustomMetrics = map[string]kt.MetricInfo{}

			dst.CustomBigInt["Sent"] = int64(client.Usage["sent"])
			dst.CustomMetrics["Sent"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

			dst.CustomBigInt["Recv"] = int64(client.Usage["recv"])
			dst.CustomMetrics["Recv"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

			c.conf.SetUserTags(dst.CustomStr)
			res = append(res, dst)
		}
	}

	return res, nil
}
