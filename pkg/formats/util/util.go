package util

import (
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	DroppedAttrs = map[string]bool{
		"timestamp":               true,
		"sampled_packet_size":     true,
		"Lat/Long Dest":           true,
		"MEMBER_ID":               true,
		"dst_eth_mac":             true,
		"src_eth_mac":             true,
		"Manufacturer":            true,
		"Error Cause/Trace Route": true,
		"Hop Data":                true,
		"STR01":                   true,
		"ULT_EXIT_PORT":           true,
		"Task ID":                 true,
		"APP_PROTOCOL":            true,
		"Agent ID":                true,
		"ULT_EXIT_DEVICE_ID":      true,
		"device_id":               true,
		"kt_functional_testing":   true,
		"CLIENT_NW_LATENCY_MS":    true,
		"APPL_LATENCY_MS":         true,
		"SERVER_NW_LATENCY_MS":    true,
		"CONNECTION_ID":           true,
	}
)

func SetAttr(attr map[string]interface{}, in *kt.JCHF, metrics map[string]bool, lastMetadata *kt.LastMetadata) {
	mapr := in.Flatten()
	for k, v := range mapr {
		if DroppedAttrs[k] {
			continue // Skip because we don't want this messing up cardinality.
		}

		switch vt := v.(type) {
		case string:
			if !metrics[k] && vt != "" {
				attr[k] = vt
			}
		case int64:
			if !metrics[k] && vt > 0 {
				attr[k] = int(vt)
			}
		case int32:
			if !metrics[k] && vt > 0 {
				attr[k] = int(vt)
			}
		}
	}

	if lastMetadata != nil {
		for k, v := range lastMetadata.DeviceInfo {
			attr[k] = v
		}

		if in.OutputPort != in.InputPort {
			if ii, ok := lastMetadata.InterfaceInfo[in.InputPort]; ok {
				for k, v := range ii {
					if v != "" {
						attr["input_if_"+k] = v
					}
				}
			}
			if ii, ok := lastMetadata.InterfaceInfo[in.OutputPort]; ok {
				for k, v := range ii {
					if v != "" {
						attr["output_if_"+k] = v
					}
				}
			}
		} else {
			if ii, ok := lastMetadata.InterfaceInfo[in.OutputPort]; ok {
				for k, v := range ii {
					if v != "" {
						attr["if_"+k] = v
					}
				}
			}
		}
	}
}
