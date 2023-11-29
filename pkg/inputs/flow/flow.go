package flow

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/resolv"

	"github.com/netsampler/goflow2/v2/decoders/netflow"
	"github.com/netsampler/goflow2/v2/format"
	"github.com/netsampler/goflow2/v2/metrics"
	protoproducer "github.com/netsampler/goflow2/v2/producer/proto"
	"github.com/netsampler/goflow2/v2/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

type FlowSource string

const (
	Ipfix    FlowSource = "ipfix"
	Sflow               = "sflow"
	Netflow5            = "netflow5"
	Netflow9            = "netflow9"
	NBar                = "nbar"
	ASA                 = "asa"
	PAN                 = "pan"
	JFlow               = "jflow"
	CFlow               = "cflow"
)

var (
	addr        string
	port        int
	reuse       bool
	workers     int
	fields      string
	promListen  string
	mappingFile string
)

func init() {
	flag.StringVar(&addr, "nf.addr", "0.0.0.0", "Sflow/NetFlow/IPFIX listening address")
	flag.IntVar(&port, "nf.port", 9995, "Sflow/NetFlow/IPFIX listening port")
	flag.BoolVar(&reuse, "nf.reuserport", false, "Enable so_reuseport for Sflow/NetFlow/IPFIX")
	flag.IntVar(&workers, "nf.workers", 1, "Number of workers per flow collector")
	flag.StringVar(&fields, "nf.message.fields", ktranslate.FlowDefaultFields, "The list of fields to include in flow messages. Can be any of "+ktranslate.FlowFields)
	flag.StringVar(&promListen, "nf.prom.listen", "", "Run a promethues metrics collector here")
	flag.StringVar(&mappingFile, "nf.mapping", "", "Configuration file for custom netflow mappings")
}

func NewFlowSource(ctx context.Context, proto FlowSource, maxBatchSize int, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi, resolv *resolv.Resolver, cfg *ktranslate.FlowInputConfig) (*KentikDriver, error) {

	defer func() {
		if v := cfg.PrometheusListenAddr; v != "" {
			http.Handle("/metrics", promhttp.Handler())
			go http.ListenAndServe(v, nil)
		}
	}()

	// Allow processing of custom ipfix templates here.
	var config *protoproducer.ProducerConfig

	// Load a special config if there's a known enriched flow.
	switch proto {
	case ASA:
		config = loadASA(cfg)
	case NBar:
		config = loadNBar(cfg)
	case PAN:
		config = loadPAN(cfg)
	default:
		config = loadDefault(cfg)
	}

	kt := NewKentikDriver(ctx, proto, maxBatchSize, log, registry, jchfChan, apic, cfg.MessageFields, resolv, cfg)

	// Or pull up a special file if needed.
	if v := cfg.MappingFile; v != "" {
		f, err := os.Open(v)
		if err != nil {
			kt.Errorf("Cannot load netflow mapping file: %v", err)
			return nil, err
		}
		pc, err := loadMapping(f)
		f.Close()
		if err != nil {
			kt.Errorf("Invalid yaml for netflow mapping file: %v", err)
			return nil, err
		}
		config = pc
	}

	kt.SetConfig(config)

	flowProducer, err := protoproducer.CreateProtoProducer(config, protoproducer.CreateSamplingSystem)
	if err != nil {
		return nil, err
	}

	flowProducer = metrics.WrapPromProducer(flowProducer)
	udpCfg := &utils.UDPReceiverConfig{
		Sockets: cfg.Workers,
		Workers: cfg.Workers,
	}
	recv, err := utils.NewUDPReceiver(udpCfg)
	if err != nil {
		return nil, err
	}

	format.RegisterFormatDriver("chf", kt) // Let goflow know about kt.
	formatter, err := format.FindFormat("chf")
	if err != nil {
		return nil, err
	}

	cfgPipe := &utils.PipeConfig{
		Format:           formatter,
		Transport:        nil,
		Producer:         flowProducer,
		NetFlowTemplater: metrics.NewDefaultPromTemplateSystem, // wrap template system to get Prometheus info
	}

	var decodeFunc utils.DecoderFunc

	switch proto {
	case Ipfix, Netflow9, ASA, NBar, PAN, Netflow5, JFlow, CFlow:
		switch proto {
		case Ipfix, ASA, NBar, PAN:
			for _, v := range config.IPFIX.Mapping {
				kt.Infof("Custom IPFIX Field Mapping: Field=%v, Pen=%v -> %v", v.Type, v.Pen, v.Destination)
			}
		case Netflow9:
			for _, v := range config.NetFlowV9.Mapping {
				kt.Infof("Custom Netflow9 Field Mapping: Field=%v -> %v", v.Type, v.Destination)
			}
		}
		kt.pipe = utils.NewNetFlowPipe(cfgPipe)
		decodeFunc = metrics.PromDecoderWrapper(kt.pipe.DecodeFlow, string(proto))
	case Sflow:
		for _, v := range config.SFlow.Mapping {
			kt.Infof("Custom SFlow Field Mapping: Layer=%s, Offset=%d, Length=%d -> %v", v.Layer, v.Offset, v.Length, v.Destination)
		}
		kt.pipe = utils.NewSFlowPipe(cfgPipe)
		decodeFunc = metrics.PromDecoderWrapper(kt.pipe.DecodeFlow, string(proto))
	default:
		return nil, fmt.Errorf("Unknown flow format %v", proto)
	}

	kt.producer = flowProducer
	kt.receiver = recv
	if err := kt.receiver.Start(cfg.ListenIP, cfg.ListenPort, decodeFunc); err != nil {
		return nil, err
	} else {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case err := <-recv.Errors():
					if errors.Is(err, netflow.ErrorTemplateNotFound) {
						kt.Warnf("template error: %v", err)
					} else if errors.Is(err, net.ErrClosed) {
						kt.Infof("closed receiver")
					} else {
						kt.Warnf("error: %v", err)
					}

				}
			}
		}()
	}

	kt.Infof("Netflow listener running on %s:%d for format %s and a batch size of %d", cfg.ListenIP, cfg.ListenPort, proto, maxBatchSize)
	kt.Infof("Netflow listener sending fields %s", cfg.MessageFields)

	return kt, nil
}

func loadMapping(f io.Reader) (*protoproducer.ProducerConfig, error) {
	config := &protoproducer.ProducerConfig{}
	dec := yaml.NewDecoder(f)
	err := dec.Decode(config)
	return config, err
}
