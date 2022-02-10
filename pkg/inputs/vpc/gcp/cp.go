package gcp

import (
	"context"
	"flag"
	"fmt"
	"time"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"cloud.google.com/go/pubsub"
)

type GcpVpc struct {
	logger.ContextL
	metrics      *GcpMetric
	recs         chan *GCELogLine
	client       *pubsub.Client
	jchfChan     chan []*kt.JCHF
	maxBatchSize int
}

type GcpMetric struct {
	FlowsIn     go_metrics.Meter
	RateInvalid go_metrics.Meter
	RateError   go_metrics.Meter
}

var (
	ProjectID  = flag.String("gcp.project", "", "Google ProjectID to listen for flows on")
	SourceSub  = flag.String("gcp.sub", "", "Google Sub to listen for flows on")
	SampleRate = flag.Float64("gcp.sample", 1, "Sample rate of the vpc export (as defined in the VPC setup)")

	ERROR_SLEEP_TIME = 20 * time.Second
)

func NewVpc(ctx context.Context, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi, maxBatchSize int) (*GcpVpc, error) {
	vpc := &GcpVpc{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "gcpVpc"}, log),
		recs:         make(chan *GCELogLine, 1000),
		jchfChan:     jchfChan,
		maxBatchSize: maxBatchSize,
		metrics: &GcpMetric{
			FlowsIn:     go_metrics.GetOrRegisterMeter("flows_in", registry),
			RateInvalid: go_metrics.GetOrRegisterMeter("rate_invalid", registry),
			RateError:   go_metrics.GetOrRegisterMeter("rate_error", registry),
		},
	}

	if *ProjectID == "" || *SourceSub == "" {
		return nil, fmt.Errorf("Flags gcp.project and gcp.sub must be set for a GCP flow export")
	}

	client, err := pubsub.NewClient(ctx, *ProjectID)
	if err != nil {
		return nil, err
	}

	sub := client.Subscription(*SourceSub)
	if sub == nil {
		return nil, fmt.Errorf("GCP Subscription not found: %s", *SourceSub)
	}

	// Calulate how much to adjust sample rate. This needs to be an int for us, from the % based input which google uses.
	flowSample := getSampleRate(*SampleRate)

	go vpc.runSubscription(ctx, sub, client)
	go vpc.checkQOut(ctx, flowSample)

	vpc.Infof("Running GCP subscription on project: %s, sub: %s, sampleRate %d", *ProjectID, *SourceSub, flowSample)

	return vpc, nil
}

func (vpc *GcpVpc) Close() {}

func (vpc *GcpVpc) HttpInfo() map[string]float64 {
	return map[string]float64{
		"FlowsIn":     vpc.metrics.FlowsIn.Rate1(),
		"RateInvalid": vpc.metrics.RateInvalid.Rate1(),
		"RateError":   vpc.metrics.RateError.Rate1(),
	}
}

// Runs the subscription and reads messages.
func (vpc *GcpVpc) runSubscription(ctx context.Context, sub *pubsub.Subscription, client *pubsub.Client) {
	defer client.Close()

	sleepTime := 60 * time.Second
	nextSeek := time.Now().Add(sleepTime)

	for {
		if err := ctx.Err(); err != nil {
			vpc.Infof("subscription call exiting (%v)", err)
			return
		}

		vpc.Infof("subscription call running")

		if err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			m.Ack()
			var data GCELogLine
			if err := json.Unmarshal(m.Data, &data); err != nil {
				vpc.metrics.RateError.Mark(1)
				vpc.Errorf("Error reading log line: %v", err)
			} else {
				vpc.metrics.FlowsIn.Mark(1)
				if data.IsValid() {
					vpc.recs <- &data
				} else {
					vpc.metrics.RateInvalid.Mark(1)
					now := time.Now()
					if now.After(nextSeek) {
						nextSeek = now.Add(sleepTime)
						if err := sub.SeekToTime(ctx, now.Add(-1*30*time.Second)); err != nil {
							vpc.Warnf("Could not seek to window: %v", err)
						} else {
							vpc.Warnf("Was falling behind, skipping ahead to now: %v", data)
						}
					}
				}
			}
		}); err != nil {
			vpc.Warnf("Err0r on sub system receive, waiting %v -- %v", sleepTime, err)
			time.Sleep(sleepTime)
		}
	}
}

func (vpc *GcpVpc) checkQOut(ctx context.Context, flowSample uint32) {
	sendTicker := time.NewTicker(kt.SendBatchDuration)
	defer sendTicker.Stop()
	batch := make([]*kt.JCHF, 0, vpc.maxBatchSize)

	vpc.Infof("kentik driver running")
	for {
		select {
		case rec := <-vpc.recs:
			log, err := rec.ToFlow(vpc, flowSample)
			if err != nil {
				vpc.Errorf("Cannot process VPC flow: %v", err)
				continue
			}
			batch = append(batch, log)
			if len(batch) >= vpc.maxBatchSize {
				vpc.jchfChan <- batch
				batch = make([]*kt.JCHF, 0, vpc.maxBatchSize)
			}

		case <-sendTicker.C:
			if len(batch) > 0 {
				vpc.jchfChan <- batch
				batch = make([]*kt.JCHF, 0, vpc.maxBatchSize)
			}

		case <-ctx.Done():
			vpc.Infof("CheckQOut Done")
			return
		}
	}
}

// Convert from a % based sampling number to a 1:X style one.
func getSampleRate(sr float64) uint32 {
	if sr > 1 { // Doesn't make sense if more than 1 but just assume 1.
		return 1
	}
	if sr < 0.00001 { // Lower bound check to make math easier.
		return 100000
	}

	// x / y = sr, solve for x. Then reduce down to smallest y.
	top := sr * 10000.0
	bottom := 10000.0
	x := int(top)
	y := int(bottom)

	// Now, reduce to lowest terms
	//Calculate GCD
	c := x % y
	for {
		if c <= 0 {
			break
		}
		y = x
		x = c
		c = x % y
	}
	gcd := float64(x)
	return uint32(bottom / gcd)
}
