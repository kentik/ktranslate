package redis

import (
	"context"
	"flag"
	"fmt"
	"net"
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
	invalids     map[string]bool
	metrics      *RedisMetrics
	rdb          *redis.Client
}

var (
	redisAddr     string
	redisPassword string
	redisDB       int
	keyPrefix     string
)

const (
	aPrefix    = ":A"
	aaaaPrefix = ":AAAA"
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
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "redis"}, log),
		lastMetadata: map[string]*kt.LastMetadata{},
		invalids:     map[string]bool{},
		ctx:          ctx,
		config:       cfg,
		metrics: &RedisMetrics{
			ExportDrops: go_metrics.GetOrRegisterCounter("redis_export_drops", registry),
		},
	}

	jf.rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := jf.rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cannot connect to Redis addr %s err %v", cfg.RedisAddr, err)
	}
	jf.Infof("connected to Redis: addr=%s", cfg.RedisAddr)

	return jf, nil
}

func (f *RedisFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	res := map[string][]RedisData{}
	for _, m := range msgs {
		for fqdn, r := range f.toRedisMetric(m) {
			if _, ok := res[fqdn]; !ok {
				res[fqdn] = r
			} else {
				res[fqdn] = append(res[fqdn], r...)
			}
		}
	}

	if len(res) == 0 {
		return nil, nil
	}

	pipe := f.rdb.Pipeline()

	for fqdn, results := range res {
		key := f.config.KeyPrefix + fqdn
		membersA := make([]redis.Z, 0, len(results))
		membersAAAA := make([]redis.Z, 0, len(results))
		for _, r := range results {
			if r.Is6 {
				membersAAAA = append(membersAAAA, redis.Z{
					Score:  r.Latency,
					Member: r.IP.String(),
				})
			} else {
				membersA = append(membersA, redis.Z{
					Score:  r.Latency,
					Member: r.IP.String(),
				})
			}
		}
		// ZADD key <score> <member> ... — overwrites existing scores atomically.
		if len(membersA) > 0 {
			pipe.ZAdd(f.ctx, key+aPrefix, membersA...)
		}
		if len(membersAAAA) > 0 {
			pipe.ZAdd(f.ctx, key+aaaaPrefix, membersAAAA...)
		}
		f.Debugf("queued ZADD key %s members %d", key, len(membersA)+len(membersAAAA))
	}

	_, err := pipe.Exec(f.ctx)
	return nil, err
}

func (f *RedisFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *RedisFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

func (f *RedisFormat) toRedisMetric(in *kt.JCHF) map[string][]RedisData {
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

	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := make([]RedisData, 0, len(metrics))
	var fqdn string
	var ip net.IP

	if tn, ok := attr["test_name"].(string); ok {
		fqdn = tn
	} else {
		return nil
	}
	if da, ok := attr["dst_addr"].(string); ok {
		ip = net.ParseIP(da)
		if ip == nil {
			return nil
		}
	} else {
		return nil
	}

	for m, name := range metrics {
		switch name.Name {
		case "avg_rtt":
			ms = append(ms, RedisData{
				IP:      ip,
				Latency: float64(in.CustomInt[m]),
				Is6:     (ip.To4() == nil),
			})
		}
	}

	return map[string][]RedisData{fqdn: ms}
}

func (f *RedisFormat) fromSnmpMetadata(in *kt.JCHF) map[string][]RedisData {
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
	IP      net.IP
	Latency float64
	Is6     bool
}
