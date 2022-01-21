# KTranslate - Kentik data to the world

Listen for a feed of data to or from Kentik and pass on in a common form. Supports rollups and filtering as well.

See the [Wiki](https://github.com/kentik/ktranslate/wiki) for more details. Come visit the [Discord](https://discord.gg/kentik) if you have any questions, need any assistance, or want to talk about the development of ktranslate.

# Build:

make && make test

# Build Docker Image:

To build and use a Docker image, you must specify `MAXMIND_DOWNLOAD_KEY` as a build arg:

```bash
docker build --build-arg MAXMIND_LICENSE_KEY=xxxxx -t ktranslate:v2 .
```

To get your own MaxMind key, visit [MaxMind](https://www.maxmind.com).

# Flags:

```Usage of ./bin/ktranslate:
  -api_device_file string
    	File to sideload devices without hitting API
  -api_devices string
    	json file containing dumy devices to use for the stub Kentik API
  -api_root string
    	API url prefix. If not set, defaults to https://api.kentik.com (default "https://api.kentik.com")
  -application_map string
    	File containing custom application mappings
  -asn string
    	Asn mapping file
  -aws_lambda
    	Run as a AWS Lambda function
  -aws_local_file string
    	If set, process this local file and exit
  -aws_regions string
    	CSV list of region to run in. Will look for metadata in all regions, run SQS in first region. (default "us-east-1")
  -bootstrap.servers string
    	bootstrap.servers
  -compression string
    	compression algo to use (none|gzip|snappy|deflate|null) (default "none")
  -dns string
    	Resolve IPs at this ip:port
  -enricher string
    	Send data to this http url for enrichment.
  -file_flush_sec int
    	Create a new output file every this many seconds (default 60)
  -file_on
    	If true, start writting to file sink right away. Otherwise, wait for a USR1 signal
  -file_out string
    	Write flows seen to log to this directory if set (default "./")
  -filters value
    	Any filters to use. Format: type dimension operator value
  -flow_only
    	If true, don't poll snmp devices.
  -format string
    	Format to convert kflow to: (json|flat_json|avro|netflow|influx|prometheus|new_relic|new_relic_metric|elasticsearch|kflow) (default "flat_json")
  -format_rollup string
    	Format to convert rollups to: (json|avro|netflow|influx|prometheus|new_relic|new_relic_metric|elasticsearch|kflow)
  -gcloud_bucket string
    	GCloud Storage Bucket to write flows to
  -gcloud_content_type string
    	GCloud Storage Content Type (default "application/json")
  -gcloud_prefix string
    	GCloud Storage object prefix (default "/kentik")
  -gcp.project string
    	Google ProjectID to listen for flows on
  -gcp.sub string
    	Google Sub to listen for flows on
  -gcp_pubsub_project_id string
    	GCP PubSub Project ID to use
  -gcp_pubsub_topic string
    	GCP PubSub Topic to publish to
  -geo string
    	Geo mapping file
  -http.source
    	Listen for content sent via http.
  -http_header value
    	Any custom http headers to set on outbound requests
  -http_url string
    	URL to post to (default "http://localhost:8086/write?db=kentik")
  -iam_role string
    	IAM Role to use for processing flow
  -info_collector
    	Also send stats about this collector
  -input_threads int
    	Number of threads to run for input processing
  -kafka_topic string
    	kafka topic to produce on
  -kentik_email string
    	Kentik email to use for API calls
  -kentik_plan int
    	Kentik plan id to use for creating devices
  -kentik_relay_url string
    	If set, override the kentik api url to send flow over here.
  -listen string
    	IP:Port to listen on (default "127.0.0.1:8081")
  -log_level string
    	Logging Level (default "info")
  -mapping string
    	Mapping file to use for enums
  -max_before_sample int
    	Only sample when a set of inputs is at least this many (default 1)
  -max_flows_per_message int
    	Max number of flows to put in each emitted message (default 10000)
  -max_threads int
    	Dynamically grow threads up to this number
  -metalisten string
    	HTTP interface and port to bind on
  -metrics string
    	Metrics Configuration. none|syslog|stderr|graphite:127.0.0.1:2003 (default "none")
  -net_protocol string
    	Use this protocol for writing data (udp|tcp|unix) (default "udp")
  -net_server string
    	Write flows seen to this address (host and port)
  -netflow_version string
    	Version of netflow to produce: (netflow9|ipfix) (default "ipfix")
  -nf.addr string
    	Sflow/NetFlow/IPFIX listening address (default "0.0.0.0")
  -nf.mapping string
    	Configuration file for custom netflow mappings
  -nf.message.fields string
    	The list of fields to include in flow messages. Can be any of Type,TimeReceived,SequenceNum,SamplingRate,SamplerAddress,TimeFlowStart,TimeFlowEnd,Bytes,Packets,SrcAddr,DstAddr,Etype,Proto,SrcPort,DstPort,InIf,OutIf,SrcMac,DstMac,SrcVlan,DstVlan,VlanId,IngressVrfID,EgressVrfID,IPTos,ForwardingStatus,IPTTL,TCPFlags,IcmpType,IcmpCode,IPv6FlowLabel,FragmentId,FragmentOffset,BiFlowDirection,SrcAS,DstAS,NextHop,NextHopAS,SrcNet,DstNet,HasMPLS,MPLSCount,MPLS1TTL,MPLS1Label,MPLS2TTL,MPLS2Label,MPLS3TTL,MPLS3Label,MPLSLastTTL,MPLSLastLabel,CustomInteger1,CustomInteger2,CustomBytes1,CustomBytes2 (default "TimeReceived,SamplingRate,Bytes,Packets,SrcAddr,DstAddr,Proto,SrcPort,DstPort,InIf,OutIf,SrcVlan,DstVlan,TCPFlags,SrcAS,DstAS,Type,SamplerAddress")
  -nf.port int
    	Sflow/NetFlow/IPFIX listening port (default 9995)
  -nf.prom.listen string
    	Run a promethues metrics collector here
  -nf.reuserport
    	Enable so_reuseport for Sflow/NetFlow/IPFIX
  -nf.source string
    	Run NetFlow Ingest Directly. Valid values here are netflow5|netflow9|ipfix|sflow
  -nf.workers int
    	Number of workers per flow collector (default 1)
  -nr_account_id string
    	If set, sends flow to New Relic
  -nr_check_json
    	Verify body is valid json before sending on
  -nr_estimate_only
    	If true, record size of inputs to NR but don't actually send anything
  -nr_region string
    	NR Region to use. US|EU
  -olly_dataset string
    	Olly dataset name
  -olly_write_key string
    	Olly dataset name
  -prom_listen string
    	Bind to listen for prometheus requests on. (default ":8082")
  -prom_seen int
    	Number of flows needed inbound before we start writting to the collector (default 10)
  -rollup_and_alpha
    	Send both rollups and alpha inputs to sinks
  -rollup_interval int
    	Export timer for rollups in seconds
  -rollup_key_join string
    	Token to use to join dimension keys together (default "^")
  -rollup_top_k int
    	Export only these top values (default 10)
  -rollups value
    	Any rollups to use. Format: type, name, metric, dimension 1, dimension 2, ..., dimension n: sum,bytes,in_bytes,dst_addr
  -s3_bucket string
    	AWS S3 Bucket to write flows to
  -s3_flush_sec int
    	Create a new output file every this many seconds (default 60)
  -s3_prefix string
    	AWS S3 Object prefix (default "/kentik")
  -sample_rate int
    	Sampling rate to use. 1 -> 1:1 sampling, 2 -> 1:2 sampling and so on.
  -service_name string
    	Service identifier (default "ktranslate")
  -sinks string
    	List of sinks to send data to. Options: (kafka|stdout|new_relic|kentik|net|http|prometheus|file|s3|gcloud) (default "stdout")
  -snmp string
    	yaml file containing snmp config to use
  -snmp_discovery
    	If true, try to discover snmp devices on this network as configured.
  -snmp_do_walk string
    	If set, try to perform a snmp walk against the targeted device.
  -snmp_dump_mibs
    	If true, dump the list of possible mibs on start.
  -snmp_json2yaml string
    	If set, convert the passed in json file to a yaml profile.
  -snmp_out_file string
    	If set, write updated snmp file here.
  -snmp_poll_now string
    	If set, run one snmp poll for the specified device and then exit.
  -snmp_walk_format string
    	use this format for walked values if -snmp_do_walk is set.
  -snmp_walk_oid string
    	Walk this oid if -snmp_do_walk is set. (default ".1.3.6.1.2.1")
  -sqs_name string
    	Listen for events from this queue for new objects to look at.
  -ssl_cert_file string
    	SSL Cert file to use for serving HTTPS traffic
  -ssl_key_file string
    	SSL Key file to use for serving HTTPS traffic
  -stdout
    	Log to stdout (default true)
  -syslog.format string
    	Format to parse syslog messages with. Options are: Automatic|RFC3164|RFC5424|RFC6587. (default "Automatic")
  -syslog.source string
    	Run Syslog Server at this IP:Port or unix socket.
  -syslog.tcp
    	Listen on TCP for syslog messages. (default true)
  -syslog.threads int
    	Number of threads to use to process messages. (default 1)
  -syslog.udp
    	Listen on UDP for syslog messages. (default true)
  -syslog.unix
    	Listen on a Unix socket for syslog messages.
  -tag_map string
    	CSV file mapping tag ids to strings
  -tag_map_type string
    	type of mapping to use for tag values. file|null
  -tee_logs
    	Tee log messages to sink
  -threads int
    	Number of threads to run for processing
  -udrs string
    	UDR mapping file
  -v	Show version and build information
  -vpc string
    	Run VPC Flow Ingest
```

# pprof
To expose profiling endpoints, use the `-metalisten` flag. This can be used with tools such as
`go tool pprof` to capture and view the data. For example, if `ktranslate` was started with
`-metalisten :6060`:

```
go tool pprof -http :8080 http://127.0.0.1:6060/debug/pprof/profile
```

To view all available profiles, open http://localhost:6060/debug/pprof/ in your browser.

This product includes GeoLite2 data created by MaxMind, available from
<a href="https://www.maxmind.com">https://www.maxmind.com</a>.
