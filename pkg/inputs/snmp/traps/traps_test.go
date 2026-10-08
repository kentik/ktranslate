package traps

import (
	"testing"

	"github.com/gosnmp/gosnmp"
	"github.com/stretchr/testify/assert"

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
