package sinks

import (
	"context"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/sinks/file"
	"github.com/kentik/ktranslate/pkg/sinks/gcloud"
	"github.com/kentik/ktranslate/pkg/sinks/gcppubsub"
	"github.com/kentik/ktranslate/pkg/sinks/http"
	"github.com/kentik/ktranslate/pkg/sinks/kafka"
	"github.com/kentik/ktranslate/pkg/sinks/kentik"
	"github.com/kentik/ktranslate/pkg/sinks/net"
	"github.com/kentik/ktranslate/pkg/sinks/nr"
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

type Sink string

const (
	KafkaSink    Sink = "kafka"
	StdOutSink        = "stdout"
	NewRelicSink      = "new_relic"
	KentikSink        = "kentik"
	FileSink          = "file"
	NetSink           = "net"
	HttpSink          = "http"
	SplunkSink        = "splunk"
	PromSink          = "prometheus"
	S3Sink            = "s3"
	GCloudSink        = "gcloud"
	GCPPubSub         = "gcppubsub"
)

func NewSink(sink Sink, log logger.Underlying, registry go_metrics.Registry, tooBig chan int, conf *kt.KentikConfig, logTee chan string) (SinkImpl, error) {
	switch sink {
	case StdOutSink:
		return stdout.NewSink(log, registry, logTee)
	case FileSink:
		return file.NewSink(log, registry)
	case KafkaSink:
		return kafka.NewSink(log, registry)
	case NewRelicSink:
		return nr.NewSink(log, registry, tooBig, logTee)
	case KentikSink:
		return kentik.NewSink(log, registry, conf)
	case NetSink:
		return net.NewSink(log, registry)
	case HttpSink, SplunkSink:
		return http.NewSink(log, registry, string(sink))
	case PromSink:
		return prom.NewSink(log, registry)
	case S3Sink:
		return s3.NewSink(log, registry)
	case GCloudSink:
		return gcloud.NewSink(log, registry)
	case GCPPubSub:
		return gcppubsub.NewSink(log, registry)
	}
	return nil, fmt.Errorf("Unknown sink %v", sink)
}
