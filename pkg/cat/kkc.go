package cat

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/cat/auth"
	"github.com/kentik/ktranslate/pkg/eggs/baseserver"
	"github.com/kentik/ktranslate/pkg/eggs/kmux"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/filter"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/inputs/flow"
	ihttp "github.com/kentik/ktranslate/pkg/inputs/http"
	"github.com/kentik/ktranslate/pkg/inputs/snmp"
	"github.com/kentik/ktranslate/pkg/inputs/syslog"
	"github.com/kentik/ktranslate/pkg/inputs/vpc"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/maps"
	"github.com/kentik/ktranslate/pkg/rollup"
	ss "github.com/kentik/ktranslate/pkg/sinks"
	"github.com/kentik/ktranslate/pkg/sinks/kentik"
	"github.com/kentik/ktranslate/pkg/util/enrich"
	"github.com/kentik/ktranslate/pkg/util/gopatricia/patricia"
	"github.com/kentik/ktranslate/pkg/util/resolv"
	"github.com/kentik/ktranslate/pkg/util/rule"

	"github.com/judwhite/go-svc"
	go_metrics "github.com/kentik/go-metrics"
)

// Setting this decode limit explicitly so that we know what it is.
// By default at the time of this writing the library would have used 64MiB.
// Feel free to change if appropriate.
const (
	CHAN_SLACK              = 8000 // Up to this many messages / sec
	MetricsCheckDuration    = 60 * time.Second
	CacheInvalidateDuration = 8 * time.Hour
	MDB_NO_LOCK             = 0x400000
	MDB_PERMS               = 0666
)

var (
	RollupsSendDuration = 15 * time.Second
)

func NewKTranslate(config *ktranslate.Config, log logger.ContextL, registry go_metrics.Registry, version string, sinks string, serviceName string, logTee chan string, metricsChan chan []*kt.JCHF) (*KTranslate, error) {
	kc := &KTranslate{
		log:      log,
		registry: registry,
		config:   config,
		metrics: &KKCMetric{
			Flows:        go_metrics.GetOrRegisterMeter("flows", registry),
			FlowsOut:     go_metrics.GetOrRegisterMeter("flows_out", registry),
			DroppedFlows: go_metrics.GetOrRegisterMeter("dropped_flows", registry),
			Errors:       go_metrics.GetOrRegisterMeter("errors", registry),
			AlphaQ:       go_metrics.GetOrRegisterGauge("alphaq", registry),
			AlphaQDrop:   go_metrics.GetOrRegisterMeter("alphaq_drop", registry),
			JCHFQ:        go_metrics.GetOrRegisterGauge("jchfq", registry),
			InputQ:       go_metrics.GetOrRegisterMeter("inputq", registry),
			InputQLen:    go_metrics.GetOrRegisterGauge("inputq_len^force=true", registry),
			OutputQLen:   go_metrics.GetOrRegisterGauge("outputq_len^force=true", registry),
		},
		alphaChans:  make([]chan *Flow, config.ProcessingThreads),
		jchfChans:   make([]chan *kt.JCHF, config.ProcessingThreads),
		metricsChan: metricsChan,
		logTee:      logTee,
		msgsc:       make(chan *kt.Output, 60),
		tooBig:      make(chan int, CHAN_SLACK),
	}

	if v := config.API.DeviceFile; v != "" {
		kc.authConfig = &auth.AuthConfig{
			DevicesFile: v,
		}
	}

	if config.KentikEmail != "" && config.KentikAPIToken != "" {
		kentikConfig := &kt.KentikConfig{
			ApiEmail: config.KentikEmail,
			ApiToken: config.KentikAPIToken,
			ApiRoot:  config.APIBaseURL,
			ApiPlan:  config.KentikPlan,
		}

		kc.kentikConfig = kentikConfig
	}

	for i := 0; i < config.ProcessingThreads; i++ {
		kc.jchfChans[i] = make(chan *kt.JCHF, CHAN_SLACK)
		for j := 0; j < CHAN_SLACK; j++ {
			kc.jchfChans[i] <- kt.NewJCHF()
		}
	}

	log.Infof("Turning on %d processing threads", config.ProcessingThreads)
	for i := 0; i < config.ProcessingThreads; i++ {
		kc.alphaChans[i] = make(chan *Flow, CHAN_SLACK)
	}

	// Load any rollups we are doing
	rolls, err := rollup.GetRollups(log.GetLogger().GetUnderlyingLogger(), config.Rollup)
	if err != nil {
		return nil, err
	}
	kc.rollups = rolls
	kc.doRollups = len(rolls) > 0

	// And load any filters we are doing
	filters, err := filter.GetFilters(log.GetLogger().GetUnderlyingLogger(), config.Filters)
	if err != nil {
		return nil, err
	}
	kc.filters = filters
	kc.doFilter = len(filters) > 0

	// Grab the custom data directly from a file.
	if config.MappingFile != "" {
		m, err := NewCustomMapper(config.MappingFile)
		if err != nil {
			return nil, err
		}
		kc.mapr = m
		kc.log.Infof("Loaded %d custom mappings", len(m.Customs))
	} else { // Make this empty to we don't error out.
		kc.mapr = &CustomMapper{Customs: map[uint32]string{}}
	}

	if config.UDRSFile != "" {
		m, udrs, err := NewUDRMapper(config.UDRSFile)
		if err != nil {
			return nil, err
		}
		kc.udrMapr = m
		kc.log.Infof("Loaded %d udr and %d subtype mappings with %d udrs total", len(m.UDRs), len(m.Subtypes), udrs)
	}

	m, err := maps.LoadMapper(maps.Mapper(config.TagMapType), log.GetLogger().GetUnderlyingLogger(), config.TagMapFile)
	if err != nil {
		kc.log.Errorf("There was an error when opening the tag service: %v.", err)
		return nil, err
	}
	kc.tagMap = m

	// Load up a geo file if one is passed in.
	if config.GeoFile != "" {
		geo, err := patricia.NewMapFromMM(config.GeoFile, log)
		if err != nil {
			kc.log.Errorf("There was an error with geo service: %v.", err)
			return nil, err
		} else {
			kc.geo = geo
		}
	}

	// Load asn mapper if set.
	if config.ASNFile != "" {
		asn, err := patricia.NewMapFromMM(config.ASNFile, log)
		if err != nil {
			kc.log.Errorf("There was an error with the asn service: &v.", err)
			return nil, err
		} else {
			kc.asn = asn
		}
	}

	// Define our sinks for where to send data to.
	kc.sinks = make(map[ss.Sink]ss.SinkImpl)
	for _, sinkStr := range strings.Split(sinks, ",") {
		sink := ss.Sink(sinkStr)
		snk, err := ss.NewSink(sink, log.GetLogger().GetUnderlyingLogger(), registry, kc.tooBig, kc.kentikConfig, logTee, kc.config)
		if err != nil {
			return nil, fmt.Errorf("Invalid sink: %s, %v", sink, err)
		}
		kc.sinks[sink] = snk
		kc.log.Infof("Using sink %s", sink)

		// Kentik gets special cased
		if sink == ss.KentikSink {
			kc.kentik = snk.(*kentik.KentikSink)
		}
	}

	// IP based rules
	rule, err := rule.NewRuleSet(config.ApplicationFile, log)
	if err != nil {
		return nil, err
	}
	kc.rule = rule

	// External Enrichment.
	if config.EnricherURL != "" {
		en, err := enrich.NewEnricher(config.EnricherURL, log.GetLogger().GetUnderlyingLogger())
		if err != nil {
			return nil, err
		}
		kc.enricher = en
	}

	if len(kc.sinks) == 0 {
		return nil, fmt.Errorf("No sinks set")
	}

	// Set snmp know what the service name is:
	snmp.ServiceName = serviceName

	// Get some randomness
	rand.Seed(time.Now().UnixNano())

	return kc, nil
}

// nolint: errcheck
func (kc *KTranslate) cleanup() {
	snmp.Close()
	for _, sink := range kc.sinks {
		sink.Close()
	}
	if kc.pgdb != nil {
		kc.pgdb.Close()
	}
	if kc.geo != nil {
		kc.geo.Close()
	}
	if kc.asn != nil {
		kc.asn.Close()
	}
	if kc.vpc != nil {
		kc.vpc.Close()
	}
	if kc.nfs != nil {
		kc.nfs.Close()
	}
	if kc.syslog != nil {
		kc.syslog.Close()
	}
}

// GetStatus implements the baseserver.Service interface.
func (kc *KTranslate) GetStatus() []byte {
	return []byte("OK")
}

// RunHealthCheck implements the baseserver.Service interface.
func (kc *KTranslate) RunHealthCheck(ctx context.Context, result *baseserver.HealthCheckResult) {
}

// HttpInfo implements the baseserver.Service interface.
func (kc *KTranslate) HttpInfo(w http.ResponseWriter, r *http.Request) {
	total := 0
	for _, c := range kc.alphaChans {
		total += len(c)
	}
	kc.metrics.AlphaQ.Update(int64(total)) // Update these on demand.

	total = 0
	for _, c := range kc.jchfChans {
		total += len(c)
	}
	kc.metrics.JCHFQ.Update(int64(total))
	h := hc{
		Flows:          kc.metrics.Flows.Rate1(),
		FlowsOut:       kc.metrics.FlowsOut.Rate1(),
		DroppedFlows:   kc.metrics.DroppedFlows.Rate1(),
		Errors:         kc.metrics.Errors.Rate1(),
		AlphaQ:         kc.metrics.AlphaQ.Value(),
		JCHFQ:          kc.metrics.JCHFQ.Value(),
		AlphaQDrop:     kc.metrics.AlphaQDrop.Rate1(),
		InputQ:         kc.metrics.InputQ.Rate1(),
		InputQLen:      kc.metrics.InputQLen.Value(),
		OutputQLen:     kc.metrics.OutputQLen.Value(),
		Sinks:          map[ss.Sink]map[string]float64{},
		SnmpDeviceData: map[string]map[string]float64{},
		Inputs:         map[string]map[string]float64{},
	}

	// Now, let other sinks do their work
	for sn, sink := range kc.sinks {
		h.Sinks[sn] = sink.HttpInfo()
	}

	// And store any metrics from inputs.
	if kc.metrics.SnmpDeviceData != nil {
		kc.metrics.SnmpDeviceData.Mux.RLock()
		defer kc.metrics.SnmpDeviceData.Mux.RUnlock()
		for d, met := range kc.metrics.SnmpDeviceData.Devices {
			h.SnmpDeviceData[d] = map[string]float64{
				"DeviceMetrics":    met.DeviceMetrics.Rate1(),
				"InterfaceMetrics": met.InterfaceMetrics.Rate1(),
				"Metadata":         met.Metadata.Rate1(),
				"Errors":           met.Errors.Rate1(),
			}
		}
	}
	if kc.vpc != nil {
		h.Inputs["vpc"] = kc.vpc.HttpInfo()
	}
	if kc.nfs != nil {
		h.Inputs["flow"] = kc.nfs.HttpInfo()
	}
	if kc.syslog != nil {
		h.Inputs["syslog"] = kc.syslog.HttpInfo()
	}

	b, err := json.Marshal(h)
	if err != nil {
		kc.log.Errorf("Error in HC: %v", err)
	} else {
		w.Write(b)
	}
}

func (kc *KTranslate) doSend(ctx context.Context) {
	kc.log.Infof("do sendToKTranslate Starting")

	for {
		select {
		case ser := <-kc.msgsc:
			if ser.BodyLen() == 0 {
				continue
			}

			for _, sink := range kc.sinks {
				sink.Send(ctx, ser)
			}

		case <-ctx.Done():
			kc.log.Infof("do sendToKTranslate Done")
			return
		}
	}
}

func (kc *KTranslate) sendToSinks(ctx context.Context) error {

	metricsTicker := time.NewTicker(MetricsCheckDuration)
	defer metricsTicker.Stop()

	rollupsTicker := time.NewTicker(RollupsSendDuration)
	defer rollupsTicker.Stop()

	// This one is in charge of sending on to sinks.
	go kc.doSend(ctx)
	kc.log.Infof("sendToSinks base Online")

	// These do the actual processing now for data from kentik.
	for i := 0; i < kc.config.ProcessingThreads; i++ {
		go kc.monitorAlphaChan(ctx, i, kc.format.To)
	}

	for {
		select {
		case <-metricsTicker.C:
			total := 0
			for _, c := range kc.alphaChans {
				total += len(c)
			}
			kc.metrics.AlphaQ.Update(int64(total))

			total = 0
			for _, c := range kc.jchfChans {
				total += len(c)
			}
			kc.metrics.JCHFQ.Update(int64(total))

		case <-rollupsTicker.C:
			for _, r := range kc.rollups {
				export := r.Export()
				if len(export) > 0 {
					res, err := kc.formatRollup.Rollup(export)
					if err != nil {
						kc.log.Errorf("There was an error when handling rollup: %v.", err)
					} else {
						kc.msgsc <- res
					}
				}
			}

		case <-kc.tooBig:
			// We need to dynamically shrink the size of data being sent in based on feedback from one of our sinks.
			os := kc.config.MaxFlowsPerMessage
			kc.config.MaxFlowsPerMessage = int(math.Max((float64(kc.config.MaxFlowsPerMessage) * .75), 1))
			kc.log.Infof("Updating MaxFlowsPerMessage to %d from %d based on errors sending", kc.config.MaxFlowsPerMessage, os)

		case <-ctx.Done():
			kc.log.Infof("sendToSinks base Done")
			return nil
		}
	}
}

// This processes data from the non-kentik input sets.
func (kc *KTranslate) handleInput(ctx context.Context, msgs []*kt.JCHF, serBuf []byte, citycache map[uint32]string, regioncache map[uint32]string, cb func(error), seri func([]*kt.JCHF, []byte) (*kt.Output, error)) {
	if kc.geo != nil || kc.asn != nil {
		kc.doEnrichments(ctx, citycache, regioncache, msgs)
	}

	// If we are filtering, cut any out here.
	if kc.doFilter {
		msgs = kc.reduce(msgs)
	}

	// If we have any rollups defined, send here instead of directly to the output format.
	if kc.doRollups {
		rv := make([]map[string]interface{}, len(msgs))
		for i, msg := range msgs {
			rv[i] = msg.ToMap()
		}
		for _, r := range kc.rollups {
			r.Add(rv)
		}
	}

	// Turn into a binary format here, using the passed in encoder.
	if !kc.doRollups || kc.config.RollupAndAlpha {
		// Compute and sample rate stuff here.
		keep := len(msgs)
		if kc.config.SampleRate > 1 && keep > kc.config.SampleMin {
			rand.Shuffle(len(msgs), func(i, j int) {
				msgs[i], msgs[j] = msgs[j], msgs[i]
			})
			keep = int(math.Max(float64(len(msgs))/float64(kc.config.SampleRate), 1))
			for _, msg := range msgs {
				msg.SampleRate = msg.SampleRate * uint32(kc.config.SampleRate)
			}
			kc.log.Debugf("Reduced input from %d to %d", len(msgs), keep)
		}

		// Ship all the logs out, according to max flows per message.
		last := 0
		for next := kc.config.MaxFlowsPerMessage; next < keep+kc.config.MaxFlowsPerMessage; next += kc.config.MaxFlowsPerMessage {
			batch := next
			if batch > keep {
				batch = keep
			}
			ser, err := seri(msgs[last:batch], serBuf)
			if err != nil {
				kc.log.Errorf("There was an error when converting to native: %v.", err)
			} else if ser != nil {
				ser.CB = cb
				kc.msgsc <- ser
			}
			last = next

			if batch == keep { // We're done here, no need to send more.
				break
			}
		}
	}

	kc.metrics.InputQ.Mark(int64(len(msgs)))
}

func (kc *KTranslate) watchInput(ctx context.Context, seri func([]*kt.JCHF, []byte) (*kt.Output, error)) {
	kc.log.Infof("watchInput running")
	checkTicker := time.NewTicker(60 * time.Second)
	defer checkTicker.Stop()

	for {
		select {
		case _ = <-checkTicker.C:
			if kc.config.InputThreads < kc.config.MaxThreads {
				if len(kc.inputChan) > CHAN_SLACK-10 { // We're filling up our channel here. Try launching another thread.
					kc.log.Infof("watchInput launching another input channel. input at %d", len(kc.inputChan))
					go kc.monitorInput(ctx, kc.config.InputThreads, seri)
					kc.config.InputThreads++
				}
			}
			kc.metrics.InputQLen.Update(int64(len(kc.inputChan)))
			kc.metrics.OutputQLen.Update(int64(len(kc.msgsc)))
		case <-ctx.Done():
			kc.log.Infof("watchInput Done")
			return
		}
	}
}

func (kc *KTranslate) monitorInput(ctx context.Context, num int, seri func([]*kt.JCHF, []byte) (*kt.Output, error)) {
	kc.log.Infof("monitorInput %d Starting", num)
	serBuf := make([]byte, 0)
	citycache := map[uint32]string{}
	regioncache := map[uint32]string{}

	for {
		select {
		case msgs := <-kc.inputChan:
			kc.handleInput(ctx, msgs, serBuf, citycache, regioncache, nil, seri)
		case <-ctx.Done():
			kc.log.Infof("monitorInput %d Done", num)
			return
		}
	}
}

func (kc *KTranslate) monitorMetricsInput(ctx context.Context, seri func([]*kt.JCHF, []byte) (*kt.Output, error)) {
	kc.log.Infof("monitorMetricsInput Starting")
	serBuf := make([]byte, 0)
	citycache := map[uint32]string{}
	regioncache := map[uint32]string{}

	for {
		select {
		case msgs := <-kc.metricsChan:
			kc.handleInput(ctx, msgs, serBuf, citycache, regioncache, nil, seri)
		case <-ctx.Done():
			kc.log.Infof("monitorMetricsInput Done")
			return
		}
	}
}

// Removes any flows which don't pass the filters.
// This is On*f -- is there a better way?
func (kc *KTranslate) reduce(in []*kt.JCHF) []*kt.JCHF {
	out := make([]*kt.JCHF, 0, len(in))
	for _, msg := range in {
		keep := true
		for _, f := range kc.filters {
			if !f.Filter(msg) {
				keep = false
				break
			}
		}
		if keep {
			out = append(out, msg)
		}
	}

	return out
}

func (kc *KTranslate) getRouter() http.Handler {
	r := kmux.NewRouter()
	r.HandleFunc(HttpAlertInboundPath, kc.handleFlow)
	r.HandleFunc(HttpHealthCheckPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n") // nolint: errcheck
	})
	if kc.auth != nil {
		kc.auth.RegisterRoutes(r)
	}
	if kc.http != nil {
		kc.http.RegisterRoutes(r)
	}

	return r
}

func (kc *KTranslate) listenHTTP() {
	if kc.config.ListenAddr == "off" {
		kc.log.Infof("Turning off HTTP server.")
		return
	}

	server := &http.Server{Addr: kc.config.ListenAddr, Handler: kc.getRouter()}
	var err error
	if kc.config.SSLCertFile != "" {
		kc.log.Infof("Setting up HTTPS system on %s%s", kc.config.ListenAddr, HttpAlertInboundPath)
		err = server.ListenAndServeTLS(kc.config.SSLCertFile, kc.config.SSLKeyFile)
	} else {
		kc.log.Infof("Setting up HTTP system on %s%s", kc.config.ListenAddr, HttpAlertInboundPath)
		err = server.ListenAndServe()
	}

	// err is always non-nil -- the http server stopped.
	if err != http.ErrServerClosed {
		kc.log.Errorf("There was an error when bringing up the HTTP system on %s: %v.", kc.config.ListenAddr, err)
		panic(err)
	}
	kc.log.Infof("HTTP server shut down on %s -- %v", kc.config.ListenAddr, err)
}

func (kc *KTranslate) Run(ctx context.Context) error {
	defer kc.cleanup()

	format := formats.Format(kc.config.Format)
	formatRollup := formats.Format(kc.config.FormatRollup)
	compression := kt.Compression(kc.config.Compression)

	// DNS mapper if set.
	if kc.config.DNS != "" {
		res, err := resolv.NewResolver(ctx, kc.log.GetLogger().GetUnderlyingLogger(), kc.config.DNS)
		if err != nil {
			return err
		}
		kc.resolver = res
		kc.log.Infof("Enabled DNS resolution at: %s", kc.config.DNS)
	}

	// Set up formatter
	fmtr, err := formats.NewFormat(format, kc.log.GetLogger().GetUnderlyingLogger(), kc.registry, compression, kc.config)
	if err != nil {
		return err
	}
	kc.format = fmtr

	if kc.config.FormatRollup != "" { // Rollups default to using the same format as main, but can be seperated out.
		fmtr, err := formats.NewFormat(formatRollup, kc.log.GetLogger().GetUnderlyingLogger(), kc.registry, compression, kc.config)
		if err != nil {
			return err
		}
		kc.formatRollup = fmtr
	} else {
		kc.formatRollup = fmtr
	}

	// Connect our sinks.
	for _, sink := range kc.sinks {
		err := sink.Init(ctx, format, compression, kc.format)
		if err != nil {
			return err
		}
	}

	// Set up api auth system if this is set. Allows kproxy|kprobe|kappa|ksynth and others to use this without phoneing home to kentik.
	if kc.authConfig != nil {
		authr, err := auth.NewServer(kc.authConfig, kc.config.SNMPInput.SNMPFile, kc.log)
		if err != nil {
			return err
		}
		kc.auth = authr
	}

	// Api system for talking to kentik.
	if kc.kentikConfig != nil && kc.kentikConfig.ApiEmail != "" {
		apic, err := api.NewKentikApi(ctx, kc.kentikConfig, kc.log, kc.config.API)
		if err != nil {
			return err
		}
		kc.apic = apic
	} else {
		kc.apic = api.NewKentikApiFromLocalDevices(kc.auth.GetDeviceMap(), kc.log)
	}

	assureInput := func() { // Start up input processing if any is asked of us.
		if kc.inputChan == nil {
			kc.inputChan = make(chan []*kt.JCHF, CHAN_SLACK)
			for i := 0; i < kc.config.InputThreads; i++ {
				go kc.monitorInput(ctx, i, kc.format.To)
			}
			if kc.config.InputThreads < kc.config.MaxThreads {
				go kc.watchInput(ctx, kc.format.To)
			}
		}
	}

	// If SNMP is configured, start this system too. Poll for metrics and metadata, also handle traps.
	if kc.config.SNMPInput.Enable {
		if kc.config.EnableSNMPDiscovery { // Here, we're just returning the list of devices on the network which might speak snmp.
			_, err := snmp.Discover(ctx, kc.log, 0, kc.config.SNMPInput)
			return err
		}
		assureInput()
		kc.metrics.SnmpDeviceData = kt.NewSnmpMetricSet(kc.registry)
		err := snmp.StartSNMPPolls(ctx, kc.inputChan, kc.metrics.SnmpDeviceData, kc.registry, kc.apic, kc.log, kc.config.SNMPInput)
		if err != nil {
			return err
		}
	}

	// If we're looking for vpc flows coming in
	if kc.config.GCPVPCInput.Enable || kc.config.AWSVPCInput.Enable {
		assureInput()
		serBufInput := make([]byte, 0)
		citycacheInput := map[uint32]string{}
		regioncacheInput := map[uint32]string{}
		handler := func(msgs []*kt.JCHF, cb func(error)) { // Capture this in a closure.
			kc.handleInput(ctx, msgs, serBufInput, citycacheInput, regioncacheInput, cb, kc.format.To)
		}
		var vpcSource vpc.CloudSource
		if kc.config.GCPVPCInput.Enable && kc.config.AWSVPCInput.Enable {
			return fmt.Errorf("cannot enable both GCP and VPC input sources")
		}
		if kc.config.GCPVPCInput.Enable {
			vpcSource = vpc.Gcp
		}
		if kc.config.AWSVPCInput.Enable {
			vpcSource = vpc.Aws
		}
		vpci, err := vpc.NewVpc(ctx, vpcSource, kc.log.GetLogger().GetUnderlyingLogger(), kc.registry, kc.inputChan, kc.apic, kc.config.MaxFlowsPerMessage, handler, kc.config)
		if err != nil {
			return err
		}
		kc.vpc = vpci
	}

	// If we're looking for netflow direct flows coming in
	if kc.config.FlowInput.Enable {
		assureInput()
		nfs, err := flow.NewFlowSource(ctx, flow.FlowSource(kc.config.FlowInput.Protocol), kc.config.MaxFlowsPerMessage, kc.log.GetLogger().GetUnderlyingLogger(), kc.registry, kc.inputChan, kc.apic, kc.resolver, kc.config.FlowInput)
		if err != nil {
			return err
		}
		kc.nfs = nfs
	}

	// If we're looking for syslog flows coming in
	if kc.config.SyslogInput.Enable {
		assureInput()
		ss, err := syslog.NewSyslogSource(ctx, kc.log.GetLogger().GetUnderlyingLogger(), kc.logTee, kc.registry, kc.apic, kc.resolver, kc.config.SyslogInput)
		if err != nil {
			return err
		}
		kc.syslog = ss
	}

	// If we're looking for json over http
	if kc.config.EnableHTTPInput {
		assureInput()
		sh, err := ihttp.NewHttpListener(ctx, kc.config.ListenAddr, kc.log.GetLogger().GetUnderlyingLogger(), kc.registry, kc.inputChan, kc.apic)
		if err != nil {
			return err
		}
		kc.http = sh
	}

	// If we're sending self metrics via a chan to sinks. This one always get sent via nrm.
	if kc.metricsChan != nil {
		// Set up formatter
		format := formats.Format(formats.FORMAT_NRM)
		if kc.config.FormatMetric != "" {
			format = formats.Format(kc.config.FormatMetric)
		}

		fmtr, err := formats.NewFormat(format, kc.log.GetLogger().GetUnderlyingLogger(), kc.registry, compression, kc.config)
		if err != nil {
			return err
		}
		go kc.monitorMetricsInput(ctx, fmtr.To)
	}

	kc.log.Infof("System running with format %s, compression %s, max flows: %d, sample rate %d:1 after %d", kc.config.Format, kc.config.Compression, kc.config.MaxFlowsPerMessage, kc.config.SampleRate, kc.config.SampleMin)
	go kc.listenHTTP()
	return kc.sendToSinks(ctx)
}

// These are needed in case we are running under windows.
func (kc *KTranslate) Init(env svc.Environment) error {
	return nil
}

func (kc *KTranslate) Start() error {
	go kc.Run(context.Background())
	return nil
}

func (kc *KTranslate) Stop() error {
	kc.cleanup()
	return nil
}
