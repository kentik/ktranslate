package flow

import (
	"context"
	"flag"
	"fmt"
	"strings"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/cloudflare/goflow/v3/utils"
)

type FlowSource string

const (
	Ipfix    FlowSource = "ipfix"
	Sflow               = "sflow"
	Netflow5            = "netflow5"
	Netflow9            = "netflow9"

	defaultFields = "Type,TimeReceived,SequenceNum,SamplingRate,SamplerAddress,TimeFlowStart,TimeFlowEnd,Bytes,Packets,SrcAddr,DstAddr,Etype,Proto,SrcPort,DstPort,InIf,OutIf,SrcMac,DstMac,SrcVlan,DstVlan,VlanId,IngressVrfID,EgressVrfID,IPTos,ForwardingStatus,IPTTL,TCPFlags,IcmpType,IcmpCode,IPv6FlowLabel,FragmentId,FragmentOffset,BiFlowDirection,SrcAS,DstAS,NextHop,NextHopAS,SrcNet,DstNet,HasEncap,SrcAddrEncap,DstAddrEncap,ProtoEncap,EtypeEncap,IPTosEncap,IPTTLEncap,IPv6FlowLabelEncap,FragmentIdEncap,FragmentOffsetEncap,HasMPLS,MPLSCount,MPLS1TTL,MPLS1Label,MPLS2TTL,MPLS2Label,MPLS3TTL,MPLS3Label,MPLSLastTTL,MPLSLastLabel,HasPPP,PPPAddressControl"
)

var (
	Addr          = flag.String("nf.addr", "", "Sflow/NetFlow/IPFIX listening address")
	Port          = flag.Int("nf.port", 2055, "Sflow/NetFlow/IPFIX listening port")
	Reuse         = flag.Bool("nf.reuserport", false, "Enable so_reuseport for Sflow/NetFlow/IPFIX")
	Workers       = flag.Int("nf.workers", 1, "Number of workers per flow collector")
	MessageFields = flag.String("nf.message.fields", defaultFields, "The list of fields to include in flow messages")
)

func NewFlowSource(ctx context.Context, proto FlowSource, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi) (*KentikTransport, error) {
	kt := KentikTransport{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "flow"}, log),
		ctx:      ctx,
		jchfChan: jchfChan,
		apic:     apic,
		metrics: &FlowMetric{
			Flows: go_metrics.GetOrRegisterMeter("netflow.flows^fmt="+string(proto), registry),
		},
		fields: strings.Split(*MessageFields, ","),
	}

	kt.Infof("Netflow listener running on %s:%d for format %s", *Addr, *Port, proto)

	switch proto {
	case Ipfix, Netflow9:
		sNF := &utils.StateNetFlow{
			Transport: &kt,
			Logger:    &KentikLog{ContextL: kt},
		}
		go func() { // Let this run, returning flow into the kentik transport struct
			err := sNF.FlowRoutine(*Workers, *Addr, *Port, *Reuse)
			if err != nil {
				sNF.Logger.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()
		return &kt, nil
	case Sflow:
		sSF := &utils.StateSFlow{
			Transport: &kt,
			Logger:    &KentikLog{ContextL: kt},
		}
		go func() { // Let this run, returning flow into the kentik transport struct
			err := sSF.FlowRoutine(*Workers, *Addr, *Port, *Reuse)
			if err != nil {
				sSF.Logger.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()
		return &kt, nil
	case Netflow5:
		sNFL := &utils.StateNFLegacy{
			Transport: &kt,
			Logger:    &KentikLog{ContextL: kt},
		}
		go func() { // Let this run, returning flow into the kentik transport struct
			err := sNFL.FlowRoutine(*Workers, *Addr, *Port, *Reuse)
			if err != nil {
				sNFL.Logger.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()
		return &kt, nil
	}
	return nil, fmt.Errorf("Unknown flow format %v", proto)
}
