package prom

import (
	"fmt"
	"sort"
	"testing"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"
)

// TestRemotePromSynthOutcome verifies the prometheus remote-write format emits the
// enumerated synth outcome (and other synth metrics) for error/timeout/ok records.
// Run with:
//
//	go test -tags dynamic -run TestRemotePromSynthOutcome -v ./pkg/formats/prom/
func TestRemotePromSynthOutcome(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	cfg := ktranslate.DefaultConfig().PrometheusFormat
	f, err := NewRemoteFormat(l, kt.CompressionSnappy, cfg)
	assert.NoError(err)

	mk := func(rt int32, rtStr, testType string) *kt.JCHF {
		j := kt.NewJCHF()
		j.EventType = kt.KENTIK_EVENT_SYNTH
		j.Timestamp = 1
		j.DeviceName = "synthetic-agent"
		j.CustomInt = map[string]int32{"result_type": rt}
		j.CustomBigInt = map[string]int64{}
		j.CustomStr = map[string]string{
			"result_type_str": rtStr,
			"test_type":       testType,
			"test_name":       "my-test",
		}
		return j
	}

	batch := []*kt.JCHF{
		mk(0, "error", "http"),
		mk(1, "timeout", "http"),
		mk(2, "ping", "ping"),
	}

	out, err := f.To(batch, make([]byte, 0))
	assert.NoError(err)
	assert.NotNil(out)

	// The remote-write payload is snappy-compressed protobuf; decode it.
	raw, err := snappy.Decode(nil, out.Body)
	assert.NoError(err)
	var wr prompb.WriteRequest
	assert.NoError(proto.Unmarshal(raw, &wr))

	// Collect the outcome value per result_type_str label.
	outcomes := map[string]float64{}
	names := map[string]bool{}
	for _, ts := range wr.Timeseries {
		var name, rtStr string
		for _, lb := range ts.Labels {
			if lb.Name == "name" {
				name = lb.Value
			}
			if lb.Name == "result_type_str" {
				rtStr = lb.Value
			}
		}
		names[name] = true
		if name == "kentik.synth.outcome" && len(ts.Samples) > 0 {
			outcomes[rtStr] = ts.Samples[0].Value
		}
	}

	printed := make([]string, 0, len(names))
	for n := range names {
		printed = append(printed, n)
	}
	sort.Strings(printed)
	fmt.Println("---- remote-write series names ----")
	for _, n := range printed {
		fmt.Println(n)
	}
	fmt.Printf("---- outcome by result_type_str ----\n%v\n", outcomes)

	assert.Equal(float64(util.SynthOutcomeError), outcomes["error"], "error -> 2")
	assert.Equal(float64(util.SynthOutcomeTimeout), outcomes["timeout"], "timeout -> 1")
	assert.Equal(float64(util.SynthOutcomeOK), outcomes["ping"], "ok -> 0")
}

// TestRemotePromSyngest verifies the remote-write format also emits KSynthgest (mesh) metrics.
//
//	go test -tags dynamic -run TestRemotePromSyngest -v ./pkg/formats/prom/
func TestRemotePromSyngest(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	cfg := ktranslate.DefaultConfig().PrometheusFormat
	f, err := NewRemoteFormat(l, kt.CompressionSnappy, cfg)
	assert.NoError(err)

	j := kt.NewJCHF()
	j.EventType = kt.KENTIK_EVENT_SYNTH_GEST
	j.Timestamp = 1
	j.DeviceName = "synthetic-agent"
	j.CustomInt = map[string]int32{"avg_latency": 27366, "avg_jitter": 934}
	j.CustomBigInt = map[string]int64{}
	j.CustomStr = map[string]string{"test_type": "agent", "test_name": "mesh"}

	out, err := f.To([]*kt.JCHF{j}, make([]byte, 0))
	assert.NoError(err)
	assert.NotNil(out)

	raw, err := snappy.Decode(nil, out.Body)
	assert.NoError(err)
	var wr prompb.WriteRequest
	assert.NoError(proto.Unmarshal(raw, &wr))

	got := map[string]float64{}
	for _, ts := range wr.Timeseries {
		for _, lb := range ts.Labels {
			if lb.Name == "name" && len(ts.Samples) > 0 {
				got[lb.Value] = ts.Samples[0].Value
			}
		}
	}
	fmt.Printf("---- syngest series ----\n%v\n", got)
	assert.Equal(float64(27366), got["kentik.syngest.avg_latency"])
	assert.Equal(float64(934), got["kentik.syngest.avg_jitter"])
}
