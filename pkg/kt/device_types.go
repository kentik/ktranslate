package kt

// Devices is a map of device ids to devices for a company.
type Devices map[DeviceID]Device

// A Device represents a device, corresponding to a row in mn_device.
// It also has all its interfaces attached.
type Device struct {
	ID         DeviceID
	Name       string
	Interfaces map[IfaceID]Interface
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

	DeviceID   DeviceID `json:"device_id"`
	DeviceName string   `json:"device_name"`
	DeviceType string   `json:"device_type"`
	SiteID     int      `json:"site_id"`

	SnmpID               string `json:"snmp_id"`
	SnmpSpeedMbps        int64  `json:"snmp_speed"` // unit? TODO: switch to uint64, rename to SnmpSpeedMbps
	SnmpType             int    `json:"snmp_type"`
	SnmpAlias            string `json:"snmp_alias"`
	InterfaceIP          string `json:"interface_ip"`
	InterfaceDescription string `json:"interface_description"`
	Provider             string `json:"provider"`
	VrfID                int64  `json:"vrf_id"`

	SiteTitle   string `json:"site_title"`
	SiteCountry string `json:"site_country"`
}
