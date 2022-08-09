package gcppubsub

import (
	"context"
	"flag"
	"fmt"

	"cloud.google.com/go/pubsub"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	gcpPubSubProjectID string
	gcpPubSubTopic     string
)

func init() {
	flag.StringVar(&gcpPubSubProjectID, "gcp_pubsub_project_id", "", "GCP PubSub Project ID to use")
	flag.StringVar(&gcpPubSubTopic, "gcp_pubsub_topic", "", "GCP PubSub Topic to publish to")
}

type GCPPubSub struct {
	logger.ContextL
	client    *pubsub.Client
	registry  go_metrics.Registry
	metrics   *GCPPubSubMetric
	projectID string
	topicID   string
	topic     *pubsub.Topic
	config    *ktranslate.GCloudPubSubSinkConfig
}

type GCPPubSubMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.GCloudPubSubSinkConfig) (*GCPPubSub, error) {
	return &GCPPubSub{
		registry:  registry,
		projectID: cfg.ProjectID,
		topicID:   cfg.Topic,
		ContextL:  logger.NewContextLFromUnderlying(logger.SContext{S: "gcppubsub"}, log),
		metrics: &GCPPubSubMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_gcppubsub", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_gcppubsub", registry),
		},
		config: cfg,
	}, nil
}

func (s *GCPPubSub) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	if s.projectID == "" {
		return fmt.Errorf("No project ID was set for GCP PubSub. Use the -gcp_pubsub_project_id option to set one up.")
	}
	if s.topicID == "" {
		return fmt.Errorf("No topic was set for GCP PubSub. Use the -gcp_pubsub_topic option to set one up.")
	}
	client, err := pubsub.NewClient(ctx, s.projectID)
	if err != nil {
		return fmt.Errorf("There was an error when creating the gcp pubsub client: %v.", err)
	}
	s.client = client
	s.topic = s.client.Topic(s.topicID)

	s.Infof("connected to GCP PubSub for project %s on topic %s", s.projectID, s.topicID)

	return nil
}

func (s *GCPPubSub) Send(ctx context.Context, payload *kt.Output) {
	go s.send(ctx, payload.Body)
}

func (s *GCPPubSub) send(ctx context.Context, payload []byte) {
	res := s.topic.Publish(ctx, &pubsub.Message{Data: payload})
	if _, err := res.Get(ctx); err != nil {
		s.Errorf("Error publishing to GCP: %v", err)
		s.metrics.DeliveryErr.Mark(1)
		return
	}
	s.metrics.DeliveryWin.Mark(1)
}

func (s *GCPPubSub) Close() {
	if s.topic != nil {
		s.topic.Stop()
	}

	if s.client != nil {
		s.client.Close()
	}
}

func (s *GCPPubSub) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}
