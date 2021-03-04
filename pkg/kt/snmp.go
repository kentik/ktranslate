package kt

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	go_metrics "github.com/kentik/go-metrics"
)

// DeviceData holds information about a device, sent via ST, and sent
// to the portal in writeSNMPWrapper.
type DeviceData struct {
	Manufacturer          string                    `json:"manufacturer"`
	InterfaceData         map[string]*InterfaceData `json:"interface_data"`
	DeviceMetricsMetadata *DeviceMetricsMetadata    `json:"device_metrics_metadata,omitempty"`
}

type DeviceMetricsMetadata struct {
	SysName     string `json:"sys_name,omitempty"`
	SysObjectID string `json:"sys_object_id,omitempty"`
	SysDescr    string `json:"sys_descr,omitempty"`
	SysLocation string `json:"sys_location,omitempty"`
	SysContact  string `json:"sys_contact,omitempty"`
	SysServices int    `json:"sys_services,omitempty"`
}

const (
	LOGICAL_INTERFACE  = "logical"
	PHYSICAL_INTERFACE = "physical"
)

var (
	baseCharOnlyRegexp = regexp.MustCompile("[^a-zA-Z0-9]+")
)

// InterfaceData is the metadata we've discovered about an interface during SNMP
// polling, or from Streaming Telemetry.  It's sent *to* the API server as part of
// a DeviceData.
type InterfaceData struct {
	IPAddr                        // embedded
	AliasAddr   []IPAddr          `json:"alias_address"`
	Index       string            `json:"index"`
	Speed       uint64            `json:"speed"`
	Description string            `json:"desc"`
	Alias       string            `json:"alias"`
	Type        uint64            `json:"type"`
	VrfName     string            `json:"vrf_name"`
	VrfDescr    string            `json:"vrf_desc"`
	VrfRD       string            `json:"vrf_rd"`
	VrfExtRD    uint64            `json:"vrf_ext_rd"`
	VrfRT       string            `json:"vrf_rt"`
	ExtraInfo   map[string]string `json:"extra_info"`
}

type IPAddr struct {
	Address string `json:"address"`
	Netmask string `json:"netmask"`
}

type TopValue struct {
	Asn     uint32
	Packets uint64
}

type V3SNMPConfig struct {
	UserName                 string `json:"UserName"`
	AuthenticationProtocol   string `json:"AuthenticationProtocol"`
	AuthenticationPassphrase string `json:"AuthenticationPassphrase"`
	PrivacyProtocol          string `json:"PrivacyProtocol"`
	PrivacyPassphrase        string `json:"PrivacyPassphrase"`
	ContextEngineID          string `json:"ContextEngineID"`
	ContextName              string `json:"ContextName"`
}

type SnmpDeviceConfig struct {
	DeviceName             string            `json:"device_name"`
	DeviceIP               string            `json:"device_ip"`
	Community              string            `json:"snmp_comm"`
	V3                     *V3SNMPConfig     `json:"snmp_v3"`
	Debug                  bool              `json:"debug"`
	RateMultiplier         int64             `json:"rate_multiplier"`
	Port                   uint16            `json:"port"`
	OID                    string            `json:"oid"`
	Description            string            `json:"description"`
	Checked                time.Time         `json:"last_checked"`
	InterfaceMetricsOidMap map[string]string `json:"interface_metrics_oids"`
	DeviceOids             map[string]*Mib   `json:"device_oids"`
	MibProfile             string            `json:"mib_profile"`
}

type SnmpTrapConfig struct {
	Listen    string `json:"listen"`
	Community string `json:"community"`
	Version   string `json:"version"`
	Transport string `json:"transport"`
}

type SnmpDiscoConfig struct {
	Cidrs              []string      `json:"cidrs"`
	Debug              bool          `json:"debug"`
	TimeoutMS          int           `json:"timeout_ms"`
	Ports              []int         `json:"ports"`
	DefaultCommunities []string      `json:"default_communities"`
	DefaultV3          *V3SNMPConfig `json:"default_v3"`
	AddDevices         bool          `json:"add_devices"`
	Retries            int           `json:"retries"`
	Threads            int           `json:"threads"`
	MibDB              string        `json:"mibs_db"`
	MibProfileDir      string        `json:"mib_profile_dir"`
	CheckAll           bool          `json:"check_all"`
	ReplaceDevices     bool          `json:"replace_devices"`
}

type SnmpGlobalConfig struct {
	PollTimeSec   int  `json:"poll_time_sec"`
	DropIfOutside bool `json:"drop_if_outside_poll"`
}

type SnmpConfig struct {
	Devices map[string]*SnmpDeviceConfig `json:"devices"`
	Trap    *SnmpTrapConfig              `json:"trap"`
	Disco   *SnmpDiscoConfig             `json:"discovery"`
	Global  *SnmpGlobalConfig            `json:"global"`
}

type SnmpMetricSet struct {
	Devices map[string]*SnmpDeviceMetric
	Mux     sync.RWMutex
	Traps   go_metrics.Meter
}

func NewSnmpMetricSet(registry go_metrics.Registry) *SnmpMetricSet {
	return &SnmpMetricSet{
		Devices: map[string]*SnmpDeviceMetric{},
		Traps:   go_metrics.GetOrRegisterMeter("snmp_traps", registry),
	}
}

type SnmpDeviceMetric struct {
	DeviceMetrics    go_metrics.Meter
	InterfaceMetrics go_metrics.Meter
	Metadata         go_metrics.Meter
	Errors           go_metrics.Meter
}

func NewSnmpDeviceMetric(registry go_metrics.Registry, deviceName string) *SnmpDeviceMetric {
	sm := SnmpDeviceMetric{
		DeviceMetrics:    go_metrics.GetOrRegisterMeter("device_metrics^device_name="+deviceName, registry),
		InterfaceMetrics: go_metrics.GetOrRegisterMeter("interface_metrics^device_name="+deviceName, registry),
		Metadata:         go_metrics.GetOrRegisterMeter("metadata^device_name="+deviceName, registry),
		Errors:           go_metrics.GetOrRegisterMeter("snmp_errors^device_name="+deviceName, registry),
	}
	return &sm
}

type Oidtype int

const (
	ObjID     Oidtype = 1
	String            = 2
	INTEGER           = 3
	NetAddr           = 4
	IpAddr            = 5
	Counter           = 6
	Gauge             = 7
	TimeTicks         = 8
	Counter64         = 11
	BitString         = 12
	Index             = 15
	Integer32         = 16
)

type Mib struct {
	Oid   string
	Name  string
	Type  Oidtype
	Extra string
}

func (mb Mib) String() string {
	return fmt.Sprintf("Name: %s, Oid: %s: Type: %d, Extra: %s", mb.Name, mb.Oid, mb.Type, mb.Extra)
}
