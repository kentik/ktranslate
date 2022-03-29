package arista

import (
	"encoding/json"
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestParseBGP(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	testIn := []byte(`
    {
      "vrfs": {
        "default": {
          "routerId": "10.192.255.3",
          "peers": {
            "10.99.2.107": {
              "description": "AWS_PRIVATE",
              "msgSent": 1411756,
              "inMsgQueue": 1,
              "prefixReceived": 10,
              "upDownTime": 1647491584,
              "version": 4,
              "prefixAccepted": 10,
              "msgReceived": 1222121,
              "peerState": "Established",
              "outMsgQueue": 0,
              "underMaintenance": false,
              "asn": "1"
            },
            "10.192.250.19": {
              "description": "FW_INSIDE",
              "msgSent": 4135220,
              "inMsgQueue": 2,
              "prefixReceived": 16,
              "upDownTime": 1629291904,
              "version": 4,
              "prefixAccepted": 9,
              "msgReceived": 1377152,
              "peerState": "Established",
              "outMsgQueue": 0,
              "underMaintenance": false,
              "asn": "2"
            },
            "10.192.250.4": {
              "description": "LEAF_UNDERLAY",
              "msgSent": 902580,
              "inMsgQueue": 3,
              "prefixReceived": 15,
              "upDownTime": 1612622720,
              "version": 4,
              "prefixAccepted": 15,
              "msgReceived": 924834,
              "peerState": "Established",
              "outMsgQueue": 0,
              "underMaintenance": false,
              "asn": "3"
            }
          },
          "vrf": "default",
          "asn": "100"
        },
        "internet": {
          "routerId": "2.78.112.2",
          "peers": {
            "1.14.99.165": {
              "description": "DESC",
              "msgSent": 1406038,
              "inMsgQueue": 4,
              "prefixReceived": 872997,
              "upDownTime": 1630422912,
              "version": 4,
              "prefixAccepted": 1,
              "msgReceived": 110745391,
              "peerState": "Established",
              "outMsgQueue": 0,
              "underMaintenance": false,
              "asn": "5"
            },
            "1.78.112.3": {
              "description": "INET_IBGP",
              "msgSent": 702988,
              "inMsgQueue": 5,
              "prefixReceived": 2,
              "upDownTime": 1612622720,
              "version": 4,
              "prefixAccepted": 2,
              "msgReceived": 703019,
              "peerState": "Established",
              "outMsgQueue": 0,
              "underMaintenance": false,
              "asn": "6"
            }
          },
          "vrf": "internet",
          "asn": "100"
        }
      }
    }`)

	sv := ShowBGP{}
	err := json.Unmarshal(testIn, &sv)
	assert.NoError(t, err)

	c := &EAPIClient{
		log:  l,
		conf: &kt.SnmpDeviceConfig{Provider: "provider"},
	}
	res, err := c.parseBGP(&sv)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(res)) // 6 results, because 6 peers.
	for _, r := range res {
		assert.True(t, r.CustomBigInt["InMsgQueue"] > 0, r.CustomBigInt)
		assert.True(t, r.CustomBigInt["OutMsgQueue"] == 0, r.CustomBigInt)
		assert.True(t, r.CustomStr["asn"] == "100", r.CustomStr)
	}

}

func TestParseMLAG(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	testIn := []byte(`
    {
      "localInterface": "",
      "state": "connected",
      "reloadDelay": 300,
      "peerLink": "",
      "dualPrimaryDetectionState": "disabled",
      "reloadDelayNonMlag": 300,
      "mlagPorts": {
        "Disabled": 0,
        "Active-partial": 0,
        "Inactive": 0,
        "Configured": 0,
        "Active-full": 2
      },
      "portsErrdisabled": false,
      "dualPrimaryPortsErrdisabled": false,
      "detail": {
        "failoverCauseList": [
          "Unknown"
        ],
        "udpHeartbeatsSent": 0,
        "lacpStandby": false,
        "mlagState": "disabled",
        "heartbeatInterval": 4000,
        "effectiveHeartbeatInterval": 0,
        "heartbeatTimeout": 0,
        "stateChanges": 0,
        "fastMacRedirectionEnabled": false,
        "peerPortsErrdisabled": false,
        "mlagHwReady": true,
        "udpHeartbeatAlive": false,
        "failoverInitiated": false,
        "peerMlagState": "unknown",
        "secondaryFromFailover": false,
        "primaryPriority": 32767,
        "failover": false,
        "enabled": false,
        "peerMacRoutingSupported": false,
        "peerPrimaryPriority": 0,
        "udpHeartbeatsReceived": 0,
        "peerMacAddress": "00:00:00:00:00:00",
        "mountChanges": 0,
        "heartbeatTimeoutsSinceReboot": 0
      }
    }`)

	sv := ShowMlag{}
	err := json.Unmarshal(testIn, &sv)
	assert.NoError(t, err)

	c := &EAPIClient{
		log:  l,
		conf: &kt.SnmpDeviceConfig{Provider: "provider"},
	}
	res, err := c.parseMLAG(&sv)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res)) // 6 results, because 6 peers.
	assert.Equal(t, 1, int(res[0].CustomBigInt["State"]), res[0].CustomBigInt)
	assert.Equal(t, 2, int(res[0].CustomBigInt["PortsActiveFull"]), res[0].CustomBigInt)
}
