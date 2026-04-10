package redis

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"sync"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/redis/go-redis/v9"
)

type RedisFormat struct {
	logger.ContextL
	lastMetadata map[string]*kt.LastMetadata
	mux          sync.RWMutex
	config       *ktranslate.RedisFormatConfig
	ctx          context.Context
	metrics      *RedisMetrics
}

var (
	redisAddr     string
	redisPassword string
	redisDB       int
	keyPrefix     string
)

func init() {
	flag.StringVar(&redisAddr, "redis.addr", "localhost:6379", "Where to connect to redis.")
	flag.StringVar(&redisPassword, "redis.password", "", "Password for redis")
	flag.IntVar(&redisDB, "redis.db", 0, "Use this redis DB.")
	flag.StringVar(&keyPrefix, "redis.key_prefix", "", "Use this key prefix.")
}

type RedisMetrics struct {
	ExportDrops go_metrics.Counter
}

func NewFormat(ctx context.Context, log logger.Underlying, cfg *ktranslate.RedisFormatConfig, registry go_metrics.Registry) (*RedisFormat, error) {
	jf := &RedisFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "otel"}, log),
		lastMetadata: map[string]*kt.LastMetadata{},
		invalids:     map[string]bool{},
		ctx:          ctx,
		config:       cfg,
		inputs:       map[string]chan OtelData{},
		metrics: &OtelMetrics{
			ExportDrops: go_metrics.GetOrRegisterCounter("otel_export_drops", registry),
		},
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cannot connect to Redis", "addr", cfg.RedisAddr, "err", err)
	}
	fj.Infof("connected to Redis", "addr", cfg.RedisAddr)

	return jf, nil
}

func (f *OtelFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	res := make([]OtelData, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, f.toOtelMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	f.mux.RLock()
	for _, m := range res {
		if _, ok := f.vecs[m.Name]; !ok {
			lm := m
			cv, err := otelm.Float64ObservableGauge(
				lm.Name,
				metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
					f.Debugf("Exporting data from otel for %s", lm.Name)
					for {
						select {
						case mm := <-f.inputs[lm.Name]:
							o.Observe(mm.Value, metric.WithAttributeSet(mm.GetTagValues()))
						default:
							return nil
						}
					}
				}),
			)
			if err != nil {
				f.mux.RUnlock()
				return nil, err
			}
			f.mux.RUnlock()
			f.mux.Lock()
			f.vecs[lm.Name] = cv
			f.inputs[lm.Name] = make(chan OtelData, CHAN_SLACK)
			f.mux.Unlock()
			f.mux.RLock()
			f.Infof("Creating otel export for %s", m.Name)
		}
		// Save this for later, for the next time the async callback is run.
		ch := f.inputs[m.Name]
		queueDepth := len(ch)
		if queueDepth >= CHAN_SLACK {
			f.Debugf("Channel queue at CHAN_SLACK limit for %s: %d/%d (100%%)", m.Name, queueDepth, CHAN_SLACK)
		}
		if f.config.NoBlockExport {
			select {
			case ch <- m:
			default:
				f.metrics.ExportDrops.Inc(1)
			}
		} else {
			ch <- m
		}
	}

	f.mux.RUnlock()
	return nil, nil
}

func (f *OtelFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *OtelFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

func (f *OtelFormat) toOtelMetric(in *kt.JCHF) []OtelData {
	switch in.EventType {
	case kt.KENTIK_EVENT_SYNTH:
		return f.fromKSynth(in)
	case kt.KENTIK_EVENT_SNMP_METADATA:
		return f.fromSnmpMetadata(in)
	default:
		f.mux.Lock()
		defer f.mux.Unlock()
		if !f.invalids[in.EventType] {
			f.Warnf("Invalid EventType: %s", in.EventType)
			f.invalids[in.EventType] = true
		}
	}

	return nil
}

func (f *RedisFormat) fromKSynth(in *kt.JCHF) map[string][]RedisData {
	if in.CustomInt["result_type"] <= 1 {
		return nil // Don't worry about timeouts and errors for now.
	}

	rawStr := in.CustomStr["error_cause/trace_route"] // Pull this out early.
	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := make([]OtelData, 0, len(metrics))

	for m, name := range metrics {
		switch name.Name {
		case "avg_rtt":
			ms = append(ms, RedisData{
				IP:      attr["dst_addr"],
				Latency: float64(in.CustomInt[m]),
			})
		}
	}

	return ms
}

func (f *OtelFormat) fromSnmpMetadata(in *kt.JCHF) []OtelData {
	if in.DeviceName == "" { // Only run if this is set.
		return nil
	}

	lm := util.SetMetadata(in)

	f.mux.Lock()
	defer f.mux.Unlock()
	if f.lastMetadata[in.DeviceName] == nil || lm.Size() >= f.lastMetadata[in.DeviceName].Size() {
		f.Infof("New Metadata for %s", in.DeviceName)
		f.lastMetadata[in.DeviceName] = lm
	} else {
		f.Infof("The metadata for %s was not updated since the attribute size is smaller. New = %d < Old = %d, Size difference = %v.",
			in.DeviceName, lm.Size(), f.lastMetadata[in.DeviceName].Size(), f.lastMetadata[in.DeviceName].Missing(lm))
	}

	return nil
}

type RedisData struct {
	IP      string
	Latency float64
}
