package s3

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"math/rand"
	"sync"
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
	suffix   string
	buf      *bytes.Buffer
	mux      sync.RWMutex
}

type S3Metric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

var (
	S3Bucket    = flag.String("s3_bucket", "", "AWS S3 Bucket to write flows to")
	S3Prefix    = flag.String("s3_prefix", "/kentik", "AWS S3 Object prefix")
	FlushDurSec = flag.Int("s3_flush_sec", 60, "Create a new output file every this many seconds")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*S3Sink, error) {
	rand.Seed(time.Now().UnixNano())
	return &S3Sink{
		registry: registry,
		Bucket:   *S3Bucket,
		prefix:   *S3Prefix,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "s3Sink"}, log),
		metrics: &S3Metric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_s3", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_s3", registry),
		},
		buf: bytes.NewBuffer(make([]byte, 0)),
	}, nil
}

func (s *S3Sink) getName() string {
	now := time.Now()
	return fmt.Sprintf("%s/%s/%d_%d%s", s.prefix, now.Format("2006/01/02/15/04"), now.Unix(), rand.Intn(100000), s.suffix)
}

func (s *S3Sink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	rand.Seed(time.Now().UnixNano())
	if s.Bucket == "" {
		return fmt.Errorf("Not writing to s3 -- no bucket set, use -s3_bucket flag")
	}
	sess := session.Must(session.NewSession())
	s.client = s3manager.NewUploader(sess)

	switch compression {
	case kt.CompressionNone, kt.CompressionNull:
		s.suffix = ""
	case kt.CompressionGzip:
		s.suffix = ".gz"
	default:
		s.suffix = "." + string(compression)
	}

	s.Infof("System connected to s3, bucket is %s, dumping on %v", s.Bucket, time.Duration(*FlushDurSec)*time.Second)

	go func() {
		dumpTick := time.NewTicker(time.Duration(*FlushDurSec) * time.Second)
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
				s.Infof("s3Sink Done")
				return
			}
		}
	}()

	return nil
}

func (s *S3Sink) Send(ctx context.Context, payload *kt.Output) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.buf.Write(payload.Body)
}

func (s *S3Sink) send(ctx context.Context, payload []byte) {
	_, err := s.client.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Body:   bytes.NewBuffer(payload),
		Key:    aws.String(s.getName()),
	})
	if err != nil {
		s.Errorf("Cannot upload to s3: %v", err)
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
