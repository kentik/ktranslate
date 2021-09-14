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
		"service_name":            true,
	}
)

func SetAttr(attr map[string]interface{}, in *kt.JCHF, metrics map[string]kt.MetricInfo, lastMetadata *kt.LastMetadata) {
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
			if vt != "" && vt != "-" && vt != "--" {
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

		idx := in.CustomStr[kt.IndexVar]
		if idx != "" {
			if idx[0:1] == "." {
				idx = idx[1:]
			}
			if table, ok := lastMetadata.Tables[idx]; ok {
				for k, v := range table.Customs {
					attr[k] = v
				}
				for k, v := range table.CustomInts {
					attr[k] = v
				}
			}
			// If the index is longer, see if there's a parent table to look into also.
			pts := strings.Split(idx, ".")
			if len(pts) > 1 {
				parent := strings.Join(pts[0:len(pts)-1], ".")
				if table, ok := lastMetadata.Tables[parent]; ok {
					for k, v := range table.Customs {
						attr[k] = v
					}
					for k, v := range table.CustomInts {
						attr[k] = v
					}
				}
			}
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

		// Finally, drop any values which do not match the whitelist.
		if lastMetadata.MatchAttr != nil {
			dropOnAdminStatus := false // Default this false.
			keepForOtherMatch := false // We use inverse of this, so default is to drop flow. BUT, need to have a match set else all flow passes.
			seenAdminStatus := false
			for k, re := range lastMetadata.MatchAttr {
				if v, ok := attr[k]; ok {
					if strv, ok := v.(string); ok {
						if k == kt.AdminStatus { // If admin status is causing us to drop, drop right away.
							seenAdminStatus = true
							dropOnAdminStatus = !re.MatchString(strv)
							if dropOnAdminStatus == true {
								break
							}
						} else { // Otherwise, OR all the matches together. Keep trying until we find an RE which matches.
							if !keepForOtherMatch {
								keepForOtherMatch = re.MatchString(strv)
							}
						}
					}
				}
			}

			// This is special cased as an AND to any other attributes.
			if dropOnAdminStatus {
				attr[kt.DropMetric] = true
			} else {
				if seenAdminStatus && len(lastMetadata.MatchAttr) > 1 {
					attr[kt.DropMetric] = !keepForOtherMatch
				} else if !seenAdminStatus && len(lastMetadata.MatchAttr) > 0 {
					attr[kt.DropMetric] = !keepForOtherMatch
				} else {
					attr[kt.DropMetric] = false
				}
			}
		}
	}
}

func SetMetadata(in *kt.JCHF) *kt.LastMetadata {
	lm := kt.LastMetadata{
		DeviceInfo:    map[string]interface{}{},
		InterfaceInfo: map[kt.IfaceID]map[string]interface{}{},
		Tables:        map[string]kt.DeviceTableMetadata{},
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
	lm.Tables = in.CustomTables
	lm.MatchAttr = in.MatchAttr

	return &lm
}

var (
	synMetrics = map[int32]map[string]kt.MetricInfo{
		0: map[string]kt.MetricInfo{"error": kt.MetricInfo{Name: "error"}, "fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "sent"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "lost"},
			"fetch_size_|_ping_min_rtt": kt.MetricInfo{Name: "min_rtt"}, "ping_max_rtt": kt.MetricInfo{Name: "max_rtt"}, "ping_avg_rtt": kt.MetricInfo{Name: "avg_rtt"}, "ping_std_rtt": kt.MetricInfo{Name: "std_rtt"}, "ping_jit_rtt": kt.MetricInfo{Name: "jit_rtt"}},

		1: map[string]kt.MetricInfo{"timeout": kt.MetricInfo{Name: "timeout"}, "fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "sent"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "lost"},
			"fetch_size_|_ping_min_rtt": kt.MetricInfo{Name: "min_rtt"}, "ping_max_rtt": kt.MetricInfo{Name: "max_rtt"}, "ping_avg_rtt": kt.MetricInfo{Name: "avg_rtt"}, "ping_std_rtt": kt.MetricInfo{Name: "std_rtt"},
			"ping_jit_rtt": kt.MetricInfo{Name: "jit_rtt"}},

		2: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "sent"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "lost"},
			"fetch_size_|_ping_min_rtt": kt.MetricInfo{Name: "min_rtt"}, "ping_max_rtt": kt.MetricInfo{Name: "max_rtt"}, "ping_avg_rtt": kt.MetricInfo{Name: "avg_rtt"},
			"ping_std_rtt": kt.MetricInfo{Name: "std_rtt"}, "ping_jit_rtt": kt.MetricInfo{Name: "jit_rtt"}},

		3: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "status"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "ttlb"},
			"fetch_size_|_ping_min_rtt": kt.MetricInfo{Name: "size"}},

		4: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "time"}},

		5: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "sent"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "lost"},
			"fetch_size_|_ping_min_rtt": kt.MetricInfo{Name: "min_rtt"}, "ping_max_rtt": kt.MetricInfo{Name: "max_rtt"}, "ping_avg_rtt": kt.MetricInfo{Name: "avg_rtt"},
			"ping_std_rtt": kt.MetricInfo{Name: "std_rtt"}, "ping_jit_rtt": kt.MetricInfo{Name: "jit_rtt"}},

		6: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "time"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "code"}},

		7: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "time"}, "lat/long_dest": kt.MetricInfo{Name: "port"}},
	}
)

func GetSynMetricNameSet(rt int32) map[string]kt.MetricInfo {
	return synMetrics[rt]
}
