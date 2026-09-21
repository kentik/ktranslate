# KTranslate - Kentik data to the world

Listen for a feed of data to or from Kentik and pass on in a common form. Supports rollups and filtering as well.

See the [Wiki](https://github.com/kentik/ktranslate/wiki) for more details. Come visit the [Discord](https://discord.gg/XGDNRj528C) if you have any questions, need any assistance, or want to talk about the development of ktranslate.

# Build:

make && make test

# Build Docker Image:

To build and use a Docker image, you must specify `MAXMIND_LICENSE_KEY` and `YOUR_ACCOUNT_ID` as build args:

```bash
docker build --build-arg YOUR_ACCOUNT_ID=xxxxx --build-arg MAXMIND_LICENSE_KEY=xxxxx -t ktranslate:v2 .
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
  -assume_role_or_instance_profile_interval_seconds int
    	Refresh credentials of Assume Role or Instance Profile (whichever is earliest) after this many seconds (default 900)
  -aws_lambda
    	Run as a AWS Lambda function
  -aws_local_file string
    	If set, process this local file and exit
  -aws_regions string
    	CSV list of region to run in. Will look for metadata in all regions, run SQS in first region. (default "us-east-1")
  -compression string
    	compression algo to use (none|gzip|snappy|deflate|null) (default "none")
  -config string
    	path to ktranslate config
  -config_provider string
    	Implementation of which provider controls the config process. Can be one of (new_relic,local)
  -dns string
    	Resolve IPs at this ip:port
  -ec2_instance_profile
    	EC2 Instance Profile
  -elastic.action string
    	Use this action when sending to elastic. (default "index")
  -enricher string
    	Send data to this http url for enrichment.
  -filters value
    	Any filters to use. Format: type dimension operator value
  -flow_only
    	If true, don't poll snmp devices.
  -format string
    	Format to convert kflow to: (json|flat_json|avro|netflow|influx|carbon|prometheus|new_relic|new_relic_metric|elasticsearch|kflow|otel|snmp) (default "flat_json")
  -format_metric string
    	Format to convert metrics to: (json|avro|netflow|influx|prometheus|new_relic|new_relic_metric|elasticsearch|kflow)
  -format_rollup string
    	Format to convert rollups to: (json|avro|netflow|influx|prometheus|new_relic|new_relic_metric|elasticsearch|kflow)
  -gcp.project string
    	Google ProjectID to listen for flows on
  -gcp.sample float
    	Sample rate of the vpc export (as defined in the VPC setup) (default 1)
  -gcp.sub string
    	Google Sub to listen for flows on
  -generate-config
    	generate ktranslate config and exit
  -geo string
    	Geo mapping file
  -geo_city_map string
    	CSV file mapping geo city ids to strings
  -geo_region_map string
    	CSV file mapping geo region ids to strings
  -http.remote_ip string
    	If set, ignore actual remote IP and use this for device mapping.
  -http.source
    	Listen for content sent via http.
  -http_header value
    	Any custom http headers to set on outbound requests
  -http_insecure
    	Allow insecure urls.
  -http_log_url string
    	URL to post logs to (default "http://localhost:8088/services/collector/event")
  -http_timeout_sec int
    	Timeout each request after this long. (default 30)
  -http_url string
    	URL to post to (default "http://localhost:8086/write?db=kentik")
  -iam_role string
    	IAM Role to use for processing flow
  -influxdb_measurement_prefix string
    	Prefix metric names with this
  -influxdb_namespace_token string
    	Use this token to seperate namespaces (default ":")
  -info_collector
    	Also send stats about this collector
  -input_threads int
    	Number of threads to run for input processing (default 1)
  -kentik_email string
    	Kentik email to use for API calls
  -kentik_plan int
    	Kentik plan id to use for creating devices
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
    	Dynamically grow threads up to this number (default 1)
  -metalisten string
    	HTTP interface and port to bind on (default "localhost:0")
  -metrics string
    	Metrics Configuration. none|syslog|stderr|graphite:127.0.0.1:2003 (default "none")
  -net_protocol string
    	Use this protocol for writing data (udp|tcp|unix) (default "udp")
  -net_server string
    	Write flows seen to this address (host and port). Comma seperate to send to multiple servers.
  -netflow_version string
    	Version of netflow to produce: (netflow9|ipfix) (default "ipfix")
  -nf.addr string
    	Sflow/NetFlow/IPFIX listening address (default "0.0.0.0")
  -nf.mapping string
    	Configuration file for custom netflow mappings
  -nf.message.fields string
    	The list of fields to include in flow messages. Can be any of Type,TimeReceived,SequenceNum,SamplingRate,FlowDirection,SamplerAddress,TimeFlowStart,TimeFlowEnd,Bytes,Packets,SrcAddr,DstAddr,Etype,Proto,SrcPort,DstPort,InIf,OutIf,SrcMac,DstMac,SrcVlan,DstVlan,VlanId,IPTos,ForwardingStatus,IPTTL,TCPFlags,IcmpType,IcmpCode,IPv6FlowLabel,FragmentId,FragmentOffset,SrcAS,DstAS,NextHop,NextHopAS,SrcNet,DstNet,MPLSCount (default "TimeReceived,SamplingRate,Bytes,Packets,SrcAddr,DstAddr,Proto,SrcPort,DstPort,InIf,OutIf,SrcVlan,DstVlan,TCPFlags,SrcAS,DstAS,Type,SamplerAddress,FlowDirection")
  -nf.port int
    	Sflow/NetFlow/IPFIX listening port (default 9995)
  -nf.prom.listen string
    	Run a promethues metrics collector here
  -nf.queuesize int
    	How big of a queue to hold for incomming flow packets. (default 10000)
  -nf.reuserport
    	Enable so_reuseport for Sflow/NetFlow/IPFIX
  -nf.source string
    	Run NetFlow Ingest Directly. Valid values here are netflow5|netflow9|ipfix|sflow|nbar|asa|pan|auto
  -nf.workers int
    	Number of workers per flow collector (default 2)
  -nr_account_id string
    	If set, sends flow to New Relic
  -nr_check_json
    	Verify body is valid json before sending on
  -nr_estimate_only
    	If true, record size of inputs to NR but don't actually send anything
  -nr_region string
    	NR Region to use. US|EU|GOV|JP. If not set, this is auto-detected from the NEW_RELIC_API_KEY license key prefix (EU/JP only; unrecognized keys default to US).
  -olly_dataset string
    	Olly dataset name
  -olly_write_key string
    	Olly dataset name
  -otel.endpoint string
    	Send data to this endpoint.
  -otel.no_block
    	If set, drop metrics when the sending chan is full.
  -otel.protocol string
    	Send data using this protocol. (grpc,http,https,stdout) (default "stdout")
  -otel.root_ca string
    	Load TLS root CA from file.
  -otel.tls_cert string
    	Load TLS client cert from file.
  -otel.tls_key string
    	Load TLS client key from file.
  -prom_seen int
    	Number of flows needed inbound before we start writting to the collector (default 4)
  -redis.addr string
    	Where to connect to redis. (default "localhost:6379")
  -redis.db int
    	Use this redis DB.
  -redis.key_prefix string
    	Use this key prefix.
  -redis.password string
    	Password for redis
  -redis.ttl.sec int
    	Expire measurements if they are not refreshed within this number of sec. (default 60)
  -rollup_and_alpha
    	Send both rollups and alpha inputs to sinks
  -rollup_interval int
    	Export timer for rollups in seconds
  -rollup_keep_undefined
    	If set, mark undefined values with the string undefined.
  -rollup_key_join string
    	Token to use to join dimension keys together (default "^")
  -rollup_top_k int
    	Export only these top values (default 10)
  -rollups value
    	Any rollups to use. Format: type, name, metric, dimension 1, dimension 2, ..., dimension n: sum,bytes,in_bytes,dst_addr
  -s3_assume_role_arn string
    	AWS assume role ARN which has permissions to write to S3 bucket
  -s3_bucket string
    	AWS S3 Bucket to write flows to
  -s3_endpoint string
    	S3 Endpoint
  -s3_flush_sec int
    	Create a new output file every this many seconds (default 60)
  -s3_prefix string
    	AWS S3 Object prefix (default "/kentik")
  -s3_region string
    	S3 Bucket region where S3 bucket is created (default "us-east-1")
  -s3_signing_region string
    	S3 endpoint signing region
  -sample_rate int
    	Sampling rate to use. 1 -> 1:1 sampling, 2 -> 1:2 sampling and so on. (default 1)
  -service_name string
    	Service identifier
  -sinks string
    	List of sinks to send data to. Options: (stdout|new_relic|new_relic_multi|otel|http|net) (default "stdout")
  -snmp string
    	yaml file containing snmp config to use
  -snmp.format.conf string
    	Parse this file for the snmp format option. Same format as -snmp flag.
  -snmp_discovery
    	If true, try to discover snmp devices on this network as configured.
  -snmp_discovery_min int
    	If set, run snmp discovery on this interval (in minutes).
  -snmp_discovery_on_start
    	If set, run snmp discovery on application start.
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
  -snmp_validate
    	If true, validate mib profiles and exit.
  -snmp_walk_file string
    	If set, use the walk file instead of polling.
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
    	Log to stdout
  -stitch.buffer.len int
    	How large of a buffer of flows to try and stitch together. (default 10000)
  -stitch.enable
    	Turn on flow stitching.
  -syslog.format string
    	Format to parse syslog messages with. Options are: Automatic|RFC3164|RFC5424|RFC6587|NoFormat. (default "Automatic")
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
  -tee_flow string
    	If set, tee flow to another ktranslate instance here.
  -tee_logs
    	Tee log messages to sink
  -threads int
    	Number of threads to run for processing (default 1)
  -udrs string
    	UDR mapping file
  -vpc string
    	Run VPC Flow Ingest
```


# Further documentation

The flag list above is a snapshot. `ktranslate -h` on a current binary is authoritative for names.

Newer operator guides live on the [wiki](https://github.com/kentik/ktranslate/wiki):

* [Sending data with OTLP](https://github.com/kentik/ktranslate/wiki/Sending-Data-with-OTLP) (`-format=otel`, Grafana Cloud / Alloy)
* [NetBox discovery](https://github.com/kentik/ktranslate/wiki/NetBox-Discovery)
* [Sinks, formats, and rollups](https://github.com/kentik/ktranslate/wiki/Sinks-Formats-and-Rollups) (Kafka SASL/TLS, S3 endpoints, parquet, rollup flags)
* [Advanced Configuration](https://github.com/kentik/ktranslate/wiki/Advanced-Ktranslate-Configuration) (`snmp-base.yaml`, profile git URL)

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
