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
	Manufacturer          string                    `yaml:"manufacturer"`
	InterfaceData         map[string]*InterfaceData `yaml:"interface_data"`
	DeviceMetricsMetadata *DeviceMetricsMetadata    `yaml:"device_metrics_metadata,omitempty"`
}

type DeviceMetricsMetadata struct {
	SysName     string `yaml:"sys_name,omitempty"`
	SysObjectID string `yaml:"sys_object_id,omitempty"`
	SysDescr    string `yaml:"sys_descr,omitempty"`
	SysLocation string `yaml:"sys_location,omitempty"`
	SysContact  string `yaml:"sys_contact,omitempty"`
	SysServices int    `yaml:"sys_services,omitempty"`
	Customs     map[string]string
	CustomInts  map[string]int64
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
	AliasAddr   []IPAddr          `yaml:"alias_address"`
	Index       string            `yaml:"index"`
	Speed       uint64            `yaml:"speed"`
	Description string            `yaml:"desc"`
	Alias       string            `yaml:"alias"`
	Type        uint64            `yaml:"type"`
	VrfName     string            `yaml:"vrf_name"`
	VrfDescr    string            `yaml:"vrf_desc"`
	VrfRD       string            `yaml:"vrf_rd"`
	VrfExtRD    uint64            `yaml:"vrf_ext_rd"`
	VrfRT       string            `yaml:"vrf_rt"`
	ExtraInfo   map[string]string `yaml:"extra_info"`
}

type IPAddr struct {
	Address string `yaml:"address"`
	Netmask string `yaml:"netmask"`
}

type TopValue struct {
	Asn     uint32
	Packets uint64
}

type V3SNMPConfig struct {
	UserName                 string `yaml:"user_name"`
	AuthenticationProtocol   string `yaml:"authentication_protocol"`
	AuthenticationPassphrase string `yaml:"authentication_passphrase"`
	PrivacyProtocol          string `yaml:"privacy_protocol"`
	PrivacyPassphrase        string `yaml:"privacy_passphrase"`
	ContextEngineID          string `yaml:"context_engine_id"`
	ContextName              string `yaml:"context_name"`
}

type SnmpDeviceConfig struct {
	DeviceName             string            `yaml:"device_name"`
	DeviceIP               string            `yaml:"device_ip"`
	Community              string            `yaml:"snmp_comm"`
	UseV1                  bool              `yaml:"use_snmp_v1"`
	V3                     *V3SNMPConfig     `yaml:"snmp_v3"`
	Debug                  bool              `yaml:"debug"`
	RateMultiplier         int64             `yaml:"rate_multiplier"`
	SampleRate             int64             `yaml:"sample_rate"`
	Port                   uint16            `yaml:"port"`
	OID                    string            `yaml:"oid"`
	Description            string            `yaml:"description"`
	Checked                time.Time         `yaml:"last_checked"`
	InterfaceMetricsOidMap map[string]string `yaml:"interface_metrics_oids"`
	DeviceOids             map[string]*Mib   `yaml:"device_oids"`
	MibProfile             string            `yaml:"mib_profile"`
	Provider               Provider          `yaml:"provider"`
	FlowOnly               bool              `yaml:"flow_only"`
}

type SnmpTrapConfig struct {
	Listen    string `yaml:"listen"`
	Community string `yaml:"community"`
	Version   string `yaml:"version"`
	Transport string `yaml:"transport"`
}

type SnmpDiscoConfig struct {
	Cidrs              StringArray   `yaml:"cidrs"`
	Debug              bool          `yaml:"debug"`
	Ports              []int         `yaml:"ports"`
	DefaultCommunities []string      `yaml:"default_communities"`
	UseV1              bool          `yaml:"use_snmp_v1"`
	DefaultV3          *V3SNMPConfig `yaml:"default_v3"`
	AddDevices         bool          `yaml:"add_devices"`
	Threads            int           `yaml:"threads"`
	CheckAll           bool          `yaml:"check_all"`
	ReplaceDevices     bool          `yaml:"replace_devices"`
	AddFromMibDB       bool          `yaml:"add_from_mibdb"`
	CidrOrig           string        `yaml:"-"`
}

type SnmpGlobalConfig struct {
	PollTimeSec     int      `yaml:"poll_time_sec"`
	DropIfOutside   bool     `yaml:"drop_if_outside_poll"`
	MibProfileDir   string   `yaml:"mib_profile_dir"`
	PyMibProfileDir string   `yaml:"pymib_profile_dir"`
	MibDB           string   `yaml:"mibs_db"`
	MibsEnabled     []string `yaml:"mibs_enabled"`
	TimeoutMS       int      `yaml:"timeout_ms"`
	Retries         int      `yaml:"retries"`
}

type SnmpConfig struct {
	Devices    DeviceMap         `yaml:"devices"`
	Trap       *SnmpTrapConfig   `yaml:"trap"`
	Disco      *SnmpDiscoConfig  `yaml:"discovery"`
	Global     *SnmpGlobalConfig `yaml:"global"`
	DeviceOrig string            `yaml:"-"`
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

type LastMetadata struct {
	DeviceInfo    map[string]interface{}
	InterfaceInfo map[IfaceID]map[string]interface{}
}

type DeviceMap map[string]*SnmpDeviceConfig

func (a *DeviceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi map[string]*SnmpDeviceConfig
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = map[string]*SnmpDeviceConfig{"file": &SnmpDeviceConfig{DeviceName: single}}
	} else {
		*a = multi
	}
	return nil
}

type StringArray []string

func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}
