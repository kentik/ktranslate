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

	"github.com/netsampler/goflow2/v2/producer"
	pp "github.com/netsampler/goflow2/v2/producer/proto"
	"github.com/netsampler/goflow2/v2/utils"
	"google.golang.org/protobuf/encoding/protowire"
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
	config       *pp.ProducerConfig
	cfg          *ktranslate.FlowInputConfig
	receiver     *utils.UDPReceiver
	pipe         utils.FlowPipe
	producer     producer.ProducerInterface
	pb2ixd       map[int32]pbInfo
}

// For now, don't handle arrays
type pbInfo struct {
	Name string
	Type string
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

func (t *KentikDriver) SetConfig(c *pp.ProducerConfig) {
	t.config = c
	t.pb2ixd = map[int32]pbInfo{}
	for _, pf := range t.config.Formatter.Protobuf {
		t.pb2ixd[pf.Index] = pbInfo{
			Name: pf.Name,
			Type: pf.Type,
		}
	}
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

// Noop for now. Possibly add later?
func (t *KentikDriver) Send(key, data []byte) error {
	return nil
}

func (t *KentikDriver) Format(data interface{}) ([]byte, []byte, error) {
	msg, ok := data.(*pp.ProtoProducerMessage)
	if !ok {
		return nil, nil, fmt.Errorf("message is not protobuf")
	}
	t.inputs <- t.toJCHF(msg) // Pull out into our own system here.
	return nil, nil, nil
}

func (t *KentikDriver) Close() error {
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

	return nil
}

func (t *KentikDriver) HttpInfo() map[string]float64 {
	flows := map[string]float64{}
	for d, f := range t.metrics {
		flows[d+"_rate"] = f.Flows.Rate1()
		flows[d+"_count"] = float64(f.Flows.Count())
	}
	return flows
}

func (t *KentikDriver) mapCustoms(m *pp.ProtoProducerMessage, in *kt.JCHF) {
	if t.pb2ixd == nil {
		return
	}

	fmr := m.ProtoReflect()
	unk := fmr.GetUnknown()
	var offset int
	for offset < len(unk) {
		num, dataType, length := protowire.ConsumeTag(unk[offset:])
		offset += length
		length = protowire.ConsumeFieldValue(num, dataType, unk[offset:])
		data := unk[offset : offset+length]
		offset += length

		// we check if the index is listed in the config
		if pbField, ok := t.pb2ixd[int32(num)]; ok {
			if dataType == protowire.VarintType {
				v, _ := protowire.ConsumeVarint(data)
				if strings.HasPrefix(pbField.Name, "ip_") {
					in.CustomStr[pbField.Name] = kt.Int2ip(uint32(v)).String()
				} else {
					in.CustomBigInt[pbField.Name] = int64(v)
				}
			} else if dataType == protowire.BytesType {
				if strings.HasPrefix(pbField.Name, "ip_") {
					v, _ := protowire.ConsumeBytes(data)
					if len(v) == 4 {
						in.CustomStr[pbField.Name] = net.IPv4(v[3], v[2], v[1], v[0]).String()
					} else {
						in.CustomStr[pbField.Name] = net.IP(v).String()
					}
				} else {
					v, _ := protowire.ConsumeString(data)
					in.CustomStr[pbField.Name] = v
				}
			} else {
				continue
			}
		}
	}
}

func (t *KentikDriver) toJCHF(fmsg *pp.ProtoProducerMessage) *kt.JCHF {
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

	t.mapCustoms(fmsg, in)
	ocflowDir := in.CustomBigInt["flow_direction"]
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
		case "FlowDirection":
			if v, ok := in.CustomBigInt["flow_direction"]; ok {
				delete(in.CustomBigInt, "flow_direction")
				if v == 1 {
					in.CustomStr["FlowDirection"] = "egress"
				} else {
					in.CustomStr["FlowDirection"] = "ingress"
				}
			}
		case "SamplerAddress":
			in.CustomStr[field] = net.IP(fmsg.SamplerAddress).String()
		case "TimeFlowStart":
			in.CustomBigInt[field] = int64(fmsg.TimeFlowStartNs)
		case "TimeFlowEnd":
			in.CustomBigInt[field] = int64(fmsg.TimeFlowEndNs)
		case "Bytes":
			if ocflowDir == 1 {
				in.OutBytes = fmsg.Bytes
			} else {
				in.InBytes = fmsg.Bytes
			}
		case "Packets":
			if ocflowDir == 1 {
				in.OutPkts = fmsg.Packets
			} else {
				in.InPkts = fmsg.Packets
			}
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
		}
	}

	// Default in direction if its not already set. Assume its ingress unless otherwise told.
	if _, ok := in.CustomStr["FlowDirection"]; !ok {
		in.CustomStr["FlowDirection"] = "ingress"
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
