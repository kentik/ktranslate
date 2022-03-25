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

func (c *EAPIClient) getBGP() ([]*kt.JCHF, error) {
	bgp := module.Show(c.client)
	sum, err := bgp.ShowIPBGPSummary()
	if err != nil {
		return nil, err
	}

	// Testing.
	// sum = getfake()

	// Each VFR + peer combo is a unique point to record for a metric.
	res := make([]*kt.JCHF, 0)
	for _, vrf := range sum.VRFs {
		for peer, state := range vrf.Peers {
			dst := kt.NewJCHF()
			dst.CustomStr = map[string]string{
				"router_id":              vrf.RouterID,
				"vrf":                    vrf.VRF,
				"peer":                   peer,
				"peer_state":             state.PeerState,
				"under_maintenance":      boolMap[state.UnderMaintenance],
				"peer_state_idle_reason": state.PeerStateIdleReason,
			}
			dst.CustomInt = map[string]int32{}
			dst.CustomBigInt = map[string]int64{
				"asn":      vrf.ASN,
				"peer_asn": state.ASN,
			}
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
	// Old system, using the mlag module.
	/*
		mlag := module.Mlag(c.client)
		config := mlag.Get()
		if config == nil {
			return nil, fmt.Errorf("Could not get a mlag config")
		}

			// This is just a dumb system right now. @TODO, is there more info to pull out?
			dst := kt.NewJCHF()
			dst.CustomStr = map[string]string{
				"domain_id":       config.DomainID(),
				"local_interface": config.LocalInterface(),
				"peer_address":    config.PeerAddress(),
				"peer_link":       config.PeerLink(),
				"shutdown":        config.Shutdown(),
			}
			dst.CustomInt = map[string]int32{}
			dst.CustomBigInt = map[string]int64{}
			dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
			dst.Provider = c.conf.Provider
			dst.DeviceName = c.conf.DeviceName
			dst.SrcAddr = c.conf.DeviceIP
			dst.Timestamp = time.Now().Unix()
			dst.CustomMetrics = map[string]kt.MetricInfo{}

			dst.CustomBigInt["MLAGShutdown"] = shutdownMap[config.Shutdown()]
			dst.CustomMetrics["MLAGShutdown"] = kt.MetricInfo{Oid: "eapi", Mib: "eapi", Profile: "eapi", Type: "eapi"}
	*/
}

func getfake() module.ShowIPBGPSummary {
	return module.ShowIPBGPSummary{
		VRFs: map[string]module.VRF{
			"one": module.VRF{
				RouterID: "id.one",
				VRF:      "one",
				ASN:      1212,
				Peers: map[string]module.BGPNeighborSummary{
					"p1": module.BGPNeighborSummary{
						MsgSent:    10,
						InMsgQueue: 1,
						PeerState:  "active",
						ASN:        65555,
					},
					"p2": module.BGPNeighborSummary{
						MsgSent:             101,
						InMsgQueue:          2,
						PeerState:           "idle",
						ASN:                 65556,
						PeerStateIdleReason: "lazy",
						UnderMaintenance:    true,
					},
				},
			},
			"two": module.VRF{
				RouterID: "id.two",
				VRF:      "two",
				ASN:      1215,
				Peers: map[string]module.BGPNeighborSummary{
					"p1": module.BGPNeighborSummary{
						MsgSent:    999,
						InMsgQueue: 1111,
						PeerState:  "active",
						ASN:        65555,
					},
				},
			},
		},
	}
}
