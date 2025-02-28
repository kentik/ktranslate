package elasticsearch

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	jsoniter "github.com/json-iterator/go"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

const (
	actionEntryFormat = `{"%s":{}}`
)

var (
	actionEntrySet string
)

func init() {
	flag.StringVar(&actionEntrySet, "elastic.action", "index", "Use this action when sending to elastic.")
}

var json = jsoniter.ConfigFastest

type ElasticsearchFormat struct {
	logger.ContextL
	compression kt.Compression
	useGzip     bool
	action      string
}

func NewFormat(log logger.Underlying, compression kt.Compression, cfg *ktranslate.ElasticFormatConfig) (*ElasticsearchFormat, error) {
	ef := &ElasticsearchFormat{
		ContextL:    logger.NewContextLFromUnderlying(logger.SContext{S: "elasticsearchFormat"}, log),
		compression: compression,
		action:      fmt.Sprintf(actionEntryFormat, cfg.Action),
	}

	switch compression {
	case kt.CompressionNone:
		ef.useGzip = false
	case kt.CompressionGzip:
		ef.useGzip = true
	default:
		return nil, fmt.Errorf("Invalid compression (%s): format json only supports none|gzip", compression)
	}

	ef.Infof("Using action %s", ef.action)

	return ef, nil
}

func (f *ElasticsearchFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	msgsNew := make([]map[string]interface{}, 0, len(msgs))
	for _, msg := range msgs {
		switch msg.EventType {
		case kt.KENTIK_EVENT_SYNTH, kt.KENTIK_EVENT_TRACE:
			msgsNew = append(msgsNew, f.handleSynth(msg))
		default:
			mm := msg.Flatten()
			strip(mm)
			msgsNew = append(msgsNew, mm)
		}
	}

	esBulkData := []string{}
	for _, m := range msgsNew {
		data, err := serialize(m, f.action)
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
		if v == f.action {
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
		data, err := serialize(m, f.action)
		if err != nil {
			return nil, err
		}
		esBulkData = append(esBulkData, data)
	}
	return kt.NewOutputWithProvider([]byte(strings.Join(esBulkData, "")), rolls[0].Provider, kt.RollupOutput), nil
}

func (f *ElasticsearchFormat) handleSynth(in *kt.JCHF) map[string]interface{} {
	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := in.Flatten()

	// Map metrics to better names.
	for m, name := range metrics {
		if _, ok := attr[m]; ok {
			attr[name.Name] = attr[m]
			delete(attr, m)
		}
	}

	// If there's a traceroute, try to unserialize this one.
	if raw, ok := attr["error_cause/trace_route"]; ok {
		trace := []interface{}{}
		if err := json.Unmarshal([]byte(raw.(string)), &trace); err == nil {
			attr["trace_route"] = trace
			delete(attr, "error_cause/trace_route")
		}
	}

	strip(attr) // Take out the rest here.

	return attr
}

func serialize(o interface{}, action string) (string, error) {
	s := action + "\n"
	data, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	s += string(data) + "\n"
	return s, nil
}

var (
	DroppedAttrs = map[string]bool{
		"sampled_packet_size": true,
		"lat/long_dest":       true,
		"member_id":           true,
		"dst_eth_mac":         true,
		"src_eth_mac":         true,
		"ult_exit_port":       true,
		"app_protocol":        true,
		"dst_route_prefix":    true,
		"src_route_prefix":    true,
		"trf_termination":     true,
		"simple_trf_prof":     true,
	}

	KeepAttrs = map[string]bool{
		"lost":           true,
		"sent":           true,
		"src_geo_city":   true,
		"src_geo_region": true,
		"dst_geo_city":   true,
		"dst_geo_region": true,
	}
)

func strip(in map[string]interface{}) {
	for k, v := range in {
		if DroppedAttrs[k] {
			delete(in, k) // Skip.
			continue
		}
		if KeepAttrs[k] {
			continue // Always pass along even if 0.
		}
		if k == "timestamp" { // Convert to the ES native timestamp format.
			if ints, ok := v.(int64); ok {
				ts := time.Unix(ints, 0)
				in["timestamp_str"] = ts.Format(time.RFC3339)
			}
			continue
		}
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
	}
	in["instrumentation.provider"] = kt.InstProvider // Let them know who sent this.
	in["collector.name"] = kt.CollectorName
}
