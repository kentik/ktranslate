package filter

import (
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type WLFilter struct {
	logger.ContextL
	wl map[string]bool
}

func newWhitelistFilter(log logger.Underlying, fd FilterDef) (*WLFilter, error) {
	wf := &WLFilter{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "whitelistFilter"}, log),
		wl:       map[string]bool{},
	}

	for _, d := range fd.Whitelist {
		wf.wl[d] = true
	}

	return wf, nil
}

func (f *WLFilter) Filter(in *kt.JCHF) bool {
	for k, _ := range in.CustomStr {
		if !f.wl[k] {
			delete(in.CustomStr, k)
		}
	}
	for k, _ := range in.CustomInt {
		if !f.wl[k] {
			delete(in.CustomInt, k)
		}
	}
	for k, _ := range in.CustomBigInt {
		if !f.wl[k] {
			delete(in.CustomBigInt, k)
		}
	}
	if !f.wl["timestamp"] {
		in.Timestamp = 0
	}
	if !f.wl["dst_as"] {
		in.DstAs = 0
	}
	if !f.wl["dst_geo"] {
		in.DstGeo = ""
	}
	if !f.wl["header_len"] {
		in.HeaderLen = 0
	}
	if !f.wl["in_bytes"] {
		in.InBytes = 0
	}
	if !f.wl["in_pkts"] {
		in.InPkts = 0
	}
	if !f.wl["input_port"] {
		in.InputPort = 0
	}
	if !f.wl["ip_size"] {
		in.IpSize = 0
	}
	if !f.wl["dst_addr"] {
		in.DstAddr = ""
	}
	if !f.wl["src_addr"] {
		in.SrcAddr = ""
	}
	if !f.wl["l4_dst_port"] {
		in.L4DstPort = 0
	}
	if !f.wl["l4_src_port"] {
		in.L4SrcPort = 0
	}
	if !f.wl["output_port"] {
		in.OutputPort = 0
	}
	if !f.wl["protocol"] {
		in.Protocol = ""
	}
	if !f.wl["sampled_packet_size"] {
		in.SampledPacketSize = 0
	}

	if !f.wl["src_as"] {
		in.SrcAs = 0
	}
	if !f.wl["src_geo"] {
		in.SrcGeo = ""
	}
	if !f.wl["tcp_flags"] {
		in.TcpFlags = 0
	}
	if !f.wl["tos"] {
		in.Tos = 0
	}
	if !f.wl["vlan_in"] {
		in.VlanIn = 0
	}
	if !f.wl["vlan_out"] {
		in.VlanOut = 0
	}
	if !f.wl["next_hop"] {
		in.NextHop = ""
	}
	if !f.wl["mpls_type"] {
		in.MplsType = 0
	}
	if !f.wl["out_bytes"] {
		in.OutBytes = 0
	}
	if !f.wl["out_pkts"] {
		in.OutPkts = 0
	}
	if !f.wl["tcp_rx"] {
		in.TcpRetransmit = 0
	}
	if !f.wl["src_flow_tags"] {
		in.SrcFlowTags = ""
	}
	if !f.wl["dst_flow_tags"] {
		in.DstFlowTags = ""
	}
	if !f.wl["sample_rate"] {
		in.SampleRate = 0
	}
	if !f.wl["device_id"] {
		in.DeviceId = 0
	}
	if !f.wl["device_name"] {
		in.DeviceName = ""
	}
	if !f.wl["company_id"] {
		in.CompanyId = 0
	}
	if !f.wl["dst_bgp_as_path"] {
		in.DstBgpAsPath = ""
	}
	if !f.wl["dst_bgp_comm"] {
		in.DstBgpCommunity = ""
	}
	if !f.wl["src_bpg_as_path"] {
		in.SrcBgpAsPath = ""
	}
	if !f.wl["src_bgp_comm"] {
		in.SrcBgpCommunity = ""
	}
	if !f.wl["src_nexthop_as"] {
		in.SrcNextHopAs = 0
	}
	if !f.wl["dst_nexthop_as"] {
		in.DstNextHopAs = 0
	}
	if !f.wl["src_geo_region"] {
		in.SrcGeoRegion = ""
	}
	if !f.wl["dst_geo_region"] {
		in.DstGeoRegion = ""
	}
	if !f.wl["src_geo_city"] {
		in.SrcGeoCity = ""
	}
	if !f.wl["dst_geo_city"] {
		in.DstGeoCity = ""
	}
	if !f.wl["dst_nexthop"] {
		in.DstNextHop = ""
	}
	if !f.wl["src_nexthop"] {
		in.SrcNextHop = ""
	}
	if !f.wl["src_route_prefix"] {
		in.SrcRoutePrefix = ""
	}
	if !f.wl["dst_route_prefix"] {
		in.DstRoutePrefix = ""
	}
	if !f.wl["src_second_asn"] {
		in.SrcSecondAsn = 0
	}
	if !f.wl["dst_second_asn"] {
		in.DstSecondAsn = 0
	}
	if !f.wl["src_third_asn"] {
		in.SrcThirdAsn = 0
	}
	if !f.wl["dst_third_asn"] {
		in.DstThirdAsn = 0
	}
	if !f.wl["src_eth_mac"] {
		in.SrcEthMac = ""
	}
	if !f.wl["dst_eth_mac"] {
		in.DstEthMac = ""
	}
	if !f.wl["input_int_desc"] {
		in.InputIntDesc = ""
	}
	if !f.wl["output_int_desc"] {
		in.OutputIntDesc = ""
	}
	if !f.wl["input_int_alias"] {
		in.InputIntAlias = ""
	}
	if !f.wl["output_int_alias"] {
		in.OutputIntAlias = ""
	}
	if !f.wl["input_int_capacity"] {
		in.InputInterfaceCapacity = 0
	}
	if !f.wl["output_int_capacity"] {
		in.OutputInterfaceCapacity = 0
	}
	if !f.wl["input_int_ip"] {
		in.InputInterfaceIP = ""
	}
	if !f.wl["output_int_ip"] {
		in.OutputInterfaceIP = ""
	}
	return true
}

func (f *WLFilter) FilterMap(mapr map[string]interface{}) bool {
	return true
}

func (f *WLFilter) GetName() string {
	return ""
}
