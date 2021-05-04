package aws

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type AwsVpc struct {
	logger.ContextL
	metrics  *OrangeMetric
	recs     chan *FlowSet
	sqsCli   *sqs.SQS
	awsQUrl  string
	client   *s3.S3
	jchfChan chan []*kt.JCHF
	topo     *AWSTopology
	regions  []string
	mux      sync.RWMutex
}

type OrangeMetric struct {
	ObjectsSeen       go_metrics.Meter
	Flows             go_metrics.Meter
	DroppedFlows      go_metrics.Meter
	RateSent          go_metrics.Meter
	DispatchCount     go_metrics.Counter
	DispatchRecsCount go_metrics.Counter
}

var (
	IamRole = flag.String("iam_role", "", "IAM Role to use for processing flow")
	SqsName = flag.String("sqs_name", "", "Listen for events from this queue for new objects to look at.")
	Regions = flag.String("aws_regions", "us-east-1", "CSV list of region to run in. Will look for metadata in all regions, run SQS in first region.")

	ERROR_SLEEP_TIME     = 20 * time.Second
	MappingCheckDuration = 30 * 60 * time.Second
)

func NewVpc(ctx context.Context, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi) (*AwsVpc, error) {
	vpc := &AwsVpc{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "awsVpc"}, log),
		recs:     make(chan *FlowSet, 1000),
		jchfChan: jchfChan,
		metrics: &OrangeMetric{
			ObjectsSeen:  go_metrics.GetOrRegisterMeter("objects_seen", registry),
			Flows:        go_metrics.GetOrRegisterMeter("flows", registry),
			DroppedFlows: go_metrics.GetOrRegisterMeter("dropped_flows", registry),
		},
		awsQUrl: *SqsName,
	}

	if vpc.awsQUrl == "" {
		return nil, fmt.Errorf("Flag --sqs_name required")
	}
	if *IamRole == "" {
		return nil, fmt.Errorf("Flag --iam_role required")
	}

	regions := strings.Split(*Regions, ",")
	if len(regions) == 0 {
		return nil, fmt.Errorf("Flag --regions required")
	}
	vpc.regions = regions

	vpc.Infof("Running with role %s in region %s looking at q %s", *IamRole, vpc.regions[0], *SqsName)
	sess := session.Must(session.NewSession())
	conf := aws.NewConfig().
		WithRegion(vpc.regions[0]).
		WithCredentials(stscreds.NewCredentials(sess, *IamRole))
	vpc.client = s3.New(sess, conf)
	vpc.sqsCli = sqs.New(sess, conf)

	go vpc.checkQIn(ctx)
	go vpc.checkQOut(ctx)

	return vpc, nil
}

func (vpc *AwsVpc) Close() {}

func (vpc *AwsVpc) HttpInfo() map[string]float64 {
	return map[string]float64{
		"ObjectsSeen":  vpc.metrics.ObjectsSeen.Rate1(),
		"Flows":        vpc.metrics.Flows.Rate1(),
		"DroppedFlows": vpc.metrics.Flows.Rate1(),
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
				dst[i] = l.ToFlow()
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
