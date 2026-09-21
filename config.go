package ktranslate

import (
	"context"
	yaml "gopkg.in/yaml.v3"
	"io/fs"

	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
)

const (
	// FlowFields are fields for the Flow input
	FlowFields = "Type,TimeReceived,SequenceNum,SamplingRate,FlowDirection,SamplerAddress,TimeFlowStart,TimeFlowEnd,Bytes,Packets,SrcAddr,DstAddr,Etype,Proto,SrcPort,DstPort,InIf,OutIf,SrcMac,DstMac,SrcVlan,DstVlan,VlanId,IPTos,ForwardingStatus,IPTTL,TCPFlags,IcmpType,IcmpCode,IPv6FlowLabel,FragmentId,FragmentOffset,SrcAS,DstAS,NextHop,NextHopAS,SrcNet,DstNet,MPLSCount"
	// FlowDefaultFields are the default fields for flow
	FlowDefaultFields = "TimeReceived,SamplingRate,Bytes,Packets,SrcAddr,DstAddr,Proto,SrcPort,DstPort,InIf,OutIf,SrcVlan,DstVlan,TCPFlags,SrcAS,DstAS,Type,SamplerAddress,FlowDirection"
	// KentikAPITokenEnvVar is the environment variables used to get the Kentik API Token
	KentikAPITokenEnvVar = "KENTIK_API_TOKEN"
	// MaxNetflowsPerMessage is the maximum flowsets in a single netflow message
	// Must be less than ( 65535[nf_max_size] - 16[nf_header_size] - 96(nf_template_flowset_size] ) / 128[nf_data_flowset_size]
	MaxNetflowsPerMessage = 510
)

// NetflowFormatConfig is the config format for netflow
type NetflowFormatConfig struct {
	Version string
}

// InfluxDBFormatConfig is the config format for influxdb
type InfluxDBFormatConfig struct {
	MeasurementPrefix string
	NamespaceToken    string
}

// PrometheusFormatConfig is the config for the prometheus format
type PrometheusFormatConfig struct {
	EnableCollectorStats bool
	FlowsNeeded          int
}

// ElasticFormatConfig is the config for the elastic format
type ElasticFormatConfig struct {
	Action string
}

// OtelFormatConfig is the config for the otel format
type OtelFormatConfig struct {
	Endpoint      string
	Protocol      string
	ClientCert    string
	ClientKey     string
	RootCA        string
	NoBlockExport bool
}

// RedisFormatConfig is the config for the redis format
type RedisFormatConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	KeyPrefix     string
	KeyTTLSeconds int
}

// SnmpFormatConfig is the config for the snmp format
type SnmpFormatConfig struct {
	ConfigFile string
}

// NRMFormatConfig is the config for the new_relic_metric format
type NRMFormatConfig struct {
	// CustomAttributes are stamped into the common attributes block sent with every
	// metric batch, regardless of device/metric type -- see -nr_custom_attributes.
	CustomAttributes map[string]string
}

// S3SinkConfig is the config for the S3-backed object store, used to
// dereference large out-of-band objects (e.g. HAR files) referenced by path.
// Not a selectable --sinks destination — see pkg/cat/kkc.go's objmgr wiring.
type S3SinkConfig struct {
	Bucket                                     string
	Prefix                                     string
	FlushIntervalSeconds                       int
	AssumeRoleARN                              string
	Region                                     string
	EC2InstanceProfile                         bool
	AssumeRoleOrInstanceProfileIntervalSeconds int
	CheckDangling                              bool
	Endpoint                                   string
	SigningRegion                              string
}

// HTTPSinkConfig is the config for the HTTP sink
type HTTPSinkConfig struct {
	Target             string
	TargetLogs         string
	Headers            []string
	InsecureSkipVerify bool
	TimeoutInSeconds   int
}

// NetSinkConfig is the config for the net sink
type NetSinkConfig struct {
	Endpoint string
	Protocol string
}

// NewRelicSinkConfig is the config for the NewRelic sink
type NewRelicSinkConfig struct {
	Account      string
	EstimateOnly bool
	Region       string
	ValidateJSON bool
}

// NewRelicMultiSinkConfig is the config for Multi New Relic
type NewRelicMultiSinkConfig struct {
	CredMap map[int]NRCred
}

// NRCred exposes a list of NR creds.
type NRCred struct {
	NRAccount  string
	NRApiToken string
}

// RollupConfig is the config for rollups
type RollupConfig struct {
	JoinKey          string
	TopK             int
	Formats          []string
	KeepUndefined    bool
	MaxMemoryMB      int
	MaxKeys          int
	EmergencyCleanup bool
}

// KMuxConfig is the config for the mux server
type KMuxConfig struct {
	Dir string
}

// ServerConfig is the config for the meta server
type ServerConfig struct {
	ServiceName     string
	LogLevel        string
	LogToStdout     bool
	MetricsEndpoint string
	MetaListenAddr  string
	CfgPath         string `yaml:"-"` // We don't want to read this directly because it comes from a flag but saved here for internal use.
}

// APIConfig is the config for the API service
type APIConfig struct {
	DeviceFile string
}

// SyslogInputConfig is the config for the syslog input
type SyslogInputConfig struct {
	Enable     bool
	ListenAddr string
	EnableTCP  bool
	EnableUDP  bool
	EnableUnix bool
	Format     string
	Threads    int
}

// SNMPInputConfig is the config for SNMP input
type SNMPInputConfig struct {
	Enable                   bool
	SNMPFile                 string
	DumpMIBs                 bool
	FlowOnly                 bool
	JSONToYAML               string
	WalkTarget               string
	WalkOID                  string
	WalkFormat               string
	OutputFile               string
	DiscoveryIntervalMinutes int
	DiscoveryOnStart         bool
	ValidateMIBs             bool
	RunFromWalkFile          string
	PollNowTarget            string
}

// GCPVPCInputConfig is the config for GCP VPC
type GCPVPCInputConfig struct {
	Enable     bool
	ProjectID  string
	Subject    string
	SampleRate float64
}

// AWSVPCInputConfig is the config for AWS VPC
type AWSVPCInputConfig struct {
	Enable    bool
	IAMRole   string
	SQSName   string
	Regions   []string
	IsLambda  bool
	LocalFile string
}

// FlowInputConfig is the config for flow input
type FlowInputConfig struct {
	Enable               bool
	Protocol             string
	ListenIP             string
	ListenPort           int
	EnableReusePort      bool
	Workers              int
	QueueSize            int
	MessageFields        string
	PrometheusListenAddr string
	MappingFile          string
}

// KentikCred is information needed to auth the Kentik API.
type KentikCred struct {
	APIEmail string
	APIToken string
}

// ConfigManager is the config for how to manage configs.
type ConfigManager struct {
	ConfigImpl  string
	PollTimeSec int
}

// StitchConfig is the config on how to manage stitching flows together
type StitchConfig struct {
	Enable bool
	BufLen int
}

// Config is the ktranslate configuration
type Config struct {
	// ktranslate
	ListenAddr          string
	MappingFile         string
	UDRSFile            string
	KMUdr               string
	GeoFile             string
	ASNFile             string
	ApplicationFile     string
	DNS                 string
	ProcessingThreads   int
	InputThreads        int
	MaxThreads          int
	Format              string
	FormatRollup        string
	FormatMetric        string
	Compression         string
	Sinks               []string
	MaxFlowsPerMessage  int
	RollupInterval      int
	RollupAndAlpha      bool
	SampleRate          int
	SampleMin           int
	EnableSNMPDiscovery bool
	KentikCreds         []KentikCred
	KentikPlan          int
	APIBaseURL          string
	GRPCBaseURL         string
	SSLCertFile         string
	SSLKeyFile          string
	TagMapType          string
	EnableTeeLogs       bool
	EnableHTTPInput     bool
	HttpRemoteIp        string
	EnricherURL         string
	TeeFlow             string
	EnricherSource      string
	EnricherScript      string

	// pkg/maps/file
	TagMapFile string
	// pkg/maps/file
	TagMapCity string
	// pkg/maps/file
	TagMapRegion string
	// pkg/formats/netflow
	NetflowFormat *NetflowFormatConfig
	// pkg/formats/prom
	PrometheusFormat *PrometheusFormatConfig
	// pkg/formats/influxdb
	InfluxDBFormat *InfluxDBFormatConfig
	// pkg/formats/otel
	OtelFormat *OtelFormatConfig
	// pkg/formats/redis
	RedisFormat *RedisFormatConfig
	// pkg/formats/snmp
	SnmpFormat *SnmpFormatConfig
	// pkg/formats/elasticsearch
	ElasticFormat *ElasticFormatConfig
	// pkg/formats/nrm
	NRMFormat *NRMFormatConfig

	// pkg/sinks/s3 (object store backing for out-of-band object dereferencing, not a --sinks option)
	S3Sink *S3SinkConfig
	// pkg/sinks/http
	HTTPSink *HTTPSinkConfig
	// pkg/sinks/net
	NetSink *NetSinkConfig
	// pkg/sinks/nr
	NewRelicSink *NewRelicSinkConfig
	// pkg/sinks/nrmulti
	NewRelicMultiSink *NewRelicMultiSinkConfig

	// pkg/rollup
	Rollup *RollupConfig
	// pkg/eggs/kmux
	KMux *KMuxConfig
	// pkg/eggs/baseserver
	Server *ServerConfig
	// pkg/api
	API *APIConfig
	// pkg/filter
	Filters []string

	// pkg/inputs/syslog
	SyslogInput *SyslogInputConfig
	// pkg/inputs/snmp
	SNMPInput *SNMPInputConfig
	// pkg/inputs/vpc/gcp
	GCPVPCInput *GCPVPCInputConfig
	// pkg/inputs/vpc/aws
	AWSVPCInput *AWSVPCInputConfig
	// pkg/inputs/flow
	FlowInput *FlowInputConfig
	// pkg/config
	CfgManager *ConfigManager
	// pkg/stitch
	Lilo *StitchConfig
}

// DefaultConfig returns a ktranslate configuration with defaults applied
func DefaultConfig() *Config {
	return &Config{
		ListenAddr:          "127.0.0.1:8081",
		MappingFile:         "",
		UDRSFile:            "",
		KMUdr:               "",
		GeoFile:             "",
		ASNFile:             "",
		ApplicationFile:     "",
		DNS:                 "",
		ProcessingThreads:   1,
		InputThreads:        1,
		MaxThreads:          1,
		Format:              "flat_json",
		FormatRollup:        "",
		Compression:         "none",
		Sinks:               []string{"stdout"},
		MaxFlowsPerMessage:  10000,
		RollupInterval:      0,
		RollupAndAlpha:      false,
		SampleRate:          1,
		SampleMin:           1,
		EnableSNMPDiscovery: false,
		KentikCreds:         nil,
		KentikPlan:          0,
		APIBaseURL:          "https://api.kentik.com",
		SSLCertFile:         "",
		SSLKeyFile:          "",
		TagMapType:          "",
		EnableTeeLogs:       false,
		EnableHTTPInput:     false,
		EnricherURL:         "",
		TagMapFile:          "",
		NetflowFormat: &NetflowFormatConfig{
			Version: "ipfix",
		},
		PrometheusFormat: &PrometheusFormatConfig{
			EnableCollectorStats: false,
			FlowsNeeded:          10,
		},
		OtelFormat: &OtelFormatConfig{
			Endpoint:      "",
			Protocol:      "stdout",
			ClientKey:     "",
			ClientCert:    "",
			RootCA:        "",
			NoBlockExport: false,
		},
		RedisFormat: &RedisFormatConfig{
			RedisAddr:     "localhost:6379",
			RedisPassword: "",
			RedisDB:       0,
			KeyPrefix:     "",
			KeyTTLSeconds: 60,
		},
		ElasticFormat: &ElasticFormatConfig{
			Action: "index",
		},
		SnmpFormat: &SnmpFormatConfig{
			ConfigFile: "",
		},
		NRMFormat: &NRMFormatConfig{
			CustomAttributes: map[string]string{},
		},
		InfluxDBFormat: &InfluxDBFormatConfig{
			MeasurementPrefix: "",
			NamespaceToken:    ":",
		},
		S3Sink: &S3SinkConfig{
			Bucket:               "",
			Prefix:               "/kentik",
			FlushIntervalSeconds: 60,
			AssumeRoleARN:        "",
			Region:               "us-east-1",
			EC2InstanceProfile:   false,
			AssumeRoleOrInstanceProfileIntervalSeconds: 900,
			CheckDangling: false,
			Endpoint:      "",
			SigningRegion: "",
		},
		HTTPSink: &HTTPSinkConfig{
			Target:             "http://localhost:8086/write?db=kentik",
			TargetLogs:         "http://localhost:8088/services/collector/event",
			Headers:            []string{},
			InsecureSkipVerify: false,
			TimeoutInSeconds:   30,
		},
		NetSink: &NetSinkConfig{
			Endpoint: "",
			Protocol: "udp",
		},
		NewRelicSink: &NewRelicSinkConfig{
			Account:      "",
			EstimateOnly: false,
			Region:       "",
			ValidateJSON: false,
		},
		NewRelicMultiSink: &NewRelicMultiSinkConfig{
			CredMap: nil,
		},
		Rollup: &RollupConfig{
			JoinKey:          "^",
			TopK:             10,
			Formats:          []string{},
			KeepUndefined:    false,
			MaxMemoryMB:      100,
			MaxKeys:          5000,
			EmergencyCleanup: true,
		},
		KMux: &KMuxConfig{
			Dir: ".",
		},
		Server: &ServerConfig{
			ServiceName:     "",
			LogLevel:        "info",
			LogToStdout:     false,
			MetricsEndpoint: "none",
			MetaListenAddr:  "localhost:0",
			CfgPath:         "",
		},
		API: &APIConfig{
			DeviceFile: "",
		},
		Filters: []string{},
		SyslogInput: &SyslogInputConfig{
			Enable:     false,
			EnableTCP:  true,
			EnableUDP:  true,
			EnableUnix: false,
			Format:     "Automatic",
			Threads:    1,
		},
		SNMPInput: &SNMPInputConfig{
			Enable:                   false,
			DumpMIBs:                 false,
			SNMPFile:                 "",
			FlowOnly:                 false,
			JSONToYAML:               "",
			WalkTarget:               "",
			WalkOID:                  ".1.3.6.1.2.1",
			WalkFormat:               "",
			OutputFile:               "",
			DiscoveryIntervalMinutes: 0,
			DiscoveryOnStart:         false,
			ValidateMIBs:             false,
			RunFromWalkFile:          "",
		},
		GCPVPCInput: &GCPVPCInputConfig{
			Enable:     false,
			ProjectID:  "",
			Subject:    "",
			SampleRate: float64(1.0),
		},
		AWSVPCInput: &AWSVPCInputConfig{
			Enable:    false,
			IAMRole:   "",
			SQSName:   "",
			Regions:   []string{"us-east-1"},
			IsLambda:  false,
			LocalFile: "",
		},
		FlowInput: &FlowInputConfig{
			Enable:               false,
			Protocol:             "",
			ListenIP:             "0.0.0.0",
			ListenPort:           9995,
			EnableReusePort:      false,
			Workers:              2,
			QueueSize:            10000,
			MessageFields:        FlowDefaultFields,
			PrometheusListenAddr: "",
			MappingFile:          "",
		},
		CfgManager: &ConfigManager{
			ConfigImpl:  "",
			PollTimeSec: 1200,
		},
		Lilo: &StitchConfig{
			Enable: false,
			BufLen: 10000,
		},
	}
}

// LoadConfig returns a ktranslate configuration from the specified path
func LoadConfig(ctx context.Context, configPath string) (*Config, error) {
	confBytes, err := snmp_util.LoadFile(ctx, configPath)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err := yaml.Unmarshal(confBytes, &cfg); err != nil {
		return nil, err
	}

	if cfg.Format == "netflow" && (cfg.MaxFlowsPerMessage <= 0 || cfg.MaxFlowsPerMessage > MaxNetflowsPerMessage) {
		cfg.MaxFlowsPerMessage = MaxNetflowsPerMessage
	}

	return &cfg, nil
}

// SaveConfig saves the ktranslate configuration to the specified path
func (c *Config) SaveConfig() error {
	t, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	permissions := fs.FileMode(0644)
	return snmp_util.WriteFile(context.Background(), c.Server.CfgPath, t, permissions)
}
