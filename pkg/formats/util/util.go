package util

import (
	"strconv"
	"strings"

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

func SetMetadata(in *kt.JCHF) *kt.LastMetadata {
	lm := kt.LastMetadata{
		DeviceInfo:    map[string]interface{}{},
		InterfaceInfo: map[kt.IfaceID]map[string]interface{}{},
	}
	for k, v := range in.CustomStr {
		if DroppedAttrs[k] {
			continue // Skip because we don't want this messing up cardinality.
		}
		if strings.HasPrefix(k, "if.") {
			pts := strings.SplitN(k, ".", 3)
			if len(pts) == 3 {
				if ifint, err := strconv.Atoi(pts[1]); err == nil {
					if _, ok := lm.InterfaceInfo[kt.IfaceID(ifint)]; !ok {
						lm.InterfaceInfo[kt.IfaceID(ifint)] = map[string]interface{}{}
					}
					if v != "" {
						lm.InterfaceInfo[kt.IfaceID(ifint)][pts[2]] = v
					}
				}
			}
		} else {
			if v != "" {
				lm.DeviceInfo[k] = v
			}
		}
	}
	for k, v := range in.CustomInt {
		if DroppedAttrs[k] {
			continue // Skip because we don't want this messing up cardinality.
		}
		if strings.HasPrefix(k, "if.") {
			pts := strings.SplitN(k, ".", 3)
			if len(pts) == 3 {
				if ifint, err := strconv.Atoi(pts[1]); err == nil {
					if _, ok := lm.InterfaceInfo[kt.IfaceID(ifint)]; !ok {
						lm.InterfaceInfo[kt.IfaceID(ifint)] = map[string]interface{}{}
					}
					lm.InterfaceInfo[kt.IfaceID(ifint)][pts[2]] = v
				}
			}
		} else {
			lm.DeviceInfo[k] = v
		}
	}

	return &lm
}
