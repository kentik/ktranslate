package gcloud

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	gcloudBucket      string
	gcloudPrefix      string
	gcloudContentType string
	flushDurSec       int
)

func init() {
	flag.StringVar(&gcloudBucket, "gcloud_bucket", "", "GCloud Storage Bucket to write flows to")
	flag.StringVar(&gcloudPrefix, "gcloud_prefix", "/kentik", "GCloud Storage object prefix")
	flag.StringVar(&gcloudContentType, "gcloud_content_type", "application/json", "GCloud Storage Content Type")
	flag.IntVar(&flushDurSec, "gcloud_flush_sec", 60, "Create a new output file every this many seconds")
}

type GCloudSink struct {
	logger.ContextL
	Bucket      string
	client      *storage.Client
	registry    go_metrics.Registry
	metrics     *GCloudMetric
	prefix      string
	suffix      string
	contentType string
	buf         *bytes.Buffer
	mux         sync.RWMutex
	config      *ktranslate.GCloudSinkConfig
}

type GCloudMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.GCloudSinkConfig) (*GCloudSink, error) {
	return &GCloudSink{
		registry:    registry,
		Bucket:      cfg.Bucket,
		prefix:      cfg.Prefix,
		contentType: cfg.ContentType,
		ContextL:    logger.NewContextLFromUnderlying(logger.SContext{S: "gcloudSink"}, log),
		metrics: &GCloudMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_gcloud", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_gcloud", registry),
		},
		buf:    bytes.NewBuffer(make([]byte, 0)),
		config: cfg,
	}, nil
}

func (s *GCloudSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	rand.Seed(time.Now().UnixNano())
	if s.Bucket == "" {
		return fmt.Errorf("No bucket was set for gcloud. Use the -gcloud_bucket option to set one up.")
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("There was an error when creating the gcloud client: %v.", err)
	}
	s.client = client

	switch compression {
	case kt.CompressionNone, kt.CompressionNull:
		s.suffix = ""
	case kt.CompressionGzip:
		s.suffix = ".gz"
	default:
		s.suffix = "." + string(compression)
	}

	s.Infof("System connected to gcloud, bucket is %s, dumping on %v", s.Bucket, time.Duration(s.config.FlushIntervalSeconds)*time.Second)

	go func() {
		dumpTick := time.NewTicker(time.Duration(s.config.FlushIntervalSeconds) * time.Second)
		defer dumpTick.Stop()

		for {
			select {
			case _ = <-dumpTick.C:
				s.mux.Lock()
				if s.buf.Len() == 0 {
					s.mux.Unlock()
					continue
				}
				ob := s.buf
				s.buf = bytes.NewBuffer(make([]byte, 0, ob.Len()))
				go s.send(ctx, ob.Bytes())
				s.mux.Unlock()

			case <-ctx.Done():
				s.Infof("gcloudSink Done")
				return
			}
		}
	}()

	return nil
}

func (s *GCloudSink) Send(ctx context.Context, payload *kt.Output) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.buf.Write(payload.Body)
}

func (s *GCloudSink) getName() string {
	now := time.Now()
	return fmt.Sprintf("%s/%s/%d_%d%s", s.prefix, now.Format("2006/01/02/15/04"), now.Unix(), rand.Intn(100000), s.suffix)
}

func (s *GCloudSink) send(ctx context.Context, payload []byte) {
	wc := s.client.Bucket(s.Bucket).Object(s.getName()).NewWriter(ctx)
	wc.ContentType = s.contentType
	if _, err := wc.Write(payload); err != nil {
		s.Errorf("There was an error when uploading to gcloud: %v.", err)
		s.metrics.DeliveryErr.Mark(1)
	}
	if err := wc.Close(); err != nil {
		s.Errorf("There was an error when uploading to gcloud: %v.", err)
		s.metrics.DeliveryErr.Mark(1)
	}
	s.metrics.DeliveryWin.Mark(1)
}

func (s *GCloudSink) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

func (s *GCloudSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}
