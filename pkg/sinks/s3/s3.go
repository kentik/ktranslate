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
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	s3Bucket                     string
	s3Prefix                     string
	flushDurSec                  int
	s3assumeRoleARN              string
	s3Region                     string
	ec2InstanceProfile           bool
	ec2assumeRoleIntervalSeconds int
)

// var wg sync.WaitGroup

func init() {
	flag.StringVar(&s3Bucket, "s3_bucket", "", "AWS S3 Bucket to write flows to")
	flag.StringVar(&s3Prefix, "s3_prefix", "/kentik", "AWS S3 Object prefix")
	flag.IntVar(&flushDurSec, "s3_flush_sec", 60, "Create a new output file every this many seconds")
	flag.StringVar(&s3assumeRoleARN, "s3_assume_role_arn", "", "AWS assume role ARN which has permissions to write to S3 bucket")
	flag.StringVar(&s3Region, "s3_region", "us-east-1", "S3 Bucket region where S3 bucket is created")
	flag.BoolVar(&ec2InstanceProfile, "ec2_instance_profile", false, "EC2 Instance Profile")
	flag.IntVar(&ec2assumeRoleIntervalSeconds, "assume_role_interval_seconds", 900, "Refresh credentials after this many seconds")
}

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
	config   *ktranslate.S3SinkConfig
}

type S3Metric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.S3SinkConfig) (*S3Sink, error) {
	rand.Seed(time.Now().UnixNano())
	return &S3Sink{
		registry: registry,
		Bucket:   cfg.Bucket,
		prefix:   cfg.Prefix,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "s3Sink"}, log),
		metrics: &S3Metric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_s3", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_s3", registry),
		},
		buf:    bytes.NewBuffer(make([]byte, 0)),
		config: cfg,
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
	if s.config.AssumeRoleARN != "" {
		if err := s.get_tmp_credentials(ctx); err != nil {
			return err
		}
	} else {
		sess := session.Must(session.NewSession())
		s.Infof("Session is created using default settings")
		s.client = s3manager.NewUploader(sess)
	}

	switch compression {
	case kt.CompressionNone, kt.CompressionNull:
		s.suffix = ""
	case kt.CompressionGzip:
		s.suffix = ".gz"
	default:
		s.suffix = "." + string(compression)
	}

	s.Infof("System connected to s3, bucket is %s, dumping on %v", s.Bucket, time.Duration(s.config.FlushIntervalSeconds)*time.Second)

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

func (s *S3Sink) get_tmp_credentials(ctx context.Context) error {
	// First, make sure we can get some credentials:
	creds, err := s.tmp_credentials(ctx)
	if err != nil {
		return err
	}
	s.client = creds

	// Now loop forever getting new creds.
	newCredTick := time.NewTicker(time.Duration(s.config.AssumeRoleIntervalSeconds) * time.Second)
	defer newCredTick.Stop()

	for {
		select {
		case _ = <-newCredTick.C:
			creds, err := s.tmp_credentials(ctx)
			if err != nil {
				// In this case, keep trying while not replacing the old creds, maybe next time around will work.
				s.Errorf("Cannot get new AWS creds: %v", err)
			} else {
				s.client = creds
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func (s *S3Sink) tmp_credentials(ctx context.Context) (*s3manager.Uploader, error) {
	if s.config.EC2InstanceProfile && s.config.AssumeRoleARN != "" {

		svc := ec2metadata.New(session.Must(session.NewSession()))
		ec2_role_creds := ec2rolecreds.NewCredentialsWithClient(svc)
		sess_tmp := session.Must(
			session.NewSession(&aws.Config{
				Region:      aws.String(s.config.Region),
				Credentials: ec2_role_creds,
			}),
		)
		_, err_role := ec2_role_creds.Get()
		if err_role != nil {
			s.Errorf("Not able to retrieve credentials via Instance Profile. ARN: %v. ERROR: %v", s.config.AssumeRoleARN, err_role)
			return nil, err_role
		}

		creds := stscreds.NewCredentials(sess_tmp, s.config.AssumeRoleARN)
		_, err_creds := creds.Get()
		if err_creds != nil {
			s.Errorf("Assume Role ARN doesn't work. ARN: %v. ERROR: %v", s.config.AssumeRoleARN, err_creds)
			return nil, err_creds
		}

		// Creating a new session from assume role
		sess, err := session.NewSession(
			&aws.Config{
				Region:      aws.String(s.config.Region),
				Credentials: creds,
			},
		)
		if err != nil {
			s.Errorf("Session is not created ERROR: %v", err)
			return nil, err
		} else {
			s.Infof("Session is created using assume role based on EC2 Instance Profile")
		}

		return s3manager.NewUploader(sess), nil

	} else {
		// Getting credentials from assume role ARN
		sess_tmp := session.Must(
			session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}),
		)
		creds := stscreds.NewCredentials(sess_tmp, s.config.AssumeRoleARN)
		_, err := creds.Get()
		if err != nil {
			s.Errorf("Assume Role ARN doesn't work. ARN: %v, %v", s.config.AssumeRoleARN, err)
			return nil, err
		}
		// Creating a new session from assume role
		sess, err := session.NewSession(
			&aws.Config{
				Region:      aws.String(s.config.Region),
				Credentials: creds,
			},
		)
		if err != nil {
			s.Errorf("Session is not created with region %v, %v", s.config.Region, err)
			return nil, err
		} else {
			s.Infof("Session is created using assume role via shared configuration")
		}

		return s3manager.NewUploader(sess), nil
	}
}
