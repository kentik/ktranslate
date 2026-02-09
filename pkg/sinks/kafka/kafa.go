package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	topic                        string
	bootstrapServers             string
	kafkaSecurityProtocol        string
	kafkaSASLMechanism           string
	kafkaSASLUsername            string
	kafkaSASLPassword            string
	kafkaKerberosServiceName     string
	kafkaKerberosRealm           string
	kafkaKerberosConfigPath      string
	kafkaKerberosKeytabPath      string
	kafkaKerberosPrincipal       string
	kafkaKerberosDisablePAFXFAST bool
	kafkaSSLCAFile               string
	kafkaSSLCertFile             string
	kafkaSSLKeyFile              string
	kafkaSSLKeyPassword          string
	kafkaSSLInsecure             bool
	kafkaRequiredAcks            int
	kafkaCompression             string
	kafkaMaxMessageBytes         int
	kafkaRetryMax                int
	kafkaFlushFrequency          int
	kafkaFlushMessages           int
	kafkaFlushBytes              int
)

func init() {
	flag.StringVar(&topic, "kafka_topic", "", "kafka topic to produce on")
	flag.StringVar(&bootstrapServers, "bootstrap.servers", "", "bootstrap.servers")
	flag.StringVar(&kafkaSecurityProtocol, "kafka_security_protocol", "PLAINTEXT", "Security protocol (PLAINTEXT, SASL_PLAINTEXT, SASL_SSL, SSL)")
	flag.StringVar(&kafkaSASLMechanism, "kafka_sasl_mechanism", "", "SASL mechanism (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, GSSAPI)")
	flag.StringVar(&kafkaSASLUsername, "kafka_sasl_username", "", "SASL username")
	flag.StringVar(&kafkaSASLPassword, "kafka_sasl_password", "", "SASL password")
	flag.StringVar(&kafkaKerberosServiceName, "kafka_kerberos_service_name", "kafka", "Kerberos service name")
	flag.StringVar(&kafkaKerberosRealm, "kafka_kerberos_realm", "", "Kerberos realm")
	flag.StringVar(&kafkaKerberosConfigPath, "kafka_kerberos_config_path", "/etc/krb5.conf", "Path to Kerberos config file")
	flag.StringVar(&kafkaKerberosKeytabPath, "kafka_kerberos_keytab_path", "", "Path to Kerberos keytab file")
	flag.StringVar(&kafkaKerberosPrincipal, "kafka_kerberos_principal", "", "Kerberos principal")
	flag.BoolVar(&kafkaKerberosDisablePAFXFAST, "kafka_kerberos_disable_pafx_fast", false, "Disable PA-FX-FAST")
	flag.StringVar(&kafkaSSLCAFile, "kafka_ssl_ca_file", "", "SSL CA certificate file")
	flag.StringVar(&kafkaSSLCertFile, "kafka_ssl_cert_file", "", "SSL client certificate file")
	flag.StringVar(&kafkaSSLKeyFile, "kafka_ssl_key_file", "", "SSL client private key file")
	flag.StringVar(&kafkaSSLKeyPassword, "kafka_ssl_key_password", "", "SSL private key password")
	flag.BoolVar(&kafkaSSLInsecure, "kafka_ssl_insecure", false, "Skip SSL certificate verification")
	flag.IntVar(&kafkaRequiredAcks, "kafka_required_acks", 1, "Required acks (0=NoResponse, 1=WaitForLocal, -1=WaitForAll)")
	flag.StringVar(&kafkaCompression, "kafka_compression", "none", "Compression codec (none, gzip, snappy, lz4, zstd)")
	flag.IntVar(&kafkaMaxMessageBytes, "kafka_max_message_bytes", 1000000, "Maximum message size in bytes")
	flag.IntVar(&kafkaRetryMax, "kafka_retry_max", 3, "Maximum number of retries")
	flag.IntVar(&kafkaFlushFrequency, "kafka_flush_frequency", 100, "Flush frequency in milliseconds")
	flag.IntVar(&kafkaFlushMessages, "kafka_flush_messages", 100, "Flush after this many messages")
	flag.IntVar(&kafkaFlushBytes, "kafka_flush_bytes", 65536, "Flush after this many bytes")
}

/**
Config options for Sarama Kafka client with full Kerberos support
*/

type KafkaSink struct {
	logger.ContextL
	Topic        string
	producer     sarama.AsyncProducer
	registry     go_metrics.Registry
	metrics      *KafkaMetric
	config       *ktranslate.KafkaSinkConfig
	saramaConfig *sarama.Config
}

type KafkaMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.KafkaSinkConfig) (*KafkaSink, error) {
	return &KafkaSink{
		registry: registry,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kafkaSink"}, log),
		metrics: &KafkaMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_kafka", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_kafka", registry),
		},
		config: cfg,
	}, nil
}

func (s *KafkaSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {

	if s.config.Topic == "" {
		return fmt.Errorf("Not writing to kafka -- no topic set, use -kafka_topic flag or KafkaSink.Topic")
	}

	if s.config.BootstrapServers == "" {
		return fmt.Errorf("Not writing to kafka -- no bootstrap servers set, use -bootstrap.servers flag")
	}

	// Create Sarama configuration
	config := sarama.NewConfig()

	// Basic producer settings
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = s.config.RetryMax
	config.Producer.RequiredAcks = sarama.RequiredAcks(s.config.RequiredAcks)
	config.Producer.MaxMessageBytes = s.config.MaxMessageBytes

	// Compression
	switch strings.ToLower(s.config.Compression) {
	case "none":
		config.Producer.Compression = sarama.CompressionNone
	case "gzip":
		config.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		config.Producer.Compression = sarama.CompressionSnappy
	case "lz4":
		config.Producer.Compression = sarama.CompressionLZ4
	case "zstd":
		config.Producer.Compression = sarama.CompressionZSTD
	default:
		config.Producer.Compression = sarama.CompressionNone
	}

	// Flush settings
	config.Producer.Flush.Frequency = time.Duration(s.config.FlushFrequency) * time.Millisecond
	config.Producer.Flush.Messages = s.config.FlushMessages
	config.Producer.Flush.Bytes = s.config.FlushBytes

	// Security configuration
	switch strings.ToUpper(s.config.SecurityProtocol) {
	case "PLAINTEXT":
		// No additional configuration needed
	case "SSL":
		config.Net.TLS.Enable = true
		if err := s.configureTLS(config); err != nil {
			return fmt.Errorf("failed to configure SSL: %v", err)
		}
	case "SASL_PLAINTEXT":
		config.Net.SASL.Enable = true
		if err := s.configureSASL(config); err != nil {
			return fmt.Errorf("failed to configure SASL: %v", err)
		}
	case "SASL_SSL":
		config.Net.TLS.Enable = true
		config.Net.SASL.Enable = true
		if err := s.configureTLS(config); err != nil {
			return fmt.Errorf("failed to configure SSL: %v", err)
		}
		if err := s.configureSASL(config); err != nil {
			return fmt.Errorf("failed to configure SASL: %v", err)
		}
	default:
		return fmt.Errorf("unsupported security protocol: %s", s.config.SecurityProtocol)
	}

	// Version configuration - use a recent version that supports all features
	version, err := sarama.ParseKafkaVersion("2.8.0")
	if err != nil {
		return fmt.Errorf("failed to parse kafka version: %v", err)
	}
	config.Version = version

	s.saramaConfig = config

	// Create producer
	brokers := strings.Split(s.config.BootstrapServers, ",")
	for i, broker := range brokers {
		brokers[i] = strings.TrimSpace(broker)
	}

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %v", err)
	}

	s.producer = producer
	s.Topic = s.config.Topic

	// Start goroutines to handle success/error responses
	go s.handleResponses(ctx)

	s.Infof("System connected to kafka, topic is %s, brokers: %s, security: %s",
		s.Topic, s.config.BootstrapServers, s.config.SecurityProtocol)

	return nil
}

func (s *KafkaSink) configureTLS(config *sarama.Config) error {
	tlsConfig := &tls.Config{}

	if s.config.SSLInsecure {
		tlsConfig.InsecureSkipVerify = true
	}

	// Load CA certificate
	if s.config.SSLCAFile != "" {
		caCert, err := ioutil.ReadFile(s.config.SSLCAFile)
		if err != nil {
			return fmt.Errorf("failed to read CA certificate: %v", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	// Load client certificate and key
	if s.config.SSLCertFile != "" && s.config.SSLKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(s.config.SSLCertFile, s.config.SSLKeyFile)
		if err != nil {
			return fmt.Errorf("failed to load client certificate: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	config.Net.TLS.Config = tlsConfig
	return nil
}

func (s *KafkaSink) configureSASL(config *sarama.Config) error {
	switch strings.ToUpper(s.config.SASLMechanism) {
	case "PLAIN":
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.User = s.config.SASLUsername
		config.Net.SASL.Password = s.config.SASLPassword

	case "SCRAM-SHA-256":
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		config.Net.SASL.User = s.config.SASLUsername
		config.Net.SASL.Password = s.config.SASLPassword

	case "SCRAM-SHA-512":
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.User = s.config.SASLUsername
		config.Net.SASL.Password = s.config.SASLPassword

	case "GSSAPI":
		// Kerberos configuration
		config.Net.SASL.Mechanism = sarama.SASLTypeGSSAPI
		config.Net.SASL.GSSAPI.ServiceName = s.config.KerberosServiceName
		config.Net.SASL.GSSAPI.Realm = s.config.KerberosRealm
		config.Net.SASL.GSSAPI.Username = s.config.KerberosPrincipal
		config.Net.SASL.GSSAPI.DisablePAFXFAST = s.config.KerberosDisablePAFXFAST

		// Set Kerberos config file path
		if s.config.KerberosConfigPath != "" {
			if _, err := os.Stat(s.config.KerberosConfigPath); err != nil {
				return fmt.Errorf("kerberos config file not found: %s", s.config.KerberosConfigPath)
			}
			os.Setenv("KRB5_CONFIG", s.config.KerberosConfigPath)
			config.Net.SASL.GSSAPI.KerberosConfigPath = s.config.KerberosConfigPath
		}

		// Set keytab file if provided
		if s.config.KerberosKeytabPath != "" {
			if _, err := os.Stat(s.config.KerberosKeytabPath); err != nil {
				return fmt.Errorf("kerberos keytab file not found: %s", s.config.KerberosKeytabPath)
			}
			config.Net.SASL.GSSAPI.KeyTabPath = s.config.KerberosKeytabPath
			config.Net.SASL.GSSAPI.AuthType = sarama.KRB5_KEYTAB_AUTH
		} else {
			// Use user authentication (requires kinit)
			config.Net.SASL.GSSAPI.AuthType = sarama.KRB5_USER_AUTH
		}

	case "OAUTHBEARER":
		config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		// OAuth configuration would need additional implementation
		return fmt.Errorf("OAUTHBEARER mechanism not fully implemented yet")

	default:
		if s.config.SASLMechanism != "" {
			return fmt.Errorf("unsupported SASL mechanism: %s", s.config.SASLMechanism)
		}
	}

	return nil
}

func (s *KafkaSink) handleResponses(ctx context.Context) {
	for {
		select {
		case success := <-s.producer.Successes():
			if success != nil {
				s.metrics.DeliveryWin.Mark(1)
			}
		case err := <-s.producer.Errors():
			if err != nil {
				s.Errorf("Failed to produce message: %v", err.Err)
				s.metrics.DeliveryErr.Mark(1)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *KafkaSink) Send(ctx context.Context, payload *kt.Output) {
	if s.producer == nil {
		s.Errorf("Kafka producer not initialized")
		s.metrics.DeliveryErr.Mark(1)
		return
	}

	message := &sarama.ProducerMessage{
		Topic: s.config.Topic,
		Value: sarama.ByteEncoder(payload.Body),
	}

	select {
	case s.producer.Input() <- message:
		// Message sent to producer successfully
	case <-ctx.Done():
		s.Warnf("Context cancelled while sending message")
		return
	default:
		s.Errorf("Producer input channel is full")
		s.metrics.DeliveryErr.Mark(1)
	}
}

func (s *KafkaSink) Close() {
	if s.producer != nil {
		s.Infof("Closing Kafka producer")
		if err := s.producer.Close(); err != nil {
			s.Errorf("Error closing Kafka producer: %v", err)
		}
	}
}

func (s *KafkaSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}
