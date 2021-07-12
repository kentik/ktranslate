package avro

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/linkedin/goavro/v2"
)

type AvroFormat struct {
	logger.ContextL
	codec       *goavro.Codec
	codecRollup *goavro.Codec
	cfg         goavro.OCFConfig
	compression kt.Compression
	compName    string
}

const chfAvroSchema = `
{
"namespace": "kentik.com",
"type": "record",
"name": "kflow",
"fields": [
{"name": "timestamp", "type": "long", "default": 0},
{"name": "dst_as", "type": "long", "default": 0},
{"name": "dst_geo", "type": "string", "default": ""},
{"name": "header_len", "type": "long", "default": 0},
{"name": "in_bytes", "type": "long", "default": 0},
{"name": "in_pkts", "type": "long", "default": 0},
{"name": "input_port", "type": "long", "default": 0},
{"name": "ip_size", "type": "long", "default": 0},
{"name": "dst_addr", "type": "string", "default": ""},
{"name": "src_addr", "type": "string", "default": ""},
{"name": "l4_dst_port", "type": "long", "default": 0},
{"name": "l4_src_port", "type": "long", "default": 0},
{"name": "output_port", "type": "long", "default": 0},
{"name": "protocol", "type": "string", "default": ""},
{"name": "sampled_packet_size", "type": "long", "default": 0},
{"name": "src_as", "type": "long", "default": 0},
{"name": "src_geo", "type": "string", "default": ""},
{"name": "tcp_flags", "type": "long", "default": 0},
{"name": "tos", "type": "long", "default": 0},
{"name": "vlan_in", "type": "long", "default": 0},
{"name": "vlan_out", "type": "long", "default": 0},
{"name": "out_bytes", "type": "long", "default": 0},
{"name": "out_pkts", "type": "long", "default": 0},
{"name": "tcp_rx", "type": "long", "default": 0},
{"name": "src_flow_tags", "type": "string", "default": ""},
{"name": "dst_flow_tags", "type": "string", "default": ""},
{"name": "sample_rate", "type": "long", "default": 0},
{"name": "device_id", "type": "long", "default": 0},
{"name": "device_name", "type": "string", "default": ""},
{"name": "company_id", "type": "long", "default": 0},
{"name": "dst_bgp_as_path", "type": "string", "default": ""},
{"name": "dst_bgp_comm", "type": "string", "default": ""},
{"name": "src_bpg_as_path", "type": "string", "default": ""},
{"name": "src_bgp_comm", "type": "string", "default": ""},
{"name": "src_nexthop_as", "type": "long", "default": 0},
{"name": "dst_nexthop_as", "type": "long", "default": 0},
{"name": "src_geo_region", "type": "string", "default": ""},
{"name": "dst_geo_region", "type": "string", "default": ""},
{"name": "src_geo_city", "type": "string", "default": ""},
{"name": "dst_geo_city", "type": "string", "default": ""},
{"name": "dst_nexthop", "type": "string", "default": ""},
{"name": "src_nexthop", "type": "string", "default": ""},
{"name": "src_route_prefix", "type": "string", "default": ""},
{"name": "dst_route_prefix", "type": "string", "default": ""},
{"name": "src_second_asn", "type": "long", "default": 0},
{"name": "dst_second_asn", "type": "long", "default": 0},
{"name": "src_third_asn", "type": "long", "default": 0},
{"name": "dst_third_asn", "type": "long", "default": 0},
{"name": "src_eth_mac", "type": "string", "default": ""},
{"name": "dst_eth_mac", "type": "string", "default": ""},
{"name": "input_int_desc", "type": "string", "default": ""},
{"name": "output_int_desc", "type": "string", "default": ""},
{"name": "input_int_alias", "type": "string", "default": ""},
{"name": "output_int_alias", "type": "string", "default": ""},
{"name": "input_interface_capacity", "type": "long", "default": 0},
{"name": "output_interface_capacity", "type": "long", "default": 0},
{"name": "input_interface_ip", "type": "string", "default": ""},
{"name": "output_interface_ip", "type": "string", "default": ""},
{"name": "custom_str", "type": {"type": "map", "values":"string"}},
{"name": "custom_int", "type": {"type": "map", "values":"int"}},
{"name": "custom_bigint", "type": {"type": "map", "values":"long"}}
]
}
`

const chfAvroRollupSchema = `
{
"namespace": "kentik.com",
"type": "record",
"name": "rollup",
"fields": [
{"name": "dimension", "type": "string", "default": ""},
{"name": "metric", "type": "long", "default": 0}
]
}
`

func NewFormat(log logger.Underlying, comp kt.Compression) (*AvroFormat, error) {
	af := &AvroFormat{
		compression: comp,
		ContextL:    logger.NewContextLFromUnderlying(logger.SContext{S: "avroFormat"}, log),
	}

	codec, err := goavro.NewCodec(chfAvroSchema)
	if err != nil {
		return nil, err
	}
	af.codec = codec

	codecRollup, err := goavro.NewCodec(chfAvroRollupSchema)
	if err != nil {
		return nil, err
	}
	af.codecRollup = codecRollup

	switch string(comp) {
	case goavro.CompressionSnappyLabel:
		af.compName = goavro.CompressionSnappyLabel
	case goavro.CompressionDeflateLabel:
		af.compName = goavro.CompressionDeflateLabel
	case goavro.CompressionNullLabel:
		af.compName = goavro.CompressionNullLabel
	case "none":
		// Noop
	default:
		return nil, fmt.Errorf("Invalid compression option: %s", comp)
	}

	return af, nil
}

func (f *AvroFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()

	values := make([]map[string]interface{}, len(msgs))
	for i, m := range msgs {
		values[i] = m.ToMap()
	}

	cfg := goavro.OCFConfig{
		W:               buf,
		Codec:           f.codec,
		CompressionName: f.compName,
	}

	ocfw, err := goavro.NewOCFWriter(cfg)
	if err != nil {
		return nil, err
	}
	if err = ocfw.Append(values); err != nil {
		return nil, err
	}
	return kt.NewOutput(buf.Bytes()), nil
}

func (f *AvroFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	ior := bytes.NewBuffer(raw.Body)
	br := bufio.NewReader(ior)
	values := make([]map[string]interface{}, 0)
	ocfr, err := goavro.NewOCFReader(br)
	if err != nil {
		return nil, err
	}
	for ocfr.Scan() {
		datum, err := ocfr.Read()
		if err != nil {
			return nil, err
		}
		values = append(values, datum.(map[string]interface{}))
	}
	return values, ocfr.Err()
}

func (f *AvroFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	values := make([]map[string]interface{}, len(rolls))
	for i, r := range rolls {
		values[i] = map[string]interface{}{
			"dimension": r.Dimension,
			"metric":    r.Metric,
		}
	}

	cfg := goavro.OCFConfig{
		W:               buf,
		Codec:           f.codecRollup,
		CompressionName: f.compName,
	}

	ocfw, err := goavro.NewOCFWriter(cfg)
	if err != nil {
		return nil, err
	}
	if err = ocfw.Append(values); err != nil {
		return nil, err
	}
	return kt.NewOutput(buf.Bytes()), nil
}
