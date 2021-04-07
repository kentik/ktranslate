package json

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
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

func (f *JsonFormat) To(msgs []*kt.JCHF, serBuf []byte) ([]byte, error) {
	var target []byte
	if f.doFlatten {
		msgsNew := make([]map[string]interface{}, len(msgs))
		for i, msg := range msgs {
			msgsNew[i] = msg.Flatten()
			strip(msgsNew[i])
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
		return target, nil
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

	return buf.Bytes(), nil
}

func (f *JsonFormat) From(raw []byte) ([]map[string]interface{}, error) {
	msgs := []*kt.JCHF{}
	var err error

	if !f.doGz {
		err = json.Unmarshal(raw, &msgs)
	} else {
		r, err := gzip.NewReader(bytes.NewBuffer(raw))
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

func (f *JsonFormat) Rollup(rolls []rollup.Rollup) ([]byte, error) {
	if !f.doGz {
		return json.Marshal(rolls)
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

	return buf.Bytes(), nil
}

func strip(in map[string]interface{}) {
	for k, v := range in {
		switch tv := v.(type) {
		case string:
			if tv == "" {
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
	}
}
