package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	DroppedAttrs = map[string]bool{
		"result_type":             true,
		"timestamp":               true,
		"sampled_packet_size":     true,
		"lat/long_dest":           true,
		"member_id":               true,
		"dst_eth_mac":             true,
		"src_eth_mac":             true,
		"Manufacturer":            true,
		"error_cause/trace_route": true,
		"hop_data":                true,
		"str01":                   true,
		"ult_exit_port":           true,
		"app_protocol":            true,
		"kt_functional_testing":   true,
		"client_nw_latency_ms":    true,
		"appl_latency_ms":         true,
		"server_nw_latency_ms":    true,
		"connection_id":           true,
	}
)

func SetAttr(attr map[string]interface{}, in *kt.JCHF, metrics map[string]string, lastMetadata *kt.LastMetadata) {
	mapr := in.Flatten()
	for k, v := range mapr {
		if DroppedAttrs[k] {
			continue // Skip because we don't want this messing up cardinality.
		}

		if _, ok := metrics[k]; ok { // Skip because this one is a metric, not an attribute.
			continue
		}

		switch vt := v.(type) {
		case string, kt.Provider:
			if vt != "" {
				attr[k] = vt
			}
		case int64:
			if vt > 0 {
				attr[k] = int(vt)
			}
		case int32:
			if vt > 0 {
				attr[k] = int(vt)
			}
		default:
			panic(fmt.Sprintf("Unknown type: %v", v.(int32)))
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

/**
jit_rtt

application
application-type
company_id
dst_threat_host
end_timestamp
result-type-str
simple_trf_prod
task_id
trf-origination
trf-term
net-bound
utl-exit
dst-cloud
dst-cloud-service
src-cloud
src-cloud-service

add -- url for info

*/

var (
	synMetrics = map[int32]map[string]string{
		0: map[string]string{"error": "error", "fetch_status_|_ping_sent_|_trace_time": "sent", "fetch_ttlb_|_ping_lost": "lost",
			"fetch_size_|_ping_min_rtt": "min_rtt", "ping_max_rtt": "max_rtt", "ping_avg_rtt": "avg_rtt", "ping_std_rtt": "std_rtt", "ping_jit_rtt": "jit_rtt"},

		1: map[string]string{"timeout": "timeout", "fetch_status_|_ping_sent_|_trace_time": "sent", "fetch_ttlb_|_ping_lost": "lost",
			"fetch_size_|_ping_min_rtt": "min_rtt", "ping_max_rtt": "max_rtt", "ping_avg_rtt": "avg_rtt", "ping_std_rtt": "std_rtt", "ping_jit_rtt": "jit_rtt"},

		2: map[string]string{"fetch_status_|_ping_sent_|_trace_time": "sent", "fetch_ttlb_|_ping_lost": "lost",
			"fetch_size_|_ping_min_rtt": "min_rtt", "ping_max_rtt": "max_rtt", "ping_avg_rtt": "avg_rtt", "ping_std_rtt": "std_rtt", "ping_jit_rtt": "jit_rtt"},

		3: map[string]string{"fetch_status_|_ping_sent_|_trace_time": "status", "fetch_ttlb_|_ping_lost": "ttlb",
			"fetch_size_|_ping_min_rtt": "size"},

		4: map[string]string{"fetch_status_|_ping_sent_|_trace_time": "time"},

		5: map[string]string{"fetch_status_|_ping_sent_|_trace_time": "sent", "fetch_ttlb_|_ping_lost": "lost",
			"fetch_size_|_ping_min_rtt": "min_rtt", "ping_max_rtt": "max_rtt", "ping_avg_rtt": "avg_rtt", "ping_std_rtt": "std_rtt", "ping_jit_rtt": "jit_rtt"},

		6: map[string]string{"fetch_status_|_ping_sent_|_trace_time": "time", "fetch_ttlb_|_ping_lost": "code"},

		7: map[string]string{"fetch_status_|_ping_sent_|_trace_time": "time", "lat/long_dest": "port"},
	}
)

func GetSynMetricNameSet(rt int32) map[string]string {
	return synMetrics[rt]
}
