module github.com/kentik/ktranslate

go 1.21

toolchain go1.21.6

require (
	cloud.google.com/go/pubsub v1.36.1
	cloud.google.com/go/secretmanager v1.11.5
	cloud.google.com/go/storage v1.36.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.6.0
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets v0.7.1
	github.com/DataDog/datadog-api-client-go/v2 v2.5.0
	github.com/agoda-com/opentelemetry-go/otelslog v0.1.1
	github.com/agoda-com/opentelemetry-logs-go v0.5.0
	github.com/aristanetworks/goeapi v0.6.0
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.44.302
	github.com/dchest/siphash v1.2.2
	github.com/elliotchance/orderedmap v1.4.0
	github.com/fsnotify/fsnotify v1.6.0
	github.com/go-openapi/runtime v0.26.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/snappy v0.0.4
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/gosnmp/gosnmp v1.35.0
	github.com/honeycombio/libhoney-go v1.15.6
	github.com/imdario/mergo v0.3.16
	github.com/influxdata/line-protocol/v2 v2.2.1
	github.com/json-iterator/go v1.1.12
	github.com/judwhite/go-svc v1.2.1
	github.com/kentik/api-schema-public v0.0.0-20220322181339-896729e59945
	github.com/kentik/dashboard-api-golang v0.0.0-20240305065153-1c88e8477770
	github.com/kentik/go-metrics v0.0.0-20200109025407-4bfd4a9b42c5
	github.com/kentik/go-syslog v0.0.0-20240302065219-cea330e269ea
	github.com/kentik/patricia v0.0.0-20210909164817-21603333b70e
	github.com/liamg/furious v0.0.0-20191231090757-c295c872d6c1
	github.com/linkedin/goavro/v2 v2.10.1
	github.com/montanaflynn/stats v0.7.0
	github.com/netsampler/goflow2/v2 v2.1.3
	github.com/oschwald/geoip2-golang v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus-community/pro-bing v0.3.0
	github.com/prometheus/client_golang v1.19.0
	github.com/prometheus/prometheus v0.46.0
	github.com/sasha-s/go-hll v0.0.0-20180522065212-c6eb27aee351
	github.com/segmentio/kafka-go v0.4.23
	github.com/stretchr/testify v1.9.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/xitongsys/parquet-go v1.6.2
	github.com/xitongsys/parquet-go-source v0.0.0-20240122235623-d6294584ab18
	go.opentelemetry.io/otel v1.25.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.25.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.25.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.25.0
	go.opentelemetry.io/otel/metric v1.25.0
	go.opentelemetry.io/otel/sdk v1.25.0
	go.opentelemetry.io/otel/sdk/metric v1.25.0
	go.starlark.net v0.0.0-20220926145019-14b050677505
	google.golang.org/genproto v0.0.0-20240227224415-6ceb2ff114de
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	zombiezen.com/go/capnproto2 v2.18.2+incompatible
)

require (
	cloud.google.com/go v0.112.0 // indirect
	cloud.google.com/go/compute v1.24.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.6 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.11.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.8.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal v0.5.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/apache/arrow/go/arrow v0.0.0-20200730104253-651201b0f516 // indirect
	github.com/apache/thrift v0.14.2 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-bits v0.0.0-20180113010104-bd8a69a71dc2 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/facebookgo/limitgroup v0.0.0-20150612190941-6abd8d71ec01 // indirect
	github.com/facebookgo/muster v0.0.0-20150708232844-fd3d7953fd52 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.4 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/strfmt v0.21.7 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-openapi/validate v0.22.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gopacket v1.1.17 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/libp2p/go-reuseport v0.4.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostlygeek/arp v0.0.0-20170424181311-541a2129847a // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/onsi/ginkgo v1.16.2 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/oschwald/maxminddb-golang v1.11.0 // indirect
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.4 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.mongodb.org/mongo-driver v1.12.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.47.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0 // indirect
	go.opentelemetry.io/otel/trace v1.25.0 // indirect
	go.opentelemetry.io/proto/otlp v1.2.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/oauth2 v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.162.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415180920-8c6c420018be // indirect
	gopkg.in/alexcesaro/statsd.v2 v2.0.0 // indirect
)
