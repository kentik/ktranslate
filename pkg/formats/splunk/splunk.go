package splunk

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type SplunkFormat struct {
	logger.ContextL
	compression  kt.Compression
	doGz         bool
	lastMetadata map[string]*kt.LastMetadata
	invalids     map[string]bool
	mux          sync.RWMutex

	EventChan chan []byte
}

type SplunkMetric struct {
	Host       string                 `json:"host"`
	SourceType string                 `json:"sourcetype"`
	Timestamp  int64                  `json:"time"`
	Event      map[string]interface{} `json:"event"`
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*SplunkFormat, error) {
	jf := &SplunkFormat{
		compression:  compression,
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "nrmFormat"}, log),
		doGz:         false,
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
		EventChan:    make(chan []byte, 100), // Used for sending events to the event API.
	}

	switch compression {
	case kt.CompressionNone:
		jf.doGz = false
	case kt.CompressionGzip:
		jf.doGz = true
	default:
		return nil, fmt.Errorf("Invalid compression (%s): format nrm only supports none|gzip", compression)
	}

	return jf, nil
}

func (f *SplunkFormat) To(msgs []*kt.JCHF, serBuf []byte) ([]byte, error) {
	ms := make([]SplunkMetric, 0, len(msgs)*4)
	ct := time.Now().Unix()
	for _, m := range msgs {
		ms = append(ms, f.toSplunkMetric(m, ct)...)
	}

	if len(ms) == 0 {
		return nil, nil
	}

	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	for _, m := range ms {
		target, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		buf.Write(target)
		buf.WriteString("\n") // Because splunk
	}

	if !f.doGz {
		return buf.Bytes(), nil
	}

	var buf2 bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf2, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return buf2.Bytes(), nil
}

func (f *SplunkFormat) From(raw []byte) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *SplunkFormat) Rollup(rolls []rollup.Rollup) ([]byte, error) {
	return nil, nil // Not supported for now.
}

func (f *SplunkFormat) toSplunkMetric(in *kt.JCHF, ts int64) []SplunkMetric {
	out := in.Flatten()
	host := in.DeviceName
	if host == "" {
		host = strconv.Itoa(int(in.DeviceId))
	}
	for k, v := range out {
		switch tv := v.(type) {
		case string:
			if tv == "" {
				delete(out, k)
			}
		case int:
			if tv == 0 {
				delete(out, k)
			}
		case int32:
			if tv == 0 {
				delete(out, k)
			}
		case int64:
			if tv == 0 {
				delete(out, k)
			} else if k == "timestamp" {
				out["timestamp"] = time.Unix(tv, 0).Format(time.RFC3339)

			}
		}
	}
	return []SplunkMetric{
		SplunkMetric{
			Host:       host,
			SourceType: string(in.Provider),
			Timestamp:  ts,
			Event:      out,
		},
	}
}
