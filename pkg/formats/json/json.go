package json

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

const (
	InstNameNetflowEvent = "netflow-events"
	InstNameVPCEvent     = "vpc-logs-events"
)

type JsonFormat struct {
	logger.ContextL
	compression kt.Compression
	doGz        bool
	doFlatten   bool
}

func NewFormat(log logger.Underlying, compression kt.Compression, doFlatten bool) (*JsonFormat, error) {
	jf := &JsonFormat{
		compression: compression,
		ContextL:    logger.NewContextLFromUnderlying(logger.SContext{S: "jsonFormat"}, log),
		doGz:        false,
		doFlatten:   doFlatten,
	}

	switch compression {
	case kt.CompressionNone:
		jf.doGz = false
	case kt.CompressionGzip:
		jf.doGz = true
	default:
		return nil, fmt.Errorf("Invalid compression (%s): format json only supports none|gzip", compression)
	}

	return jf, nil
}

func (f *JsonFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	var target []byte
	if f.doFlatten {
		msgsNew := make([]map[string]interface{}, 0, len(msgs))
		for _, msg := range msgs {
			if msg.EventType == kt.KENTIK_EVENT_SNMP {
				continue
			}
			mm := msg.Flatten()
			strip(mm)
			msgsNew = append(msgsNew, mm)
		}
		t, err := json.Marshal(msgsNew)
		if err != nil {
			return nil, err
		}
		target = t
	} else {
		t, err := json.Marshal(msgs)
		if err != nil {
			return nil, err
		}
		target = t
	}

	if !f.doGz {
		return kt.NewOutputWithProvider(target, msgs[0].Provider, kt.EventOutput), nil
	}

	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(target)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return kt.NewOutputWithProvider(buf.Bytes(), msgs[0].Provider, kt.EventOutput), nil
}

func (f *JsonFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	msgs := []*kt.JCHF{}
	var err error

	if !f.doGz {
		err = json.Unmarshal(raw.Body, &msgs)
	} else {
		r, err := gzip.NewReader(bytes.NewBuffer(raw.Body))
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(r).Decode(&msgs)
	}
	if err != nil {
		return nil, err
	}

	values := make([]map[string]interface{}, len(msgs))
	for i, m := range msgs {
		m.SetMap()
		values[i] = m.ToMap()
	}

	return values, err
}

func (f *JsonFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	if !f.doGz {
		res, err := json.Marshal(rolls)
		return kt.NewOutputWithProvider(res, rolls[0].Provider, kt.RollupOutput), err
	}

	serBuf := make([]byte, 0)
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(rolls)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(b)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return kt.NewOutputWithProvider(buf.Bytes(), rolls[0].Provider, kt.RollupOutput), nil
}

func strip(in map[string]interface{}) {
	for k, v := range in {
		switch tv := v.(type) {
		case string:
			if tv == "" || tv == "-" || tv == "--" {
				delete(in, k)
			}
		case int32:
			if tv == 0 {
				delete(in, k)
			}
		case int64:
			if tv == 0 {
				delete(in, k)
			}
		}
		if _, ok := util.DroppedAttrs[k]; ok {
			delete(in, k)
		}
	}
	in["instrumentation.provider"] = kt.InstProvider  // Let them know who sent this.
	in["instrumentation.name"] = InstNameNetflowEvent // @TODO -- think about how to handle VPC
	in["collector.name"] = kt.CollectorName
}
