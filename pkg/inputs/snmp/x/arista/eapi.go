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
)

func (c *EAPIClient) getBGP() ([]*kt.JCHF, error) {
	bgp := module.Show(c.client)
	sum, err := bgp.ShowIPBGPSummary()
	if err != nil {
		return nil, err
	}

	res := make([]*kt.JCHF, 0)

	// Testing.
	sum = getfake()

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
			dst.CustomMetrics["MsgSent"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			dst.CustomBigInt["InMsgQueue"] = int64(state.InMsgQueue)
			dst.CustomMetrics["InMsgQueue"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			dst.CustomBigInt["PrefixReceived"] = int64(state.PrefixReceived)
			dst.CustomMetrics["PrefixReceived"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			dst.CustomBigInt["UpDownTime"] = int64(state.UpDownTime)
			dst.CustomMetrics["UpDownTime"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			dst.CustomBigInt["Version"] = int64(state.Version)
			dst.CustomMetrics["Version"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			dst.CustomBigInt["MsgReceived"] = int64(state.MsgReceived)
			dst.CustomMetrics["MsgReceived"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			dst.CustomBigInt["PrefixAccepted"] = int64(state.PrefixAccepted)
			dst.CustomMetrics["PrefixAccepted"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			dst.CustomBigInt["PeerState"] = stateMap[state.PeerState]
			dst.CustomMetrics["PeerState"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			dst.CustomBigInt["OutMsgQueue"] = int64(state.OutMsgQueue)
			dst.CustomMetrics["OutMsgQueue"] = kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: "eapi", Type: "eapi"}

			c.conf.SetUserTags(dst.CustomStr)
			res = append(res, dst)
		}
	}

	return res, nil
}

func (c *EAPIClient) getMLAG() ([]*kt.JCHF, error) {
	return nil, nil
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
