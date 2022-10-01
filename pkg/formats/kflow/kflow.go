package kflow

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"net"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"
	patricia "github.com/kentik/ktranslate/pkg/util/gopatricia/patricia"
	"github.com/kentik/ktranslate/pkg/util/ic"
	model "github.com/kentik/ktranslate/pkg/util/kflow2"

	capn "zombiezen.com/go/capnproto2"
)

const (
	MSG_KEY_PREFIX                    = 80 // This many bytes in every rcv message are for the key.
	KTRANSLATE_PROTO                  = 0
	KTRANSLATE_MAP_PROTO              = 101
	kentikDefaultCapnprotoDecodeLimit = 128 << 20 // 128 MiB
)

type KflowFormat struct {
	logger.ContextL
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*KflowFormat, error) {
	kf := &KflowFormat{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kflowFormat"}, log),
	}

	// For now, we force gzip in the kflow case.
	if compression != kt.CompressionGzip {
		return nil, fmt.Errorf("Invalid compression (%s): format kflow only supports gzip", compression)
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

	msgs, err := root.NewMsgs(int32(len(flows) + 1)) // +1 for the encoded mapping flow.
	if err != nil {
		return nil, err
	}

	ids, err := f.getIds(flows, msgs.At(0), seg) // Set the first msg to be the mapping one.
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

		f.pack(flow, msgs.At(i+1), list, ids) // Offset by 1 here because the first flow is the map.
	}

	root.SetMsgs(msgs)

	key := fmt.Sprintf("%d:%s:%d^", flows[0].CompanyId, flows[0].DeviceName, flows[0].DeviceId)
	cid := make([]byte, MSG_KEY_PREFIX)
	if len(key) < MSG_KEY_PREFIX {
		copy(cid, key)
	}

	buf := bytes.NewBuffer(serBuf)
	z := gzip.NewWriter(buf)
	z.Reset(buf)
	z.Write(cid)

	err = capn.NewPackedEncoder(z).Encode(msg)
	if err != nil {
		return nil, err
	}

	f.Infof("Sending to %s", key)
	z.Close()
	return kt.NewOutputWithProviderAndCompany(buf.Bytes(), flows[0].Provider, flows[0].CompanyId, kt.EventOutput), nil
}

func (f *KflowFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	z, err := gzip.NewReader(bytes.NewBuffer(raw.Body))
	if err != nil {
		return nil, err
	}
	defer z.Close()

	var bodyBufferBytes []byte
	bodyBuffer := bytes.NewBuffer(bodyBufferBytes)
	_, err = bodyBuffer.ReadFrom(z)
	if err != nil {
		return nil, err
	}
	evt := bodyBuffer.Bytes()

	keyP := bytes.Split(evt[0:MSG_KEY_PREFIX], []byte("^"))
	if len(keyP) < 2 {
		return nil, fmt.Errorf("Invalid prefix found for kflow: %s", string(evt[0:MSG_KEY_PREFIX]))
	}

	decoder := capn.NewPackedDecoder(bytes.NewBuffer(evt[MSG_KEY_PREFIX:]))
	decoder.MaxMessageSize = kentikDefaultCapnprotoDecodeLimit
	capnprotoMessage, err := decoder.Decode()
	if err != nil {
		return nil, err
	}

	// unpack flow messages and pass them down
	packedCHF, err := model.ReadRootPackedCHF(capnprotoMessage)
	if err != nil {
		return nil, err
	}

	messages, err := packedCHF.Msgs()
	if err != nil {
		return nil, err
	}

	out := []map[string]interface{}{}
	customMap := map[uint32]string{}
	for i := 0; i < messages.Len(); i++ {
		msg := messages.At(i)
		switch msg.AppProtocol() {
		case KTRANSLATE_PROTO:
			flow := map[string]interface{}{
				"timestamp": msg.Timestamp(),
				"protocol":  ic.PROTO_NAMES[uint16(msg.Protocol())],
				"src_geo":   fmt.Sprintf("%c%c", msg.SrcGeo()>>8, msg.SrcGeo()&0xFF),
				"server_id": keyP[0],
			}

			// Now the addresses.
			var addr net.IP
			if msg.Ipv4DstAddr() > 0 {
				addr = int2ip(msg.Ipv4DstAddr())
			} else {
				ipr, _ := msg.Ipv6DstAddr()
				addr = net.IP(ipr)
			}
			flow["dst_addr"] = addr.String()

			if msg.Ipv4SrcAddr() > 0 {
				addr = int2ip(msg.Ipv4SrcAddr())
			} else {
				ipr, _ := msg.Ipv6SrcAddr()
				addr = net.IP(ipr)
			}
			flow["src_addr"] = addr.String()

			customs, _ := msg.Custom()
			for i, customsLen := 0, customs.Len(); i < customsLen; i++ {
				cust := customs.At(i)
				val := cust.Value()
				if key, ok := customMap[cust.Id()]; !ok {
					continue // Skip because we don't have a key for this.
				} else {
					switch val.Which() {
					case model.Custom_value_Which_uint32Val:
						flow[key] = val.Uint32Val()
					case model.Custom_value_Which_uint64Val:
						flow[key] = val.Uint64Val()
					case model.Custom_value_Which_strVal:
						sv, _ := val.StrVal()
						flow[key] = sv
					}
				}
			}
			out = append(out, flow)

		case KTRANSLATE_MAP_PROTO:
			customs, _ := msg.Custom()
			for i, customsLen := 0, customs.Len(); i < customsLen; i++ {
				cust := customs.At(i)
				val := cust.Value()
				sv, _ := val.StrVal()
				customMap[cust.Id()] = sv
			}
		}
	}

	return out, nil
}

func (f *KflowFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

func (ff *KflowFormat) pack(f *kt.JCHF, kflow model.CHF, list model.Custom_List, ids map[string]uint32) error {
	kflow.SetAppProtocol(KTRANSLATE_PROTO)
	kflow.SetTimestamp(f.Timestamp)
	kflow.SetDstAs(f.DstAs)
	kflow.SetDstGeo(patricia.PackGeo([]byte(f.DstGeo)))
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
	kflow.SetSrcGeo(patricia.PackGeo([]byte(f.SrcGeo)))
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

	sip := net.ParseIP(f.SrcAddr)
	dip := net.ParseIP(f.DstAddr)
	if dip != nil {
		if dip.To4() != nil {
			kflow.SetIpv4DstAddr(binary.BigEndian.Uint32(dip.To4()))
		} else {
			kflow.SetIpv6DstAddr(dip)
		}
	}

	if sip != nil {
		if sip.To4() != nil {
			kflow.SetIpv4SrcAddr(binary.BigEndian.Uint32(sip.To4()))
		} else {
			kflow.SetIpv6SrcAddr(sip)
		}
	}

	next := 0
	for key, val := range f.CustomStr {
		kc := list.At(next)
		kc.SetId(ids[key])
		kc.Value().SetStrVal(val)
		next++
	}
	for key, val := range f.CustomInt {
		kc := list.At(next)
		kc.SetId(ids[key])
		kc.Value().SetUint32Val(uint32(val))
		next++
	}
	for key, val := range f.CustomBigInt {
		kc := list.At(next)
		kc.SetId(ids[key])
		kc.Value().SetUint64Val(uint64(val))
		next++
	}

	if next > 0 {
		kflow.SetCustom(list)
	}

	return nil
}

func (ff *KflowFormat) getIds(flows []*kt.JCHF, kflow model.CHF, seg *capn.Segment) (map[string]uint32, error) {
	ids := map[string]uint32{}

	for _, flow := range flows {
		for key, _ := range flow.CustomStr {
			if _, ok := ids[key]; ok {
				continue
			}
			ids[key] = crc32.ChecksumIEEE([]byte(key))
		}
		for key, _ := range flow.CustomInt {
			if _, ok := ids[key]; ok {
				continue
			}
			ids[key] = crc32.ChecksumIEEE([]byte(key))
		}
		for key, _ := range flow.CustomBigInt {
			if _, ok := ids[key]; ok {
				continue
			}
			ids[key] = crc32.ChecksumIEEE([]byte(key))
		}
	}

	// Now, set up our mapping flow.
	list, err := model.NewCustom_List(seg, int32(len(ids)))
	if err != nil {
		return nil, err
	}

	kflow.SetAppProtocol(KTRANSLATE_MAP_PROTO)
	next := 0
	for k, id := range ids {
		kc := list.At(next)
		kc.SetId(id)
		kc.Value().SetStrVal(k)
		next++
	}
	kflow.SetCustom(list) // Don't forget to save this into the kflow itself.

	// And return the map we used.
	return ids, nil
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
