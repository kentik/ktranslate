package kt

import ()

var (
	InputTesting = []*JCHF{
		&JCHF{CompanyId: 10, SrcAddr: "10.2.2.1", Protocol: "TCP", DstAddr: "2001:db8::68", Timestamp: 1, L4DstPort: 80, SrcAs: 1111, SrcGeo: "US", OutputPort: IfaceID(20), EventType: KENTIK_EVENT_TYPE, CustomStr: map[string]string{"foo": "bar"}, CustomInt: map[string]int32{"fooI": 1}, CustomBigInt: map[string]int64{"fooII": 12}, InBytes: 12121, InPkts: 12, OutBytes: 13, OutPkts: 1, SrcEthMac: "90:61:ae:fb:c2:19", avroSet: map[string]interface{}{}},
		&JCHF{CompanyId: 10, SrcAddr: "3.2.2.2", InBytes: 1, OutBytes: 12, InPkts: 12, OutPkts: 1, Protocol: "UDP", SrcAs: 222223, SrcGeo: "CA", DstAddr: "2001:db8::69", SrcEthMac: "90:61:ae:fb:c2:20", Timestamp: 2, CustomStr: map[string]string{"tar": "far"}, EventType: KENTIK_EVENT_TYPE, avroSet: map[string]interface{}{}},
	}

	InputTestingSynth = []*JCHF{
		&JCHF{CompanyId: 10, SrcAddr: "10.2.2.1", Protocol: "TCP", DstAddr: "2001:db8::68", Timestamp: 1, L4DstPort: 80, OutputPort: IfaceID(20), EventType: KENTIK_EVENT_SYNTH, CustomStr: map[string]string{"foo": "bar"}, CustomInt: map[string]int32{"foo": 1}, CustomBigInt: map[string]int64{"foo": 12}, InBytes: 12121, InPkts: 12, OutBytes: 13, OutPkts: 1, SrcEthMac: "90:61:ae:fb:c2:19", avroSet: map[string]interface{}{}},
	}

	InputTestingSnmp = []*JCHF{
		&JCHF{CompanyId: 10, SrcAddr: "10.2.2.1", Protocol: "UDP", DstAddr: "2001:db8::68", Timestamp: 1, L4DstPort: 80, OutputPort: IfaceID(20), EventType: KENTIK_EVENT_SNMP_INT_METRIC, CustomStr: map[string]string{"foo": "bar"}, CustomInt: map[string]int32{"foo": 1}, CustomBigInt: map[string]int64{"foo": 12}, InBytes: 12121, InPkts: 12, OutBytes: 13, OutPkts: 1, SrcEthMac: "90:61:ae:fb:c2:19", avroSet: map[string]interface{}{}},
		&JCHF{CompanyId: 10, SrcAddr: "10.2.2.1", Protocol: "UDP", DstAddr: "2001:db8::68", Timestamp: 1, L4DstPort: 90, OutputPort: IfaceID(40), EventType: KENTIK_EVENT_SNMP_INT_METRIC, CustomStr: map[string]string{"foo": "sas"}, CustomInt: map[string]int32{"foo": 2}, CustomBigInt: map[string]int64{"foo": 22}, InBytes: 12121, InPkts: 12, OutBytes: 13, OutPkts: 1, SrcEthMac: "90:61:ae:fb:c2:19", avroSet: map[string]interface{}{}},
	}
)
