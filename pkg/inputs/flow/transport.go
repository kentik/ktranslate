package flow

import (
	"context"
	"encoding/binary"
	"net"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/ic"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type KentikTransport struct {
	logger.ContextL
	ctx      context.Context
	jchfChan chan []*kt.JCHF
	apic     *api.KentikApi
	metrics  *FlowMetric
	fields   []string
	devices  map[string]*kt.Device
}

type FlowMetric struct {
	Flows go_metrics.Meter
}

func (t *KentikTransport) Publish(msgs []*flowmessage.FlowMessage) {
	t.metrics.Flows.Mark(int64(len(msgs)))
	jchfs := make([]*kt.JCHF, len(msgs))
	for i, m := range msgs {
		jchfs[i] = t.toJCHF(m)
	}

	if len(jchfs) > 0 {
		t.jchfChan <- jchfs
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
			in.SrcEthMac = net.HardwareAddr(srcmac).String()
		case "DstMac":
			in.DstEthMac = net.HardwareAddr(dstmac).String()
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
