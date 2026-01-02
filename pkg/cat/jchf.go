package cat

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/cdn"
	patricia "github.com/kentik/ktranslate/pkg/util/gopatricia/patricia"
	"github.com/kentik/ktranslate/pkg/util/ic"
	model "github.com/kentik/ktranslate/pkg/util/kflow2"
)

var (
	DEFAULT_GEO        = []byte("--")
	DEFAULT_GEO_PACKED = patricia.PackGeo(DEFAULT_GEO)

	defaultProvider = kt.ProviderFlowDevice
)

func (kc *KTranslate) getEventType(dst *kt.JCHF) string {

	// if app_proto is 12, this is snmp and return as such.
	if dst.CustomInt[APP_PROTOCOL_COL] == 12 {
		return kt.KENTIK_EVENT_SNMP
	}

	// Else, if its synth, split out into traceroute and other.
	if dst.CustomInt[APP_PROTOCOL_COL] == 10 {
		if dst.CustomInt["result_type"] == 4 {
			return kt.KENTIK_EVENT_TRACE
		} else {
			return kt.KENTIK_EVENT_SYNTH
		}
	}

	if dst.CustomInt[APP_PROTOCOL_COL] == 18 {
		return kt.KENTIK_EVENT_SYNTH_GEST
	}

	return kt.KENTIK_EVENT_TYPE
}

func (kc *KTranslate) getProviderType(dst *kt.JCHF) kt.Provider {

	udr, ok := dst.CustomStr[UDR_TYPE]
	if !ok { // Return this right away.
		return defaultProvider
	}

	// Or maybe its a host.
	if udr == "kprobe" || udr == "kappa" {
		return defaultProvider
	}

	// Else, if its synth, return this.
	if dst.CustomInt[APP_PROTOCOL_COL] == 10 || dst.CustomInt[APP_PROTOCOL_COL] == 11 {
		return kt.ProviderSynth
	}

	// Cloud subnet here.
	if strings.HasSuffix(udr, "_subnet") {
		return kt.ProviderVPC
	}

	// Default to provider.
	return defaultProvider
}

func (kc *KTranslate) flowToJCHF(ctx context.Context, dst *kt.JCHF, src *Flow, currentTS int64, tagcache map[uint64]string) error {

	dst.CustomStr = make(map[string]string)
	dst.CustomInt = make(map[string]int32)
	dst.CustomBigInt = make(map[string]int64)

	// dst.Timestamp = src.CHF.Timestamp() This is being strage, use current timestamp for now.
	dst.Timestamp = currentTS
	dst.DstAs = src.CHF.DstAs()
	if src.CHF.DstGeo() > 0 {
		dst.DstGeo = fmt.Sprintf("%c%c", src.CHF.DstGeo()>>8, src.CHF.DstGeo()&0xFF)
		if dst.DstGeo[0] == '-' {
			dst.DstGeo = "--"
		}
	}
	dst.HeaderLen = src.CHF.HeaderLen()
	dst.InBytes = src.CHF.InBytes()
	dst.InPkts = src.CHF.InPkts()
	dst.InputPort = kt.IfaceID(src.CHF.InputPort())
	dst.IpSize = src.CHF.IpSize()
	dst.L4DstPort = src.CHF.L4DstPort()
	dst.L4SrcPort = src.CHF.L4SrcPort()
	dst.OutputPort = kt.IfaceID(src.CHF.OutputPort())
	dst.Protocol = ic.PROTO_NAMES[uint16(src.CHF.Protocol())]
	dst.SampledPacketSize = src.CHF.SampledPacketSize()
	dst.SrcAs = src.CHF.SrcAs()
	if src.CHF.SrcGeo() > 0 {
		dst.SrcGeo = fmt.Sprintf("%c%c", src.CHF.SrcGeo()>>8, src.CHF.SrcGeo()&0xFF)
		if dst.SrcGeo[0] == '-' {
			dst.SrcGeo = "--"
		}
	}
	dst.TcpFlags = src.CHF.TcpFlags()
	dst.Tos = src.CHF.Tos()
	dst.VlanIn = src.CHF.VlanIn()
	dst.VlanOut = src.CHF.VlanOut()
	dst.MplsType = src.CHF.MplsType()
	dst.OutBytes = src.CHF.OutBytes()
	dst.OutPkts = src.CHF.OutPkts()
	dst.TcpRetransmit = src.CHF.TcpRetransmit()
	dst.SampleRate = src.CHF.SampleRate() / 100 // Reduce by 100 to get actual rate.
	dst.DeviceId = kt.DeviceID(src.CHF.DeviceId())
	dst.DeviceName = src.DeviceName
	dst.CompanyId = kt.Cid(src.CompanyId)
	dst.SrcNextHopAs = src.CHF.SrcNextHopAs()
	dst.DstNextHopAs = src.CHF.DstNextHopAs()
	dst.SrcSecondAsn = src.CHF.SrcSecondAsn()
	dst.DstSecondAsn = src.CHF.DstSecondAsn()
	dst.SrcThirdAsn = src.CHF.SrcThirdAsn()
	dst.DstThirdAsn = src.CHF.DstThirdAsn()

	// Do we have info about this device?
	custColNames := map[uint32]string{}
	if d := kc.apic.GetDevice(dst.CompanyId, dst.DeviceId); d != nil {
		dst.DeviceName = d.Name
		dst.CustomStr[UDR_TYPE] = d.DeviceSubtype
		if len(d.SendingIps) > 0 {
			dst.CustomStr["SamplerAddress"] = d.SendingIps[0].String()
		}
		dst.CustomStr["device_site"] = d.Site.SiteName
		if i, ok := d.Interfaces[dst.InputPort]; ok {
			dst.InputIntDesc = i.Description
			dst.InputIntAlias = i.Alias
			dst.InputInterfaceCapacity = i.SnmpSpeedMbps
			dst.InputInterfaceIP = i.Address
			dst.CustomStr["input_provider"] = i.Provider
			dst.CustomStr["input_network_boundary"] = i.NetworkBoundary
			dst.CustomStr["input_connectivity_type"] = i.ConnectivityType
		}
		if i, ok := d.Interfaces[dst.OutputPort]; ok {
			dst.OutputIntDesc = i.Description
			dst.OutputIntAlias = i.Alias
			dst.OutputInterfaceCapacity = i.SnmpSpeedMbps
			dst.OutputInterfaceIP = i.Address
			dst.CustomStr["output_provider"] = i.Provider
			dst.CustomStr["output_network_boundary"] = i.NetworkBoundary
			dst.CustomStr["output_connectivity_type"] = i.ConnectivityType
		}
		for _, v := range d.Customs {
			custColNames[v.ID] = v.Name
		}
		if d.FullSite != nil {
			dst.CustomStr["device_site_market"] = d.FullSite.SiteMarket.Name
			dst.CustomStr["device_site_country"] = d.FullSite.PostalAddress.Country
		}
	}

	// Now the strings.
	smac := make([]byte, 8)
	binary.BigEndian.PutUint64(smac, src.CHF.SrcEthMac())
	dst.SrcEthMac = net.HardwareAddr(smac).String()
	binary.BigEndian.PutUint64(smac, src.CHF.DstEthMac())
	dst.DstEthMac = net.HardwareAddr(smac).String()

	if sft, err := src.CHF.SrcFlowTags(); err != nil {
		dst.SrcFlowTags = sft
	}
	if sft, err := src.CHF.DstFlowTags(); err != nil {
		dst.DstFlowTags = sft
	}
	if sft, err := src.CHF.SrcBgpAsPath(); err != nil {
		dst.SrcBgpAsPath = sft
	}
	if sft, err := src.CHF.DstBgpAsPath(); err != nil {
		dst.DstBgpAsPath = sft
	}
	if sft, err := src.CHF.SrcBgpCommunity(); err != nil {
		dst.SrcBgpCommunity = sft
	}
	if sft, err := src.CHF.DstBgpCommunity(); err != nil {
		dst.DstBgpCommunity = sft
	}

	// Now the addresses.
	var addr net.IP

	// start with the basic src and dst.
	if src.CHF.Ipv4DstAddr() > 0 {
		addr = kt.Int2ip(src.CHF.Ipv4DstAddr())
	} else {
		ipr, _ := src.CHF.Ipv6DstAddr()
		addr = net.IP(ipr)
	}
	dst.DstAddr = addr.String()

	// Resolve any hostnames if a resolver is set up.
	if kc.resolver != nil {
		dst.CustomStr["dst_host"] = kc.resolver.Resolve(ctx, dst.DstAddr, false)
	}

	if src.CHF.Ipv4SrcAddr() > 0 {
		addr = kt.Int2ip(src.CHF.Ipv4SrcAddr())
	} else {
		ipr, _ := src.CHF.Ipv6SrcAddr()
		addr = net.IP(ipr)
	}
	dst.SrcAddr = addr.String()

	// These are ipv4 addresses.
	addr = kt.Int2ip(src.CHF.SrcRoutePrefix())
	dst.SrcRoutePrefix = addr.String()
	addr = kt.Int2ip(src.CHF.DstRoutePrefix())
	dst.DstRoutePrefix = addr.String()

	if kc.resolver != nil {
		dst.CustomStr["src_host"] = kc.resolver.Resolve(ctx, dst.SrcAddr, false)
	}

	// next hops
	if src.CHF.Ipv4SrcNextHop() > 0 {
		addr = kt.Int2ip(src.CHF.Ipv4SrcNextHop())
	} else {
		ipr, _ := src.CHF.Ipv6SrcNextHop()
		addr = net.IP(ipr)
	}
	dst.SrcNextHop = addr.String()

	if src.CHF.Ipv4DstNextHop() > 0 {
		addr = kt.Int2ip(src.CHF.Ipv4DstNextHop())
	} else {
		ipr, _ := src.CHF.Ipv6DstNextHop()
		addr = net.IP(ipr)
	}
	dst.DstNextHop = addr.String()

	customs, _ := src.CHF.Custom()
	for i, customsLen := 0, customs.Len(); i < customsLen; i++ {
		cust := customs.At(i)
		val := cust.Value()
		name, ok := kc.mapr.Customs[cust.Id()]

		isInt := false
		if !ok {
			if k, ok := custColNames[cust.Id()]; ok {
				name = k
			} else {
				name = strconv.Itoa(int(cust.Id()))
				isInt = true
			}
		}
		switch val.Which() {
		case model.Custom_value_Which_uint16Val:
			dst.CustomInt[name] = int32(val.Uint16Val())
		case model.Custom_value_Which_uint32Val:
			v := val.Uint32Val()
			switch name {
			case "src_cdn_int", "dst_cdn_int":
				dst.CustomStr[name] = cdn.NameByCDN(v)
			case "trf_origination", "trf_termination", "host_direction":
				dst.CustomStr[name] = ic.NETWORK_CLASS_INT_TO_NAME[v]
			case "src_network_bndry", "dst_network_bndry", "ult_exit_network_bndry":
				dst.CustomStr[name] = ic.NameFromNBInt(int(v))
			case "src_connect_type", "dst_connect_type", "ult_exit_connect_type":
				dst.CustomStr[name] = ic.NameFromCTInt(int(v))
			case "dst_rpki":
				if v > ic.RPKI_MAX_NUM || v == ic.RPKI_INVALID {
					dst.CustomStr["i_dst_rpki_name"] = fmt.Sprintf(ic.RPKI_INVALID_NAME, v)
					dst.CustomStr["i_dst_rpki_min_name"] = ic.RPKI_INVALID_MIN_NAME
				} else {
					dst.CustomStr["i_dst_rpki_name"] = ic.RPKI_INT_TO_NAME[v]
					dst.CustomStr["i_dst_rpki_min_name"] = ic.RPKI_INT_TO_MIN_NAME[v]
				}
			case "ult_exit_device_id":
				dst.CustomInt[name] = int32(v)
				if d := kc.apic.GetDevice(dst.CompanyId, kt.DeviceID(v)); d != nil {
					dst.CustomStr["ult_exit_device"] = d.Name
					dst.CustomStr["ult_exit_site"] = d.Site.SiteName
				}
			default:
				if tk, tv, ok := kc.tagMap.LookupTagValue(dst.CompanyId, v, name); ok {
					dst.CustomStr[tk] = tv
				} else if !isInt {
					dst.CustomInt[name] = int32(v) // We don't know anything more about this so best to leave it as it is.
				}
			}
		case model.Custom_value_Which_uint64Val:
			iv := int64(val.Uint64Val())
			if tk, tv, ok := kc.tagMap.LookupTagValueBig(dst.CompanyId, iv, name); ok {
				dst.CustomStr[tk] = tv
			} else {
				dst.CustomBigInt[name] = iv
			}
		case model.Custom_value_Which_strVal:
			sv, _ := val.StrVal()
			dst.CustomStr[name] = sv
		case model.Custom_value_Which_addrVal:
			sv, _ := val.AddrVal()
			if len(sv) > 1 {
				var addr net.IP
				if sv[0] == 4 && len(sv) == 5 {
					addr = net.IP(sv[1:5])
				} else {
					addr = net.IP(sv[1:])
				}
				dst.CustomStr[name] = addr.String()
			}
		}
	}

	// Finally, update any udr based columns with the correct mapping
	if kc.udrMapr != nil {
		var mapr map[string]*UDR
		if dst.CustomStr[UDR_TYPE] != "" && kc.udrMapr.Subtypes[dst.CustomStr[UDR_TYPE]] != nil {
			mapr = kc.udrMapr.Subtypes[dst.CustomStr[UDR_TYPE]]
		} else if ap, ok := dst.CustomInt[APP_PROTOCOL_COL]; ok {
			if maprr, ok := kc.udrMapr.UDRs[ap]; ok {
				mapr = maprr
			}
		}
		for col, udr := range mapr {
			switch udr.Type {
			case UDR_TYPE_INT:
				if val, ok := dst.CustomInt[col]; ok {
					dst.CustomInt[udr.ColumnName] = val
					delete(dst.CustomInt, col)
				}
			case UDR_TYPE_STRING:
				if val, ok := dst.CustomStr[col]; ok {
					dst.CustomStr[udr.ColumnName] = val
					delete(dst.CustomStr, col)
				}
			case UDR_TYPE_BIGINT:
				if val, ok := dst.CustomBigInt[col]; ok {
					dst.CustomBigInt[udr.ColumnName] = val
					delete(dst.CustomBigInt, col)
				} else { // Sometimes this is mis-labeled.
					if val, ok := dst.CustomInt[col]; ok {
						dst.CustomInt[udr.ColumnName] = val
						delete(dst.CustomInt, col)
					}
				}
			}
			if _, ok := dst.CustomStr[UDR_TYPE]; !ok {
				dst.CustomStr[UDR_TYPE] = udr.ApplicationName
			}
			switch udr.ColumnName { // Fill these in directly if they are set.
			case "result_type":
				dst.CustomStr["result_type_str"] = synResultTypes[dst.CustomInt[udr.ColumnName]]
			case "test_id":
				dst.CustomStr["test_url"] = fmt.Sprintf("https://portal.kentik.com/v4/synthetics/tests/%d/results", dst.CustomBigInt[udr.ColumnName])
				if t := kc.apic.GetTest(kt.TestId(dst.CustomBigInt[udr.ColumnName])); t != nil {
					dst.CustomStr["test_name"] = t.GetName()
					dst.CustomStr["test_type"] = t.GetType()
				} else {
					dst.CustomStr["test_name"] = ""
					dst.CustomStr["test_type"] = ""
					kc.apic.UpdateTests(ctx) // On demand, check to see if this test is new.
				}
			case "agent_id":
				if a := kc.apic.GetAgent(kt.AgentId(dst.CustomBigInt[udr.ColumnName])); a != nil {
					dst.CustomStr["agent_name"] = a.GetAlias()
					if dst.SrcAddr == "" || dst.SrcAddr == "<nil>" { // Try getting via agent info.
						lip := a.GetLocalIp()
						if lip != "" {
							dst.SrcAddr = lip
						} else {
							dst.SrcAddr = a.GetIp()
						}
						if kc.resolver != nil {
							dst.CustomStr["src_host"] = kc.resolver.Resolve(ctx, dst.SrcAddr, false)
						}
					}
					dst.SrcAs = a.GetAsn()
					dst.SrcGeoRegion = a.GetRegion()
					dst.SrcGeoCity = a.GetCity()
					dst.SrcGeo = a.GetCountry()
					dst.CustomStr["src_cloud_region"] = a.GetCloudRegion()
					dst.CustomStr["src_cloud_provider"] = a.GetCloudProvider()
					dst.CustomStr["src_site"] = a.GetSiteName()

					// Try getting dest agent info also, but we have to use IP to look up.
					if dst.DstAddr != "" && dst.DstAddr != "<nil>" {
						if da := kc.apic.GetAgentByIP(dst.DstAddr); da != nil {
							dst.CustomStr["dst_agent_name"] = da.GetAlias()
							dst.CustomStr["dst_agent_id"] = da.GetId()
							dst.DstAs = da.GetAsn()
							dst.DstGeoRegion = da.GetRegion()
							dst.DstGeoCity = da.GetCity()
							dst.DstGeo = da.GetCountry()
							dst.CustomStr["dst_cloud_region"] = da.GetCloudRegion()
							dst.CustomStr["dst_cloud_provider"] = da.GetCloudProvider()
							dst.CustomStr["dst_site"] = da.GetSiteName()
						}
					}
				} else {
					dst.CustomStr["agent_name"] = ""
				}
			}
		}
	}

	// Check if there is ultimate exit data interface and pull this in also.
	if udid, ok := dst.CustomInt["ult_exit_device_id"]; ok {
		if d := kc.apic.GetDevice(dst.CompanyId, kt.DeviceID(udid)); d != nil {
			if ui, ok := dst.CustomInt["ult_exit_port"]; ok {
				if i, ok := d.Interfaces[kt.IfaceID(ui)]; ok {
					dst.CustomStr["ult_exit_port_alias"] = i.Alias
					dst.CustomStr["ult_exit_port_description"] = i.Description
					dst.CustomStr["ult_exit_port_provider"] = i.Provider
				}
			}
		}
	}

	// Do we need to remap any of the custom strings?
	for k, v := range dst.CustomStr {
		switch dst.CustomStr[UDR_TYPE] { // Kick out any cross contaminated tags.
		case "gcp_subnet":
			if strings.HasPrefix(k, "kt_aws") || strings.HasPrefix(k, "kt_az") {
				delete(dst.CustomStr, k)
				continue
			}
		case "aws_subnet":
			if strings.HasPrefix(k, "kt_az") {
				delete(dst.CustomStr, k)
				continue
			}
		case "azure_subnet":
			if strings.HasPrefix(k, "kt_aws") {
				delete(dst.CustomStr, k)
				continue
			}
		case "ibm_subnet":
			if strings.HasPrefix(k, "kt_aws") || strings.HasPrefix(k, "kt_az") {
				delete(dst.CustomStr, k)
				continue
			}
		}

		// Now remap to common fields.
		if n, ok := remapCustomStrings[k]; ok {
			dst.CustomStr[n] = v
			delete(dst.CustomStr, k)
		}
	}

	// Fill in as needed.
	if _, ok := dst.CustomStr["Type"]; !ok {
		dst.CustomStr["Type"] = "kflow"
	}

	// Now add some combo fields.
	dst.CustomStr["src_endpoint"] = dst.SrcAddr + ":" + strconv.Itoa(int(dst.L4SrcPort))
	dst.CustomStr["dst_endpoint"] = dst.DstAddr + ":" + strconv.Itoa(int(dst.L4DstPort))

	// Set the type dynamically here to help out processing.
	dst.EventType = kc.getEventType(dst)
	dst.Provider = kc.getProviderType(dst)

	return nil
}

var (
	synResultTypes = map[int32]string{
		0:  "error",
		1:  "timeout",
		2:  "ping",
		3:  "fetch",
		4:  "trace",
		5:  "knock",
		6:  "query",
		7:  "shake",
		8:  "pageload",
		9:  "transaction",
		10: "dnssec",
	}

	remapCustomStrings = map[string]string{
		"kt_aws_dst_acc_id":       "dest_account",
		"kt_aws_src_acc_id":       "source_account",
		"kt_az_dst_sub_id":        "dest_account",
		"kt_az_src_sub_id":        "source_account",
		"destination_project_id":  "dest_account",
		"source_project_id":       "source_account",
		"account":                 "source_account",
		"kt_aws_src_region":       "source_region",
		"kt_aws_dst_region":       "dest_region",
		"kt_az_src_region":        "source_region",
		"kt_az_dst_region":        "dest_region",
		"src_region":              "source_region",
		"dst_region":              "dest_region",
		"source_region":           "source_region",
		"destination_region":      "dest_region",
		"kt_aws_src_vpc_id":       "source_vpc",
		"kt_aws_dst_vpc_id":       "dest_vpc",
		"kt_az_dst_rsrc_group":    "source_vpc",
		"kt_az_src_rsrc_group":    "dest_vpc",
		"source_subnet_name":      "source_vpc",
		"destination_subnet_name": "dest_vpc",
		"src_vpc":                 "source_vpc",
		"dst_vpc":                 "dest_vpc",
	}
)

// Updates asn and geo if set for any of these inputs.
func (kc *KTranslate) doEnrichments(ctx context.Context, msgs []*kt.JCHF) []*kt.JCHF {
	for _, msg := range msgs {
		sip := net.ParseIP(msg.SrcAddr)
		dip := net.ParseIP(msg.DstAddr)
		setSip := false
		setDip := false

		// Internal ips get special handling
		if kc.rule.IsInternal(sip, msg.SrcAs) {
			msg.SrcGeo = kt.PrivateIP
			msg.CustomStr["src_as_name"] = kt.PrivateIP
			setSip = true
		}
		if kc.rule.IsInternal(dip, msg.DstAs) {
			msg.DstGeo = kt.PrivateIP
			msg.CustomStr["dst_as_name"] = kt.PrivateIP
			setDip = true
		}

		// Fetch our own geo if not already set.
		if kc.geo != nil {
			if sip != nil && !setSip {
				if geo, err := kc.geo.SearchBestFromHostGeo(sip); err == nil {
					msg.SrcGeo = geo
				}
			}
			if dip != nil && !setDip {
				if geo, err := kc.geo.SearchBestFromHostGeo(dip); err == nil {
					msg.DstGeo = geo
				}
			}
		}

		// And set our own asn also if not set.
		if kc.asn != nil {
			if sip != nil && !setSip {
				if asn, name, err := kc.asn.SearchBestFromHostAsn(sip); err == nil {
					msg.SrcAs = asn
					msg.CustomStr["src_as_name"] = name
				}
			}
			if dip != nil && !setDip {
				if asn, name, err := kc.asn.SearchBestFromHostAsn(dip); err == nil {
					msg.DstAs = asn
					msg.CustomStr["dst_as_name"] = name
				}
			}
		}

		// See if we know what service this is based on proto and port.
		if msg.CustomStr["application"] == "" {
			if msg.L4SrcPort > 0 && msg.L4SrcPort < msg.L4DstPort {
				if app, ok := kc.rule.GetService(sip, msg.L4SrcPort, ic.PROTO_NUMS[msg.Protocol]); ok {
					msg.CustomStr["application"] = app
				}
			} else {
				if app, ok := kc.rule.GetService(dip, msg.L4DstPort, ic.PROTO_NUMS[msg.Protocol]); ok {
					msg.CustomStr["application"] = app
				}
			}
		}

		// If there's a resolver, try to resolve to hostnames.
		if kc.resolver != nil {
			msg.CustomStr["src_host"] = kc.resolver.Resolve(ctx, msg.SrcAddr, false)
			msg.CustomStr["dst_host"] = kc.resolver.Resolve(ctx, msg.DstAddr, false)
		}

		// This data is typically json, try to parse and pull out.
		if msg.EventType == kt.KENTIK_EVENT_SYNTH {
			if rawStr, ok := msg.CustomStr["error_cause/trace_route"]; ok {
				if rawStr != "" {
					strData := []interface{}{}
					if err := json.Unmarshal([]byte(rawStr), &strData); err == nil {
						if len(strData) > 0 {
							switch sd := strData[0].(type) {
							case map[string]interface{}:
								for key, val := range sd {
									switch tv := val.(type) {
									case string:
										msg.CustomStr[key] = tv
									case int:
										msg.CustomInt[key] = int32(tv)
									case int64:
										msg.CustomBigInt[key] = tv
									case float64:
										msg.CustomBigInt[key] = int64(tv)
									case map[string]interface{}:
										if hv, ok := tv["har"]; ok {
											switch av := hv.(type) {
											case []interface{}:
												if len(av) > 0 {
													switch iner := av[0].(type) {
													case map[string]interface{}:
														if path, ok := iner["path"]; ok {
															switch pt := path.(type) {
															case string:
																kc.getHar(ctx, pt, msg)
															}
														}
													}
												}
											}
										}
									case nil:
										// Noop here.
									default:
										// And noop here.
									}
								}
								delete(msg.CustomStr, "error_cause/trace_route")
							}
						}
					}
				}
			}
		}
	}

	// If there's an outside enrichment service, send over here.
	if kc.enricher != nil {
		new, err := kc.enricher.Enrich(ctx, msgs)
		if err != nil {
			kc.log.Errorf("Cannot enrich: %v", err)
		} else {
			msgs = new
		}
	}

	return msgs
}

// Pulls in a har file if possible.
func (kc *KTranslate) getHar(ctx context.Context, path string, msg *kt.JCHF) {
	if kc.objmgr != nil {
		data, err := kc.objmgr.Get(ctx, path)
		if err != nil {
			kc.log.Errorf("Cannot get path %s %v", path, err)
			return
		}

		var har kt.HarFile
		if err := json.Unmarshal(data, &har); err == nil {
			msg.Har = &har
		}
	}
}
