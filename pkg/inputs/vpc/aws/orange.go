package aws

import (
	"context"
	"flag"
	"fmt"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	iamRole   string
	sqsName   string
	regions   string
	isLambda  bool
	localFile string

	ERROR_SLEEP_TIME     = 20 * time.Second
	MappingCheckDuration = 30 * 60 * time.Second
)

func init() {
	flag.StringVar(&iamRole, "iam_role", "", "IAM Role to use for processing flow")
	flag.StringVar(&sqsName, "sqs_name", "", "Listen for events from this queue for new objects to look at.")
	flag.StringVar(&regions, "aws_regions", "us-east-1", "CSV list of region to run in. Will look for metadata in all regions, run SQS in first region.")
	flag.BoolVar(&isLambda, "aws_lambda", kt.LookupEnvBool("AWS_IS_LAMBDA", false), "Run as a AWS Lambda function")
	flag.StringVar(&localFile, "aws_local_file", "", "If set, process this local file and exit")
}

type AwsVpc struct {
	logger.ContextL
	metrics       *OrangeMetric
	recs          chan *FlowSet
	sqsCli        *sqs.SQS
	awsQUrl       string
	client        *s3.S3
	jchfChan      chan []*kt.JCHF
	topo          *AWSTopology
	regions       []string
	lambdaHandler func([]*kt.JCHF, func(error))
	config        *ktranslate.AWSVPCInputConfig
}

type OrangeMetric struct {
	ObjectsSeen       go_metrics.Meter
	Flows             go_metrics.Meter
	DroppedFlows      go_metrics.Meter
	RateSent          go_metrics.Meter
	DispatchCount     go_metrics.Counter
	DispatchRecsCount go_metrics.Counter
}

func NewVpc(ctx context.Context, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi, lambdaHandler func([]*kt.JCHF, func(error)), cfg *ktranslate.AWSVPCInputConfig) (*AwsVpc, error) {
	vpc := &AwsVpc{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "awsVpc"}, log),
		recs:     make(chan *FlowSet, 1000),
		jchfChan: jchfChan,
		metrics: &OrangeMetric{
			ObjectsSeen:  go_metrics.GetOrRegisterMeter("objects_seen", registry),
			Flows:        go_metrics.GetOrRegisterMeter("flows", registry),
			DroppedFlows: go_metrics.GetOrRegisterMeter("dropped_flows", registry),
		},
		awsQUrl: cfg.SQSName,
		config:  cfg,
	}

	if cfg.IsLambda {
		sess := session.Must(session.NewSession())
		conf := aws.NewConfig()
		vpc.client = s3.New(sess, conf)
		vpc.lambdaHandler = lambdaHandler
		go lambda.Start(vpc.handleLamdba)
		vpc.Infof("Running as a lamdba function")
	} else if cfg.LocalFile != "" {
		vpc.Infof("Running on %s directly.", cfg.LocalFile)
		go vpc.handleLocal(cfg.LocalFile)
	} else {
		if vpc.awsQUrl == "" {
			return nil, fmt.Errorf("Flag --sqs_name (or AWSVPCInput.SQSName) required")
		}
		if cfg.IAMRole == "" {
			return nil, fmt.Errorf("Flag --iam_role (or AWSVPCInput.IAMRole) required")
		}

		if len(cfg.Regions) == 0 {
			return nil, fmt.Errorf("Flag --regions (or AWSVPCInput.Regions) required")
		}
		vpc.regions = cfg.Regions

		vpc.Infof("Running with role %s in region %s looking at q %s and metadata in %v", cfg.IAMRole, vpc.regions[0], cfg.SQSName, vpc.regions)
		sess := session.Must(session.NewSession())
		conf := aws.NewConfig().
			WithRegion(vpc.regions[0]).
			WithCredentials(stscreds.NewCredentials(sess, cfg.IAMRole))
		vpc.client = s3.New(sess, conf)
		vpc.sqsCli = sqs.New(sess, conf)

		go vpc.checkQIn(ctx)
		go vpc.checkQOut(ctx)
		go vpc.checkMappings(ctx)
	}

	return vpc, nil
}

func (vpc *AwsVpc) Close() {}

func (vpc *AwsVpc) HttpInfo() map[string]float64 {
	return map[string]float64{
		"ObjectsSeen":  vpc.metrics.ObjectsSeen.Rate1(),
		"Flows":        vpc.metrics.Flows.Rate1(),
		"DroppedFlows": vpc.metrics.DroppedFlows.Rate1(),
	}
}

func (vpc *AwsVpc) checkQIn(ctx context.Context) {
	for {
		err := vpc.checkQueue(ctx) // Will block waiting for input from sqs queue.
		if err != nil {
			vpc.Errorf("Cannot check queue: %s -> %v", vpc.awsQUrl, err)
			time.Sleep(ERROR_SLEEP_TIME)
		}
	}
}

func (vpc *AwsVpc) checkQOut(ctx context.Context) {
	for {
		select {
		case rec := <-vpc.recs:
			dst := make([]*kt.JCHF, len(rec.Lines))
			for i, l := range rec.Lines {
				dst[i] = l.ToFlow(vpc, vpc.topo)
			}

			if len(dst) > 0 {
				vpc.jchfChan <- dst
			}

		case <-ctx.Done():
			vpc.Infof("CheckQOut Done")
			return
		}
	}
}
