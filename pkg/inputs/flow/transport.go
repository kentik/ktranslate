package flow

import (
	"context"
	"encoding/binary"
	"net"
	"strings"
	"time"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/ic"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

const (
	DeviceUpdateDuration = 1 * time.Hour
)

type KentikTransport struct {
	logger.ContextL
	jchfChan     chan []*kt.JCHF
	apic         *api.KentikApi
	metrics      *FlowMetric
	fields       []string
	devices      map[string]*kt.Device
	maxBatchSize int
	inputs       chan []*kt.JCHF
}

type FlowMetric struct {
	Flows go_metrics.Meter
}

func NewKentikTransport(ctx context.Context, proto FlowSource, maxBatchSize int, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi, fields string) *KentikTransport {
	kt := KentikTransport{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "flow"}, log),
		jchfChan: jchfChan,
		apic:     apic,
		metrics: &FlowMetric{
			Flows: go_metrics.GetOrRegisterMeter("netflow.flows^fmt="+string(proto), registry),
		},
		fields:       strings.Split(fields, ","),
		devices:      apic.GetDevicesAsMap(0),
		maxBatchSize: maxBatchSize,
		inputs:       make(chan []*kt.JCHF, 10),
	}
	go kt.run(ctx) // Process flows and send them on.
	return &kt
}

func (t *KentikTransport) Publish(msgs []*flowmessage.FlowMessage) {
	t.metrics.Flows.Mark(int64(len(msgs)))
	local := make([]*kt.JCHF, len(msgs))
	for i, m := range msgs {
		local[i] = t.toJCHF(m)
	}

	if len(local) >= t.maxBatchSize {
		t.jchfChan <- local
	} else {
		t.inputs <- local
	}
}

func (t *KentikTransport) Close() {}

func (t *KentikTransport) HttpInfo() map[string]float64 {
	return map[string]float64{
		"Flows": t.metrics.Flows.Rate1(),
	}
}

func (t *KentikTransport) toJCHF(fmsg *flowmessage.FlowMessage) *kt.JCHF {
	srcmac := make([]byte, 8)
	dstmac := make([]byte, 8)
	binary.BigEndian.PutUint64(srcmac, fmsg.SrcMac)
	binary.BigEndian.PutUint64(dstmac, fmsg.DstMac)
	srcmac = srcmac[2:8]
	dstmac = dstmac[2:8]

	in := kt.NewJCHF()
	in.CustomStr = make(map[string]string)
	in.CustomInt = make(map[string]int32)
	in.CustomBigInt = make(map[string]int64)
	in.EventType = kt.KENTIK_EVENT_TYPE
	in.Provider = kt.ProviderFlowDevice
	in.SampleRate = 1
	if dev, ok := t.devices[net.IP(fmsg.SamplerAddress).String()]; ok {
		in.DeviceName = dev.Name
		in.DeviceId = dev.ID
		in.CompanyId = dev.CompanyID
	} else {
		in.DeviceName = net.IP(fmsg.SamplerAddress).String()
	}

	for _, field := range t.fields {
		switch field {
		case "Type":
			in.CustomStr[field] = fmsg.Type.String()
		case "TimeReceived":
			in.Timestamp = int64(fmsg.TimeReceived)
		case "SequenceNum":
			in.CustomBigInt[field] = int64(fmsg.SequenceNum)
		case "SamplingRate":
			if fmsg.SamplingRate > 0 {
				in.SampleRate = uint32(fmsg.SamplingRate)
			}
		case "SamplerAddress":
			in.CustomStr[field] = net.IP(fmsg.SamplerAddress).String()
		case "TimeFlowStart":
			in.CustomBigInt[field] = int64(fmsg.TimeFlowStart)
		case "TimeFlowEnd":
			in.CustomBigInt[field] = int64(fmsg.TimeFlowEnd)
		case "Bytes":
			in.InBytes = fmsg.Bytes
		case "Packets":
			in.InPkts = fmsg.Packets
		case "SrcAddr":
			in.SrcAddr = net.IP(fmsg.SrcAddr).String()
		case "DstAddr":
			in.DstAddr = net.IP(fmsg.DstAddr).String()
		case "Etype":
			in.CustomBigInt[field] = int64(fmsg.Etype)
		case "Proto":
			in.Protocol = ic.PROTO_NAMES[uint16(fmsg.Proto)]
		case "SrcPort":
			in.L4SrcPort = fmsg.SrcPort
		case "DstPort":
			in.L4DstPort = fmsg.DstPort
		case "InIf":
			in.InputPort = kt.IfaceID(fmsg.InIf)
		case "OutIf":
			in.OutputPort = kt.IfaceID(fmsg.OutIf)
		case "SrcMac":
			if fmsg.SrcMac > 0 {
				in.SrcEthMac = net.HardwareAddr(srcmac).String()
			}
		case "DstMac":
			if fmsg.DstMac > 0 {
				in.DstEthMac = net.HardwareAddr(dstmac).String()
			}
		case "SrcVlan":
			in.VlanIn = fmsg.SrcVlan
		case "DstVlan":
			in.VlanOut = fmsg.DstVlan
		case "VlanId":
			in.CustomBigInt[field] = int64(fmsg.VlanId)
		case "IngressVrfID":
			in.CustomBigInt[field] = int64(fmsg.IngressVrfID)
		case "EgressVrfID":
			in.CustomBigInt[field] = int64(fmsg.EgressVrfID)
		case "IPTos":
			in.Tos = fmsg.IPTos
		case "ForwardingStatus":
			in.CustomBigInt[field] = int64(fmsg.ForwardingStatus)
		case "IPTTL":
			in.CustomBigInt[field] = int64(fmsg.IPTTL)
		case "TCPFlags":
			in.TcpFlags = fmsg.TCPFlags
		case "IcmpType":
			in.CustomBigInt[field] = int64(fmsg.IcmpType)
		case "IcmpCode":
			in.CustomBigInt[field] = int64(fmsg.IcmpCode)
		case "IPv6FlowLabel":
			in.CustomBigInt[field] = int64(fmsg.IPv6FlowLabel)
		case "FragmentId":
			in.CustomBigInt[field] = int64(fmsg.FragmentId)
		case "FragmentOffset":
			in.CustomBigInt[field] = int64(fmsg.FragmentOffset)
		case "BiFlowDirection":
			in.CustomBigInt[field] = int64(fmsg.BiFlowDirection)
		case "SrcAS":
			in.SrcAs = fmsg.SrcAS
		case "DstAS":
			in.DstAs = fmsg.DstAS
		case "NextHop":
			in.NextHop = net.IP(fmsg.NextHop).String()
		case "NextHopAS":
			in.DstNextHopAs = fmsg.NextHopAS
		case "SrcNet":
			in.CustomBigInt[field] = int64(fmsg.SrcNet)
		case "DstNet":
			in.CustomBigInt[field] = int64(fmsg.DstNet)
		case "HasEncap":
			if fmsg.HasEncap {
				in.CustomInt[field] = 1
			} else {
				in.CustomInt[field] = 0
			}
		case "SrcAddrEncap":
			if ip := net.IP(fmsg.SrcAddrEncap); ip != nil {
				in.CustomStr[field] = ip.String()
			}
		case "DstAddrEncap":
			if ip := net.IP(fmsg.DstAddrEncap); ip != nil {
				in.CustomStr[field] = ip.String()
			}
		case "ProtoEncap":
			in.CustomBigInt[field] = int64(fmsg.ProtoEncap)
		case "EtypeEncap":
			in.CustomBigInt[field] = int64(fmsg.EtypeEncap)
		case "IPTosEncap":
			in.CustomBigInt[field] = int64(fmsg.IPTosEncap)
		case "IPTTLEncap":
			in.CustomBigInt[field] = int64(fmsg.IPTTLEncap)
		case "IPv6FlowLabelEncap":
			in.CustomBigInt[field] = int64(fmsg.IPv6FlowLabelEncap)
		case "FragmentIdEncap":
			in.CustomBigInt[field] = int64(fmsg.FragmentIdEncap)
		case "FragmentOffsetEncap":
			in.CustomBigInt[field] = int64(fmsg.FragmentOffsetEncap)
		case "HasMPLS":
			if fmsg.HasMPLS {
				in.CustomInt[field] = 1
			} else {
				in.CustomInt[field] = 0
			}
		case "MPLSCount":
			in.CustomBigInt[field] = int64(fmsg.MPLSCount)
		case "MPLS1TTL":
			in.CustomBigInt[field] = int64(fmsg.MPLS1TTL)
		case "MPLS1Label":
			in.CustomBigInt[field] = int64(fmsg.MPLS1Label)
		case "MPLS2TTL":
			in.CustomBigInt[field] = int64(fmsg.MPLS2TTL)
		case "MPLS2Label":
			in.CustomBigInt[field] = int64(fmsg.MPLS2Label)
		case "MPLS3TTL":
			in.CustomBigInt[field] = int64(fmsg.MPLS3TTL)
		case "MPLS3Label":
			in.CustomBigInt[field] = int64(fmsg.MPLS3Label)
		case "MPLSLastTTL":
			in.CustomBigInt[field] = int64(fmsg.MPLSLastTTL)
		case "MPLSLastLabel":
			in.CustomBigInt[field] = int64(fmsg.MPLSLastLabel)
		case "HasPPP":
			if fmsg.HasPPP {
				in.CustomInt[field] = 1
			} else {
				in.CustomInt[field] = 0
			}
		case "PPPAddressControl":
			in.CustomBigInt[field] = int64(fmsg.PPPAddressControl)
		}
	}

	return in
}

func (t *KentikTransport) run(ctx context.Context) {
	sendTicker := time.NewTicker(kt.SendBatchDuration)
	defer sendTicker.Stop()
	deviceTicker := time.NewTicker(DeviceUpdateDuration)
	defer deviceTicker.Stop()
	batch := make([]*kt.JCHF, 0, t.maxBatchSize)

	t.Infof("flow transport running")
	for {
		select {
		case local := <-t.inputs:
			batch = append(batch, local...)
			if len(batch) >= t.maxBatchSize {
				t.jchfChan <- batch
				batch = make([]*kt.JCHF, 0, t.maxBatchSize)
			}
		case <-sendTicker.C:
			if len(batch) > 0 {
				t.jchfChan <- batch
				batch = make([]*kt.JCHF, 0, t.maxBatchSize)
			}
		case <-deviceTicker.C:
			go func() {
				t.Infof("updating device list for flow")
				t.devices = t.apic.GetDevicesAsMap(0)
			}()
		case <-ctx.Done():
			t.Infof("flow transport done")
			return
		}
	}
}
