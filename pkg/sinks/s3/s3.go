package s3

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

type S3Sink struct {
	logger.ContextL
	Bucket   string
	client   *s3manager.Uploader
	registry go_metrics.Registry
	metrics  *S3Metric
	prefix   string
}

type S3Metric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

var (
	S3Bucket = flag.String("s3_bucket", "", "AWS S3 Bucket to write flows to")
	S3Prefix = flag.String("s3_prefix", "/kentik", "AWS S3 Object prefix")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*S3Sink, error) {
	return &S3Sink{
		registry: registry,
		Bucket:   *S3Bucket,
		prefix:   *S3Prefix,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "s3Sink"}, log),
		metrics: &S3Metric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_s3", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_s3", registry),
		},
	}, nil
}

func (s *S3Sink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	rand.Seed(time.Now().UnixNano())
	if s.Bucket == "" {
		return fmt.Errorf("Not writing to s3 -- no bucket set, use -s3_bucket flag")
	}
	sess := session.Must(session.NewSession())
	s.client = s3manager.NewUploader(sess)

	s.Infof("System connected to s3, bucket is %s", s.Bucket)

	return nil
}

func (s *S3Sink) Send(ctx context.Context, payload []byte) {
	go s.send(ctx, payload)
}

func (s *S3Sink) send(ctx context.Context, payload []byte) {
	_, err := s.client.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%d_%d", s.prefix, time.Now().Unix(), rand.Intn(100000))),
		Body:   bytes.NewBuffer(payload),
	})
	if err != nil {
		s.Errorf("Cannot close upload to s3: %v", err)
		s.metrics.DeliveryErr.Mark(1)
	} else {
		s.metrics.DeliveryWin.Mark(1)
	}
}

func (s *S3Sink) Close() {

}

func (s *S3Sink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}
