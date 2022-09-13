package nrmulti

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	jsoniter "github.com/json-iterator/go"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/sinks/nr"
)

var json = jsoniter.ConfigFastest

type NRMultiSink struct {
	logger.ContextL
	sync.RWMutex
	sinks       map[kt.Cid]*nr.NRSink
	registry    go_metrics.Registry
	tooBig      chan int
	configMult  *ktranslate.NewRelicMultiSinkConfig
	config      *ktranslate.NewRelicSinkConfig
	confKentik  *kt.KentikConfig
	ctx         context.Context
	format      formats.Format
	compression kt.Compression
	fmtr        formats.Formatter
	creds       map[kt.Cid]NRInfo
}

type NRInfo struct {
	NRAccount      string `json:"nr_account_id"`
	NRApiToken     string `json:"nr_api_token"`
	KentikApiEmail string `json:"kentik_api_email"`
	KentikApiToken string `json:"kentik_api_token"`
}

var (
	nrMultiFile string
)

func init() {
	flag.StringVar(&nrMultiFile, "nr_multi_config_file", "", "Used to send multiple accounts to NR")
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, tooBig chan int, logTee chan string, confK *kt.KentikConfig, cfg *ktranslate.NewRelicSinkConfig, cfgMult *ktranslate.NewRelicMultiSinkConfig) (*NRMultiSink, error) {
	return &NRMultiSink{
		ContextL:   logger.NewContextLFromUnderlying(logger.SContext{S: "nrMultiSink"}, log),
		sinks:      map[kt.Cid]*nr.NRSink{},
		registry:   registry,
		tooBig:     tooBig,
		configMult: cfgMult,
		config:     cfg,
		confKentik: confK,
	}, nil
}

func (s *NRMultiSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	s.ctx = ctx
	s.format = format
	s.compression = compression
	s.fmtr = fmtr

	// Load the config map.
	m := map[kt.Cid]NRInfo{}
	by, err := os.ReadFile(s.configMult.CredFile)
	if err != nil {
		return fmt.Errorf("Cannot load nrMulti file from %s -- %v", s.configMult.CredFile, err)
	}
	err = json.Unmarshal(by, &m)
	if err != nil {
		return err
	}
	s.creds = m

	if s.confKentik == nil { // Set this up if nil.
		base := ktranslate.DefaultConfig()
		s.confKentik = &kt.KentikConfig{
			ApiRoot: base.APIBaseURL,
			ApiPlan: base.KentikPlan,
		}
	}
	kentikInfo := map[kt.Cid]*kt.KentikConfig{}
	for cid, i := range m {
		kentikInfo[cid] = &kt.KentikConfig{
			ApiEmail: i.KentikApiEmail,
			ApiToken: i.KentikApiToken,
		}
	}
	s.confKentik.CredMap = kentikInfo // Save this one for the api system also.
	s.Infof("Online with %d accounts", len(s.creds))

	return nil
}

func (s *NRMultiSink) Send(ctx context.Context, payload *kt.Output) {

	place := func(ss *nr.NRSink) { // Do the lock dance.
		s.RUnlock()
		s.Lock()
		s.sinks[payload.Ctx.CompanyId] = ss
		s.Unlock()
		s.RLock()
	}

	s.RLock()
	defer s.RUnlock()
	if _, ok := s.sinks[payload.Ctx.CompanyId]; !ok {
		sink, err := nr.NewSink(s.GetLogger().GetUnderlyingLogger(), s.registry, s.tooBig, nil, s.config)
		if err != nil {
			s.Errorf("Cannot create NR sink for %d %v", payload.Ctx.CompanyId, err)
			place(nil) // Nil means skip this cid.
			return
		}
		account, token, err := s.getCreds(payload.Ctx.CompanyId)
		if err != nil {
			s.Errorf("Cannot get NR creds for %d %v", payload.Ctx.CompanyId, err)
			place(nil) // Nil means skip this cid.
			return
		}

		sink.NRAccount = account // Update these values per cid.
		sink.NRApiKey = token
		err = sink.Init(s.ctx, s.format, s.compression, s.fmtr)
		if err != nil {
			s.Errorf("Cannot init NR sink for %d %v", payload.Ctx.CompanyId, err)
			place(nil) // Same, skip out.
			return
		}

		// Sink is good, go ahead and run.
		place(sink)
	}

	sink := s.sinks[payload.Ctx.CompanyId]
	if sink != nil {
		s.sinks[payload.Ctx.CompanyId].Send(ctx, payload)
	}
}

func (s *NRMultiSink) Close() {}

func (s *NRMultiSink) HttpInfo() map[string]float64 {
	s.RLock()
	defer s.RUnlock()

	res := map[string]float64{}
	for cid, sink := range s.sinks {
		si := sink.HttpInfo()
		for k, v := range si {
			res[fmt.Sprintf("%d.%s", cid, k)] = v
		}
	}

	return res
}

func (s *NRMultiSink) getCreds(cid kt.Cid) (string, string, error) {
	if info, ok := s.creds[cid]; ok {
		return info.NRAccount, info.NRApiToken, nil
	}
	return "", "", fmt.Errorf("Cannot find credencial info for cid %d", cid)
}
