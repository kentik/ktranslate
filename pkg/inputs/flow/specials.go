package flow

import (
	"strings"

	"github.com/kentik/ktranslate"

	"github.com/netsampler/goflow2/producer"
)

func loadASA(cfg *ktranslate.FlowInputConfig) EntConfig {
	config := EntConfig{
		FlowConfig: producer.ProducerConfig{
			IPFIX: producer.IPFIXProducerConfig{
				Mapping: []producer.NetFlowMapField{
					producer.NetFlowMapField{
						Type:        231,
						Destination: "CustomInteger1",
					},
					producer.NetFlowMapField{
						Type:        232,
						Destination: "CustomInteger2",
					},
					producer.NetFlowMapField{
						Type:        298,
						Destination: "CustomInteger3",
					},
					producer.NetFlowMapField{
						Type:        299,
						Destination: "CustomInteger4",
					},
					producer.NetFlowMapField{
						Type:        233,
						Destination: "CustomInteger5",
					},
				},
			},
			NetFlowV9: producer.NetFlowV9ProducerConfig{
				Mapping: []producer.NetFlowMapField{
					producer.NetFlowMapField{
						Type:        231,
						Destination: "CustomInteger1",
					},
					producer.NetFlowMapField{
						Type:        232,
						Destination: "CustomInteger2",
					},
					producer.NetFlowMapField{
						Type:        298,
						Destination: "CustomInteger3",
					},
					producer.NetFlowMapField{
						Type:        299,
						Destination: "CustomInteger4",
					},
					producer.NetFlowMapField{
						Type:        233,
						Destination: "CustomInteger5",
					},
				},
			},
		},
		NameMap: map[string]string{
			"CustomInteger1": "in_bytes",
			"CustomInteger2": "out_bytes",
			"CustomInteger3": "in_pkts",
			"CustomInteger4": "out_pkts",
			"CustomInteger5": "Firewall Event",
		},
	}

	for field, _ := range config.NameMap {
		if !strings.Contains(cfg.MessageFields, field) {
			cfg.MessageFields = cfg.MessageFields + "," + field
		}
	}

	return config
}

func loadNBar(cfg *ktranslate.FlowInputConfig) EntConfig {
	config := EntConfig{
		FlowConfig: producer.ProducerConfig{
			IPFIX: producer.IPFIXProducerConfig{
				Mapping: []producer.NetFlowMapField{
					producer.NetFlowMapField{
						Type:        96,
						Destination: "CustomBytes1",
					},
					producer.NetFlowMapField{
						Type:        12232,
						Destination: "CustomBytes2",
						PenProvided: true,
						Pen:         9,
					},
					producer.NetFlowMapField{
						Type:        12233,
						Destination: "CustomBytes3",
						PenProvided: true,
						Pen:         9,
					},
					producer.NetFlowMapField{
						Type:        12234,
						Destination: "CustomBytes4",
						PenProvided: true,
						Pen:         9,
					},
					producer.NetFlowMapField{
						Type:        12243,
						Destination: "CustomBytes5",
						PenProvided: true,
						Pen:         9,
					},
				},
			},
			NetFlowV9: producer.NetFlowV9ProducerConfig{
				Mapping: []producer.NetFlowMapField{
					producer.NetFlowMapField{
						Type:        96,
						Destination: "CustomBytes1",
					},
				},
			},
		},
		NameMap: map[string]string{
			"CustomBytes1": "application",
			"CustomBytes2": "Application Category",
			"CustomBytes3": "Application Subcategory",
			"CustomBytes4": "Application Group",
			"CustomBytes5": "Application Traffic Class",
		},
	}

	for field, _ := range config.NameMap {
		if !strings.Contains(cfg.MessageFields, field) {
			cfg.MessageFields = cfg.MessageFields + "," + field
		}
	}

	return config
}

func loadPAN(cfg *ktranslate.FlowInputConfig) EntConfig {
	config := EntConfig{
		FlowConfig: producer.ProducerConfig{
			IPFIX: producer.IPFIXProducerConfig{
				Mapping: []producer.NetFlowMapField{
					producer.NetFlowMapField{
						Type:        32,
						Destination: "CustomInteger1",
					},
					producer.NetFlowMapField{
						Type:        148,
						Destination: "CustomInteger2",
					},
					producer.NetFlowMapField{
						Type:        233,
						Destination: "CustomInteger3",
					},
					producer.NetFlowMapField{
						Type:        61,
						Destination: "CustomInteger4",
					},
					producer.NetFlowMapField{
						Type:        56701,
						Destination: "CustomBytes1",
					},
					producer.NetFlowMapField{
						Type:        56702,
						Destination: "CustomBytes2",
					},
				},
			},
		},
		NameMap: map[string]string{
			"CustomInteger1": "ICMP Type",
			"CustomInteger2": "Flow ID",
			"CustomInteger3": "Firewall Event",
			"CustomInteger4": "Direction",
			"CustomBytes1":   "Application ID",
			"CustomBytes2":   "User ID",
		},
	}

	for field, _ := range config.NameMap {
		if !strings.Contains(cfg.MessageFields, field) {
			cfg.MessageFields = cfg.MessageFields + "," + field
		}
	}

	return config
}
