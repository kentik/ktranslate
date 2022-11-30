package kt

import (
	sfmt "gopkg.in/mcuadros/go-syslog.v2/format"

	"net"
)

// Devices is a map of device ids to devices for a company.
type Devices map[DeviceID]*Device

// A Device represents a device, corresponding to a row in mn_device.
// It also has all its interfaces attached.
type Device struct {
	ID            DeviceID              `json:"id,string"`
	Name          string                `json:"device_name"`
	CompanyID     Cid                   `json:"company_id,string"`
	DeviceType    string                `json:"device_type"`
	DeviceSubtype string                `json:"device_subtype"`
	Description   string                `json:"device_description"`
	IP            net.IP                `json:"ip"`
	Interfaces    map[IfaceID]Interface `json:"-"`
	AllInterfaces []Interface           `json:"all_interfaces"`
	SendingIps    []net.IP              `json:"sending_ips"`
	SampleRate    uint32                `json:"device_sample_rate,string"`
	BgpType       string                `json:"device_bgp_type"`
	Plan          Plan                  `json:"plan"`
	CdnAttr       string                `json:"cdn_attr"`
	MaxFlowRate   int                   `json:"max_flow_rate"`
	Customs       []Column              `json:"custom_column_data,omitempty"`
	CustomStr     string                `json:"custom_columns"`
	SnmpCommunity string                `json:"device_snmp_community"`
	SnmpIp        string                `json:"device_snmp_ip"`
	SnmpV3        *V3SNMPConfig         `json:"device_snmp_v3_conf"`
	Labels        []DeviceLabel         `json:"labels"`
	Site          DeviceSite            `json:"site"`
	allUserTags   map[string]string
}

type DeviceLabel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

type DeviceSite struct {
	ID       int    `json:"id"`
	SiteName string `json:"site_name"`
}

// A CustomColumn corresponds a row in mn_kflow_field, which represents
// a field in the "Custom" field of a kflow record. This could either be
// a field we add to most kflow (e.g. i_ult_exit_network_bndry_name), or
// a "custom dimension" for a company (e.g. c_customer, c_kentik_services, etc.).
// The CustomMapValues field is filled in with data from mn_flow_tag_kv if
// it exists, but in a post HSCD (hyper-scale custom dimensions aka hippo tagging)
// world, these generally won't be there.
type CustomColumn struct {
	ID              uint32
	Name            string
	Type            string // kt.FORMAT_UINT32 or kt.FORMAT_STRING or kt.FORMAT_ADDR
	CustomMapValues map[uint32]string
}

// InterfaceCapacityBPS denotes capacity of interface in bits per second
type InterfaceCapacityBPS = uint64

// An Interface is everything we know about a device's interfaces.
// It corresponds to a row in mn_interface, joined with information
// from mn_device and mn_site.
type Interface struct {
	ID int64 `json:"id"`

	DeviceID   DeviceID `json:"device_id,string"`
	DeviceName string   `json:"device_name"`
	DeviceType string   `json:"device_type"`
	SiteID     int      `json:"site_id"`

	SnmpID               IfaceID `json:"snmp_id,string"`
	SnmpSpeedMbps        int64   `json:"snmp_speed,string"` // unit? TODO: switch to uint64, rename to SnmpSpeedMbps
	SnmpType             int     `json:"snmp_type"`
	SnmpAlias            string  `json:"snmp_alias"`
	InterfaceIP          string  `json:"interface_ip"`
	InterfaceDescription string  `json:"interface_description"`
	Provider             string  `json:"provider"`
	VrfID                int64   `json:"vrf_id"`
	Netmask              string  `json:"interface_ip_netmask"`
	Addrs                []Addr  `json:"secondary_ips"`

	SiteTitle   string `json:"site_title"`
	SiteCountry string `json:"site_country"`
}

type Addr struct {
	Address string `json:"address"`
	Netmask string `json:"netmask"`
}

type DeviceList struct {
	Devices []Device `json:"devices"`
}

type DeviceMapper struct {
	Devices map[DeviceID]map[IfaceID]*InterfaceRow
}

type InterfaceRow struct {
	DeviceId             uint32 `json:"device_id"`
	DeviceName           string `json:"device_name"`
	DeviceType           string `json:"device_type"`
	SiteId               uint32 `json:"site_id"`
	SnmpId               string `json:"snmp_id"`
	SnmpSpeed            int64  `json:"snmp_speed"`
	SnmpType             uint32 `json:"snmp_type"`
	SnmpAlias            string `json:"snmp_alias"`
	InterfaceIp          string `json:"interface_ip"`
	InterfaceDescription string `json:"interface_description"`
	Provider             string `json:"provider"`
	VrfId                uint32 `json:"vrf_id"`
	SiteTitle            string `json:"site_title"`
	SiteCountry          string `json:"site_country"`
}

type Plan struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Column struct {
	ID   uint64 `json:"field_id,string"`
	Name string `json:"col_name"`
	Type string `json:"col_type"`
}

func (d *Device) InitUserTags(serviceName string, tags map[string]string) {
	d.allUserTags = tags
	if serviceName != "ktranslate" {
		d.allUserTags["tags.container_service"] = serviceName
	}
}

func (d *Device) SetMsgUserTags(in sfmt.LogParts) {
	for k, v := range d.allUserTags {
		in[k] = v
	}
}

func (d *Device) SetUserTags(in map[string]string) {
	for k, v := range d.allUserTags {
		in[k] = v
	}
}
