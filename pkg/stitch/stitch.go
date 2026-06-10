package stitch

import (
	"flag"
	"fmt"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/jellydator/ttlcache/v3"
)

var (
	enable bool
	ttlSec int
)

func init() {
	flag.IntVar(&ttlSec, "stitch.ttl.sec", 1, "TTL for holding flows to try and stitch together.")
	flag.BoolVar(&enable, "stitch.enable", false, "Turn on flow stitching.")
}

type Stitcher struct {
	logger.ContextL
	cache    *ttlcache.Cache[string, *kt.JCHF]
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

	cache := ttlcache.New[string, *kt.JCHF](
		ttlcache.WithTTL[string, *kt.JCHF](time.Duration(cfg.TTLSec) * time.Second),
	)

	go cache.Start() // starts automatic expired item deletion

	s := &Stitcher{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "flowStitch"}, log),
		cache:    cache,
		registry: registry,
		metrics: &StitchMetric{
			FlowsIn:      go_metrics.GetOrRegisterMeter(fmt.Sprintf("stitch.in^force=true"), registry),
			FlowsMatched: go_metrics.GetOrRegisterMeter(fmt.Sprintf("stitch.matched^force=true"), registry),
		},
	}

	s.Infof("Starting flow unification system with a ttl of %v", time.Duration(cfg.TTLSec)*time.Second)
	return s, nil
}

/*
If there's a matching ingress / egress flow, record it here.
*/
func (s *Stitcher) Stitch(msg *kt.JCHF) bool {
	key := msg.GetKey()

	s.metrics.FlowsIn.Mark(1)
	item, retrieved := s.cache.GetOrSet(key, msg, ttlcache.WithTTL[string, *kt.JCHF](ttlcache.DefaultTTL))
	if retrieved {
		msg.Pair = item.Value()
		s.cache.Delete(item.Key())
		s.metrics.FlowsMatched.Mark(1)
		return true
	}

	return false
}

func (s *Stitcher) Stop() {
	if s == nil {
		return
	}

	s.cache.Stop()
}

func (s *Stitcher) HttpInfo() map[string]float64 {
	if s == nil {
		return nil
	}

	return map[string]float64{
		"flows_in":      s.metrics.FlowsIn.Rate1(),
		"flows_matched": s.metrics.FlowsMatched.Rate1(),
	}
}
