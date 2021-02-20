# KTranslate - Kentik data to the world 
 
Listen for a feed of data to or from Kentik and pass on in a common form. Supports rollups and filtering as well. 

Input format: kflow\
Output formats: json, avro, ipfix, prometheus\
Output destinations: stdout, file, New Relic, Splunk, Elastic, Prometheus, Kafka, Kentik

Run with:

`docker run -p 8082:8082 ktranslate`

Flags:

```
Usage of ./bin/ktranslate:
  -asn4 string
    	Asn ipv6 mapping file
  -asn6 string
    	Asn ipv6 mapping file
  -bootstrap.servers string
    	bootstrap.servers
  -city string
    	City mapping file
  -compression string
    	compression algo to use (none|gzip|snappy|deflate|null) (default "none")
  -dns string
    	Resolve IPs at this ip:port
  -file_out string
    	Write flows seen to log to this directory if set (default "./")
  -filters value
    	Any filters to use. Format: type dimension operator value
  -format string
    	Format to convert kflow to: (json|avro|netflow|influx|prometheus) (default "json")
  -geo string
    	Geo mapping file
  -healthcheck string
    	Bind to this interface to allow healthchecks
  -http_header value
    	Any custom http headers to set on outbound requests
  -http_url string
    	URL to post to (default "http://localhost:8086/write?db=kentik")
  -info_collector
    	Also send stats about this collector
  -interfaces string
    	Interface mapping file
  -kafka.debug string
    	debug contexts to enable for kafka
  -kafka_topic string
    	kafka topic to produce on
  -kentik_email string
    	Email to use for sending flow on to Kentik
  -kentik_url string
    	URL to use for sending flow on to Kentik (default "https://flow.kentik.com/chf")
  -listen string
    	IP:Port to listen on (default "127.0.0.1:8081")
  -log_level string
    	Logging Level (default "debug")
  -mapping string
    	Mapping file to use for enums (default "config.json")
  -max_flows_per_message int
    	Max number of flows to put in each emitted message (default 10000)
  -max_sql_conns int
    	Max concurrent SQL connections (default 16)
  -measurement string
    	Measurement to use for rollups. (default "kflow")
  -metalisten string
    	HTTP port to bind on (default "localhost:0")
  -metrics string
    	Metrics Configuration. none|syslog|stderr|graphite:127.0.0.1:2003 (default "syslog")
  -net_protocol string
    	Use this protocol for writing data (udp|tcp|unix) (default "udp")
  -net_server string
    	Write flows seen to this address (host and port)
  -netflow_version string
    	Version of netflow to produce: (netflow9|ipfix) (default "ipfix")
  -nr_account_id string
    	If set, sends flow to New Relic
  -nr_url string
    	URL to use to send into NR (default "https://insights-collector.newrelic.com/v1/accounts/%s/events")
  -olly_dataset string
    	Olly dataset name
  -olly_write_key string
    	Olly dataset name
  -prom_listen string
    	Bind to listen for prometheus requests on. (default ":8082")
  -region string
    	Region mapping file
  -rollup_and_alpha
    	Send both rollups and alpha inputs to sinks
  -rollup_export int
    	Export timer for rollups in seconds
  -rollup_key_join string
    	Token to use to join dimension keys together (default "^")
  -rollup_top_k int
    	Export only these top values (default 10)
  -rollups value
    	Any rollups to use. Format: type, metric, dimension 1, dimension 2, ..., dimension n: sum,in_bytes,dst_addr
  -sample_rate int
    	Sampling rate to use. 1 -> 1:1 sampling, 2 -> 1:2 sampling and so on. (default 1)
  -sasl.kerberos.keytab string
    	sasl.kerberos.keytab
  -sasl.kerberos.kinit.cmd string
    	sasl.kerberos.kinit.cmd (default "kinit -R -t \"%{sasl.kerberos.keytab}\" -k %{sasl.kerberos.principal} || kinit -t \"%{sasl.kerberos.keytab}\" -k %{sasl.kerberos.principal}")
  -sasl.kerberos.principal string
    	sasl.kerberos.principal
  -sasl.kerberos.service.name string
    	sasl.kerberos.service.name (default "kafka")
  -sasl.mechanism string
    	sasl.mechanism
  -security.protocol string
    	security.protocol
  -service_name string
    	Service identifier (default "ktranslate")
  -sinks string
    	List of sinks to send data to. Options: (kafka|stdout|new_relic|kentik|net|http|splunk|prometheus) (default "stdout")
  -ssl.ca.location string
    	ssl.ca.location
  -stdout
    	Log to stdout (default true)
  -tag_lookup string
    	Tag service port to run lookups on
  -threads int
    	Number of threads to run for processing
  -v	Show version and build information
```
