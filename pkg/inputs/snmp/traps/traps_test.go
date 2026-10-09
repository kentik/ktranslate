package traps

import (
	"context"
	"net"
	"testing"

	"github.com/gosnmp/gosnmp"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/kt"
)

func pdu(oid string, typ gosnmp.Asn1BER, value interface{}) gosnmp.SnmpPDU {
	return gosnmp.SnmpPDU{Name: oid, Type: typ, Value: value}
}

func TestResolveSenderNoConfig(t *testing.T) {
	assert := assert.New(t)

	packet := &gosnmp.SnmpPacket{Version: gosnmp.Version2c}
	cfg := &kt.SnmpTrapConfig{}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("", sender)
	assert.Equal(senderSourceUDP, source)
}

func TestResolveSenderVarOid(t *testing.T) {
	assert := assert.New(t)

	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.3.6.1.4.1.9999.1.1.0", gosnmp.OctetString, []byte("10.1.1.10")),
		},
	}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.3.6.1.4.1.9999.1.1.0"}}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal("varoid:.1.3.6.1.4.1.9999.1.1.0", source)
}

func TestResolveSenderVarOidPriority(t *testing.T) {
	assert := assert.New(t)

	// Second configured OID is the one present; first should still be tried first and found absent.
	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.2.3.4.0", gosnmp.OctetString, []byte("10.1.1.11")),
			pdu(".1.2.3.5.0", gosnmp.OctetString, []byte("10.1.1.10")),
		},
	}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.5.0", ".1.2.3.4.0"}}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal("varoid:.1.2.3.5.0", source)
}

func TestResolveSenderSkipsEmptyVarOidValue(t *testing.T) {
	assert := assert.New(t)

	// The first configured OID is present but decodes to "", so it names nothing and
	// resolution has to keep going rather than settle on it.
	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.2.3.4.0", gosnmp.OctetString, []byte{}),
			pdu(".1.2.3.5.0", gosnmp.OctetString, []byte("10.1.1.10")),
		},
	}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.4.0", ".1.2.3.5.0"}}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal("varoid:.1.2.3.5.0", source)
}

func TestResolveSenderSkipsNulOnlyVarOidValueForStdSender(t *testing.T) {
	assert := assert.New(t)

	// A NUL-only value trims to "" in ReadOctetString; the standard sender still applies.
	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.2.3.4.0", gosnmp.OctetString, []byte{0x00, 0x00}),
			pdu(snmpTrapAddressOID, gosnmp.IPAddress, "10.1.1.10"),
		},
	}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.4.0"}, UseStdSender: true}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal(senderSourceStdAddr, source)
}

func TestResolveSenderV1PrefersAgentAddressOverStdAddr(t *testing.T) {
	assert := assert.New(t)

	// snmpTrapAddress is what RFC 3584 has a proxy add when translating a v1 trap to
	// v2c/v3. A v1 trap carrying both must still trust its own AgentAddress field.
	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version1,
		Variables: []gosnmp.SnmpPDU{
			pdu(snmpTrapAddressOID, gosnmp.IPAddress, "10.1.1.99"),
		},
		SnmpTrap: gosnmp.SnmpTrap{AgentAddress: "10.1.1.10"},
	}
	cfg := &kt.SnmpTrapConfig{UseStdSender: true}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal(senderSourceAgent, source)
}

func TestResolveSenderStdSenderSnmpTrapAddress(t *testing.T) {
	assert := assert.New(t)

	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(snmpTrapAddressOID, gosnmp.IPAddress, "10.1.1.10"),
		},
	}
	cfg := &kt.SnmpTrapConfig{UseStdSender: true}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal(senderSourceStdAddr, source)
}

func TestResolveSenderStdSenderRequiresFlag(t *testing.T) {
	assert := assert.New(t)

	// snmpTrapAddress is present, but UseStdSender is off: must not be adopted.
	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(snmpTrapAddressOID, gosnmp.IPAddress, "10.1.1.10"),
		},
	}
	cfg := &kt.SnmpTrapConfig{}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("", sender)
	assert.Equal(senderSourceUDP, source)
}

func TestResolveSenderAgentAddressV1Only(t *testing.T) {
	assert := assert.New(t)
	cfg := &kt.SnmpTrapConfig{UseStdSender: true}

	v1Packet := &gosnmp.SnmpPacket{Version: gosnmp.Version1, SnmpTrap: gosnmp.SnmpTrap{AgentAddress: "10.1.1.10"}}
	sender, source := resolveSender(v1Packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal(senderSourceAgent, source)

	// AgentAddress is a v1-only PDU field; must not be consulted for v2c/v3.
	v2cPacket := &gosnmp.SnmpPacket{Version: gosnmp.Version2c, SnmpTrap: gosnmp.SnmpTrap{AgentAddress: "10.1.1.10"}}
	sender, source = resolveSender(v2cPacket, "10.0.0.5", cfg, nil)
	assert.Equal("", sender)
	assert.Equal(senderSourceUDP, source)
}

func TestResolveSenderVarOidTakesPriorityOverStdSender(t *testing.T) {
	assert := assert.New(t)

	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(snmpTrapAddressOID, gosnmp.IPAddress, "10.1.1.99"),
			pdu(".1.2.3.4.0", gosnmp.OctetString, []byte("10.1.1.10")),
		},
	}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.4.0"}, UseStdSender: true}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal("varoid:.1.2.3.4.0", source)
}

func TestResolveSenderTrustedRelaysBlocksUntrustedSource(t *testing.T) {
	assert := assert.New(t)

	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.2.3.4.0", gosnmp.OctetString, []byte("10.1.1.10")),
		},
	}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.4.0"}}
	trustedRelays := map[string]bool{"10.0.0.5": true}

	// Untrusted relay: payload is ignored entirely.
	sender, source := resolveSender(packet, "10.0.0.99", cfg, trustedRelays)
	assert.Equal("", sender)
	assert.Equal(senderSourceUDP, source)

	// Trusted relay: payload is honored.
	sender, source = resolveSender(packet, "10.0.0.5", cfg, trustedRelays)
	assert.Equal("10.1.1.10", sender)
	assert.Equal("varoid:.1.2.3.4.0", source)
}

func TestResolveSenderInvalidValueFallsThrough(t *testing.T) {
	assert := assert.New(t)

	// A varbind of a type we don't treat as an address/hostname carrier (e.g. a counter)
	// at the configured OID must not be adopted; resolution falls through to std sender.
	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version1,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.2.3.4.0", gosnmp.Counter32, uint(42)),
		},
		SnmpTrap: gosnmp.SnmpTrap{AgentAddress: "10.1.1.10"},
	}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.4.0"}, UseStdSender: true}

	sender, source := resolveSender(packet, "10.0.0.5", cfg, nil)
	assert.Equal("10.1.1.10", sender)
	assert.Equal(senderSourceAgent, source)
}

func TestLookupDeviceByIPThenName(t *testing.T) {
	assert := assert.New(t)

	byIP := &kt.SnmpDeviceConfig{DeviceName: "router-a", DeviceIP: "10.1.1.10"}
	byName := &kt.SnmpDeviceConfig{DeviceName: "router-b.example.com", DeviceIP: "10.1.1.11"}

	st := &SnmpTrap{
		deviceMap:     map[string]*kt.SnmpDeviceConfig{"10.1.1.10": byIP},
		deviceMapName: map[string]*kt.SnmpDeviceConfig{"router-b.example.com": byName},
	}

	assert.Equal(byIP, st.lookupDevice("10.1.1.10"))
	assert.Equal(byName, st.lookupDevice("router-b.example.com"))
	assert.Nil(st.lookupDevice("10.1.1.12"))
	assert.Nil(st.lookupDevice(""))
}

// newTestTrapListener builds a listener wired for handle() without binding a socket.
// A zero-value MibDB resolves nothing, which is what an unprofiled trap sees anyway.
func newTestTrapListener(t *testing.T, cfg *kt.SnmpTrapConfig, devices ...*kt.SnmpDeviceConfig) (*SnmpTrap, chan []*kt.JCHF) {
	jchfChan := make(chan []*kt.JCHF, 1)
	st := &SnmpTrap{
		log:           lt.NewTestContextL(logger.NilContext, t),
		jchfChan:      jchfChan,
		metrics:       kt.NewSnmpMetricSet(go_metrics.NewRegistry()),
		conf:          &kt.SnmpConfig{Trap: cfg},
		mibdb:         &mibs.MibDB{},
		ctx:           context.Background(),
		deviceMap:     map[string]*kt.SnmpDeviceConfig{},
		deviceMapName: map[string]*kt.SnmpDeviceConfig{},
		baseTags:      map[string]string{},
	}
	for _, dev := range devices {
		st.deviceMap[dev.DeviceIP] = dev
		if dev.DeviceName != "" {
			st.deviceMapName[dev.DeviceName] = dev
		}
	}
	if len(cfg.TrustedRelays) > 0 {
		st.trustedRelays = map[string]bool{}
		for _, relay := range cfg.TrustedRelays {
			st.trustedRelays[relay] = true
		}
	}
	return st, jchfChan
}

func TestHandleSenderMatchedByNameReportsDeviceIP(t *testing.T) {
	assert := assert.New(t)

	dev := &kt.SnmpDeviceConfig{DeviceName: "router-b.example.com", DeviceIP: "10.1.1.11"}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.4.0"}}
	st, jchfChan := newTestTrapListener(t, cfg, dev)

	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.2.3.4.0", gosnmp.OctetString, []byte("router-b.example.com")),
		},
	}
	st.handle(packet, &net.UDPAddr{IP: net.ParseIP("10.0.0.5")})

	out := <-jchfChan
	assert.Len(out, 1)
	dst := out[0]
	// SrcAddr has to stay parseable as an IP, so the hostname only appears in SenderAddr.
	assert.Equal("10.1.1.11", dst.SrcAddr)
	assert.NotNil(net.ParseIP(dst.SrcAddr))
	assert.Equal("router-b.example.com", dst.CustomStr["SenderAddr"])
	assert.Equal("10.0.0.5", dst.CustomStr["RelayAddr"])
	assert.Equal("router-b.example.com", dst.DeviceName)
}

func TestHandleSenderMatchedByIPKeepsSenderAddr(t *testing.T) {
	assert := assert.New(t)

	dev := &kt.SnmpDeviceConfig{DeviceName: "router-a", DeviceIP: "10.1.1.10"}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.4.0"}}
	st, jchfChan := newTestTrapListener(t, cfg, dev)

	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.2.3.4.0", gosnmp.OctetString, []byte("10.1.1.10")),
		},
	}
	st.handle(packet, &net.UDPAddr{IP: net.ParseIP("10.0.0.5")})

	out := <-jchfChan
	dst := out[0]
	assert.Equal("10.1.1.10", dst.SrcAddr)
	assert.Equal("10.1.1.10", dst.CustomStr["SenderAddr"])
	assert.Equal("varoid:.1.2.3.4.0", dst.CustomStr["SenderSource"])
	assert.Equal("10.0.0.5", dst.CustomStr["RelayAddr"])
	assert.Equal("router-a", dst.DeviceName)
}

func TestHandleUnresolvedSenderKeepsUdpSource(t *testing.T) {
	assert := assert.New(t)

	dev := &kt.SnmpDeviceConfig{DeviceName: "router-a", DeviceIP: "10.1.1.10"}
	cfg := &kt.SnmpTrapConfig{SenderVarOids: []string{".1.2.3.4.0"}}
	st, jchfChan := newTestTrapListener(t, cfg, dev)

	// The varbind names a device that isn't configured, so none of the sender fields
	// are reported and the UDP source stands.
	packet := &gosnmp.SnmpPacket{
		Version: gosnmp.Version2c,
		Variables: []gosnmp.SnmpPDU{
			pdu(".1.2.3.4.0", gosnmp.OctetString, []byte("10.9.9.9")),
		},
	}
	st.handle(packet, &net.UDPAddr{IP: net.ParseIP("10.0.0.5")})

	out := <-jchfChan
	dst := out[0]
	assert.Equal("10.0.0.5", dst.SrcAddr)
	assert.NotContains(dst.CustomStr, "SenderAddr")
	assert.NotContains(dst.CustomStr, "RelayAddr")
}

func TestHandleNoSenderConfigIsUnchanged(t *testing.T) {
	assert := assert.New(t)

	dev := &kt.SnmpDeviceConfig{DeviceName: "router-a", DeviceIP: "10.1.1.10"}
	st, jchfChan := newTestTrapListener(t, &kt.SnmpTrapConfig{}, dev)

	// A trap straight from the device, with no sender config at all: the existing
	// UDP-source lookup still identifies it and no sender fields are added.
	packet := &gosnmp.SnmpPacket{Version: gosnmp.Version2c}
	st.handle(packet, &net.UDPAddr{IP: net.ParseIP("10.1.1.10")})

	out := <-jchfChan
	dst := out[0]
	assert.Equal("10.1.1.10", dst.SrcAddr)
	assert.Equal("router-a", dst.DeviceName)
	assert.NotContains(dst.CustomStr, "SenderAddr")
	assert.NotContains(dst.CustomStr, "SenderSource")
	assert.NotContains(dst.CustomStr, "RelayAddr")
}
