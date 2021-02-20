package gcloud

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/storage"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

type GCloudSink struct {
	logger.ContextL
	Bucket      string
	client      *storage.Client
	registry    go_metrics.Registry
	metrics     *GCloudMetric
	prefix      string
	contentType string
}

type GCloudMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

var (
	GCloudBucket      = flag.String("gcloud_bucket", "", "GCloud Storage Bucket to write flows to")
	GCloudPrefix      = flag.String("gcloud_prefix", "/kentik", "GCloud Storage object prefix")
	GCloudContentType = flag.String("gcloud_content_type", "application/json", "GCloud Storage Content Type")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*GCloudSink, error) {
	return &GCloudSink{
		registry:    registry,
		Bucket:      *GCloudBucket,
		prefix:      *GCloudPrefix,
		contentType: *GCloudContentType,
		ContextL:    logger.NewContextLFromUnderlying(logger.SContext{S: "gcloudSink"}, log),
		metrics: &GCloudMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_gcloud", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_gcloud", registry),
		},
	}, nil
}

func (s *GCloudSink) Init(ctx context.Context, format formats.Format, compression kt.Compression) error {
	rand.Seed(time.Now().UnixNano())
	if s.Bucket == "" {
		return fmt.Errorf("Not writing to gcloud -- no bucket set, use -gcloud_bucket flag")
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Cannot create gcloud client -- %v", err)
	}
	s.client = client

	s.Infof("System connected to gcloud, bucket is %s", s.Bucket)

	return nil
}

func (s *GCloudSink) Send(ctx context.Context, payload []byte) {
	go s.send(ctx, payload)
}

func (s *GCloudSink) send(ctx context.Context, payload []byte) {
	wc := s.client.Bucket(s.Bucket).Object(fmt.Sprintf("%s/%d_%d", s.prefix, time.Now().Unix(), rand.Intn(100000))).NewWriter(ctx)
	wc.ContentType = s.contentType
	if _, err := wc.Write(payload); err != nil {
		s.Errorf("Cannot upload to gcloud: %v", err)
		s.metrics.DeliveryErr.Mark(1)
	}
	if err := wc.Close(); err != nil {
		s.Errorf("Cannot close upload to gcloud: %v", err)
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
