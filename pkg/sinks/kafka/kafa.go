package kafka

import (
	"context"
	"flag"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
	kafka "github.com/segmentio/kafka-go"
)

/**
Config options at https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
*/

type KafkaSink struct {
	logger.ContextL
	Topic    string
	kp       *kafka.Writer
	registry go_metrics.Registry
	metrics  *KafkaMetric
}

type KafkaMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

var (
	KafkaTopic            = flag.String("kafka_topic", "", "kafka topic to produce on")
	KafkaBootstrapServers = flag.String("bootstrap.servers", "", "bootstrap.servers")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*KafkaSink, error) {
	return &KafkaSink{
		registry: registry,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kafkaSink"}, log),
		metrics: &KafkaMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_kafka", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_kafka", registry),
		},
	}, nil
}

func (s *KafkaSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {

	s.Topic = *KafkaTopic

	if s.Topic == "" {
		return fmt.Errorf("Not writing to kafka -- no topic set, use -kafka_topic flag")
	}

	s.kp = &kafka.Writer{
		Addr:     kafka.TCP(*KafkaBootstrapServers),
		Topic:    s.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	s.Infof("System connected to kafka, topic is %s", s.Topic)

	return nil
}

func (s *KafkaSink) Send(ctx context.Context, payload *kt.Output) {
	err := s.kp.WriteMessages(ctx, kafka.Message{
		Value: payload.Body,
	})
	if err != nil {
		s.Errorf("There was an error with the delivery: %v.", err)
		s.metrics.DeliveryErr.Mark(1)
	} else {
		s.metrics.DeliveryWin.Mark(1)
	}
}

func (s *KafkaSink) Close() {
	if s.kp != nil {
		s.kp.Close()
	}
}

func (s *KafkaSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}
