package metrics

import (
	"strconv"
	"strings"
	"sync"

	"github.com/kentik/gosnmp"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/kt/counters"
	snmp_util "github.com/kentik/ktranslate/pkg/snmp/util"
)

// SNMP keys used in various places
const (
	SNMP_ifHCInOctets     = "ifHCInOctets"
	SNMP_ifHCInUcastPkts  = "ifHCInUcastPkts"
	SNMP_ifHCOutOctets    = "ifHCOutOctets"
	SNMP_ifHCOutUcastPkts = "ifHCOutUcastPkts"
	SNMP_ifInErrors       = "ifInErrors"
	SNMP_ifOutErrors      = "ifOutErrors"

	SNMP_ifInDiscards         = "ifInDiscards"
	SNMP_ifOutDiscards        = "ifOutDiscards"
	SNMP_ifHCOutMulticastPkts = "ifHCOutMulticastPkts"
	SNMP_ifHCOutBroadcastPkts = "ifHCOutBroadcastPkts"
	SNMP_ifHCInMulticastPkts  = "ifHCInMulticastPkts"
	SNMP_ifHCInBroadcastPkts  = "ifHCInBroadcastPkts"

	AllDeviceInterface = "device"
	Uptime             = "Uptime"
)

type InterfaceMetrics struct {
	log     logger.ContextL
	gconf   *kt.SnmpGlobalConfig
	conf    *kt.SnmpDeviceConfig
	metrics *kt.SnmpDeviceMetric

	// guards interfaceTracker and intValues.
	// As of 10/2020, only two goroutines should ever touch this struct:  a short-lived
	// goroutine at startup and then the persistent metric-poller.  But this stuff will
	// evolve, so let's be careful.
	mux sync.Mutex

	intValues map[string]*counters.CounterSet
	oidMap    map[string]string
}

func NewInterfaceMetrics(gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, profileMetrics map[string]*kt.Mib, log logger.ContextL) *InterfaceMetrics {
	oidMap := conf.InterfaceMetricsOidMap
	if oidMap == nil {
		oidMap = make(map[string]string)
	}
	for oid, m := range profileMetrics {
		noid := oid
		if !strings.HasPrefix(noid, ".") {
			noid = "." + noid
		}
		log.Infof("Adding interface metric %s -> %s", noid, m.Name)
		oidMap[noid] = m.Name
	}

	if len(oidMap) == 0 {
		oidMap = defaultOidMap
		log.Infof("Using default interface metric set")
	} else {
		log.Infof("Using custom interface metric set")
	}

	return &InterfaceMetrics{
		log:       log,
		gconf:     gconf,
		conf:      conf,
		metrics:   metrics,
		oidMap:    oidMap,
		intValues: make(map[string]*counters.CounterSet),
	}
}

func (im *InterfaceMetrics) DiscardDeltaState() {
	im.mux.Lock()
	defer im.mux.Unlock()

	im.intValues = make(map[string]*counters.CounterSet)
}

var (
	MAX_COUNTER_INTS = 250

	// TODO: ideally, this guy is the source of truth for SNMP_ifHCInOctets, etc,
	// currently defined in common.go

	// See https://tools.ietf.org/html/rfc2863.html and (for example)
	// http://www.oid-info.com/get/1.3.6.1.2.1.31.1.1.1 for explanations about
	// these and other snmp OIDs.
	defaultOidMap = map[string]string{
		"1.3.6.1.2.1.31.1.1.1.6":  SNMP_ifHCInOctets,     // 64 bit
		"1.3.6.1.2.1.31.1.1.1.7":  SNMP_ifHCInUcastPkts,  // 64 bit
		"1.3.6.1.2.1.31.1.1.1.10": SNMP_ifHCOutOctets,    // 64 bit
		"1.3.6.1.2.1.31.1.1.1.11": SNMP_ifHCOutUcastPkts, // 64 bit
		"1.3.6.1.2.1.2.2.1.14":    SNMP_ifInErrors,
		"1.3.6.1.2.1.2.2.1.20":    SNMP_ifOutErrors,
		"1.3.6.1.2.1.2.2.1.13":    SNMP_ifInDiscards,         // 32 bit in SNMP, 64 in ST; using 64 bit flex column
		"1.3.6.1.2.1.2.2.1.19":    SNMP_ifOutDiscards,        // same
		"1.3.6.1.2.1.31.1.1.1.12": SNMP_ifHCOutMulticastPkts, // 64 bit
		"1.3.6.1.2.1.31.1.1.1.13": SNMP_ifHCOutBroadcastPkts, // 64 bit
		"1.3.6.1.2.1.31.1.1.1.8":  SNMP_ifHCInMulticastPkts,  // 64 bit
		"1.3.6.1.2.1.31.1.1.1.9":  SNMP_ifHCInBroadcastPkts,  // 64 bit
	}
)

// PollSNMPCounter polls SNMP for counter statistics like # bytes and packets transferred.
func (im *InterfaceMetrics) Poll(server *gosnmp.GoSNMP, lastDeviceMetrics []*kt.JCHF) ([]*kt.JCHF, error) {
	im.mux.Lock()
	defer im.mux.Unlock()

	deltas := map[string]map[string]uint64{}

	for oid, varName := range im.oidMap {
		results, err := snmp_util.WalkOID(oid, server, im.log, "Counter")
		if err != nil {
			im.metrics.Errors.Mark(1)
			return nil, err
		}

		for _, variable := range results {
			parts := strings.Split(variable.Name, oid)
			if len(parts) != 2 || len(parts[1]) == 0 {
				continue
			}

			// variable.Name looks like this: .<oidVal>.<intVal>, e.g.
			// .1.3.6.1.2.1.31.1.1.1.10.219, where .1.3.6.1.2.1.31.1.1.1.10 is
			// the oid and 219 is the intVal.  So splitting on oidVal gives us
			// .intVal.
			intId := parts[1][1:]
			if _, ok := im.intValues[intId]; !ok {
				im.intValues[intId] = counters.NewCounterSetWithId(intId)
			}

			value := gosnmp.ToBigInt(variable.Value).Uint64()

			// Calculate the different of this counter here.
			delta, ok := deltas[intId]
			if !ok {
				delta = map[string]uint64{}
				deltas[intId] = delta
			}
			delta[varName] = im.intValues[intId].SetValueAndReturnDelta(varName, value)
		}
	}

	// See if we have a uptime delta to work with
	for _, dm := range lastDeviceMetrics {
		intId := AllDeviceInterface
		if _, ok := im.intValues[intId]; !ok {
			im.intValues[intId] = counters.NewCounterSetWithId(intId)
		}
		delta, ok := deltas[intId]
		if !ok {
			delta = map[string]uint64{}
			deltas[intId] = delta
		}
		if dm.CustomBigInt[Uptime] > 0 {
			delta[Uptime] = im.intValues[intId].SetValueAndReturnDelta(Uptime, uint64(dm.CustomBigInt[Uptime]))
			break // We got what we need here.
		}
	}

	// send this off encoded as chf as well as via tsdb
	flows := im.convertToCHF(deltas)

	im.log.Infof("SNMP interface metric poll - found metrics for %d interfaces", len(deltas))
	im.metrics.InterfaceMetrics.Mark(int64(len(flows)))

	return flows, nil
}

func (im *InterfaceMetrics) convertToCHF(deltas map[string]map[string]uint64) []*kt.JCHF {

	uptimeDelta := uint64(0)
	if deltas[AllDeviceInterface] != nil {
		uptimeDelta = deltas[AllDeviceInterface][Uptime]
	}

	flows := make([]*kt.JCHF, 0, len(deltas))
	for strint, cs := range deltas {
		if strint == AllDeviceInterface { // Don't put these here now.
			continue
		}
		intr, _ := strconv.Atoi(strint)
		dst := kt.NewJCHF()
		dst.CustomStr = make(map[string]string)
		dst.CustomInt = make(map[string]int32)
		dst.CustomBigInt = make(map[string]int64)
		dst.EventType = kt.KENTIK_EVENT_SNMP_INT_METRIC
		dst.Provider = im.conf.Provider
		dst.InputPort = kt.IfaceID(intr)
		dst.OutputPort = kt.IfaceID(intr)
		dst.DeviceName = im.conf.DeviceName
		dst.SrcAddr = im.conf.DeviceIP

		metrics := map[string]bool{}
		for k, v := range cs {
			dst.CustomBigInt[k] = int64(v)
			switch k {
			case SNMP_ifHCInUcastPkts, SNMP_ifHCOutUcastPkts, SNMP_ifHCInOctets, SNMP_ifHCOutOctets:
				dst.CustomBigInt[k] = dst.CustomBigInt[k] * im.conf.RateMultiplier
			}
			metrics[k] = true
		}
		if uptimeDelta > 0 {
			dst.CustomBigInt[Uptime] = int64(uptimeDelta)
		}

		dst.CustomMetrics = metrics // Add this in so that we know what metrics to pull out down the road.
		flows = append(flows, dst)
	}

	return flows
}
