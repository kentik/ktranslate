package flow

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/resolv"

	"github.com/netsampler/goflow2/producer"
	"github.com/netsampler/goflow2/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	var config EntConfig

	// Load a special config if there's a known enriched flow.
	switch proto {
	case ASA:
		config = loadASA(cfg)
	case NBar:
		config = loadNBar(cfg)
	case PAN:
		config = loadPAN(cfg)
	}

	kt := NewKentikDriver(ctx, proto, maxBatchSize, log, registry, jchfChan, apic, cfg.MessageFields, resolv, cfg)

	// Or pull up a special file if needed.
	if v := cfg.MappingFile; v != "" {
		f, err := os.Open(v)
		if err != nil {
			kt.Errorf("Cannot load netflow mapping file: %v", err)
			return nil, err
		}
		config, err = loadMapping(f, proto)
		f.Close()
		if err != nil {
			kt.Errorf("Invalid yaml for netflow mapping file: %v", err)
			return nil, err
		}
	}

	kt.Infof("Netflow listener running on %s:%d for format %s and a batch size of %d", cfg.ListenIP, cfg.ListenPort, proto, maxBatchSize)
	kt.Infof("Netflow listener sending fields %s", cfg.MessageFields)
	kt.SetConfig(config)

	switch proto {
	case Ipfix, Netflow9, ASA, NBar, PAN:
		sNF := &utils.StateNetFlow{
			Format: kt,
			Logger: &KentikLog{l: kt},
			Config: &config.FlowConfig,
		}
		switch proto {
		case Ipfix, ASA, NBar, PAN:
			for _, v := range config.FlowConfig.IPFIX.Mapping {
				kt.Infof("Custom IPFIX Field Mapping: Field=%v, Pen=%v -> %v", v.Type, v.Pen, config.NameMap[v.Destination])
			}
		case Netflow9:
			for _, v := range config.FlowConfig.NetFlowV9.Mapping {
				kt.Infof("Custom Netflow9 Field Mapping: Field=%v -> %v", v.Type, config.NameMap[v.Destination])
			}
		}
		go func() { // Let this run, returning flow into the kentik transport struct
			err := sNF.FlowRoutine(cfg.Workers, cfg.ListenIP, cfg.ListenPort, cfg.EnableReusePort)
			if err != nil {
				sNF.Logger.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()
		return kt, nil
	case Sflow:
		sSF := &utils.StateSFlow{
			Format: kt,
			Logger: &KentikLog{l: kt},
			Config: &config.FlowConfig,
		}
		for _, v := range config.FlowConfig.SFlow.Mapping {
			kt.Infof("Custom SFlow Field Mapping: Layer=%d, Offset=%d, Length=%d -> %v", v.Layer, v.Offset, v.Length, config.NameMap[v.Destination])
		}
		go func() { // Let this run, returning flow into the kentik transport struct
			err := sSF.FlowRoutine(cfg.Workers, cfg.ListenIP, cfg.ListenPort, cfg.EnableReusePort)
			if err != nil {
				sSF.Logger.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()
		return kt, nil
	case Netflow5:
		sNFL := &utils.StateNFLegacy{
			Format: kt,
			Logger: &KentikLog{l: kt},
		}
		go func() { // Let this run, returning flow into the kentik transport struct
			err := sNFL.FlowRoutine(cfg.Workers, cfg.ListenIP, cfg.ListenPort, cfg.EnableReusePort)
			if err != nil {
				sNFL.Logger.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()
		return kt, nil
	}
	return nil, fmt.Errorf("Unknown flow format %v", proto)
}

type EntConfig struct {
	FlowConfig producer.ProducerConfig `json:"flow_config"`
	NameMap    map[string]string       `json:"name_map"`
}

func loadMapping(f io.Reader, proto FlowSource) (EntConfig, error) {
	config := EntConfig{}
	dec := json.NewDecoder(f)
	err := dec.Decode(&config)

	// Update any non filled in name maps to the default.
	if config.NameMap == nil {
		config.NameMap = map[string]string{}
	}
	switch proto {
	case Ipfix, ASA, NBar, PAN:
		for _, v := range config.FlowConfig.IPFIX.Mapping {
			if _, ok := config.NameMap[v.Destination]; !ok {
				config.NameMap[v.Destination] = v.Destination
			}
		}
	case Netflow9:
		for _, v := range config.FlowConfig.NetFlowV9.Mapping {
			if _, ok := config.NameMap[v.Destination]; !ok {
				config.NameMap[v.Destination] = v.Destination
			}
		}
	case Sflow:
		for _, v := range config.FlowConfig.SFlow.Mapping {
			if _, ok := config.NameMap[v.Destination]; !ok {
				config.NameMap[v.Destination] = v.Destination
			}
		}
	}

	return config, err
}
