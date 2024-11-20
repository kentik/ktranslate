package s3

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	s3Bucket                                      string
	s3Prefix                                      string
	flushDurSec                                   int
	s3assumeRoleARN                               string
	s3Region                                      string
	ec2InstanceProfile                            bool
	ec2assumeRoleOrInstanceProfileIntervalSeconds int
	checkDangling                                 bool

	TimeCheckDangling = time.Duration(1) * time.Hour
)

// var wg sync.WaitGroup

func init() {
	flag.StringVar(&s3Bucket, "s3_bucket", "", "AWS S3 Bucket to write flows to")
	flag.StringVar(&s3Prefix, "s3_prefix", "/kentik", "AWS S3 Object prefix")
	flag.IntVar(&flushDurSec, "s3_flush_sec", 60, "Create a new output file every this many seconds")
	flag.StringVar(&s3assumeRoleARN, "s3_assume_role_arn", "", "AWS assume role ARN which has permissions to write to S3 bucket")
	flag.StringVar(&s3Region, "s3_region", "us-east-1", "S3 Bucket region where S3 bucket is created")
	flag.BoolVar(&ec2InstanceProfile, "ec2_instance_profile", false, "EC2 Instance Profile")
	flag.IntVar(&ec2assumeRoleOrInstanceProfileIntervalSeconds, "assume_role_or_instance_profile_interval_seconds", 900, "Refresh credentials of Assume Role or Instance Profile (whichever is earliest) after this many seconds")
	flag.BoolVar(&checkDangling, "s3_check_dangling", false, "Check every hour for dangling folders if enabled.")
}

type S3Sink struct {
	logger.ContextL
	Bucket   string
	registry go_metrics.Registry
	metrics  *S3Metric
	prefix   string
	suffix   string
	buf      *bytes.Buffer
	mux      sync.RWMutex
	config   *ktranslate.S3SinkConfig
	client   *s3.Client
	ul       *s3manager.Uploader
	dl       *s3manager.Downloader
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

	if s.config.EC2InstanceProfile && s.config.AssumeRoleARN == "" {
		appCreds := aws.NewCredentialsCache(ec2rolecreds.New())
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(s.config.Region),
			config.WithCredentialsProvider(appCreds),
		)
		if err != nil {
			return err
		}

		value, err_role := cfg.Credentials.Retrieve(ctx)
		s.Infof("Credentials %v: ", value)
		if err_role != nil {
			s.Errorf("Not able to retrieve credentials via Instance Profile. ARN: %v. ERROR: %v", s.config.AssumeRoleARN, err_role)
			return err_role
		} else {
			s.Infof("Session is created using EC2 Instance Profile")
		}

		s.client = s3.NewFromConfig(cfg)
		s.ul = s3manager.NewUploader(s.client)
		s.dl = s3manager.NewDownloader(s.client)

	} else if s.config.AssumeRoleARN != "" || s.config.EC2InstanceProfile {
		if err := s.get_tmp_credentials(ctx); err != nil {
			return err
		}
	} else {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		s.Infof("Session is created using default settings")
		s.client = s3.NewFromConfig(cfg)
		s.ul = s3manager.NewUploader(s.client)
		s.dl = s3manager.NewDownloader(s.client)
	}

	if format == formats.FORMAT_PARQUET {
		s.suffix = ".parquet"
	} else {
		switch compression {
		case kt.CompressionNone, kt.CompressionNull:
			s.suffix = ""
		case kt.CompressionGzip:
			s.suffix = ".gz"
		default:
			s.suffix = "." + string(compression)
		}
	}

	s.Infof("System connected to s3, bucket is %s, dumping on %v", s.Bucket, time.Duration(s.config.FlushIntervalSeconds)*time.Second)

	if s.config.CheckDangling {
		go s.checkForDangling(ctx)
	}

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
	// In the un-buffered case, write this out right away.
	if payload.NoBuffer && len(payload.Body) > 0 {
		go s.send(ctx, payload.Body)
		return
	}

	// Else queue up for more efficient sending.
	s.mux.Lock()
	defer s.mux.Unlock()
	s.buf.Write(payload.Body)
}

func (s *S3Sink) Get(ctx context.Context, path string) ([]byte, error) {
	var bufb bytes.Buffer
	buf := s3manager.NewWriteAtBuffer(bufb.Bytes())
	size, err := s.dl.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes()[0:size], nil
}

func (s *S3Sink) Put(ctx context.Context, path string, data []byte) error {
	_, err := s.ul.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Body:   bytes.NewBuffer(data),
		Key:    aws.String(path),
	})
	return err
}

func (s *S3Sink) send(ctx context.Context, payload []byte) {
	_, err := s.ul.Upload(ctx, &s3.PutObjectInput{
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
	ulc, dlc, err := s.tmp_credentials(ctx)
	if err != nil {
		return err
	}
	s.ul = ulc
	s.dl = dlc

	// Now loop forever getting new creds.
	go func() {
		newCredTick := time.NewTicker(time.Duration(s.config.AssumeRoleOrInstanceProfileIntervalSeconds) * time.Second)
		defer newCredTick.Stop()

		for {
			select {
			case _ = <-newCredTick.C:
				ulc, dlc, err := s.tmp_credentials(ctx)
				if err != nil {
					// In this case, keep trying while not replacing the old creds, maybe next time around will work.
					s.Errorf("Cannot get new AWS creds: %v", err)
				} else {
					s.ul = ulc
					s.dl = dlc
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (s *S3Sink) tmp_credentials(ctx context.Context) (*s3manager.Uploader, *s3manager.Downloader, error) {

	if s.config.EC2InstanceProfile && s.config.AssumeRoleARN != "" {
		cache := aws.NewCredentialsCache(ec2rolecreds.New())
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(s.config.Region),
			config.WithCredentialsProvider(cache),
		)
		if err != nil {
			return nil, nil, err
		}

		stsc := sts.NewFromConfig(cfg)
		appCreds := stscreds.NewAssumeRoleProvider(stsc, s.config.AssumeRoleARN)
		cfgRole, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(s.config.Region),
			config.WithCredentialsProvider(appCreds),
		)
		if err != nil {
			return nil, nil, err
		}

		value, err_role := cfgRole.Credentials.Retrieve(ctx)
		s.Infof("Credentials %v: ", value)
		if err_role != nil {
			s.Errorf("Not able to retrieve credentials via Instance Profile. ARN: %v. ERROR: %v", s.config.AssumeRoleARN, err_role)
			return nil, nil, err_role
		} else {
			s.Infof("Session is created using EC2 Instance Profile")
		}

		s.client = s3.NewFromConfig(cfgRole)
		return s3manager.NewUploader(s.client), s3manager.NewDownloader(s.client), nil
	} else {
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(s.config.Region),
		)
		if err != nil {
			return nil, nil, err
		}
		value, err_role := cfg.Credentials.Retrieve(ctx)
		s.Infof("Credentials %v: ", value)
		if err_role != nil {
			s.Errorf("Not able to retrieve credentials via Instance Profile. ERROR: %v", err_role)
			return nil, nil, err_role
		} else {
			s.Infof("Session is created using EC2 Instance Profile")
		}

		s.client = s3.NewFromConfig(cfg)
		return s3manager.NewUploader(s.client), s3manager.NewDownloader(s.client), nil
	}

	return nil, nil, nil
}

// Sometimes there can be a _$folder$ object which needs to get cleaned up.
func (s *S3Sink) checkForDangling(ctx context.Context) {
	dangCheck := time.NewTicker(TimeCheckDangling)
	defer dangCheck.Stop()

	s.Infof("starting check for dangling objects.")
	for {
		select {
		case _ = <-dangCheck.C:
			go func() {
				err := s.runRemoveDangle(ctx)
				if err != nil {
					s.Errorf("Cannot check for dangling objects: %v", err)
				}
			}()

		case <-ctx.Done():
			return
		}
	}
}

func (s *S3Sink) runRemoveDangle(ctx context.Context) error {
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
	}

	// Create the Paginator for the ListObjectsV2 operation.
	p := s3.NewListObjectsV2Paginator(s.client, params, func(o *s3.ListObjectsV2PaginatorOptions) {})

	for p.HasMorePages() {
		// Next Page takes a new context for each page retrieval. This is where
		// you could add timeouts or deadlines.
		page, err := p.NextPage(ctx)
		if err != nil {
			return err
		}

		// Log the objects found
		for _, obj := range page.Contents {
			if strings.Contains(*obj.Key, "_$folder$") {
				di := &s3.DeleteObjectInput{
					Bucket: aws.String(s.Bucket),
					Key:    obj.Key,
				}
				_, err := s.client.DeleteObject(ctx, di)
				if err != nil {
					return err
				}
				s.Infof("Removed dangling: %s", *obj.Key)
			}
		}
	}

	return nil
}
