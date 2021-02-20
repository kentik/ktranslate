package influx

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

var (
	Measurement = flag.String("measurement", "kflow", "Measurement to use for rollups.")
)

type InfluxFormat struct {
	logger.ContextL
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*InfluxFormat, error) {
	jf := &InfluxFormat{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "influxFormat"}, log),
	}

	return jf, nil
}

func (f *InfluxFormat) To(msgs []*kt.JCHF, serBuf []byte) ([]byte, error) {

	res := make([]string, len(msgs))
	for i, m := range msgs {
		res[i] = fmt.Sprintf("kflow.raw,src_ip=%s,dst_ip=%s,protocol=%d,src_port=%d,dst_port=%d in_bytes=%d,in_pkts=%d,out_bytes=%d,out_pkts=%d %d",
			m.SrcAddr, m.DstAddr, m.Protocol, m.L4SrcPort, m.L4DstPort, m.InBytes, m.InPkts, m.OutBytes, m.OutPkts, m.Timestamp*1000000000) // Time to nano
	}

	return []byte(strings.Join(res, "\n")), nil
}

// Not supported.
func (f *InfluxFormat) From(raw []byte) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *InfluxFormat) Rollup(rolls []rollup.Rollup) ([]byte, error) {
	res := make([]string, len(rolls))
	ts := time.Now()
	for i, r := range rolls {
		pkts := strings.Split(r.EventType, ":")
		if len(pkts) > 2 {
			res[i] = fmt.Sprintf("%s,%s=%s %s=%d %d", *Measurement, strings.Join(pkts[2:], ":"), r.Dimension, pkts[1], uint64(r.Metric), ts.UnixNano()) // Time to nano
		}
	}

	return []byte(strings.Join(res, "\n")), nil
}
