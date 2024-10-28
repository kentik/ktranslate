package kt

import (
	"regexp"
	"strconv"
	"syscall"
	"time"
)

type Compression string

const (
	CompressionNone    Compression = "none"
	CompressionGzip                = "gzip"
	CompressionSnappy              = "snappy"
	CompressionDeflate             = "deflate"
	CompressionNull                = "null"

	KENTIK_EVENT_TYPE       = "KFlow"
	KENTIK_EVENT_SNMP       = "KSnmp"
	KENTIK_EVENT_TRACE      = "KTrace"
	KENTIK_EVENT_SYNTH      = "KSynth"
	KENTIK_EVENT_SYNTH_GEST = "KSynthgest"

	KENTIK_EVENT_SNMP_DEV_METRIC = "KSnmpDeviceMetric"
	KENTIK_EVENT_SNMP_INT_METRIC = "KSnmpInterfaceMetric"
	KENTIK_EVENT_SNMP_METADATA   = "KSnmpInterfaceMetadata"
	KENTIK_EVENT_SNMP_TRAP       = "KSnmpTrap"
	KENTIK_EVENT_EXT             = "KExtEvent"
	KENTIK_EVENT_JSON            = "KJson"
	KENTIK_EVENT_KTRANS_METRIC   = "KTranslateMetric"

	KentikAPITimeout = "KENTIK_API_TIMEOUT"
	IndexVar         = "Index"
	StringPrefix     = "ks_"
	PrivateIP        = "Private IP"
	DropMetric       = "DropMetric"
	AdminStatus      = "if_AdminStatus"
)

type OutputType string

const (
	EventOutput  OutputType = "event"
	MetricOutput            = "metric"
	RollupOutput            = "rollup"
)

type Provider string

const (
	ProviderRouter             Provider = "kentik-router"
	ProviderDefault            Provider = "kentik-default"
	ProviderVPC                Provider = "kentik-vpc"
	ProviderSynth              Provider = "kentik-network-synthetic"
	ProviderSwitch             Provider = "kentik-switch"
	ProviderFirewall           Provider = "kentik-firewall"
	ProviderUPS                Provider = "kentik-ups"
	ProviderPDU                Provider = "kentik-pdu"
	ProviderIOT                Provider = "kentik-iot"
	ProviderHost               Provider = "kentik-host"
	ProviderAlert              Provider = "kentik-alert"
	ProviderFlowDevice         Provider = "kentik-flow-device"
	ProviderWirelessController Provider = "kentik-wireless-controller"
	ProviderNas                Provider = "kentik-nas"
	ProviderSan                Provider = "kentik-san"
	ProviderAgent              Provider = "kentik-agent"
	ProviderFibreChannel       Provider = "kentik-fibre-channel"
	ProviderTrapUnknown        Provider = "kentik-trap-device"
	ProviderHttpDevice         Provider = "kentik-http"
	ProviderMerakiCloud        Provider = "meraki-cloud-controller"
	ProviderKflow              Provider = "kentik-kflow"
)

const (
	InstProvider  = "kentik"
	CollectorName = "ktranslate"
	SnmpCollector = "snmp"

	SendBatchDuration     = 1 * time.Second
	DefaultProfileMessage = "Missing matched profile. See overview page for details."
	SIGUSR2               = syscall.Signal(0xc) // Because windows doesn't have this.

	PluginSyslog   = "ktranslate-syslog"
	PluginHealth   = "ktranslate-health"
	FloatMS        = "float_ms"
	InvalidEnum    = "invalid enum"
	DeviceTagTable = "KT_Device_Tag"
	FromLambda     = "awslambda"
	GaugeMetric    = "gauge"
	CountMetric    = "count"
	FromGCP        = "gcppubsub"
)

type IntId uint64

type Cid IntId              // company id
func (id Cid) Itoa() string { return strconv.Itoa(int(id)) }

// DeviceID denotes id from mn_device table
type DeviceID IntId

// Itoa returns string formated device id
func (id DeviceID) Itoa() string {
	return strconv.FormatInt(int64(id), 10)
}

// IfaceID denotes interfce id
// note this is not mn_interface.id but snmp_id, {input,output}_port in flow
type IfaceID IntId

// Itoa returns string repr of interface id
func (id IfaceID) Itoa() string { return strconv.Itoa(int(id)) }

type JCHF struct {
	Timestamp               int64             `json:"timestamp"`
	DstAs                   uint32            `json:"dst_as"`
	DstGeo                  string            `json:"dst_geo"`
	HeaderLen               uint32            `json:"header_len"`
	InBytes                 uint64            `json:"in_bytes"`
	InPkts                  uint64            `json:"in_pkts"`
	InputPort               IfaceID           `json:"input_port"`
	IpSize                  uint32            `json:"ip_size"`
	DstAddr                 string            `json:"dst_addr"`
	SrcAddr                 string            `json:"src_addr"`
	L4DstPort               uint32            `json:"l4_dst_port"`
	L4SrcPort               uint32            `json:"l4_src_port"`
	OutputPort              IfaceID           `json:"output_port"`
	Protocol                string            `json:"protocol"`
	SampledPacketSize       uint32            `json:"sampled_packet_size"`
	SrcAs                   uint32            `json:"src_as"`
	SrcGeo                  string            `json:"src_geo"`
	TcpFlags                uint32            `json:"tcp_flags"`
	Tos                     uint32            `json:"tos"`
	VlanIn                  uint32            `json:"vlan_in"`
	VlanOut                 uint32            `json:"vlan_out"`
	NextHop                 string            `json:"next_hop"`
	MplsType                uint32            `json:"mpls_type"`
	OutBytes                uint64            `json:"out_bytes"`
	OutPkts                 uint64            `json:"out_pkts"`
	TcpRetransmit           uint32            `json:"tcp_rx"`
	SrcFlowTags             string            `json:"src_flow_tags"`
	DstFlowTags             string            `json:"dst_flow_tags"`
	SampleRate              uint32            `json:"sample_rate"`
	DeviceId                DeviceID          `json:"device_id"`
	DeviceName              string            `json:"device_name"`
	CompanyId               Cid               `json:"company_id"`
	DstBgpAsPath            string            `json:"dst_bgp_as_path"`
	DstBgpCommunity         string            `json:"dst_bgp_comm"`
	SrcBgpAsPath            string            `json:"src_bpg_as_path"`
	SrcBgpCommunity         string            `json:"src_bgp_comm"`
	SrcNextHopAs            uint32            `json:"src_nexthop_as"`
	DstNextHopAs            uint32            `json:"dst_nexthop_as"`
	SrcGeoRegion            string            `json:"src_geo_region"`
	DstGeoRegion            string            `json:"dst_geo_region"`
	SrcGeoCity              string            `json:"src_geo_city"`
	DstGeoCity              string            `json:"dst_geo_city"`
	DstNextHop              string            `json:"dst_nexthop"`
	SrcNextHop              string            `json:"src_nexthop"`
	SrcRoutePrefix          string            `json:"src_route_prefix"`
	DstRoutePrefix          string            `json:"dst_route_prefix"`
	SrcSecondAsn            uint32            `json:"src_second_asn"`
	DstSecondAsn            uint32            `json:"dst_second_asn"`
	SrcThirdAsn             uint32            `json:"src_third_asn"`
	DstThirdAsn             uint32            `json:"dst_third_asn"`
	SrcEthMac               string            `json:"src_eth_mac"`
	DstEthMac               string            `json:"dst_eth_mac"`
	InputIntDesc            string            `json:"input_int_desc"`
	OutputIntDesc           string            `json:"output_int_desc"`
	InputIntAlias           string            `json:"input_int_alias"`
	OutputIntAlias          string            `json:"output_int_alias"`
	InputInterfaceCapacity  int64             `json:"input_int_capacity"`
	OutputInterfaceCapacity int64             `json:"output_int_capacity"`
	InputInterfaceIP        string            `json:"input_int_ip"`
	OutputInterfaceIP       string            `json:"output_int_ip"`
	CustomStr               map[string]string `json:"custom_str,omitempty"`
	CustomInt               map[string]int32  `json:"custom_int,omitempty"`
	CustomBigInt            map[string]int64  `json:"custom_bigint,omitempty"`
	EventType               string            `json:"eventType"`
	Provider                Provider          `json:"provider"` // Entity type for this data.
	avroSet                 map[string]interface{}
	hasSetAvro              bool
	CustomMetrics           map[string]MetricInfo          `json:"-"`
	CustomTables            map[string]DeviceTableMetadata `json:"-"`
	MatchAttr               map[string]*regexp.Regexp      `json:"-"`
	ApplySample             bool                           `json:"-"`        // Should this value be subject to sampling?
	Har                     *HarFile                       `json:"har_file"` // Let you attatch a har file to this object if needed.
}

type MetricInfo struct {
	Oid     string          `json:"-"`
	Mib     string          `json:"-"`
	Name    string          `json:"-"`
	Profile string          `json:"-"`
	Table   string          `json:"-"`
	Format  string          `json:"-"`
	Type    string          `json:"-"`
	PollDur time.Duration   `json:"-"`
	Tables  map[string]bool `json:"-"`
}

func (m *MetricInfo) GetType() string {
	if m.Type != "" {
		return m.Type
	}
	return "snmp"
}

func NewJCHF() *JCHF {
	return &JCHF{avroSet: map[string]interface{}{}, hasSetAvro: false}
}

func (j *JCHF) Reset() {
	j.hasSetAvro = false
}

func (j *JCHF) Flatten() map[string]interface{} {
	mapr := j.ToMap()
	for k, v := range mapr {
		switch mv := v.(type) {
		case map[string]string:
			for ki, vi := range mv {
				mapr[ki] = vi
			}
			delete(mapr, k)
		case map[string]int32:
			for ki, vi := range mv {
				mapr[ki] = vi
			}
			delete(mapr, k)
		case map[string]int64:
			for ki, vi := range mv {
				mapr[ki] = vi
			}
			delete(mapr, k)
		case string:
			if mv == "<nil>" {
				mapr[k] = ""
			}
		default:
			// noop
		}
	}
	return mapr
}

func (j *JCHF) ToMap() map[string]interface{} {
	if j.hasSetAvro { //cache this also.
		return j.avroSet
	}

	j.avroSet["timestamp"] = int64(j.Timestamp)
	j.avroSet["dst_as"] = int64(j.DstAs)
	j.avroSet["dst_geo"] = j.DstGeo
	j.avroSet["header_len"] = int64(j.HeaderLen)
	j.avroSet["in_bytes"] = int64(j.InBytes)
	j.avroSet["in_pkts"] = int64(j.InPkts)
	j.avroSet["input_port"] = int64(j.InputPort)
	j.avroSet["dst_addr"] = j.DstAddr
	j.avroSet["src_addr"] = j.SrcAddr
	j.avroSet["l4_dst_port"] = int64(j.L4DstPort)
	j.avroSet["l4_src_port"] = int64(j.L4SrcPort)
	j.avroSet["output_port"] = int64(j.OutputPort)
	j.avroSet["protocol"] = j.Protocol
	j.avroSet["sampled_packet_size"] = int64(j.SampledPacketSize)
	j.avroSet["src_as"] = int64(j.SrcAs)
	j.avroSet["src_geo"] = j.SrcGeo
	j.avroSet["tcp_flags"] = int64(j.TcpFlags)
	j.avroSet["tos"] = int64(j.Tos)
	j.avroSet["vlan_in"] = int64(j.VlanIn)
	j.avroSet["vlan_out"] = int64(j.VlanOut)
	j.avroSet["out_bytes"] = int64(j.OutBytes)
	j.avroSet["out_pkts"] = int64(j.OutPkts)
	j.avroSet["tcp_rx"] = int64(j.TcpRetransmit)
	j.avroSet["src_flow_tags"] = j.SrcFlowTags
	j.avroSet["dst_flow_tags"] = j.DstFlowTags
	j.avroSet["sample_rate"] = int64(j.SampleRate)
	j.avroSet["device_id"] = int64(j.DeviceId)
	j.avroSet["device_name"] = j.DeviceName
	j.avroSet["company_id"] = int64(j.CompanyId)
	j.avroSet["dst_bgp_as_path"] = j.DstBgpAsPath
	j.avroSet["dst_bgp_comm"] = j.DstBgpCommunity
	j.avroSet["src_bpg_as_path"] = j.SrcBgpAsPath
	j.avroSet["src_bgp_comm"] = j.SrcBgpCommunity
	j.avroSet["src_nexthop_as"] = int64(j.SrcNextHopAs)
	j.avroSet["dst_nexthop_as"] = int64(j.DstNextHopAs)
	j.avroSet["src_geo_region"] = j.SrcGeoRegion
	j.avroSet["dst_geo_region"] = j.DstGeoRegion
	j.avroSet["src_geo_city"] = j.SrcGeoCity
	j.avroSet["dst_geo_city"] = j.DstGeoCity
	j.avroSet["dst_nexthop"] = j.DstNextHop
	j.avroSet["src_nexthop"] = j.SrcNextHop
	j.avroSet["src_route_prefix"] = j.SrcRoutePrefix
	j.avroSet["dst_route_prefix"] = j.DstRoutePrefix
	j.avroSet["src_second_asn"] = int64(j.SrcSecondAsn)
	j.avroSet["dst_second_asn"] = int64(j.DstSecondAsn)
	j.avroSet["src_third_asn"] = int64(j.SrcThirdAsn)
	j.avroSet["dst_third_asn"] = int64(j.DstThirdAsn)
	j.avroSet["src_eth_mac"] = j.SrcEthMac
	j.avroSet["dst_eth_mac"] = j.DstEthMac
	j.avroSet["input_int_desc"] = j.InputIntDesc
	j.avroSet["output_int_desc"] = j.OutputIntDesc
	j.avroSet["input_int_alias"] = j.InputIntAlias
	j.avroSet["output_int_alias"] = j.OutputIntAlias
	j.avroSet["input_interface_capacity"] = j.InputInterfaceCapacity
	j.avroSet["output_interface_capacity"] = j.OutputInterfaceCapacity
	j.avroSet["input_interface_ip"] = j.InputInterfaceIP
	j.avroSet["output_interface_ip"] = j.OutputInterfaceIP
	j.avroSet["custom_str"] = j.CustomStr
	j.avroSet["custom_int"] = j.CustomInt
	j.avroSet["custom_bigint"] = j.CustomBigInt
	j.avroSet["eventType"] = j.EventType
	j.avroSet["provider"] = j.Provider
	j.hasSetAvro = true
	return j.avroSet
}

func (j *JCHF) SetMap() {
	j.avroSet = map[string]interface{}{}
	j.CustomStr = map[string]string{}
	j.CustomInt = map[string]int32{}
	j.CustomBigInt = map[string]int64{}
	j.CustomMetrics = map[string]MetricInfo{}
	j.CustomTables = map[string]DeviceTableMetadata{}
	j.MatchAttr = map[string]*regexp.Regexp{}
}

func (j *JCHF) SetIFPorts(p IfaceID) *JCHF {
	j.OutputPort = p
	j.InputPort = p
	return j
}

type AgentId IntId

func NewAgentId(id string) AgentId {
	aid, _ := strconv.Atoi(id)
	return AgentId(aid)
}

type TestId IntId

func NewTestId(id string) TestId {
	tid, _ := strconv.Atoi(id)
	return TestId(tid)
}

type OutputContext struct {
	Provider  Provider
	Type      OutputType
	CompanyId Cid
	SenderId  string
}

type Output struct {
	Body     []byte
	Ctx      OutputContext
	CB       func(error) // Called when this is sent to a sink.
	NoBuffer bool        // If true, write out right away.
}

func NewOutput(body []byte) *Output {
	return &Output{Body: body}
}

func NewOutputWithProvider(body []byte, prov Provider, stype OutputType) *Output {
	return &Output{Body: body, Ctx: OutputContext{Provider: prov, Type: stype}}
}

func NewOutputWithProviderAndCompanySender(body []byte, prov Provider, cid Cid, stype OutputType, senderid string) *Output {
	return &Output{Body: body, Ctx: OutputContext{Provider: prov, Type: stype, CompanyId: cid, SenderId: senderid}}
}

func NewOutputNoBuffer(body []byte) *Output {
	return &Output{Body: body, NoBuffer: true}
}

func (o *Output) IsEvent() bool {
	return o.Ctx.Type == EventOutput
}

func (o *Output) IsMetric() bool {
	return o.Ctx.Type == RollupOutput || o.Ctx.Type == MetricOutput
}

func (o *Output) GetDataType() string {
	if o == nil {
		return ""
	}

	switch o.Ctx.Provider {
	case ProviderRouter:
		return "snmp"
	case ProviderVPC:
		return "vpc-flows"
	case ProviderSynth:
		return "synthetics"
	case ProviderSwitch:
		return "snmp"
	case ProviderFirewall:
		return "snmp"
	case ProviderUPS:
		return "snmp"
	case ProviderPDU:
		return "snmp"
	case ProviderIOT:
		return "snmp"
	case ProviderHost:
		return "snmp"
	case ProviderAlert:
		return "snmp"
	case ProviderFlowDevice:
		return "device-flows"
	}

	return "device-flows" // Default to this.
}

func (o *Output) BodyLen() int {
	if o == nil {
		return 0
	}
	return len(o.Body)
}

type HarFile struct {
	HarLog HarLog `json:"log"`
}

type HarLog struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Time     float64  `json:"time"`
}

type Request struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type Response struct {
	Status  int     `json:"status"`
	Content Content `json:"content"`
}

type Content struct {
	MimeType    string `json:"mimeType"`
	Size        int64  `json:"size"`
	Compression int64  `json:"compression"`
}
