package arista

import (
	"context"
	"time"

	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type EAPIClient struct {
	log      logger.ContextL
	jchfChan chan []*kt.JCHF
	conf     *kt.SnmpDeviceConfig
	metrics  *kt.SnmpDeviceMetric
	client   *goeapi.Node
}

func NewEAPIClient(jchfChan chan []*kt.JCHF, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, log logger.ContextL) (*EAPIClient, error) {
	c := EAPIClient{
		log:      log,
		jchfChan: jchfChan,
		conf:     conf,
		metrics:  metrics,
	}

	node, err := goeapi.Connect(conf.Ext.EAPIConfig.Transport, conf.DeviceIP, conf.Ext.EAPIConfig.Username, conf.Ext.EAPIConfig.Password, conf.Ext.EAPIConfig.Port)
	if err != nil {
		return nil, err
	}
	c.log.Infof("Testing %s to %s", c.GetName(), conf.DeviceName)
	sys := module.System(node)
	sysInfo := sys.Get() // Only here do we make the RPC call. @TODO, how long can this block for?
	c.log.Infof("%s connected to %s (%s)", c.GetName(), conf.DeviceName, sysInfo.HostName())
	c.client = node

	return &c, nil
}

func (c *EAPIClient) GetName() string {
	return "Arista eAPI"
}

func (c *EAPIClient) Run(ctx context.Context, dur time.Duration) {

	c.getBGP()
	c.getMLAG()

	poll := time.NewTicker(dur)
	defer poll.Stop()

	for {
		select {

		// Track the counters here, to convert from raw counters to differences
		case _ = <-poll.C:
			if res, err := c.getBGP(); err != nil {
				c.log.Infof("eAPI cannot get BGP Info: %v", err)
			} else if len(res) > 0 {
				c.jchfChan <- res
			}

			if res, err := c.getMLAG(); err != nil {
				c.log.Infof("eAPI cannot get MLAG Info: %v", err)
			} else if len(res) > 0 {
				c.jchfChan <- res
			}

		case <-ctx.Done():
			c.log.Infof("eAPI Poll Done")
			return
		}
	}
}

var (
	stateMap = map[string]int64{
		"established": 0,
		"openconfirm": 1,
		"opensent":    2,
		"active":      3,
		"connect":     4,
		"idle":        5,
	}

	boolMap = map[bool]string{
		true:  "true",
		false: "false",
	}

	shutdownMap = map[string]int64{
		"true":  1,
		"false": 2,
	}

	mlag_strings = map[string]int64{
		"consistent": 1,
		"up":         1,
		"connected":  1,
		"active":     1,
		"primary":    1,
		"secondary":  2,
		"disabled":   0,
		"unknown":    3,
	}

	mlag_bool = map[bool]int64{
		true:  1,
		false: 0,
	}
)

type Peer struct {
	Description         string  `json:"description"`
	MsgSent             int64   `json:"msgSent"`
	InMsgQueue          int64   `json:"inMsgQueue"`
	PrefixReceived      int64   `json:"prefixReceived"`
	UpDownTime          float64 `json:"upDownTime"`
	Version             int     `json:"version"`
	PrefixAccepted      int64   `json:"prefixAccepted"`
	MsgReceived         int64   `json:"msgReceived"`
	PeerState           string  `json:"peerState"`
	OutMsgQueue         int64   `json:"outMsgQueue"`
	UnderMaintenance    bool    `json:"underMaintenance"`
	ASN                 string  `json:"asn"`
	PeerStateIdleReason string  `json:"peerStateIdleReason"`
}

type VRF struct {
	RouterID string          `json:"routerId"`
	Peers    map[string]Peer `json:"peers"`
	VRF      string          `json:"vrf"`
	ASN      string          `json:"asn"`
}

type ShowBGP struct {
	VRFs map[string]VRF `json:"vrfs"`
}

func (s *ShowBGP) GetCmd() string {
	return "show ip bgp summary vrf all"
}

func (c *EAPIClient) getBGP() ([]*kt.JCHF, error) {
	sv := &ShowBGP{}
	handle, _ := c.client.GetHandle("json")
	handle.AddCommand(sv)
	if err := handle.Call(); err != nil {
		return nil, err
	}

	return c.parseBGP(sv)
}

func (c *EAPIClient) parseBGP(sv *ShowBGP) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)

	for v, vrf := range sv.VRFs {
		for peer, state := range vrf.Peers {
			dst := kt.NewJCHF()
			dst.CustomStr = map[string]string{
				"router_id":              vrf.RouterID,
				"vrf":                    vrf.VRF,
				"peer":                   peer,
				"peer_state":             state.PeerState,
				"under_maintenance":      boolMap[state.UnderMaintenance],
				"peer_state_idle_reason": state.PeerStateIdleReason,
				"asn":                    vrf.ASN,
				"peer_asn":               state.ASN,
			}
			dst.CustomInt = map[string]int32{}
			dst.CustomBigInt = map[string]int64{}
			dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
			dst.Provider = c.conf.Provider
			dst.DeviceName = c.conf.DeviceName
			dst.SrcAddr = c.conf.DeviceIP
			dst.Timestamp = time.Now().Unix()
			dst.CustomMetrics = map[string]kt.MetricInfo{}

			dst.CustomBigInt["MsgSent"] = int64(state.MsgSent)
			dst.CustomMetrics["MsgSent"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bpg", Type: "eapi.bgp"}

			dst.CustomBigInt["InMsgQueue"] = int64(state.InMsgQueue)
			dst.CustomMetrics["InMsgQueue"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bgp", Type: "eapi.bgp"}

			dst.CustomBigInt["PrefixReceived"] = int64(state.PrefixReceived)
			dst.CustomMetrics["PrefixReceived"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bgp", Type: "eapi.bgp"}

			dst.CustomBigInt["UpDownTime"] = int64(state.UpDownTime)
			dst.CustomMetrics["UpDownTime"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bgp", Type: "eapi.bgp"}

			dst.CustomBigInt["Version"] = int64(state.Version)
			dst.CustomMetrics["Version"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bgp", Type: "eapi.bgp"}

			dst.CustomBigInt["MsgReceived"] = int64(state.MsgReceived)
			dst.CustomMetrics["MsgReceived"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bgp", Type: "eapi.bgp"}

			dst.CustomBigInt["PrefixAccepted"] = int64(state.PrefixAccepted)
			dst.CustomMetrics["PrefixAccepted"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bgp", Type: "eapi.bgp"}

			dst.CustomBigInt["PeerState"] = stateMap[state.PeerState]
			dst.CustomMetrics["PeerState"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bgp", Type: "eapi.bgp"}

			dst.CustomBigInt["OutMsgQueue"] = int64(state.OutMsgQueue)
			dst.CustomMetrics["OutMsgQueue"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.bgp", Type: "eapi.bgp"}

			c.conf.SetUserTags(dst.CustomStr)
			res = append(res, dst)
		}
	}

	return res, nil
}

type MLAGPorts struct {
	Disabled      int64 `json:"Disabled"`
	ActivePartial int64 `json:"Active-partial"`
	Inactive      int64 `json:"Inactive"`
	Configured    int64 `json:"Configured"`
	ActiveFull    int64 `json:"ActiveFull"`
}

type MLAGDetail struct {
	FailoverCauseList            []string `json:"failoverCauseList"`
	UdpHeartbeatsSent            int64    `json:"udpHeartbeatsSent"`
	LacpStandby                  bool     `json:"lacpStandby"`
	MlagState                    string   `json:"mlagState"`
	HeartbeatInterval            int      `json:"heartbeatInterval"`
	EffectiveHeartbeatInterval   int      `json:"effectiveHeartbeatInterval"`
	HeartbeatTimeout             int      `json:"heartbeatTimeout"`
	StateChanges                 int      `json:"stateChanges"`
	FastMacRedirectionEnabled    bool     `json:"fastMacRedirectionEnabled"`
	PeerPortsErrdisabled         bool     `json:"peerPortsErrdisabled"`
	MlagHwReady                  bool     `json:"mlagHwReady"`
	UdpHeartbeatAlive            bool     `json:"udpHeartbeatAlive"`
	FailoverInitiated            bool     `json:"failoverInitiated"`
	PeerMlagState                string   `json:"peerMlagState"`
	SecondaryFromFailover        bool     `json:"secondaryFromFailover"`
	PrimaryPriority              int      `json:"primaryPriority"`
	Failover                     bool     `json:"failover"`
	Enabled                      bool     `json:"enabled"`
	PeerMacRoutingSupported      bool     `json:"peerMacRoutingSupported"`
	PeerPrimaryPriority          int      `json:"peerPrimaryPriority"`
	udpHeartbeatsReceived        int64    `json:"udpHeartbeatsReceived"`
	PeerMacAddress               string   `json:"peerMacAddress"`
	MountChanges                 int      `json:"mountChanges"`
	HeartbeatTimeoutsSinceReboot int64    `json:"heartbeatTimeoutsSinceReboot"`
}

type ShowMlag struct {
	ConfigSanity                string     `json:"configSanity"`
	DomainId                    string     `json:"domainId"`
	LocalIntfStatus             string     `json:"localIntfStatus"`
	LocalInterface              string     `json:"localInterface"`
	State                       string     `json:"state"`
	ReloadDelay                 int        `json:"reloadDelay"`
	PeerLink                    string     `json:"peerLink"`
	NegStatus                   string     `json:"negStatus"`
	PeerAddress                 string     `json:"peerAddress"`
	PeerLinkStatus              string     `json:"peerLinkStatus"`
	SystemId                    string     `json:"systemId"`
	DualPrimaryDetectionState   string     `json:"dualPrimaryDetectionState"`
	ReloadDelayNonMlag          int        `json:"reloadDelayNonMlag"`
	MlagPorts                   MLAGPorts  `json:"mlagPorts"`
	PortsErrdisabled            bool       `json:"portsErrdisabled"`
	DualPrimaryPortsErrdisabled bool       `json:"dualPrimaryPortsErrdisabled"`
	Detail                      MLAGDetail `json:"detail"`
}

func (s *ShowMlag) GetCmd() string {
	return "show mlag detail"
}

func (c *EAPIClient) getMLAG() ([]*kt.JCHF, error) {
	sv := &ShowMlag{}
	handle, _ := c.client.GetHandle("json")
	handle.AddCommand(sv)
	if err := handle.Call(); err != nil {
		return nil, err
	}

	return c.parseMLAG(sv)
}

func (c *EAPIClient) parseMLAG(sv *ShowMlag) ([]*kt.JCHF, error) {
	dst := kt.NewJCHF()
	dst.CustomStr = map[string]string{
		"config_sanity":     sv.ConfigSanity,
		"local_intf_status": sv.LocalIntfStatus,
		"neg_status":        sv.NegStatus,
		"peer_link_status":  sv.PeerLinkStatus,
		"state":             sv.State,
		"domain_id":         sv.DomainId,
		"local_interface":   sv.LocalInterface,
		"peer_address":      sv.PeerAddress,
		"peer_link":         sv.PeerLink,
	}
	dst.CustomInt = map[string]int32{}
	dst.CustomBigInt = map[string]int64{}
	dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
	dst.Provider = c.conf.Provider
	dst.DeviceName = c.conf.DeviceName
	dst.SrcAddr = c.conf.DeviceIP
	dst.Timestamp = time.Now().Unix()
	dst.CustomMetrics = map[string]kt.MetricInfo{}

	dst.CustomBigInt["ConfigSanity"] = mlag_strings[sv.ConfigSanity]
	dst.CustomMetrics["ConfigSanity"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["LocalIntfStatus"] = mlag_strings[sv.LocalIntfStatus]
	dst.CustomMetrics["LocalIntfStatus"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["NegStatus"] = mlag_strings[sv.NegStatus]
	dst.CustomMetrics["NegStatus"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["PeerLinkStatus"] = mlag_strings[sv.PeerLinkStatus]
	dst.CustomMetrics["PeerLinkStatus"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["PortsErrdisabled"] = mlag_bool[sv.PortsErrdisabled]
	dst.CustomMetrics["PortsErrdisabled"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["State"] = mlag_strings[sv.State]
	dst.CustomMetrics["State"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["PortsDisabled"] = sv.MlagPorts.Disabled
	dst.CustomMetrics["PortsDisabled"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["PortsActivePartial"] = sv.MlagPorts.ActivePartial
	dst.CustomMetrics["PortsActivePartial"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["PortsInactive"] = sv.MlagPorts.Inactive
	dst.CustomMetrics["PortsInactive"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["PortsConfigured"] = sv.MlagPorts.Configured
	dst.CustomMetrics["PortsConfigured"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	dst.CustomBigInt["PortsActiveFull"] = sv.MlagPorts.ActiveFull
	dst.CustomMetrics["PortsActiveFull"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi.mlag", Type: "eapi.mlag"}

	return []*kt.JCHF{dst}, nil
}
