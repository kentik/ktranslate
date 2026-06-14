package stitch

import (
	"flag"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/stitch/ringbuffer"
)

var (
	enable bool
	bufLen int
)

func init() {
	flag.IntVar(&bufLen, "stitch.buffer.len", 10000, "How large of a buffer of flows to try and stitch together.")
	flag.BoolVar(&enable, "stitch.enable", false, "Turn on flow stitching.")
}

type Stitcher struct {
	logger.ContextL
	cache    *ringbuffer.RingBuffer[*kt.JCHF]
	registry go_metrics.Registry
	metrics  *StitchMetric
}

type StitchMetric struct {
	FlowsIn      go_metrics.Meter
	FlowsMatched go_metrics.Meter
}

func NewStitcher(log logger.Underlying, cfg *ktranslate.StitchConfig, registry go_metrics.Registry) (*Stitcher, error) {
	if !cfg.Enable {
		return nil, nil
	}

	s := &Stitcher{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "flowStitch"}, log),
		cache:    ringbuffer.New[*kt.JCHF](cfg.BufLen),
		registry: registry,
		metrics: &StitchMetric{
			FlowsIn:      go_metrics.GetOrRegisterMeter(fmt.Sprintf("stitch.in^force=true"), registry),
			FlowsMatched: go_metrics.GetOrRegisterMeter(fmt.Sprintf("stitch.matched^force=true"), registry),
		},
	}

	s.Infof("Starting flow unification system with a buffer length of %d", cfg.BufLen)
	return s, nil
}

/*
If there's a matching ingress / egress flow, record it here.
*/
func (s *Stitcher) Stitch(msg *kt.JCHF) bool {

	key := msg.GetKey()
	s.metrics.FlowsIn.Mark(1)
	if nm, ok := s.cache.GetAndDelete(key); ok {
		msg.CustomInt["pair_tcp_flags"] = int32(nm.TcpFlags)
		msg.CustomBigInt["pair_in_bytes"] = int64(nm.InBytes)
		msg.CustomBigInt["pair_in_pkts"] = int64(nm.InPkts)
		msg.CustomInt["pair_tcp_rx"] = int32(nm.TcpRetransmit)
		msg.CustomBigInt["pair_timestamp"] = int64(nm.Timestamp)
		s.metrics.FlowsMatched.Mark(1)
		return true
	}

	s.cache.Put(key, msg)
	return false
}

// NOOP for now.
func (s *Stitcher) Stop() {}

func (s *Stitcher) HttpInfo() map[string]float64 {
	if s == nil {
		return nil
	}

	return map[string]float64{
		"flows_in":      s.metrics.FlowsIn.Rate1(),
		"flows_matched": s.metrics.FlowsMatched.Rate1(),
	}
}
