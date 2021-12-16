package kflow

import (
	"bytes"
	"compress/gzip"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"
	"github.com/kentik/ktranslate/pkg/util/ic"
	model "github.com/kentik/ktranslate/pkg/util/kflow2"

	capn "zombiezen.com/go/capnproto2"
)

const (
	MSG_KEY_PREFIX    = 80 // This many bytes in every rcv message are for the key.
	KFLOW_VERSION_ONE = '1'
)

type KflowFormat struct {
	logger.ContextL
	ids    map[string]uint32
	nextID uint32
}

func NewFormat(log logger.Underlying) (*KflowFormat, error) {
	kf := &KflowFormat{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kflowFormat"}, log),
		ids:      map[string]uint32{},
		nextID:   100,
	}

	return kf, nil
}

func (f *KflowFormat) To(flows []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	if len(flows) == 0 {
		return nil, nil
	}

	msg, seg, err := capn.NewMessage(capn.SingleSegment(nil))
	if err != nil {
		return nil, err
	}

	root, err := model.NewRootPackedCHF(seg)
	if err != nil {
		return nil, err
	}

	msgs, err := root.NewMsgs(int32(len(flows)))
	if err != nil {
		return nil, err
	}

	for i, flow := range flows {
		var list model.Custom_List
		n := len(flow.CustomStr) + len(flow.CustomInt) + len(flow.CustomBigInt)
		if n > 0 {
			if list, err = model.NewCustom_List(seg, int32(n)); err != nil {
				return nil, err
			}
		}

		f.pack(flow, msgs.At(i), list)
	}

	root.SetMsgs(msgs)

	cid := [80]byte{}
	buf := bytes.NewBuffer(serBuf)
	z := gzip.NewWriter(buf)
	z.Reset(buf)
	z.Write(cid[:])

	err = capn.NewPackedEncoder(z).Encode(msg)
	if err != nil {
		return nil, err
	}

	z.Close()
	return kt.NewOutputWithProvider(serBuf, flows[0].Provider, kt.EventOutput), nil
}

func (f *KflowFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	return nil, nil
}

func (f *KflowFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

func (ff *KflowFormat) pack(f *kt.JCHF, kflow model.CHF, list model.Custom_List) error {
	kflow.SetTimestamp(f.Timestamp)
	kflow.SetDstAs(f.DstAs)
	kflow.SetHeaderLen(f.HeaderLen)
	kflow.SetInBytes(f.InBytes)
	kflow.SetInPkts(f.InPkts)
	kflow.SetInputPort(uint32(f.InputPort))
	kflow.SetIpSize(f.IpSize)
	kflow.SetL4DstPort(f.L4DstPort)
	kflow.SetL4SrcPort(f.L4SrcPort)
	kflow.SetOutputPort(uint32(f.OutputPort))
	kflow.SetProtocol(uint32(ic.PROTO_NUMS[f.Protocol]))
	kflow.SetSampledPacketSize(f.SampledPacketSize)
	kflow.SetSrcAs(f.SrcAs)
	kflow.SetTcpFlags(f.TcpFlags)
	kflow.SetTos(f.Tos)
	kflow.SetVlanIn(f.VlanIn)
	kflow.SetVlanOut(f.VlanOut)
	kflow.SetMplsType(f.MplsType)
	kflow.SetOutBytes(f.OutBytes)
	kflow.SetOutPkts(f.OutPkts)
	kflow.SetTcpRetransmit(f.TcpRetransmit)
	kflow.SetSampleRate(f.SampleRate)
	kflow.SetDeviceId(uint32(f.DeviceId))
	kflow.SetSrcNextHopAs(f.SrcNextHopAs)
	kflow.SetDstNextHopAs(f.DstNextHopAs)
	kflow.SetSrcSecondAsn(f.SrcSecondAsn)
	kflow.SetDstSecondAsn(f.DstSecondAsn)
	kflow.SetSrcThirdAsn(f.SrcThirdAsn)
	kflow.SetDstThirdAsn(f.DstThirdAsn)

	next := 0
	for key, val := range f.CustomStr {
		kc := list.At(next)
		kc.SetId(ff.getId(key))
		kc.Value().SetStrVal(val)
		next++
	}
	for key, val := range f.CustomInt {
		kc := list.At(next)
		kc.SetId(ff.getId(key))
		kc.Value().SetUint32Val(uint32(val))
		next++
	}
	for key, val := range f.CustomBigInt {
		kc := list.At(next)
		kc.SetId(ff.getId(key))
		kc.Value().SetUint64Val(uint64(val))
		next++
	}

	return nil
}

// @TODO, make this work in a thread safe way and across restarts?
func (ff *KflowFormat) getId(key string) uint32 {
	if id, ok := ff.ids[key]; ok {
		return id
	}

	ff.ids[key] = ff.nextID
	ff.nextID++
	return ff.ids[key]
}
