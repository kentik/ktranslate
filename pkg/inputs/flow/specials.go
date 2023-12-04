package flow

import (
	"github.com/kentik/ktranslate"

	producer "github.com/netsampler/goflow2/v2/producer/proto"
)

func loadDefault(cfg *ktranslate.FlowInputConfig) *producer.ProducerConfig {
	config := &producer.ProducerConfig{
		Formatter: producer.FormatterConfig{
			Fields: []string{},
			Key:    []string{"sampler_address"},
			Protobuf: []producer.ProtobufFormatterConfig{
				producer.ProtobufFormatterConfig{
					Name:  "flow_direction",
					Index: 99,
					Type:  "varint",
				},
			},
		},
		IPFIX: producer.IPFIXProducerConfig{
			Mapping: []producer.NetFlowMapField{
				producer.NetFlowMapField{
					Type:        61,
					Destination: "flow_direction",
				},
			},
		},
		NetFlowV9: producer.NetFlowV9ProducerConfig{
			Mapping: []producer.NetFlowMapField{
				producer.NetFlowMapField{
					Type:        61,
					Destination: "flow_direction",
				},
			},
		},
	}

	return config
}

func loadASA(cfg *ktranslate.FlowInputConfig) *producer.ProducerConfig {
	config := &producer.ProducerConfig{
		Formatter: producer.FormatterConfig{
			Fields: []string{},
			Key:    []string{"sampler_address"},
			Protobuf: []producer.ProtobufFormatterConfig{
				producer.ProtobufFormatterConfig{
					Name:  "flow_direction",
					Index: 99,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "in_bytes",
					Index: 100,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "out_bytes",
					Index: 101,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "in_pkts",
					Index: 102,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "put_pkts",
					Index: 103,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "firewall_event",
					Index: 104,
					Type:  "bytes",
				},
			},
		},
		IPFIX: producer.IPFIXProducerConfig{
			Mapping: []producer.NetFlowMapField{
				producer.NetFlowMapField{
					Type:        61,
					Destination: "flow_direction",
				},
				producer.NetFlowMapField{
					Type:        231,
					Destination: "in_bytes",
				},
				producer.NetFlowMapField{
					Type:        232,
					Destination: "out_bytes",
				},
				producer.NetFlowMapField{
					Type:        298,
					Destination: "in_pkts",
				},
				producer.NetFlowMapField{
					Type:        299,
					Destination: "out_pkts",
				},
				producer.NetFlowMapField{
					Type:        233,
					Destination: "firewall_event",
				},
			},
		},
		NetFlowV9: producer.NetFlowV9ProducerConfig{
			Mapping: []producer.NetFlowMapField{
				producer.NetFlowMapField{
					Type:        61,
					Destination: "flow_direction",
				},
				producer.NetFlowMapField{
					Type:        231,
					Destination: "in_bytes",
				},
				producer.NetFlowMapField{
					Type:        232,
					Destination: "out_bytes",
				},
				producer.NetFlowMapField{
					Type:        298,
					Destination: "in_pkts",
				},
				producer.NetFlowMapField{
					Type:        299,
					Destination: "out_pkts",
				},
				producer.NetFlowMapField{
					Type:        233,
					Destination: "firewall_event",
				},
			},
		},
	}

	return config
}

func loadNBar(cfg *ktranslate.FlowInputConfig) *producer.ProducerConfig {
	config := &producer.ProducerConfig{
		Formatter: producer.FormatterConfig{
			Fields: []string{},
			Key:    []string{"sampler_address"},
			Protobuf: []producer.ProtobufFormatterConfig{
				producer.ProtobufFormatterConfig{
					Name:  "flow_direction",
					Index: 99,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "application",
					Index: 100,
					Type:  "bytes",
				},
				producer.ProtobufFormatterConfig{
					Name:  "application_category",
					Index: 101,
					Type:  "bytes",
				},
				producer.ProtobufFormatterConfig{
					Name:  "application_subcategory",
					Index: 102,
					Type:  "bytes",
				},
				producer.ProtobufFormatterConfig{
					Name:  "application_group",
					Index: 103,
					Type:  "bytes",
				},
				producer.ProtobufFormatterConfig{
					Name:  "application_traffic_class",
					Index: 104,
					Type:  "bytes",
				},
			},
		},
		IPFIX: producer.IPFIXProducerConfig{
			Mapping: []producer.NetFlowMapField{
				producer.NetFlowMapField{
					Type:        61,
					Destination: "flow_direction",
				},
				producer.NetFlowMapField{
					Type:        96,
					Destination: "application",
				},
				producer.NetFlowMapField{
					Type:        12232,
					Destination: "application_category",
					PenProvided: true,
					Pen:         9,
				},
				producer.NetFlowMapField{
					Type:        12233,
					Destination: "application_subcategory",
					PenProvided: true,
					Pen:         9,
				},
				producer.NetFlowMapField{
					Type:        12234,
					Destination: "application_group",
					PenProvided: true,
					Pen:         9,
				},
				producer.NetFlowMapField{
					Type:        12243,
					Destination: "application_traffic_class",
					PenProvided: true,
					Pen:         9,
				},
			},
		},
		NetFlowV9: producer.NetFlowV9ProducerConfig{
			Mapping: []producer.NetFlowMapField{
				producer.NetFlowMapField{
					Type:        61,
					Destination: "flow_direction",
				},
				producer.NetFlowMapField{
					Type:        96,
					Destination: "application",
				},
			},
		},
	}

	return config
}

func loadPAN(cfg *ktranslate.FlowInputConfig) *producer.ProducerConfig {
	config := &producer.ProducerConfig{
		Formatter: producer.FormatterConfig{
			Fields: []string{},
			Key:    []string{"sampler_address"},
			Protobuf: []producer.ProtobufFormatterConfig{
				producer.ProtobufFormatterConfig{
					Name:  "flow_direction",
					Index: 99,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "icmp_type",
					Index: 100,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "flow_id",
					Index: 101,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "firewall_event",
					Index: 102,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "direction",
					Index: 103,
					Type:  "varint",
				},
				producer.ProtobufFormatterConfig{
					Name:  "application_id",
					Index: 104,
					Type:  "bytes",
				},
				producer.ProtobufFormatterConfig{
					Name:  "user_id",
					Index: 105,
					Type:  "bytes",
				},
			},
		},
		IPFIX: producer.IPFIXProducerConfig{
			Mapping: []producer.NetFlowMapField{
				producer.NetFlowMapField{
					Type:        61,
					Destination: "flow_direction",
				},
				producer.NetFlowMapField{
					Type:        32,
					Destination: "icmp_type",
				},
				producer.NetFlowMapField{
					Type:        148,
					Destination: "flow_id",
				},
				producer.NetFlowMapField{
					Type:        233,
					Destination: "firewall_event",
				},
				producer.NetFlowMapField{
					Type:        61,
					Destination: "direction",
				},
				producer.NetFlowMapField{
					Type:        56701,
					Destination: "application_id",
				},
				producer.NetFlowMapField{
					Type:        56702,
					Destination: "user_id",
				},
			},
		},
	}

	return config
}
