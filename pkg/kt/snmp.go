package kt

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	go_metrics "github.com/kentik/go-metrics"
	"gopkg.in/yaml.v2"
)

const (
	UserTagPrefix = "tags."
)

// DeviceData holds information about a device, sent via ST, and sent
// to the portal in writeSNMPWrapper.
type DeviceData struct {
	Manufacturer          string                    `yaml:"manufacturer"`
	InterfaceData         map[string]*InterfaceData `yaml:"interface_data"`
	DeviceMetricsMetadata *DeviceMetricsMetadata    `yaml:"device_metrics_metadata,omitempty"`
}

type DeviceTableMetadata struct {
	Customs map[string]MetaValue
}

type MetaValue struct {
	TableName  string
	TableNames map[string]bool
	StringVal  string
	IntVal     int64
}

func (mv *MetaValue) GetValue() interface{} {
	if mv.IntVal != 0 {
		return mv.IntVal
	}
	return mv.StringVal
}

func NewDeviceTableMetadata() DeviceTableMetadata {
	return DeviceTableMetadata{
		Customs: map[string]MetaValue{},
	}
}

func NewMetaValue(mib *Mib, sv string, iv int64) MetaValue {
	if mib != nil && mib.EnumRev != nil {
		if nv, ok := mib.EnumRev[iv]; ok {
			sv = nv
		} else {
			sv = InvalidEnum
		}
		iv = 0
	}

	mv := MetaValue{
		TableName:  mib.Table,
		TableNames: map[string]bool{mib.Table: true},
		StringVal:  strings.TrimSpace(sv),
		IntVal:     iv,
	}

	for k, _ := range mib.OtherTables {
		mv.TableNames[k] = true
	}

	return mv
}

type DeviceMetricsMetadata struct {
	SysName     string `yaml:"sys_name,omitempty"`
	SysObjectID string `yaml:"sys_object_id,omitempty"`
	SysDescr    string `yaml:"sys_descr,omitempty"`
	SysLocation string `yaml:"sys_location,omitempty"`
	SysContact  string `yaml:"sys_contact,omitempty"`
	SysServices int    `yaml:"sys_services,omitempty"`
	EngineID    string `yaml:"-"`
	Customs     map[string]string
	CustomInts  map[string]int64
	Tables      map[string]DeviceTableMetadata
}

const (
	LOGICAL_INTERFACE  = "logical"
	PHYSICAL_INTERFACE = "physical"
	PollAdjustTime     = 5
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
	Type        string            `yaml:"type"`
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
	useGlobal                bool
	origStr                  string
}

type EAPIConfig struct {
	Host      string `yaml:"host"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Transport string `yaml:"transport"`
	Port      int    `yaml:"port"`
	origStr   string
}

type MerakiConfig struct {
	ApiKey string `yaml:"api_key"` // Required.
	Host   string `yaml:"host"`    // Optional, defaults to api.meraki.com
}

// Contain various extensions to snmp which can be used to get data.
type ExtensionSet struct {
	EAPIConfig   *EAPIConfig   `yaml:"eapi_config,omitempty"` // Arista eAPI.
	MerakiConfig *MerakiConfig `yaml:"meraki_config,omitempty"`
}

type SnmpDeviceConfig struct {
	DeviceName          string            `yaml:"device_name"`
	DeviceIP            string            `yaml:"device_ip"`
	Community           string            `yaml:"snmp_comm,omitempty"`
	UseV1               bool              `yaml:"use_snmp_v1"`
	V3                  *V3SNMPConfig     `yaml:"snmp_v3,omitempty"`
	Debug               bool              `yaml:"debug"`
	SampleRate          int64             `yaml:"sample_rate,omitempty"` // Used for flow.
	Port                uint16            `yaml:"port,omitempty"`
	OID                 string            `yaml:"oid"`
	Description         string            `yaml:"description"`
	Checked             time.Time         `yaml:"last_checked"`
	MibProfile          string            `yaml:"mib_profile"`
	Provider            Provider          `yaml:"provider"`
	FlowOnly            bool              `yaml:"flow_only,omitempty"`
	PingOnly            bool              `yaml:"ping_only,omitempty"`
	UserTags            map[string]string `yaml:"user_tags"`
	DiscoveredMibs      []string          `yaml:"discovered_mibs,omitempty"`
	PollTimeSec         int               `yaml:"poll_time_sec,omitempty"`
	TimeoutMS           int               `yaml:"timeout_ms,omitempty"`
	Retries             int               `yaml:"retries,omitempty"`
	EngineID            string            `yaml:"engine_id,omitempty"`
	MatchAttr           map[string]string `yaml:"match_attributes"`
	MonitorAdminShut    bool              `yaml:"monitor_admin_shut"`
	NoUseBulkWalkAll    bool              `yaml:"no_use_bulkwalkall"`
	InstrumentationName string            `yaml:"instrumentationName,omitempty"`
	RunPing             bool              `yaml:"response_time,omitempty"`
	Ext                 *ExtensionSet     `yaml:"ext,omitempty"`
	PingSec             int               `yaml:"ping_interval_sec,omitempty"`
	allUserTags         map[string]string
	walker              SNMPTestWalker
}

type SnmpTrapConfig struct {
	Listen    string        `yaml:"listen"`
	Community string        `yaml:"community"`
	Version   string        `yaml:"version"`
	Transport string        `yaml:"transport"`
	V3        *V3SNMPConfig `yaml:"v3_config"`
}

type SnmpDiscoConfig struct {
	Cidrs              StringArray     `yaml:"cidrs"`
	IgnoreList         StringArray     `yaml:"ignore_list"`
	Debug              bool            `yaml:"debug"`
	Ports              []int           `yaml:"ports"`
	DefaultCommunities []string        `yaml:"default_communities"`
	UseV1              bool            `yaml:"use_snmp_v1"`
	DefaultV3          *V3SNMPConfig   `yaml:"default_v3"`
	OtherV3s           []*V3SNMPConfig `yaml:"other_v3s,omitempty"`
	AddDevices         bool            `yaml:"add_devices"`
	AddAllMibs         bool            `yaml:"add_mibs"`
	Threads            int             `yaml:"threads"`
	ReplaceDevices     bool            `yaml:"replace_devices"`
	NoDedup            bool            `yaml:"no_dedup_engine_id,omitempty"`
	CheckAll           bool            `yaml:"check_all_ips,omitempty"`
	CidrOrig           string          `yaml:"-"`
	IgnoreOrig         string          `yaml:"-"`
}

type SnmpGlobalConfig struct {
	PollTimeSec   int               `yaml:"poll_time_sec"`
	DropIfOutside bool              `yaml:"drop_if_outside_poll"`
	MibProfileDir string            `yaml:"mib_profile_dir"`
	MibDB         string            `yaml:"mibs_db"`
	MibsEnabled   []string          `yaml:"mibs_enabled"`
	TimeoutMS     int               `yaml:"timeout_ms"`
	Retries       int               `yaml:"retries"`
	GlobalV3      *V3SNMPConfig     `yaml:"global_v3"`
	RunPing       bool              `yaml:"response_time"`
	PingSec       int               `yaml:"ping_interval_sec,omitempty"`
	UserTags      map[string]string `yaml:"user_tags"`
	MatchAttr     map[string]string `yaml:"match_attributes"`
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
	Fail             go_metrics.Gauge
	Missing          go_metrics.Gauge
	MissingMeta      go_metrics.Gauge
}

const (
	SNMP_GOOD                   = 1
	SNMP_BAD_NO_DATA            = 3
	SNMP_BAD_DEVICE_ERR         = 4
	SNMP_BAD_INTERFACE_ERR      = 5
	SNMP_BAD_POLL_TIMEOUT       = 6
	SNMP_BAD_INIT_METADATA      = 7
	SNMP_BAD_INIT_METRICS       = 8
	SNMP_BAD_FIRST_METRICS_POLL = 9
)

var (
	SNMP_STATUS_MAP = map[int64]string{
		SNMP_GOOD:                   "GOOD",
		SNMP_BAD_NO_DATA:            "BAD: No data",
		SNMP_BAD_DEVICE_ERR:         "BAD: Device Error",
		SNMP_BAD_INTERFACE_ERR:      "BAD: Interface Error",
		SNMP_BAD_POLL_TIMEOUT:       "BAD: Poll Timeout",
		SNMP_BAD_INIT_METADATA:      "BAD: Init Metadata",
		SNMP_BAD_INIT_METRICS:       "BAD: Init Metrics",
		SNMP_BAD_FIRST_METRICS_POLL: "BAD: Metrics Poll",
	}
)

func NewSnmpDeviceMetric(registry go_metrics.Registry, deviceName string) *SnmpDeviceMetric {
	sm := SnmpDeviceMetric{
		DeviceMetrics:    go_metrics.GetOrRegisterMeter("device_metrics^device_name="+deviceName, registry),
		InterfaceMetrics: go_metrics.GetOrRegisterMeter("interface_metrics^device_name="+deviceName, registry),
		Metadata:         go_metrics.GetOrRegisterMeter("metadata^device_name="+deviceName, registry),
		Errors:           go_metrics.GetOrRegisterMeter("snmp_errors^device_name="+deviceName, registry),
		Fail:             go_metrics.GetOrRegisterGauge("snmp_fail^device_name="+deviceName, registry),
		Missing:          go_metrics.GetOrRegisterGauge("snmp_missing^force=true^device_name="+deviceName, registry),
		MissingMeta:      go_metrics.GetOrRegisterGauge("snmp_missing_meta^force=true^device_name="+deviceName, registry),
	}
	sm.Fail.Update(SNMP_GOOD) // 1 means that this device is not failing.
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
	Oid          string
	Name         string
	Type         Oidtype
	Extra        string
	Tag          string
	Enum         map[string]int64
	EnumRev      map[int64]string
	Conversion   string
	Mib          string
	Table        string
	PollDur      time.Duration
	MatchAttr    map[string]*regexp.Regexp
	lastPoll     time.Time
	FromExtended bool
	OtherTables  map[string]bool
	Format       string
	AllowDup     bool
}

func (mb *Mib) String() string {
	return fmt.Sprintf("Name: %s, Oid: %s: Type: %d, Extra: %s", mb.Name, mb.Oid, mb.Type, mb.Extra)
}

func (mb *Mib) GetName() string { // Tag takes precedince over name if it is present.
	if mb == nil {
		return "missing_mib"
	}
	if mb.Tag != "" {
		return mb.Tag
	}
	return mb.Name
}

func (mb *Mib) IsPollReady() bool { // If there's a poll duration, return false if not enough time has elapsed before this next poll.
	if mb == nil || mb.PollDur == 0 { // If not set, just always return true
		return true
	}
	now := time.Now()
	ready := mb.lastPoll.Add(mb.PollDur).Before(now)
	if ready {
		mb.lastPoll = now
	}
	return ready
}

func (mb *Mib) Extend(nm *Mib) {
	if mb.OtherTables == nil {
		mb.OtherTables = map[string]bool{mb.Table: true}
	}
	mb.OtherTables[nm.Table] = true
}

type LastMetadata struct {
	DeviceInfo    map[string]interface{}
	InterfaceInfo map[IfaceID]map[string]interface{}
	Tables        map[string]DeviceTableMetadata
	MatchAttr     map[string]*regexp.Regexp
	XtraInfo      map[string]MetricInfo
}

func (lm *LastMetadata) Size() int {
	if lm == nil {
		return 0
	}

	return len(lm.DeviceInfo) + len(lm.InterfaceInfo)
}

func (lm *LastMetadata) Missing(new *LastMetadata) []string {
	missing := []string{}
	for k, _ := range lm.DeviceInfo {
		if _, ok := new.DeviceInfo[k]; !ok {
			missing = append(missing, k)
		}
	}
	for ifn, _ := range lm.InterfaceInfo {
		if _, ok := new.InterfaceInfo[ifn]; !ok {
			missing = append(missing, strconv.Itoa(int(ifn)))
		}
	}

	return missing
}

func (lm *LastMetadata) GetTableName(key string) (string, map[string]bool, bool) {
	if lm == nil {
		return "", nil, false
	}
	if i, ok := lm.XtraInfo[key]; ok {
		return i.Table, i.Tables, true
	}
	return "", nil, false
}

type DeviceMap map[string]*SnmpDeviceConfig

func (a *DeviceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi map[string]*SnmpDeviceConfig
	err := unmarshal(&multi)
	if err != nil {
		var mult []string
		err := unmarshal(&mult)
		if err != nil {
			var single string
			err := unmarshal(&single)
			if err != nil {
				return err
			}
			*a = map[string]*SnmpDeviceConfig{"file_0": &SnmpDeviceConfig{DeviceName: single}}
		} else {
			res := map[string]*SnmpDeviceConfig{}
			for i, s := range mult {
				res["file_"+strconv.Itoa(i)] = &SnmpDeviceConfig{DeviceName: s}
			}
			*a = res
		}
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

type V3SNMP V3SNMPConfig // Need a 2nd type alias to avoid stack overflow on parsing.

// Make sure that things serialize back to how they were.
func (a *V3SNMPConfig) MarshalYAML() (interface{}, error) {
	if a.origStr != "" {
		return a.origStr, nil
	}
	return a, nil
}

// This lets the config get overriden by a global_v3 string.
func (a *V3SNMPConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var conf = V3SNMP{}
	err := unmarshal(&conf)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		if single == "@global_v3" { // Should this be hard coded like this?
			conf.useGlobal = true
		} else if strings.HasPrefix(single, "${") { // get the whole yaml block out of an env var.
			raw := os.Getenv(single[2 : len(single)-1])
			if err = yaml.Unmarshal([]byte(raw), &conf); err != nil {
				return err
			}
		} else if strings.HasPrefix(single, AwsSmPrefix) { // See if we can pull these out of AWS Secret Manager directly
			raw := loadViaAWSSecrets(single[len(AwsSmPrefix):])
			if err = yaml.Unmarshal([]byte(raw), &conf); err != nil {
				return err
			}
		}
		conf.origStr = single // Let us know where this came from.
		*a = V3SNMPConfig(conf)
	} else {
		// Now, see if we need to map in any ENV vars.
		fields := reflect.VisibleFields(reflect.TypeOf(conf))
		ps := reflect.ValueOf(&conf)
		for _, field := range fields {
			if field.Type.Kind() == reflect.String {
				s := ps.Elem()
				f := s.FieldByName(field.Name)
				if f.IsValid() && f.CanSet() {
					if sval, ok := f.Interface().(string); ok {
						if strings.HasPrefix(sval, "${") { // Expecting values of the form ${V3_AUTH_PROTOCOL}
							f.SetString(os.Getenv(sval[2 : len(sval)-1]))
						} else if strings.HasPrefix(sval, AwsSmPrefix) { // See if we can pull these out of AWS Secret Manager directly
							f.SetString(loadViaAWSSecrets(sval[len(AwsSmPrefix):]))
						}
					}
				}
			}
		}
		// And pop back what we created.
		*a = V3SNMPConfig(conf)
	}
	return nil
}

func (a *V3SNMPConfig) InheritGlobal() bool {
	if a == nil {
		return false
	}
	return a.useGlobal
}

// Save any hard coded parts of this profile which might get overriten by an automatic process.
func (d *SnmpDeviceConfig) UpdateFrom(old *SnmpDeviceConfig, conf *SnmpConfig) {
	if old == nil {
		return
	}

	if strings.HasPrefix(old.MibProfile, "!") {
		d.MibProfile = old.MibProfile
	}
	if old.SampleRate > 0 {
		d.SampleRate = old.SampleRate
	}
	if old.FlowOnly {
		d.FlowOnly = old.FlowOnly
	}
	if old.PingOnly {
		d.PingOnly = old.PingOnly
	}
	if old.PollTimeSec > 0 {
		d.PollTimeSec = old.PollTimeSec
	}
	if old.TimeoutMS > 0 {
		d.TimeoutMS = old.TimeoutMS
	}
	if old.Retries > 0 {
		d.Retries = old.Retries
	}
	if old.MonitorAdminShut {
		d.MonitorAdminShut = old.MonitorAdminShut
	}
	if old.NoUseBulkWalkAll {
		d.NoUseBulkWalkAll = old.NoUseBulkWalkAll
	}
	if old.RunPing {
		d.RunPing = old.RunPing
	}
	if d.UserTags == nil && old.UserTags != nil {
		d.UserTags = map[string]string{}
	}
	for k, v := range old.UserTags {
		if _, ok := d.UserTags[k]; !ok {
			if _, ok := conf.Global.UserTags[k]; !ok {
				d.UserTags[k] = v
			}
		}
	}
	if d.MatchAttr == nil && old.MatchAttr != nil {
		d.MatchAttr = map[string]string{}
	}
	for k, v := range old.MatchAttr {
		if _, ok := d.MatchAttr[k]; !ok {
			if _, ok := conf.Global.MatchAttr[k]; !ok {
				d.MatchAttr[k] = v
			}
		}
	}
}

func (d *SnmpDeviceConfig) AddUserTag(k string, v string) {
	if d.allUserTags == nil {
		d.allUserTags = map[string]string{}
	}
	key := k
	if !strings.HasPrefix(key, UserTagPrefix) {
		key = UserTagPrefix + k
	}
	d.allUserTags[key] = v
}

func (d *SnmpDeviceConfig) InitUserTags(serviceName string) {
	d.allUserTags = map[string]string{}
	if serviceName != "ktranslate" {
		if d.UserTags == nil { // Prevent nil map assignment.
			d.UserTags = map[string]string{}
		}
		d.UserTags["container_service"] = serviceName
	}

	for k, v := range d.UserTags {
		key := k
		if !strings.HasPrefix(key, UserTagPrefix) {
			key = UserTagPrefix + k
		}
		d.allUserTags[key] = v
	}
}

func (d *SnmpDeviceConfig) SetUserTags(in map[string]string) {
	for k, v := range d.allUserTags {
		in[k] = v
	}
}

func (d *SnmpDeviceConfig) GetUserTags() map[string]string {
	if d.allUserTags == nil {
		return nil
	}

	out := map[string]string{}
	for k, v := range d.allUserTags {
		out[k] = v
	}

	return out
}

func (d *SnmpDeviceConfig) GetTestWalker() SNMPTestWalker {
	return d.walker
}

func (d *SnmpDeviceConfig) SetTestWalker(w SNMPTestWalker) {
	d.walker = w
}

type SNMPTestWalker interface {
	WalkAll(string) ([]gosnmp.SnmpPDU, error)
}

type EAPIC EAPIConfig // Need a 2nd type alias to avoid stack overflow on parsing.

// Make sure that things serialize back to how they were.
func (a *EAPIConfig) MarshalYAML() (interface{}, error) {
	if a.origStr != "" {
		return a.origStr, nil
	}
	return a, nil
}

// This lets the config get overriden by a global_v3 string.
func (a *EAPIConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var conf = EAPIC{}
	err := unmarshal(&conf)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		if strings.HasPrefix(single, "${") { // get the whole yaml block out of an env var.
			raw := os.Getenv(single[2 : len(single)-1])
			if err = yaml.Unmarshal([]byte(raw), &conf); err != nil {
				return err
			}
		} else if strings.HasPrefix(single, AwsSmPrefix) { // See if we can pull these out of AWS Secret Manager directly
			raw := loadViaAWSSecrets(single[len(AwsSmPrefix):])
			if err = yaml.Unmarshal([]byte(raw), &conf); err != nil {
				return err
			}
		}
		conf.origStr = single // Let us know where this came from.
		*a = EAPIConfig(conf)
	} else { // Just use the value directly.
		*a = EAPIConfig(conf)
	}

	return nil
}
