package kt

import ()

var (
	InputTesting = []*JCHF{
		&JCHF{CompanyId: 10, SrcAddr: "10.2.2.1", Protocol: 1, DstAddr: "2001:db8::68", Timestamp: 1, L4DstPort: 80, OutputPort: IfaceID(20), CustomStr: map[string]string{"foo": "bar"}, CustomInt: map[string]int32{"foo": 1}, CustomBigInt: map[string]int64{"foo": 12}, InBytes: 12121, InPkts: 12, OutBytes: 13, OutPkts: 1, SrcEthMac: "90:61:ae:fb:c2:19", avroSet: map[string]interface{}{}},
		&JCHF{CompanyId: 10, SrcAddr: "3.2.2.2", InBytes: 1, OutBytes: 12, InPkts: 12, OutPkts: 1, Protocol: 2, DstAddr: "2001:db8::69", SrcEthMac: "90:61:ae:fb:c2:20", Timestamp: 2, CustomStr: map[string]string{"tar": "far"}, avroSet: map[string]interface{}{}},
	}
)
