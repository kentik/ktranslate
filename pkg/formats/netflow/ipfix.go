package netflow

import (
	"bytes"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/netflow/ipfix"
)

func (f *NetflowFormat) packIpfix(msgs []*kt.JCHF, serBuf []byte) ([]byte, error) {

	// Quick guard.
	if len(msgs) == 0 {
		return nil, nil
	}

	dataBuf := &bytes.Buffer{}
	tmp, err := encodeTemplate(msgs[0], templateId) // Pass any options in here.
	if err != nil {
		return nil, err
	}
	err = f.writeWithHeader(dataBuf, templateSet, tmp)
	if err != nil {
		return nil, err
	}

	for i := range msgs {
		flow, err := encodeFlow(msgs[i])
		if err != nil {
			return nil, err
		}

		var padding [4]byte
		if n := flow.Len() % 4; n > 0 {
			flow.Write(padding[:n])
		}

		err = f.writeWithHeader(dataBuf, templateId, flow)
		if err != nil {
			return nil, err
		}
	}

	// Have to get the length before we can write the header. But now we can start doing thing.
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()

	header := ipfix.MessageHeader{
		Version: f.version,
		Length:  uint16(dataBuf.Len() + 16), // Payload + msgheader length
	}
	err = write(buf, &header)

	if err != nil {
		return nil, err
	}

	_, err = buf.Write(dataBuf.Bytes())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (f *NetflowFormat) decodeIpfix(msg *ipfix.Message, raw []byte) []*kt.JCHF {
	chfs := make([]*kt.JCHF, 0)

	for _, dt := range msg.TemplateSets {
		for _, tr := range dt.Records {
			f.addIpfixTemplate(msg.Header.ObservationDomainID, tr)
		}
	}

	for _, ds := range msg.DataSets {
		for _, dr := range ds.Records {
			chf := f.decodeIpfixRecord(dr, raw)
			if chf != nil {
				chfs = append(chfs, chf)
			}
		}
	}

	return chfs
}

func (ipf *NetflowFormat) decodeIpfixRecord(dr ipfix.DataRecord, raw []byte) *kt.JCHF {
	res := kt.NewJCHF()
	for _, f := range dr.Fields {
		if f.Translated != nil {
			// On an error decoding a field value (wrong type or value out-of-range), throw away the entire flow record
			ok := ipf.setCHFField(res, f.Translated.InformationElementID, f.Translated.Value, raw, f.Translated.EnterpriseNumber)
			if !ok {
				return nil
			}
		}
	}
	return res
}
