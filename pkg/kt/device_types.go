package kt

import (
	devicepb "github.com/kentik/api-schema-public/gen/go/kentik/device/v202504beta2"
	interfacepb "github.com/kentik/api-schema-public/gen/go/kentik/interface/v202108alpha1"
	sfmt "github.com/kentik/the-library-formally-known-as-go-syslog/format"

	"fmt"
	"github.com/kentik/ktranslate/pkg/util/ic"
	"net"
	"strconv"
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
	FullSite      *FullSite
	IDStr         string
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

// An Interface is everything we know about a device's interfaces.
// It corresponds to a row in mn_interface, joined with information
// from mn_device and mn_site.
type Interface struct {
	DeviceID    DeviceID `json:"device_id,string"`
	Address     string   `json:"interface_ip"`
	Netmask     string   `json:"interface_ip_netmask"`
	Description string   `json:"interface_description"`

	NetworkBoundary  string `json:"network_boundary"`
	ConnectivityType string `json:"connectivity_type"`
	Provider         string `json:"provider"`

	SnmpID        IfaceID `json:"snmp_id,string"`
	Alias         string  `json:"snmp_alias"`
	Type          uint64  `json:"snmp_type"`
	SnmpSpeedMbps int64   `json:"snmp_speed,string"` // unit? TODO: switch to uint64, rename to SnmpSpeedMbps
	SnmpType      int     `json:"snmp_type"`

	Addrs     []Addr            `json:"secondary_ips"`
	ExtraInfo map[string]string `json:"extra_info"`
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
	ID          uint32 `json:"field_id,string"`
	Name        string `json:"col_name"`
	Type        string `json:"col_type"`
	DeviceID    string
	FieldID     string
	Description string
	DeviceType  string
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

type SiteList struct {
	Sites []FullSite `json:"sites"`
}

type FullSite struct {
	ID            string        `json:"id": "33467"`
	Title         string        `json:"title"`
	Lat           float64       `json:"lat"`
	Lon           float64       `json:"lon"`
	PostalAddress PostalAddress `json:"postalAddress"`
	Type          string        `json:"type"`
	SiteMarket    SiteMarket    `json:"siteMarket"`
}

type PostalAddress struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	Region     string `json:"region"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type SiteMarket struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

func (d *Device) AddCustoms(dd *devicepb.DeviceDetailed) {
	customs := mapCustomColumns(dd.GetCustomColumnData())
	d.Customs = customs
	d.CustomStr = dd.GetCustomColumns()
}

func (d *Device) AddInterface(p *interfacepb.Interface) {
	if p == nil {
		return
	}

	snmpID, err := strconv.ParseInt(p.GetSnmpId(), 10, 64)
	if err != nil {
		snmpID = 0
	}

	devID, err := strconv.ParseInt(p.GetDeviceId(), 10, 64)
	if err != nil {
		devID = 0
	}

	iface := Interface{
		DeviceID:         DeviceID(devID),
		Description:      p.GetInterfaceDescription(),
		NetworkBoundary:  ic.NameFromNBInt(int(p.GetNetworkBoundary())),
		ConnectivityType: ic.NameFromCTInt(int(p.GetConnectivityType())),
		Provider:         p.GetProvider(),
		SnmpID:           IfaceID(snmpID),
		Alias:            p.GetSnmpAlias(),
		SnmpSpeedMbps:    int64(p.GetSnmpSpeed()),
		Address:          p.GetInterfaceIp(),
	}

	d.AllInterfaces = append(d.AllInterfaces, iface)
	d.Interfaces[IfaceID(snmpID)] = iface
}

func MapDeviceDetailedToDevice(dd *devicepb.DeviceDetailed) (*Device, error) {
	if dd == nil {
		return nil, fmt.Errorf("DeviceDetailed is nil")
	}

	deviceID, err := strconv.ParseInt(dd.GetId(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing device ID %q: %w", dd.GetId(), err)
	}

	companyID, err := strconv.ParseInt(dd.GetCompanyId(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing company ID %q: %w", dd.GetCompanyId(), err)
	}

	sampleRate, err := strconv.ParseUint(dd.GetDeviceSampleRate(), 10, 32)
	if err != nil {
		// Default to 0 if unparseable rather than failing hard
		sampleRate = 0
	}

	ip := net.ParseIP(dd.GetDeviceSnmpIp())

	sendingIPs := make([]net.IP, 0, len(dd.GetSendingIps()))
	for _, s := range dd.GetSendingIps() {
		if parsed := net.ParseIP(s); parsed != nil {
			sendingIPs = append(sendingIPs, parsed)
		}
	}

	ifaces, ifaceMap := mapInterfaces(dd.GetId(), dd.GetAllInterfaces())

	labels := make([]DeviceLabel, 0, len(dd.GetLabels()))
	for _, l := range dd.GetLabels() {
		labelID, err := strconv.Atoi(l.GetId())
		if err != nil {
			labelID = 0
		}
		labels = append(labels, DeviceLabel{
			ID:   labelID,
			Name: l.GetName(),
			Desc: l.GetDescription(),
		})
	}

	site := DeviceSite{}
	if s := dd.GetSite(); s != nil {
		siteID, err := strconv.Atoi(s.GetId())
		if err != nil {
			siteID = 0
		}
		site = DeviceSite{
			ID:       siteID,
			SiteName: s.GetSiteName(),
		}
	}

	plan := mapPlan(dd.GetPlan())

	customs := mapCustomColumns(dd.GetCustomColumnData())

	var snmpV3 *V3SNMPConfig
	if v3 := dd.GetDeviceSnmpV3Conf(); v3 != nil {
		snmpV3 = &V3SNMPConfig{
			UserName:                 v3.GetUsername(),
			AuthenticationProtocol:   v3.GetAuthenticationProtocol(),
			AuthenticationPassphrase: v3.GetAuthenticationPassphrase(),
			PrivacyProtocol:          v3.GetPrivacyProtocol(),
			PrivacyPassphrase:        v3.GetPrivacyPassphrase(),
		}
	}

	return &Device{
		ID:            DeviceID(deviceID),
		IDStr:         dd.GetId(),
		Name:          dd.GetDeviceName(),
		CompanyID:     Cid(companyID),
		DeviceType:    dd.GetDeviceType(),
		DeviceSubtype: dd.GetDeviceSubtype(),
		Description:   dd.GetDeviceDescription(),
		IP:            ip,
		Interfaces:    ifaceMap,
		AllInterfaces: ifaces,
		SendingIps:    sendingIPs,
		SampleRate:    uint32(sampleRate),
		BgpType:       dd.GetDeviceBgpType(),
		Plan:          plan,
		CdnAttr:       dd.GetCdnAttr(),
		MaxFlowRate:   int(dd.GetMaxFlowRate()),
		Customs:       customs,
		CustomStr:     dd.GetCustomColumns(),
		SnmpCommunity: dd.GetDeviceSnmpCommunity(),
		SnmpIp:        dd.GetDeviceSnmpIp(),
		SnmpV3:        snmpV3,
		Labels:        labels,
		Site:          site,
	}, nil
}

func mapInterfaces(deviceID string, protos []*devicepb.Interface) ([]Interface, map[IfaceID]Interface) {
	ifaces := make([]Interface, 0, len(protos))
	ifaceMap := make(map[IfaceID]Interface, len(protos))

	for _, p := range protos {
		if p == nil {
			continue
		}

		snmpID, err := strconv.ParseInt(p.GetSnmpId(), 10, 64)
		if err != nil {
			snmpID = 0
		}

		devID, err := strconv.ParseInt(deviceID, 10, 64)
		if err != nil {
			devID = 0
		}

		snmpSpeed, err := strconv.ParseInt(p.GetSnmpSpeed(), 10, 64)
		if err != nil {
			snmpSpeed = 0
		}

		iface := Interface{
			DeviceID:         DeviceID(devID),
			Description:      p.GetInterfaceDescription(),
			NetworkBoundary:  p.GetNetworkBoundary(),
			ConnectivityType: p.GetConnectivityType(),
			Provider:         p.GetProvider(),
			SnmpID:           IfaceID(snmpID),
			Alias:            p.GetSnmpAlias(),
			SnmpSpeedMbps:    snmpSpeed,
		}

		ifaces = append(ifaces, iface)
		ifaceMap[IfaceID(snmpID)] = iface
	}

	return ifaces, ifaceMap
}

func mapPlan(p *devicepb.Plan) Plan {
	if p == nil {
		return Plan{}
	}

	planID, err := strconv.Atoi(p.GetId())
	if err != nil {
		planID = 0
	}

	return Plan{
		ID:   uint64(planID),
		Name: p.GetName(),
	}
}

func mapCustomColumns(cols []*devicepb.CustomColumnData) []Column {
	if len(cols) == 0 {
		return nil
	}

	result := make([]Column, 0, len(cols))
	for _, c := range cols {
		if c == nil {
			continue
		}
		fieldIDUint, err := strconv.ParseUint(c.GetFieldId(), 10, 32)
		if err != nil {
			fieldIDUint = 0
		}
		result = append(result, Column{
			ID:          uint32(fieldIDUint),
			DeviceID:    c.GetDeviceId(),
			FieldID:     c.GetFieldId(),
			Name:        c.GetColName(),
			Description: c.GetDescription(),
			Type:        c.GetColType(),
			DeviceType:  c.GetDeviceType(),
		})
	}
	return result
}
