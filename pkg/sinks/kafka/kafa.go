package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"strings"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var (
	topic            string
	bootstrapServers string
	tlsConfig        string
	saslUser         string
	saslPass         string
	saslMech         string
	skipVerify       bool
)

func init() {
	flag.StringVar(&topic, "kafka_topic", "", "kafka topic to produce on")
	flag.StringVar(&bootstrapServers, "bootstrap.servers", "", "bootstrap.servers")
	flag.StringVar(&tlsConfig, "kafka.tls.config", "", "tls config to use for kafka. Can be basic|cert.pem,key.pem|cert.pem,key.pem,ca-cert.pem")
	flag.StringVar(&saslUser, "kafka.sasl.user", "", "kafka user")
	flag.StringVar(&saslPass, "kafka.sasl.password", "", "kafka password")
	flag.StringVar(&saslMech, "kafka.sasl.mechanism", "", "plain|scram")
	flag.BoolVar(&skipVerify, "kafka.tls.skip.verify", false, "Skip the cert verification for kafka")
}

/**
Config options at https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
*/

type KafkaSink struct {
	logger.ContextL
	kp       *kafka.Writer
	registry go_metrics.Registry
	metrics  *KafkaMetric
	config   *ktranslate.KafkaSinkConfig
	tlsConf  *tls.Config
}

type KafkaMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.KafkaSinkConfig) (*KafkaSink, error) {

	ks := KafkaSink{
		registry: registry,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kafkaSink"}, log),
		metrics: &KafkaMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_kafka", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_kafka", registry),
		},
		config: cfg,
	}

	if cfg.TlsConfig != "" {
		pts := strings.Split(cfg.TlsConfig, ",")
		switch len(pts) {
		case 1:
			if strings.ToLower(pts[0]) == "basic" {
				ks.tlsConf = &tls.Config{InsecureSkipVerify: ks.config.SkipVerify} // Just a basic config to force using tls.
			} else {
				return nil, fmt.Errorf("Invalid kafka.tls.config value: %v.", cfg.TlsConfig)
			}
		case 2:
			cert, err := tls.LoadX509KeyPair(pts[0], pts[1])
			if err != nil {
				return nil, err
			}
			ks.tlsConf = &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: ks.config.SkipVerify}
		case 3:
			// Load client cert
			cert, err := tls.LoadX509KeyPair(pts[0], pts[1])
			if err != nil {
				return nil, err
			}

			// Load CA cert
			caCert, err := os.ReadFile(pts[2])
			if err != nil {
				return nil, err
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			// Setup client
			ks.tlsConf = &tls.Config{
				Certificates:       []tls.Certificate{cert},
				RootCAs:            caCertPool,
				InsecureSkipVerify: ks.config.SkipVerify,
			}
			ks.tlsConf.BuildNameToCertificate()
		default:
			return nil, fmt.Errorf("Invalid kafka.tls.config value: %v.", cfg.TlsConfig)
		}
	}

	return &ks, nil
}

func (s *KafkaSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {

	if s.config.Topic == "" {
		return fmt.Errorf("Not writing to kafka -- no topic set, use -kafka_topic flag or KafkaSink.Topic")
	}

	s.kp = &kafka.Writer{
		Addr:     kafka.TCP(s.config.BootstrapServers),
		Topic:    s.config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	var mech sasl.Mechanism

	if s.config.SaslMech != "" {
		s.Infof("Adding Sasl Support")
		switch s.config.SaslMech {
		case "plain":
			mech = plain.Mechanism{
				Username: s.config.SaslUser,
				Password: s.config.SaslPass,
			}
		case "scram":
			mechl, err := scram.Mechanism(scram.SHA512, s.config.SaslUser, s.config.SaslPass)
			if err != nil {
				return err
			}
			mech = mechl
		default:
			return fmt.Errorf("Invalid kafka.sasl.mechanism flag=%s. Valid is plain|scram", s.config.SaslMech)
		}
	}

	if s.tlsConf != nil {
		s.Infof("Adding TLS Support")
		s.kp.Transport = &kafka.Transport{
			TLS:  s.tlsConf,
			SASL: mech,
		}
	} else if mech != nil {
		s.kp.Transport = &kafka.Transport{
			SASL: mech,
		}
	}

	s.Infof("System connected to kafka, topic is %s", s.config.Topic)

	return nil
}

func (s *KafkaSink) Send(ctx context.Context, payload *kt.Output) {
	go s.send(ctx, payload)
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

func (s *KafkaSink) send(ctx context.Context, payload *kt.Output) {
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
