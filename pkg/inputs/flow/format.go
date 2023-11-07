package flow

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/ic"
	"github.com/kentik/ktranslate/pkg/util/resolv"

	flowmessage "github.com/netsampler/goflow2/v2/pb"
	"github.com/netsampler/goflow2/v2/producer"
	"github.com/netsampler/goflow2/v2/utils"
)

const (
	DeviceUpdateDuration = 1 * time.Hour
)

type KentikDriver struct {
	logger.ContextL
	sync.RWMutex
	jchfChan     chan []*kt.JCHF
	apic         *api.KentikApi
	metrics      map[string]*FlowMetric
	fields       []string
	devices      map[string]*kt.Device
	maxBatchSize int
	inputs       chan *kt.JCHF
	proto        FlowSource
	registry     go_metrics.Registry
	resolv       *resolv.Resolver
	ctx          context.Context
	config       EntConfig
	cfg          *ktranslate.FlowInputConfig
	receiver     *utils.UDPReceiver
	pipe         utils.FlowPipe
	producer     producer.ProducerInterface
}

type FlowMetric struct {
	Flows go_metrics.Meter
}

func NewKentikDriver(ctx context.Context, proto FlowSource, maxBatchSize int, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi, fields string, resolv *resolv.Resolver, cfg *ktranslate.FlowInputConfig) *KentikDriver {
	kt := KentikDriver{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "flow"}, log),
		jchfChan:     jchfChan,
		apic:         apic,
		metrics:      map[string]*FlowMetric{},
		fields:       strings.Split(fields, ","),
		devices:      apic.GetDevicesAsMap(0),
		maxBatchSize: maxBatchSize,
		inputs:       make(chan *kt.JCHF, maxBatchSize),
		proto:        proto,
		registry:     registry,
		resolv:       resolv,
		ctx:          ctx,
		cfg:          cfg,
	}
	go kt.run(ctx) // Process flows and send them on.
	return &kt
}

func (t *KentikDriver) SetConfig(c EntConfig) {
	t.config = c
}

func (t *KentikDriver) Name() string {
	return "Kentik CHF"
}

func (t *KentikDriver) Init() error {
	return nil
}

func (t *KentikDriver) Prepare() error {
	return nil
}

func (t *KentikDriver) Format(data interface{}) ([]byte, []byte, error) {
	msg, ok := data.(*flowmessage.FlowMessage)
	if !ok {
		return nil, nil, fmt.Errorf("message is not protobuf")
	}
	t.inputs <- t.toJCHF(msg) // Pull out into our own system here.
	return nil, nil, nil
}

func (t *KentikDriver) Close() {
	if t.receiver != nil {
		if err := t.receiver.Stop(); err != nil {
			t.Errorf("Error stopping flow reciever: %v", err)
		}
	}
	if t.pipe != nil {
		t.pipe.Close()
	}
	if t.producer != nil {
		t.producer.Close()
	}
}

func (t *KentikDriver) HttpInfo() map[string]float64 {
	flows := map[string]float64{}
	for d, f := range t.metrics {
		flows[d] = f.Flows.Rate1()
	}
	return flows
}

func (t *KentikDriver) toJCHF(fmsg *flowmessage.FlowMessage) *kt.JCHF {
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
	in.ApplySample = true

	// We have enough traffic here now to require a locking primative.
	t.RLock()
	if dev, ok := t.devices[net.IP(fmsg.SamplerAddress).String()]; ok {
		in.DeviceName = dev.Name // Copy in any of these info we get
		in.DeviceId = dev.ID
		in.CompanyId = dev.CompanyID
		in.SampleRate = dev.SampleRate
		dev.SetUserTags(in.CustomStr)
	} else {
		in.DeviceName = net.IP(fmsg.SamplerAddress).String()
		if t.resolv != nil {
			dm := t.resolv.Resolve(t.ctx, in.DeviceName, true)
			if dm != "" {
				in.DeviceName = dm
			}
		}
	}

	if _, ok := t.metrics[in.DeviceName]; !ok {
		t.RUnlock() // Annoying lock dance, we assume this will not happen that much.
		t.Lock()
		t.metrics[in.DeviceName] = &FlowMetric{Flows: go_metrics.GetOrRegisterMeter(fmt.Sprintf("netflow.flows^fmt=%s^device_name=%s^force=true", string(t.proto), in.DeviceName), t.registry)}
		t.Unlock()
		t.RLock()
	}
	t.metrics[in.DeviceName].Flows.Mark(1)
	t.RUnlock()

	if in.SampleRate == 0 {
		in.SampleRate = 1
	}

	for _, field := range t.fields {
		switch field {
		case "Type":
			in.CustomStr[field] = fmsg.Type.String()
		case "TimeReceived":
			in.Timestamp = int64(fmsg.TimeReceivedNs)
		case "SequenceNum":
			in.CustomBigInt[field] = int64(fmsg.SequenceNum)
		case "SamplingRate":
			if fmsg.SamplingRate > 0 {
				in.SampleRate = uint32(fmsg.SamplingRate)
			}
		case "SamplerAddress":
			in.CustomStr[field] = net.IP(fmsg.SamplerAddress).String()
		case "TimeFlowStart":
			in.CustomBigInt[field] = int64(fmsg.TimeFlowStartNs)
		case "TimeFlowEnd":
			in.CustomBigInt[field] = int64(fmsg.TimeFlowEndNs)
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
		case "IPTos":
			in.Tos = fmsg.IpTos
		case "ForwardingStatus":
			in.CustomBigInt[field] = int64(fmsg.ForwardingStatus)
		case "IPTTL":
			in.CustomBigInt[field] = int64(fmsg.IpTtl)
		case "TCPFlags":
			in.TcpFlags = fmsg.TcpFlags
		case "IcmpType":
			in.CustomBigInt[field] = int64(fmsg.IcmpType)
		case "IcmpCode":
			in.CustomBigInt[field] = int64(fmsg.IcmpCode)
		case "IPv6FlowLabel":
			in.CustomBigInt[field] = int64(fmsg.Ipv6FlowLabel)
		case "FragmentId":
			in.CustomBigInt[field] = int64(fmsg.FragmentId)
		case "FragmentOffset":
			in.CustomBigInt[field] = int64(fmsg.FragmentOffset)
		case "SrcAS":
			in.SrcAs = fmsg.SrcAs
		case "DstAS":
			in.DstAs = fmsg.DstAs
		case "NextHop":
			in.NextHop = net.IP(fmsg.NextHop).String()
		case "NextHopAS":
			in.DstNextHopAs = fmsg.NextHopAs
		case "SrcNet":
			in.CustomBigInt[field] = int64(fmsg.SrcNet)
		case "DstNet":
			in.CustomBigInt[field] = int64(fmsg.DstNet)
		case "MPLSCount":
			in.CustomBigInt[field] = int64(len(fmsg.MplsTtl))
			/**
			case "CustomInteger1":
				in.CustomBigInt[t.config.NameMap[field]] = int64(fmsg.CustomInteger1)
			case "CustomInteger2":
				in.CustomBigInt[t.config.NameMap[field]] = int64(fmsg.CustomInteger2)
			case "CustomInteger3":
				in.CustomBigInt[t.config.NameMap[field]] = int64(fmsg.CustomInteger3)
			case "CustomInteger4":
				in.CustomBigInt[t.config.NameMap[field]] = int64(fmsg.CustomInteger4)
			case "CustomInteger5":
				in.CustomBigInt[t.config.NameMap[field]] = int64(fmsg.CustomInteger5)
			case "CustomBytes1":
				in.CustomStr[t.config.NameMap[field]] = fmt.Sprintf("%s", string(fmsg.CustomBytes1))
			case "CustomBytes2":
				in.CustomStr[t.config.NameMap[field]] = fmt.Sprintf("%s", string(fmsg.CustomBytes2))
			case "CustomBytes3":
				in.CustomStr[t.config.NameMap[field]] = fmt.Sprintf("%s", string(fmsg.CustomBytes3))
			case "CustomBytes4":
				in.CustomStr[t.config.NameMap[field]] = fmt.Sprintf("%s", string(fmsg.CustomBytes4))
			case "CustomBytes5":
				in.CustomStr[t.config.NameMap[field]] = fmt.Sprintf("%s", string(fmsg.CustomBytes5))
			case "CustomIPv41":
				in.CustomStr[t.config.NameMap[field]] = kt.Int2ip(uint32(fmsg.CustomInteger1)).String()
			case "CustomIPv42":
				in.CustomStr[t.config.NameMap[field]] = kt.Int2ip(uint32(fmsg.CustomInteger2)).String()
			case "CustomIP1":
				if len(fmsg.CustomBytes1) == 4 {
					in.CustomStr[t.config.NameMap[field]] = net.IPv4(fmsg.CustomBytes1[3], fmsg.CustomBytes1[2], fmsg.CustomBytes1[1], fmsg.CustomBytes1[0]).String()
				} else {
					in.CustomStr[t.config.NameMap[field]] = net.IP(fmsg.CustomBytes1).String()
				}
			case "CustomIP2":
				if len(fmsg.CustomBytes2) == 4 {
					in.CustomStr[t.config.NameMap[field]] = net.IPv4(fmsg.CustomBytes2[3], fmsg.CustomBytes2[2], fmsg.CustomBytes2[1], fmsg.CustomBytes2[0]).String()
				} else {
					in.CustomStr[t.config.NameMap[field]] = net.IP(fmsg.CustomBytes2).String()
				}
			*/
		}
	}

	// Now add some combo fields.
	in.CustomStr["src_endpoint"] = in.SrcAddr + ":" + strconv.Itoa(int(in.L4SrcPort))
	in.CustomStr["dst_endpoint"] = in.DstAddr + ":" + strconv.Itoa(int(in.L4DstPort))

	return in
}

func (t *KentikDriver) run(ctx context.Context) {
	sendTicker := time.NewTicker(kt.SendBatchDuration)
	defer sendTicker.Stop()
	deviceTicker := time.NewTicker(DeviceUpdateDuration)
	defer deviceTicker.Stop()
	batch := make([]*kt.JCHF, 0, t.maxBatchSize)

	t.Infof("kentik driver running")
	for {
		select {
		case local := <-t.inputs:
			batch = append(batch, local)
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
				t.Infof("Updating the network flow device list.")
				t.devices = t.apic.GetDevicesAsMap(0)
			}()
		case <-ctx.Done():
			t.Infof("kentik driver done")
			return
		}
	}
}
