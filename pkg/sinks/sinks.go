package sinks

import (
	"context"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/sinks/ddog"
	"github.com/kentik/ktranslate/pkg/sinks/file"
	"github.com/kentik/ktranslate/pkg/sinks/gcloud"
	"github.com/kentik/ktranslate/pkg/sinks/gcppubsub"
	"github.com/kentik/ktranslate/pkg/sinks/http"
	"github.com/kentik/ktranslate/pkg/sinks/kafka"
	"github.com/kentik/ktranslate/pkg/sinks/kentik"
	"github.com/kentik/ktranslate/pkg/sinks/net"
	"github.com/kentik/ktranslate/pkg/sinks/nr"
	"github.com/kentik/ktranslate/pkg/sinks/nrmulti"
	"github.com/kentik/ktranslate/pkg/sinks/prom"
	"github.com/kentik/ktranslate/pkg/sinks/s3"
	"github.com/kentik/ktranslate/pkg/sinks/stdout"
)

type SinkImpl interface {
	Init(context.Context, formats.Format, kt.Compression, formats.Formatter) error
	Send(context.Context, *kt.Output)
	Close()
	HttpInfo() map[string]float64
}

type CloudObjectManager interface {
	Init(context.Context, formats.Format, kt.Compression, formats.Formatter) error
	Put(context.Context, string, []byte) error
	Get(context.Context, string) ([]byte, error)
}

type Sink string

const (
	KafkaSink     Sink = "kafka"
	StdOutSink         = "stdout"
	NewRelicSink       = "new_relic"
	KentikSink         = "kentik"
	FileSink           = "file"
	NetSink            = "net"
	HttpSink           = "http"
	SplunkSink         = "splunk"
	PromSink           = "prometheus"
	S3Sink             = "s3"
	GCloudSink         = "gcloud"
	GCPPubSub          = "gcppubsub"
	NewRelicMulti      = "new_relic_multi"
	DDogSink           = "ddog"
)

func NewSink(sink Sink, log logger.Underlying, registry go_metrics.Registry, tooBig chan int, logTee chan string, config *ktranslate.Config) (SinkImpl, error) {
	switch sink {
	case StdOutSink:
		return stdout.NewSink(log, registry, logTee)
	case FileSink:
		return file.NewSink(log, registry, config.FileSink)
	case KafkaSink:
		return kafka.NewSink(log, registry, config.KafkaSink)
	case NewRelicSink:
		return nr.NewSink(log, registry, tooBig, logTee, config.NewRelicSink)
	case DDogSink:
		return ddog.NewSink(log, registry, config.DDogSink)
	case KentikSink:
		return kentik.NewSink(log, registry, config)
	case NetSink:
		return net.NewSink(log, registry, config.NetSink)
	case HttpSink, SplunkSink:
		return http.NewSink(log, registry, string(sink), config.HTTPSink, logTee)
	case PromSink:
		return prom.NewSink(log, registry, config.PrometheusSink)
	case S3Sink:
		return s3.NewSink(log, registry, config.S3Sink)
	case GCloudSink:
		return gcloud.NewSink(log, registry, config.GCloudSink)
	case GCPPubSub:
		return gcppubsub.NewSink(log, registry, config.GCloudPubSubSink)
	case NewRelicMulti:
		return nrmulti.NewSink(log, registry, tooBig, logTee, config.NewRelicSink, config.NewRelicMultiSink)
	}
	return nil, fmt.Errorf("Unknown sink %v", sink)
}
