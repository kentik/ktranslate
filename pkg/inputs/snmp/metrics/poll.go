package metrics

import (
	"context"
	"math/rand"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/ping"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	extension "github.com/kentik/ktranslate/pkg/inputs/snmp/x"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/tick"
)

const (
	STATUS_CHECK_TIME  = 60 * time.Second
	DefaultPingTimeout = 1000 * time.Millisecond
)

type Poller struct {
	log              logger.ContextL
	server           *gosnmp.GoSNMP
	interfaceMetrics *InterfaceMetrics
	deviceMetrics    *DeviceMetrics
	jchfChan         chan []*kt.JCHF
	metrics          *kt.SnmpDeviceMetric
	counterTimeSec   int
	jitterTimeSec    int
	dropIfOutside    bool
	pinger           *ping.Pinger
	extension        extension.Extension
	gconf            *kt.SnmpGlobalConfig
	pingSec          int
}

func NewPoller(server *gosnmp.GoSNMP, gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL, logchan chan string) *Poller {
	counterTimeSec := util.GetPollRate(gconf, conf, log)

	jitterTimeSec := 15 // This is how long to spead the polling load out across.
	if gconf.JitterTimeSec > 0 {
		jitterTimeSec = gconf.JitterTimeSec
	}

	// Default is not not drop.
	dropIfOutside := false
	if gconf != nil && gconf.PollTimeSec > 0 {
		dropIfOutside = gconf.DropIfOutside
	}

	// If there's a profile passed in, look at the mibs set for this.
	var deviceMetricMibs, interfaceMetricMibs map[string]*kt.Mib
	if profile != nil {
		minCounterTime := counterTimeSec
		deviceMetricMibs, interfaceMetricMibs, minCounterTime = profile.GetMetrics(gconf.MibsEnabled, counterTimeSec)
		if counterTimeSec != minCounterTime {
			counterTimeSec = minCounterTime
			log.Warnf("%d poll time adjusting to new one base on mibs", counterTimeSec)
		}
	}

	poller := Poller{
		jchfChan:         jchfChan,
		log:              log,
		metrics:          metrics,
		server:           server,
		interfaceMetrics: NewInterfaceMetrics(gconf, conf, metrics, interfaceMetricMibs, profile, counterTimeSec, log),
		deviceMetrics:    NewDeviceMetrics(gconf, conf, metrics, deviceMetricMibs, profile, log),
		counterTimeSec:   counterTimeSec,
		jitterTimeSec:    jitterTimeSec,
		dropIfOutside:    dropIfOutside,
		gconf:            gconf,
	}

	// If we are extending the metrics for this device in any way, set it up now.
	ext, err := extension.NewExtension(jchfChan, gconf, conf, metrics, log, logchan)
	if err != nil {
		log.Errorf("Cannot setup extension for %s -> %s: %v", err, conf.DeviceIP, conf.DeviceName)
	} else if ext != nil {
		poller.extension = ext
		log.Infof("Enabling extension %s for %s -> %s", ext.GetName(), conf.DeviceIP, conf.DeviceName)
	}

	return &poller
}

func NewPollerForPing(gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL) *Poller {
	counterTimeSec := util.GetPollRate(gconf, conf, log)

	jitterTimeSec := 15 // This is how long to spead the polling load out across.
	if gconf.JitterTimeSec > 0 {
		jitterTimeSec = gconf.JitterTimeSec
	}

	pingSec := conf.PingSec
	if pingSec == 0 { // If not per device, try per global.
		pingSec = gconf.PingSec
	}
	if pingSec == 0 { // Default to 60 (1/per min) here if not defined in either global or per device levels.
		pingSec = 60
	}

	// How long to wait for a response back.
	timeout := time.Millisecond * time.Duration(conf.TimeoutMS)
	if timeout == 0 {
		timeout = DefaultPingTimeout
	}

	poller := Poller{
		jchfChan:       jchfChan,
		log:            log,
		metrics:        metrics,
		counterTimeSec: counterTimeSec,
		jitterTimeSec:  jitterTimeSec,
		deviceMetrics:  NewDeviceMetrics(gconf, conf, metrics, nil, profile, log),
		pingSec:        pingSec,
		gconf:          gconf,
	}

	p, err := ping.NewPinger(log, conf.DeviceIP, pingSec, timeout)
	if err != nil {
		log.Errorf("Cannot setup ping service for %s -> %s: %v", err, conf.DeviceIP, conf.DeviceName)
	} else {
		poller.pinger = p
		log.Infof("Enabling response time service for %s -> %s with timeout %v", conf.DeviceIP, conf.DeviceName, timeout)
	}

	return &poller
}

func NewPollerForExtention(gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL, logchan chan string) *Poller {
	counterTimeSec := util.GetPollRate(gconf, conf, log)

	jitterTimeSec := 15 // This is how long to spead the polling load out across.
	if gconf.JitterTimeSec > 0 {
		jitterTimeSec = gconf.JitterTimeSec
	}

	poller := Poller{
		jchfChan:       jchfChan,
		log:            log,
		metrics:        metrics,
		counterTimeSec: counterTimeSec,
		jitterTimeSec:  jitterTimeSec,
		deviceMetrics:  NewDeviceMetrics(gconf, conf, metrics, nil, profile, log),
		gconf:          gconf,
	}

	// If we are extending the metrics for this device in any way, set it up now.
	ext, err := extension.NewExtension(jchfChan, gconf, conf, metrics, log, logchan)
	if err != nil {
		log.Errorf("Cannot setup extension for %s -> %s: %v", err, conf.DeviceIP, conf.DeviceName)
	} else if ext != nil {
		poller.extension = ext
		log.Infof("Enabling extension %s for %s -> %s", ext.GetName(), conf.DeviceIP, conf.DeviceName)
	}

	return &poller
}

func (p *Poller) StartLoop(ctx context.Context) {

	// Problem is, SNMP counter polls take some time, and the time varies widely from device to device, based on number of interfaces and
	// round-trip-time to the device.  So we're going to divide each aligned five minute chunk into two periods: an initial period over which
	// to jitter the devices, and the rest of the five-minute chunk to actually do the counter-polling.  For any device whose counters we can walk
	// in less than (5 minutes - jitter period), we should be able to guarantee exactly one datapoint per aligned five-minute chunk.
	counterAlignment := time.Duration(p.counterTimeSec) * time.Second
	jitterWindow := time.Duration(p.jitterTimeSec) * time.Second
	firstCollection := time.Now().Truncate(counterAlignment).Add(counterAlignment).Add(time.Duration(rand.Int63n(int64(jitterWindow))))
	counterCheck := tick.NewFixedTimer(firstCollection, counterAlignment)
	statusCheck := time.NewTicker(STATUS_CHECK_TIME)

	p.log.Infof("snmpCounterPoll: First poll will be at %v. Polling every %v, drop=%v", firstCollection, counterAlignment, p.dropIfOutside)
	go func() {
		for {
			select {

			// Track the counters here, to convert from raw counters to differences
			case scheduledTime := <-counterCheck.C:

				startTime := time.Now()
				if !startTime.Truncate(counterAlignment).Equal(scheduledTime.Truncate(counterAlignment)) {
					// This poll was supposed to occur in a previous five-minute-block, but we were delayed
					// in picking it up -- presumably because a previous poll overflowed *its* block.
					// Since we can't possibly complete this one on schedule, skip it.
					p.log.Warnf("Skipping a counter datapoint for the period %v -- poll scheduled for %v, but only dequeued at %v",
						scheduledTime.Truncate(counterAlignment), scheduledTime, startTime)
					p.interfaceMetrics.DiscardDeltaState()
					p.metrics.Fail.Update(kt.SNMP_BAD_POLL_TIMEOUT)
					continue
				}

				pollCtx, pollCancel := context.WithTimeout(ctx, STATUS_CHECK_TIME)
				flows, err := p.Poll(pollCtx)
				pollCancel()
				if err != nil {
					p.log.Warnf("There was an error when polling the SNMP counter: %v.", err)

					// We didn't collect all the metrics here, which means that our delta values are
					// off, and we have to discard them.
					p.interfaceMetrics.DiscardDeltaState()
					continue
				}

				// Send counter data as flow
				if p.dropIfOutside && !time.Now().Truncate(counterAlignment).Equal(scheduledTime.Truncate(counterAlignment)) {
					// Uggh.  calling PollSNMPCounter took us long enough that we're no longer in the five-minute block
					// we were in when we started the poll.
					p.log.Warnf("Missed a counter datapoint for the period %v -- poll scheduled for %v, started at %v, ended at %v",
						scheduledTime.Truncate(counterAlignment), scheduledTime, startTime, time.Now())

					// Because this counter poll took too long, and at least the earliest values received in the
					// poll are already over five minutes old, we can no longer use them as the basis for deltas.
					// Throw all the values away, and start over with the next polling cycle
					p.interfaceMetrics.DiscardDeltaState()
					continue
				}

				// Great!  We finished the poll in the same block we started it in!
				p.log.Debugf("Metrics Poll: In %v polled %d mibs, %d missing.",
					time.Now().Sub(startTime),
					p.deviceMetrics.GetNumberMibsTotal()+p.interfaceMetrics.GetNumberMibsTotal(),
					p.deviceMetrics.GetNumberMibsFailed())
				p.jchfChan <- flows

			case <-statusCheck.C: // Send in on a seperate timer status about how this system is working.
				p.jchfChan <- p.deviceMetrics.GetStatusFlows()

			case <-ctx.Done():
				p.log.Infof("Metrics Poll Done")
				statusCheck.Stop()
				counterCheck.Stop()
				return
			}
		}
	}()

	// If there's any extensions, start them here.
	if p.extension != nil {
		go p.extension.Run(ctx, counterAlignment)
	}
}

// PollSNMPCounter polls SNMP for counter statistics like # bytes and packets transferred.
func (p *Poller) Poll(ctx context.Context) ([]*kt.JCHF, error) {

	deviceFlows, err := p.deviceMetrics.Poll(ctx, p.server, p.pinger)
	if err != nil {
		p.log.Warnf("Cannot poll device metrics: %v", err)
		p.metrics.Fail.Update(kt.SNMP_BAD_POLL_TIMEOUT)
	}

	flows, err := p.interfaceMetrics.Poll(ctx, p.server, deviceFlows)
	if err != nil {
		p.metrics.Fail.Update(kt.SNMP_BAD_POLL_TIMEOUT)
		return nil, err
	}

	// Marshal device metrics data into flow and append them to the list.
	flows = append(flows, deviceFlows...)

	// Since we didn't error and got some flows from this, set the value to GOOD.
	if len(flows) > 0 {
		p.metrics.Fail.Update(kt.SNMP_GOOD)
	} else {
		p.metrics.Fail.Update(kt.SNMP_BAD_POLL_TIMEOUT) // Otherwise, set to bad because there's no data coming out of this device.
	}

	return flows, nil
}

// Simpler loop which only runs on ping data, no actual snmp polling.
func (p *Poller) StartPingOnlyLoop(ctx context.Context) {
	if p.pinger == nil {
		p.log.Errorf("Missing pinger in ping loop.")
		return
	}

	p.pinger.Start(ctx)
	counterAlignment := time.Duration(p.counterTimeSec) * time.Second
	jitterWindow := time.Duration(p.jitterTimeSec) * time.Second
	firstCollection := time.Now().Truncate(counterAlignment).Add(counterAlignment).Add(time.Duration(rand.Int63n(int64(jitterWindow))))
	counterCheck := tick.NewFixedTimer(firstCollection, counterAlignment)
	fastDuration := time.Duration(kt.LookupEnvInt("KENTIK_FAST_PING_DURATION_SEC", 120)) * time.Second
	fastTick := time.Duration(kt.LookupEnvInt("KENTIK_FAST_PING_TICK_SEC", 10)) * time.Second
	slowTick := time.Duration(p.pingSec) * time.Second

	p.pinger.Reset(slowTick, p.counterTimeSec/p.pingSec)
	p.log.Infof("snmpPing: First run will be at %v. Running every %v, Sending %d probes", firstCollection, counterAlignment, p.counterTimeSec/p.pingSec)

	go func() {
		seenGoodPacketLoss := true
		runningFast := false

		// Run one ping first.
		go func() {
			flows, _, err := p.deviceMetrics.GetPingStats(ctx, p.pinger)
			if err != nil {
				p.log.Warnf("There was an error when getting ping stats: %v.", err)
			}

			// Send data on.
			p.jchfChan <- flows
		}()

		// And now wait for the normal ticker to go.
		for {
			select {
			case _ = <-counterCheck.C:
				go func() {
					if runningFast {
						return
					}

					flows, isTotalLoss, err := p.deviceMetrics.GetPingStats(ctx, p.pinger)
					if err != nil {
						p.log.Warnf("There was an error when getting ping stats: %v.", err)
						return
					}

					// Send data on.
					p.jchfChan <- flows

					if !isTotalLoss { // We don't want to go back into fast polling unless we get <100% packet loss at some point.
						seenGoodPacketLoss = true
					}

					// If there's total loss, go to fast polling but only if we haven't been here before.
					if p.gconf.FastPoll && isTotalLoss && seenGoodPacketLoss && !runningFast {
						p.log.Warnf("Starting fast ping operation due to 100 percent packet loss.")
						runningFast = true
						ctxT, cancel := context.WithTimeout(ctx, fastDuration)
						p.runFastPoll(ctxT, fastTick, fastDuration, slowTick)
						cancel() // Done with fast polling.
						seenGoodPacketLoss = false
						runningFast = false
					}
				}()

			case <-ctx.Done():
				p.log.Infof("Metrics PingOnly Done")
				counterCheck.Stop()
				p.pinger.Stop()
				return
			}
		}
	}()
}

func (p *Poller) runFastPoll(ctx context.Context, fastTick time.Duration, fastDuration time.Duration, slowTick time.Duration) {

	p.log.Infof("snmpFastPoll: Running every %v for %v", fastTick, fastDuration)
	fastCheck := time.NewTicker(fastTick)

	defer func() { // When we leave this loop, return to slow polling.
		fastCheck.Stop()
		p.pinger.Reset(slowTick, p.counterTimeSec/p.pingSec)
		p.log.Infof("Fast ping done.")
	}()

	// But for now we need fast polling.
	p.pinger.Reset(1*time.Second, int(fastTick.Seconds())) // Run every second.

	for {
		select {
		case _ = <-fastCheck.C:
			flows, isTotalLoss, err := p.deviceMetrics.GetPingStats(ctx, p.pinger)
			if err != nil {
				p.log.Warnf("There was an error when getting ping stats: %v.", err)
				continue
			}

			// Send data on.
			p.jchfChan <- flows

			if !isTotalLoss { // Total loss has resolved itself so back to slow polling.
				p.log.Warnf("snmpFastPoll: FastPoll Done: not total packet loss seen.")
				return
			}

		case <-ctx.Done():
			p.log.Warnf("snmpFastPoll: FastPoll Done: %v.", ctx.Err())
			return
		}
	}
}

// Simpler loop which only runs on ext data, no actual snmp polling.
func (p *Poller) StartExtensionOnlyLoop(ctx context.Context) {
	if p.extension == nil {
		p.log.Errorf("Missing extension in Ext loop.")
		return
	}

	// Problem is, SNMP counter polls take some time, and the time varies widely from device to device, based on number of interfaces and
	// round-trip-time to the device.  So we're going to divide each aligned five minute chunk into two periods: an initial period over which
	// to jitter the devices, and the rest of the five-minute chunk to actually do the counter-polling.  For any device whose counters we can walk
	// in less than (5 minutes - jitter period), we should be able to guarantee exactly one datapoint per aligned five-minute chunk.
	counterAlignment := time.Duration(p.counterTimeSec) * time.Second

	if p.extension != nil {
		p.log.Infof("Running only extension %s", p.extension.GetName())
		go p.extension.Run(ctx, counterAlignment)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				p.log.Infof("Metrics ExtOnly Done")
				return
			}
		}
	}()
}
