package util

import (
	"fmt"
	"regexp"
	"sort"
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

const (
	FORCE_MATCH_TOKEN  = "!"
	OR_TOKEN           = "||"
	MAX_ATTR_FOR_SNMP  = 64
	NEGATE_MATCH_TOKEN = "DOES_NOT_MATCH"
)

func SetAttr(attr map[string]interface{}, in *kt.JCHF, metrics map[string]kt.MetricInfo, lastMetadata *kt.LastMetadata, stripMetrics bool) {
	mapr := in.Flatten()
	for k, v := range mapr {
		if DroppedAttrs[k] {
			continue // Skip because we don't want this messing up cardinality.
		}

		if _, ok := metrics[k]; ok { // Skip because this one is a metric, not an attribute.
			continue
		}
		if stripMetrics && strings.HasPrefix(k, kt.StringPrefix) {
			if _, ok := metrics[k[len(kt.StringPrefix):]]; ok { // Skip because this one is a metric, not an attribute.
				continue
			}
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

	// Copy this over as a deap struct.
	if in.Har != nil {
		attr["har_file"] = &in.Har
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
					attr[k] = v.GetValue()
				}
			}
			// If the index is longer, see if there's a parent table to look into also.
			pts := strings.Split(idx, ".")
			if len(pts) > 1 {
				parent := strings.Join(pts[0:len(pts)-1], ".")
				if table, ok := lastMetadata.Tables[parent]; ok {
					for k, v := range table.Customs {
						attr[k] = v.GetValue()
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

		// Now run the metadata scripts in their own loop.
		if idx != "" {
			if table, ok := lastMetadata.Tables[idx]; ok {
				for k, v := range table.Customs {
					if s := v.GetScript(); s != nil {
						s.EnrichMib(idx, k, attr, lastMetadata)
					}
				}
			}
			// If the index is longer, see if there's a parent table to look into also.
			pts := strings.Split(idx, ".")
			if len(pts) > 1 {
				parent := strings.Join(pts[0:len(pts)-1], ".")
				if table, ok := lastMetadata.Tables[parent]; ok {
					for k, v := range table.Customs {
						if s := v.GetScript(); s != nil {
							s.EnrichMib(idx, k, attr, lastMetadata)
						}
					}
				}
			}
		}
	}
}

func DropOnFilter(attr map[string]interface{}, lastMetadata *kt.LastMetadata, isIfMetric bool) bool {
	if lastMetadata != nil && lastMetadata.MatchAttr != nil {
		dropOnAdminStatus := false // Default this false.
		keepForOtherMatch := false // We use inverse of this, so default is to drop flow. BUT, need to have a match set else all flow passes.
		seenNonAdmin := 0
		negateMatch := false
		for k, re := range lastMetadata.MatchAttr {
			forceMatch := false
			if strings.HasPrefix(k, FORCE_MATCH_TOKEN) { // Handle forcing this column to exist here.
				k = k[1:]
				forceMatch = true
			}
			if k == NEGATE_MATCH_TOKEN { // If this is set, return the opposite of all matches.
				negateMatch = true
				continue
			}

			if strings.HasPrefix(k, "(") && strings.HasSuffix(k, ")") { // Handle paren groupings here.
				k = k[1 : len(k)-1]
			}
			keyCheck := strings.Split(k, OR_TOKEN)
			var cont, br bool
			seenAny := false
			if forceMatch {
				for _, key := range keyCheck {
					if _, ok := attr[key]; ok {
						seenAny = true
						break
					}
				}
			}

			for _, key := range keyCheck {
				cont, br = checkFilter(attr, key, re, forceMatch, seenAny, isIfMetric, &dropOnAdminStatus, &keepForOtherMatch, &seenNonAdmin)
			}
			if cont {
				continue
			}
			if br {
				break
			}
		}

		// This is special cased as an AND to any other attributes.
		if dropOnAdminStatus {
			return true
		} else {
			if seenNonAdmin > 0 {
				if negateMatch {
					return keepForOtherMatch
				} else {
					return !keepForOtherMatch
				}
			} else {
				if negateMatch {
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}

func checkFilter(attr map[string]interface{}, k string, re *regexp.Regexp, forceMatch bool, seenAny bool, isIfMetric bool, dropOnAdminStatus *bool, keepForOtherMatch *bool, seenNonAdmin *int) (cont bool, br bool) {
	// If this is not an interface attribute, skip interface matches.
	if !isIfMetric && (k == kt.AdminStatus || strings.HasPrefix(k, "if_") || strings.HasPrefix(k, "input_if_") || strings.HasPrefix(k, "output_if_")) {
		cont = true
		return
	}
	if v, ok := attr[k]; ok {
		if strv, ok := v.(string); ok {
			if k == kt.AdminStatus { // If admin status is causing us to drop, drop right away.
				*dropOnAdminStatus = !re.MatchString(strv)
				if *dropOnAdminStatus == true {
					br = true
					return
				}
			} else { // Otherwise, OR all the matches together. Keep trying until we find an RE which matches.
				*seenNonAdmin++
				if !*keepForOtherMatch {
					*keepForOtherMatch = re.MatchString(strv)
					if !*keepForOtherMatch && (forceMatch && !seenAny) { // In this case, we failed the force match so break right now.
						br = true
						return
					}
				}
			}
		}
	} else { // If the key doesn't exist but the match tells us to force this, drop here.
		if forceMatch && !seenAny {
			*seenNonAdmin++
			*keepForOtherMatch = false
			br = true
			return
		}
	}

	return
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
	lm.XtraInfo = in.CustomMetrics

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

		8: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "status"}, "fetch_size_|_ping_min_rtt": kt.MetricInfo{Name: "size"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "ttlb"}},

		9: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "status"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "trx_time"}},

		10: map[string]kt.MetricInfo{"fetch_status_|_ping_sent_|_trace_time": kt.MetricInfo{Name: "time"}, "fetch_ttlb_|_ping_lost": kt.MetricInfo{Name: "validation"}},
	}
)

func GetSynMetricNameSet(rt int32) map[string]kt.MetricInfo {
	return synMetrics[rt]
}

// [aggregate_interval:0 app_protocol:18 avg_jitter:934 avg_latency:27366 avg_weighted_latency:0 configured_task_type:5 health_moment_task_type:5 jitter_health:300 latency_health:300 member_id:21241 packet_loss_health:300 rolling_avg_jitter:1127 rolling_avg_latency:27172 rolling_avg_weighted_latency:0 rolling_stddev_jitter:303 rolling_stddev_latency:452 size:0 status:0
func GetSyngestMetricNameSet() map[string]kt.MetricInfo {
	return map[string]kt.MetricInfo{
		"avg_jitter":             kt.MetricInfo{Name: "avg_jitter"},
		"avg_latency":            kt.MetricInfo{Name: "avg_latency"},
		"rolling_avg_jitter":     kt.MetricInfo{Name: "rolling_avg_jitter"},
		"rolling_avg_latency":    kt.MetricInfo{Name: "rolling_avg_latency"},
		"rolling_stddev_jitter":  kt.MetricInfo{Name: "rolling_stddev_jitter"},
		"rolling_stddev_latency": kt.MetricInfo{Name: "rolling_stddev_latency"},
		"size":                   kt.MetricInfo{Name: "size"},
		"status":                 kt.MetricInfo{Name: "status"},
	}
}

// List of attributes to not pass to NR.
var removeAttrForSnmp = []string{
	"Uptime",
	"if_LastChange",
	"SysServices",
	"if_Mtu",
	"if_ConnectorPresent",
	"output_port",
	"input_port",
	"DropMetric",
	"sysoid_vendor",
}

var keepAcrossTables = map[string]bool{
	"device_name":    true,
	"eventType":      true,
	"provider":       true,
	"sysoid_profile": true,
	kt.IndexVar:      true,
	"if_Index":       true,
	"src_addr":       true,
}

var allowSysAttr = map[string]bool{
	"Uptime":        true,
	"MinRttMs":      true,
	"MaxRttMs":      true,
	"AvgRttMs":      true,
	"StdDevRtt":     true,
	"PacketsSent":   true,
	"PacketsRecv":   true,
	"PacketLossPct": true,
}

func CopyAttrForSnmp(attr map[string]interface{}, metricName string, name kt.MetricInfo, lm *kt.LastMetadata, gentleCardinality bool, dropSrcAddr bool) map[string]interface{} {
	attrNew := map[string]interface{}{
		"objectIdentifier":     name.Oid,
		"mib-name":             name.Mib,
		"instrumentation.name": name.Profile,
	}

	// If set, add this in.
	durSec := name.PollDur.Seconds()
	if durSec > 0 {
		attrNew["poll_duration_sec"] = name.PollDur.Seconds() + kt.PollAdjustTime
	}

	for k, v := range attr {
		if !allowSysAttr[metricName] { // Only allow Sys* attributes on specific metrics.
			if strings.HasPrefix(k, "Sys") || (dropSrcAddr && k == "src_addr") {
				continue
			}
		}

		newKey := k
		if strings.HasPrefix(k, kt.StringPrefix) {
			newKey = k[len(kt.StringPrefix):]
		}

		if name.Table != "" && metricName != newKey {
			if _, ok := keepAcrossTables[newKey]; !ok { // If we want this attribute in every table, list it here.
				attrNew["mib-table"] = name.Table

				// See if the metadata knows about this attribute.
				if tableName, allNames, ok := lm.GetTableName(newKey); ok && len(allNames) > 0 {
					if !allNames[name.Table] && tableName != kt.DeviceTagTable {
						continue
					}
				} else {
					// If this metric comes from a specific table, only show attributes for this table.
					if !strings.HasPrefix(newKey, name.Table) {
						if !allNames[name.Table] && tableName != kt.DeviceTagTable {
							continue
						}
					}
				}
			}
		}

		// Case where metric has no table.
		if name.Table == "" {
			if tableName, _, ok := lm.GetTableName(newKey); ok {
				if tableName != "" && tableName != kt.DeviceTagTable {
					continue
				}
			}
		}

		attrNew[newKey] = v
	}

	if gentleCardinality {
		// Delete a few attributes we don't want adding to cardinality.
		for _, key := range removeAttrForSnmp {
			delete(attrNew, key)
		}

		// These are getting dropped sometimes and we don't need to one in the other series.
		if metricName == "if_AdminStatus" {
			delete(attrNew, "if_OperStatus")
		} else if metricName == "if_OperStatus" {
			delete(attrNew, "if_AdminStatus")
		}

		if len(attrNew) > MAX_ATTR_FOR_SNMP {
			// Since NR limits us to 100 attributes, we need to prune. Take the first 100 lexographical keys.
			keys := make([]string, len(attrNew))
			i := 0
			for k, _ := range attrNew {
				keys[i] = k
				i++
			}
			sort.Strings(keys)
			for _, k := range keys[MAX_ATTR_FOR_SNMP-3:] {
				delete(attrNew, k)
			}

			// Force these to be back in.
			attrNew["objectIdentifier"] = name.Oid
			attrNew["mib-name"] = name.Mib
			attrNew["instrumentation.name"] = name.Profile
		}
	}

	if attrNew["mib-name"] == "" {
		delete(attrNew, "mib-name")
	}

	return attrNew
}
