package flow

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/resolv"

	"github.com/netsampler/goflow2/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type FlowSource string

const (
	Ipfix    FlowSource = "ipfix"
	Sflow               = "sflow"
	Netflow5            = "netflow5"
	Netflow9            = "netflow9"

	// These can be any of
	fullFieldList = "Type,TimeReceived,SequenceNum,SamplingRate,SamplerAddress,TimeFlowStart,TimeFlowEnd,Bytes,Packets,SrcAddr,DstAddr,Etype,Proto,SrcPort,DstPort,InIf,OutIf,SrcMac,DstMac,SrcVlan,DstVlan,VlanId,IngressVrfID,EgressVrfID,IPTos,ForwardingStatus,IPTTL,TCPFlags,IcmpType,IcmpCode,IPv6FlowLabel,FragmentId,FragmentOffset,BiFlowDirection,SrcAS,DstAS,NextHop,NextHopAS,SrcNet,DstNet,HasMPLS,MPLSCount,MPLS1TTL,MPLS1Label,MPLS2TTL,MPLS2Label,MPLS3TTL,MPLS3Label,MPLSLastTTL,MPLSLastLabel,CustomInteger1,CustomInteger2,CustomBytes1,CustomBytes2"
	defaultFields = "TimeReceived,SamplingRate,Bytes,Packets,SrcAddr,DstAddr,Proto,SrcPort,DstPort,InIf,OutIf,SrcVlan,DstVlan,TCPFlags,SrcAS,DstAS,Type,SamplerAddress"
)

var (
	Addr          = flag.String("nf.addr", "0.0.0.0", "Sflow/NetFlow/IPFIX listening address")
	Port          = flag.Int("nf.port", 9995, "Sflow/NetFlow/IPFIX listening port")
	Reuse         = flag.Bool("nf.reuserport", false, "Enable so_reuseport for Sflow/NetFlow/IPFIX")
	Workers       = flag.Int("nf.workers", 1, "Number of workers per flow collector")
	MessageFields = flag.String("nf.message.fields", defaultFields, "The list of fields to include in flow messages. Can be any of "+fullFieldList)
	PromPath      = flag.String("nf.prom.listen", "", "Run a promethues metrics collector here")
	MappingFile   = flag.String("nf.mapping", "", "Configuration file for custom netflow mappings")
)

func NewFlowSource(ctx context.Context, proto FlowSource, maxBatchSize int, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi, resolv *resolv.Resolver) (*KentikDriver, error) {
	kt := NewKentikDriver(ctx, proto, maxBatchSize, log, registry, jchfChan, apic, *MessageFields, resolv)
	kt.Infof("Netflow listener running on %s:%d for format %s and a batch size of %d", *Addr, *Port, proto, maxBatchSize)
	kt.Infof("Netflow listener sending fields %s", *MessageFields)

	defer func() {
		if *PromPath != "" {
			http.Handle("/metrics", promhttp.Handler())
			go http.ListenAndServe(*PromPath, nil)
		}
	}()

	// Allow processing of custom ipfix templates here.
	var config utils.ProducerConfig
	if *MappingFile != "" {
		f, err := os.Open(*MappingFile)
		if err != nil {
			kt.Errorf("Cannot load netflow mapping file: %v", err)
			return nil, err
		}
		config, err = utils.LoadMapping(f)
		f.Close()
		if err != nil {
			kt.Errorf("Invalid yaml for netflow mapping file: %v", err)
			return nil, err
		}
	}

	switch proto {
	case Ipfix, Netflow9:
		sNF := &utils.StateNetFlow{
			Format: kt,
			Logger: &KentikLog{l: kt},
			Config: config,
		}
		go func() { // Let this run, returning flow into the kentik transport struct
			err := sNF.FlowRoutine(*Workers, *Addr, *Port, *Reuse)
			if err != nil {
				sNF.Logger.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()
		return kt, nil
	case Sflow:
		sSF := &utils.StateSFlow{
			Format: kt,
			Logger: &KentikLog{l: kt},
			Config: config,
		}
		go func() { // Let this run, returning flow into the kentik transport struct
			err := sSF.FlowRoutine(*Workers, *Addr, *Port, *Reuse)
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
			err := sNFL.FlowRoutine(*Workers, *Addr, *Port, *Reuse)
			if err != nil {
				sNFL.Logger.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()
		return kt, nil
	}
	return nil, fmt.Errorf("Unknown flow format %v", proto)
}
