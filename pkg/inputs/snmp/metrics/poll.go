package metrics

import (
	"context"
	"math/rand"
	"time"

	"github.com/kentik/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/tick"
)

type Poller struct {
	log              logger.ContextL
	server           *gosnmp.GoSNMP
	interfaceMetrics *InterfaceMetrics
	deviceMetrics    *DeviceMetrics
	jchfChan         chan []*kt.JCHF
	metrics          *kt.SnmpDeviceMetric
	counterTimeSec   int
	dropIfOutside    bool
}

func NewPoller(server *gosnmp.GoSNMP, gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL) *Poller {
	// Default poll rate is 5 min. This is what a lot of SNMP billing is on.
	counterTimeSec := 5 * 60
	if conf != nil && conf.PollTimeSec > 0 {
		counterTimeSec = conf.PollTimeSec
	} else if gconf != nil && gconf.PollTimeSec > 0 {
		counterTimeSec = gconf.PollTimeSec
	}

	// Default is not not drop.
	dropIfOutside := false
	if gconf != nil && gconf.PollTimeSec > 0 {
		dropIfOutside = gconf.DropIfOutside
	}

	// If there's a profile passed in, look at the mibs set for this.
	var deviceMetricMibs, interfaceMetricMibs map[string]*kt.Mib
	if profile != nil {
		deviceMetricMibs, interfaceMetricMibs = profile.GetMetrics(gconf.MibsEnabled)
	}

	return &Poller{
		jchfChan:         jchfChan,
		log:              log,
		metrics:          metrics,
		server:           server,
		interfaceMetrics: NewInterfaceMetrics(gconf, conf, metrics, interfaceMetricMibs, log),
		deviceMetrics:    NewDeviceMetrics(gconf, conf, metrics, deviceMetricMibs, log),
		counterTimeSec:   counterTimeSec,
		dropIfOutside:    dropIfOutside,
	}
}

func (p *Poller) StartLoop(ctx context.Context) {

	// Problem is, SNMP counter polls take some time, and the time varies widely from device to device, based on number of interfaces and
	// round-trip-time to the device.  So we're going to divide each aligned five minute chunk into two periods: an initial period over which
	// to jitter the devices, and the rest of the five-minute chunk to actually do the counter-polling.  For any device whose counters we can walk
	// in less than (5 minutes - jitter period), we should be able to guarantee exactly one datapoint per aligned five-minute chunk.
	counterAlignment := time.Duration(p.counterTimeSec) * time.Second
	jitterWindow := 15 * time.Second
	firstCollection := time.Now().Truncate(counterAlignment).Add(counterAlignment).Add(time.Duration(rand.Int63n(int64(jitterWindow))))
	counterCheck := tick.NewFixedTimer(firstCollection, counterAlignment)

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
					continue
				}

				flows, err := p.Poll()
				if err != nil {
					p.log.Warnf("Issue polling SNMP Counter: %v", err)

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

				// Great!  We finished the poll in the same five-minute block we started it in!
				// send the results to Sinks.
				p.jchfChan <- flows

			case <-ctx.Done():
				p.log.Infof("Metrics Poll Done")
				return
			}
		}
	}()
}

// PollSNMPCounter polls SNMP for counter statistics like # bytes and packets transferred.
func (p *Poller) Poll() ([]*kt.JCHF, error) {

	deviceFlows, err := p.deviceMetrics.Poll(p.server)

	flows, err := p.interfaceMetrics.Poll(p.server, deviceFlows)
	if err != nil {
		return nil, err
	}

	// Marshal device metrics data into flow and append them to the list.
	flows = append(flows, deviceFlows...)

	return flows, nil
}
