package elasticsearch

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	jsoniter "github.com/json-iterator/go"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

const (
	InstNameNetflowEvent = "netflow-events"
	InstNameVPCEvent     = "vpc-flow-events"
	InstNameAWSVPCEvent  = "aws-vpc-flow-events"
	actionEntry          = `{"index":{}}`
)

var json = jsoniter.ConfigFastest

type ElasticsearchFormat struct {
	logger.ContextL
	compression kt.Compression
	useGzip     bool
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*ElasticsearchFormat, error) {
	ef := &ElasticsearchFormat{
		ContextL:    logger.NewContextLFromUnderlying(logger.SContext{S: "elasticsearchFormat"}, log),
		compression: compression,
	}

	switch compression {
	case kt.CompressionNone:
		ef.useGzip = false
	case kt.CompressionGzip:
		ef.useGzip = true
	default:
		return nil, fmt.Errorf("Invalid compression (%s): format json only supports none|gzip", compression)
	}

	return ef, nil
}

func (f *ElasticsearchFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	msgsNew := make([]map[string]interface{}, 0, len(msgs))
	for _, msg := range msgs {
		if msg.EventType == kt.KENTIK_EVENT_SNMP {
			continue
		}
		mm := msg.Flatten()
		strip(mm)
		msgsNew = append(msgsNew, mm)
	}

	esBulkData := []string{}
	for _, m := range msgsNew {
		data, err := serialize(m)
		if err != nil {
			return nil, err
		}
		esBulkData = append(esBulkData, data)
	}

	esData := []byte(strings.Join(esBulkData, ""))

	// compression
	if f.useGzip {
		buf := bytes.NewBuffer(serBuf)
		buf.Reset()
		zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
		if err != nil {
			return nil, err
		}

		if _, err := zw.Write(esData); err != nil {
			return nil, err
		}

		if err := zw.Close(); err != nil {
			return nil, err
		}
		esData = buf.Bytes()
	}
	return kt.NewOutputWithProvider(esData, msgs[0].Provider, kt.EventOutput), nil
}

func (f *ElasticsearchFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	msgs := []*kt.JCHF{}

	var r io.Reader
	r = bytes.NewReader(raw.Body)
	if f.useGzip {
		gr, err := gzip.NewReader(bytes.NewBuffer(raw.Body))
		if err != nil {
			return nil, err
		}
		r = gr
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		// check for ES action and ignore
		v := sc.Text()
		if v == actionEntry {
			continue
		}
		msg := &kt.JCHF{}
		if err := json.Unmarshal([]byte(v), &msg); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	values := make([]map[string]interface{}, len(msgs))
	for i, m := range msgs {
		m.SetMap()
		values[i] = m.ToMap()
	}

	return values, nil
}

func (f *ElasticsearchFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	// serialize rolls
	esBulkData := []string{}
	for _, m := range rolls {
		data, err := serialize(m)
		if err != nil {
			return nil, err
		}
		esBulkData = append(esBulkData, data)
	}
	return kt.NewOutputWithProvider([]byte(strings.Join(esBulkData, "")), rolls[0].Provider, kt.RollupOutput), nil
}

func serialize(o interface{}) (string, error) {
	s := actionEntry + "\n"
	data, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	s += string(data) + "\n"
	return s, nil
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
	in["instrumentation.provider"] = kt.InstProvider // Let them know who sent this.
	switch in["provider"] {
	case kt.ProviderVPC:
		switch in["kt.from"] {
		case kt.FromLambda:
			in["instrumentation.name"] = InstNameAWSVPCEvent
		default:
			in["instrumentation.name"] = InstNameVPCEvent
		}
	default:
		in["instrumentation.name"] = InstNameNetflowEvent
	}
	in["collector.name"] = kt.CollectorName
}
