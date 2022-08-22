package ktranslate

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

const (
	// FlowFields are fields for the Flow input
	FlowFields = "Type,TimeReceived,SequenceNum,SamplingRate,SamplerAddress,TimeFlowStart,TimeFlowEnd,Bytes,Packets,SrcAddr,DstAddr,Etype,Proto,SrcPort,DstPort,InIf,OutIf,SrcMac,DstMac,SrcVlan,DstVlan,VlanId,IngressVrfID,EgressVrfID,IPTos,ForwardingStatus,IPTTL,TCPFlags,IcmpType,IcmpCode,IPv6FlowLabel,FragmentId,FragmentOffset,BiFlowDirection,SrcAS,DstAS,NextHop,NextHopAS,SrcNet,DstNet,HasMPLS,MPLSCount,MPLS1TTL,MPLS1Label,MPLS2TTL,MPLS2Label,MPLS3TTL,MPLS3Label,MPLSLastTTL,MPLSLastLabel,CustomInteger1,CustomInteger2,CustomInteger3,CustomInteger4,CustomInteger5,CustomBytes1,CustomBytes2,CustomBytes3,CustomBytes4,CustomBytes5"
	// FlowDefaultFields are the default fields for flow
	FlowDefaultFields = "TimeReceived,SamplingRate,Bytes,Packets,SrcAddr,DstAddr,Proto,SrcPort,DstPort,InIf,OutIf,SrcVlan,DstVlan,TCPFlags,SrcAS,DstAS,Type,SamplerAddress"
	// KentikAPITokenEnvVar is the environment variables used to get the Kentik API Token
	KentikAPITokenEnvVar = "KENTIK_API_TOKEN"
)

// NetflowFormatConfig is the config format for netflow
type NetflowFormatConfig struct {
	Version string
}

// PrometheusFormatConfig is the config for the prometheus format
type PrometheusFormatConfig struct {
	EnableCollectorStats bool
	FlowsNeeded          int
}

// PrometheusSinkConfig is config for the prometheus sink
type PrometheusSinkConfig struct {
	ListenAddr string
}

// GCloudSinkConfig is the config for GCP
type GCloudSinkConfig struct {
	Bucket               string
	Prefix               string
	ContentType          string
	FlushIntervalSeconds int
}

// S3SinkConfig is the config for the S3 sink
type S3SinkConfig struct {
	Bucket               string
	Prefix               string
	FlushIntervalSeconds int
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

// FileSinkConfig is the config for the file sink
type FileSinkConfig struct {
	Path                 string
	EnableImmediateWrite bool
	FlushIntervalSeconds int
}

// GCloudPubSubSinkConfig is the config for GCP PubSub
type GCloudPubSubSinkConfig struct {
	ProjectID string
	Topic     string
}

// HTTPSinkConfig is the config for the HTTP sink
type HTTPSinkConfig struct {
	Target             string
	Headers            []string
	InsecureSkipVerify bool
	TimeoutInSeconds   int
}

// KafkaSinkConfig is the config for the Kafka sink
type KafkaSinkConfig struct {
	Topic            string
	BootstrapServers string
}

// KentikSinkConfig is the config for the Kentik sink
type KentikSinkConfig struct {
	RelayURL string
}

// RollupConfig is the config for rollups
type RollupConfig struct {
	JoinKey string
	TopK    int
	Formats []string
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
	OllyDataset     string
	OllyWriteKey    string
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
	MessageFields        string
	PrometheusListenAddr string
	MappingFile          string
}

// Config is the ktranslate configuration
type Config struct {
	// ktranslate
	ListenAddr          string
	MappingFile         string
	UDRSFile            string
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
	Sinks               string
	MaxFlowsPerMessage  int
	RollupInterval      int
	RollupAndAlpha      bool
	SampleRate          int
	SampleMin           int
	EnableSNMPDiscovery bool
	KentikEmail         string
	KentikAPIToken      string
	KentikPlan          int
	APIBaseURL          string
	SSLCertFile         string
	SSLKeyFile          string
	TagMapType          string
	EnableTeeLogs       bool
	EnableHTTPInput     bool
	EnricherURL         string

	// pkg/maps/file
	TagMapFile string
	// pkg/formats/netflow
	NetflowFormat *NetflowFormatConfig
	// pkg/formats/prom
	PrometheusFormat *PrometheusFormatConfig

	// pkg/sinks/prom
	PrometheusSink *PrometheusSinkConfig
	// pkg/sinks/gcloud
	GCloudSink *GCloudSinkConfig
	// pkg/sinks/s3
	S3Sink *S3SinkConfig
	// pkg/sinks/net
	NetSink *NetSinkConfig
	// pkg/sinks/nr
	NewRelicSink *NewRelicSinkConfig
	// pkg/sinks/file
	FileSink *FileSinkConfig
	// pkg/sinks/gcppubsub
	GCloudPubSubSink *GCloudPubSubSinkConfig
	// pkg/sinks/http
	HTTPSink *HTTPSinkConfig
	// pkg/sinks/kafka
	KafkaSink *KafkaSinkConfig
	// pkg/sinks/kentik
	KentikSink *KentikSinkConfig

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
}

// DefaultConfig returns a ktranslate configuration with defaults applied
func DefaultConfig() *Config {
	return &Config{
		ListenAddr:          "127.0.0.1:8081",
		MappingFile:         "",
		UDRSFile:            "",
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
		Sinks:               "stdout",
		MaxFlowsPerMessage:  10000,
		RollupInterval:      0,
		RollupAndAlpha:      false,
		SampleRate:          1,
		SampleMin:           1,
		EnableSNMPDiscovery: false,
		KentikEmail:         "",
		KentikAPIToken:      os.Getenv(KentikAPITokenEnvVar),
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
		PrometheusSink: &PrometheusSinkConfig{
			ListenAddr: ":8082",
		},
		GCloudSink: &GCloudSinkConfig{
			Bucket:               "",
			Prefix:               "/kentik",
			ContentType:          "application/json",
			FlushIntervalSeconds: 60,
		},
		S3Sink: &S3SinkConfig{
			Bucket:               "",
			Prefix:               "/kentik",
			FlushIntervalSeconds: 60,
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
		FileSink: &FileSinkConfig{
			Path:                 "./",
			EnableImmediateWrite: false,
			FlushIntervalSeconds: 60,
		},
		GCloudPubSubSink: &GCloudPubSubSinkConfig{
			ProjectID: "",
			Topic:     "",
		},
		HTTPSink: &HTTPSinkConfig{
			Target:             "http://localhost:8086/write?db=kentik",
			Headers:            []string{},
			InsecureSkipVerify: false,
			TimeoutInSeconds:   30,
		},
		KafkaSink: &KafkaSinkConfig{
			Topic:            "",
			BootstrapServers: "",
		},
		KentikSink: &KentikSinkConfig{
			RelayURL: "",
		},
		Rollup: &RollupConfig{
			JoinKey: "^",
			TopK:    10,
			Formats: []string{},
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
			OllyDataset:     "",
			OllyWriteKey:    "",
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
			Protocol:             "netflow5",
			ListenIP:             "0.0.0.0",
			ListenPort:           9995,
			EnableReusePort:      false,
			Workers:              1,
			MessageFields:        FlowDefaultFields,
			PrometheusListenAddr: "",
			MappingFile:          "",
		},
	}
}

// LoadConfig returns a ktranslate configuration from the specified path
func LoadConfig(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg *Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
