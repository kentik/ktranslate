package netflow

import (
	"bytes"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/netflow/netflow9"
)

func (f *NetflowFormat) pack9(flows []*kt.JCHF, serBuf []byte) ([]byte, error) {
	var err error

	header := netflow9.PacketHeader{
		Version: netflow9.Version,
		Count:   uint16(len(flows) + 1),
	}

	buf := bytes.NewBuffer(serBuf)
	buf.Reset()

	err = write(buf, &header)
	if err != nil {
		return nil, err
	}

	tmp, err := encodeTemplate(nil, templateId)
	if err != nil {
		return nil, err
	}

	err = f.writeWithHeader(buf, 0, tmp)
	if err != nil {
		return nil, err
	}

	for i := range flows {
		flow, err := encodeFlow(flows[i])
		if err != nil {
			return nil, err
		}

		var padding [4]byte
		if n := flow.Len() % 4; n > 0 {
			flow.Write(padding[:n])
		}

		err = f.writeWithHeader(buf, templateId, flow)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (f *NetflowFormat) decodeV9(msg *netflow9.Packet, raw []byte) []*kt.JCHF {
	chfs := make([]*kt.JCHF, 0)

	for _, dt := range msg.TemplateFlowSets {
		for _, tr := range dt.Records {
			f.addV9Template(msg.Header.SourceID, tr)
		}
	}

	for _, ds := range msg.DataFlowSets {
		for _, dr := range ds.Records {
			chf := f.decodeV9Record(dr, raw)
			if chf != nil {
				chfs = append(chfs, chf)
			}
		}
	}

	return chfs
}

func (ipf *NetflowFormat) decodeV9Record(dr netflow9.DataRecord, raw []byte) *kt.JCHF {
	res := kt.NewJCHF()
	for _, f := range dr.Fields {
		if f.Translated != nil {
			// On an error decoding a field value (wrong type or value out-of-range), throw away the entire flow record
			ok := ipf.setCHFField(res, f.Translated.Type, f.Translated.Value, raw, 0)
			if !ok {
				return nil
			}
		}
	}
	return res
}
