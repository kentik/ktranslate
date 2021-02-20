package kafka

import (
	"context"
	"flag"
	"fmt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

/**
Config options at https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
*/

type KafkaSink struct {
	logger.ContextL
	Topic    string
	conf     kafka.ConfigMap
	kp       *kafka.Producer
	registry go_metrics.Registry
	metrics  *KafkaMetric
}

type KafkaMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

var (
	KafkaTopic                      = flag.String("kafka_topic", "", "kafka topic to produce on")
	KafkaBootstrapServers           = flag.String("bootstrap.servers", "", "bootstrap.servers")
	KafkaSecurityProtocol           = flag.String("security.protocol", "", "security.protocol")
	KafkaSASLKerberosServiceName    = flag.String("sasl.kerberos.service.name", "kafka", "sasl.kerberos.service.name")
	KafkaSSLCALocation              = flag.String("ssl.ca.location", "", "ssl.ca.location")
	KafkaKerberosKeytabFileLocation = flag.String("sasl.kerberos.keytab", "", "sasl.kerberos.keytab")
	KafkaKerberosPrincipalName      = flag.String("sasl.kerberos.principal", "", "sasl.kerberos.principal")
	KafkaKerberosKinitCmd           = flag.String("sasl.kerberos.kinit.cmd", `kinit -R -t "%{sasl.kerberos.keytab}" -k %{sasl.kerberos.principal} || kinit -t "%{sasl.kerberos.keytab}" -k %{sasl.kerberos.principal}`, "sasl.kerberos.kinit.cmd")
	KafkaSASLMechanism              = flag.String("sasl.mechanism", "", "sasl.mechanism")
	KafkaDebugContext               = flag.String("kafka.debug", "", "debug contexts to enable for kafka")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*KafkaSink, error) {
	return &KafkaSink{
		registry: registry,
		conf:     kafka.ConfigMap{},
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kafkaSink"}, log),
		metrics: &KafkaMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_kafka", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_kafka", registry),
		},
	}, nil
}

func (s *KafkaSink) Init(ctx context.Context, format formats.Format, compression kt.Compression) error {

	s.Topic = *KafkaTopic

	// These are all optional
	if *KafkaBootstrapServers != "" {
		s.conf["bootstrap.servers"] = *KafkaBootstrapServers
	}
	if *KafkaSecurityProtocol != "" {
		s.conf["security.protocol"] = *KafkaSecurityProtocol
	}
	if *KafkaSASLKerberosServiceName != "" {
		s.conf["sasl.kerberos.service.name"] = *KafkaSASLKerberosServiceName
	}
	if *KafkaSSLCALocation != "" {
		s.conf["ssl.ca.location"] = *KafkaSSLCALocation
	}
	if *KafkaKerberosKeytabFileLocation != "" {
		s.conf["sasl.kerberos.keytab"] = *KafkaKerberosKeytabFileLocation
	}
	if *KafkaKerberosPrincipalName != "" {
		s.conf["sasl.kerberos.principal"] = *KafkaKerberosPrincipalName
	}
	if *KafkaDebugContext != "" {
		s.conf["debug"] = *KafkaDebugContext
	}
	if *KafkaSASLMechanism != "" {
		s.conf["sasl.mechanism"] = *KafkaSASLMechanism
	}
	if *KafkaKerberosKinitCmd != "" {
		s.conf["sasl.kerberos.kinit.cmd"] = *KafkaKerberosKinitCmd
	}

	if s.Topic == "" {
		return fmt.Errorf("Not writing to kafka -- no topic set, use -kafka_topic flag")
	}

	var err error
	s.kp, err = kafka.NewProducer(&s.conf)
	if err != nil {
		return err
	}
	s.Infof("System connected to kafka, topic is %s", s.Topic)
	go s.checkKafkaStatus()

	return nil
}

func (s *KafkaSink) Send(ctx context.Context, payload []byte) {
	s.kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.Topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, nil)
}

func (s *KafkaSink) Close() {
	s.kp.Close()
}

func (s *KafkaSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}

// Delivery report handler for produced messages
func (s *KafkaSink) checkKafkaStatus() {
	for e := range s.kp.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				s.Errorf("Delivery failed: %v\n", ev.TopicPartition)
				s.metrics.DeliveryErr.Mark(1)
			} else {
				s.metrics.DeliveryWin.Mark(1)
			}
		}
	}
}
