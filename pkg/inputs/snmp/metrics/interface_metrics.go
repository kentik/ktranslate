package metrics

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/kt/counters"
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

	SNMP_ifInErrorsPercent  = "ifInErrorPercent"
	SNMP_ifOutErrorsPercent = "ifOutErrorPercent"

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

	intValues   map[string]*counters.CounterSet
	oidMap      map[string]string
	nameOidMap  map[string]string
	oidMibMap   map[string]*kt.Mib
	profileName string
}

func NewInterfaceMetrics(gconf *kt.SnmpGlobalConfig, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, profileMetrics map[string]*kt.Mib, profile *mibs.Profile, log logger.ContextL) *InterfaceMetrics {
	oidMap := make(map[string]string)
	for oid, m := range profileMetrics {
		noid := oid
		if !strings.HasPrefix(noid, ".") {
			noid = "." + noid
		}
		oidName := m.GetName()
		log.Infof("Adding interface metric %s -> %s", noid, oidName)
		oidMap[noid] = oidName
	}

	nameOidMap := map[string]string{} // Reverse the polled oids here.
	for k, v := range oidMap {
		nameOidMap[v] = k
	}

	return &InterfaceMetrics{
		log:         log,
		gconf:       gconf,
		conf:        conf,
		metrics:     metrics,
		oidMap:      oidMap,
		nameOidMap:  nameOidMap,
		oidMibMap:   profileMetrics,
		intValues:   make(map[string]*counters.CounterSet),
		profileName: profile.GetProfileName(conf.InstrumentationName),
	}
}

func (im *InterfaceMetrics) DiscardDeltaState() {
	im.mux.Lock()
	defer im.mux.Unlock()

	im.intValues = make(map[string]*counters.CounterSet)
}

var (
	MAX_COUNTER_INTS = 250
)

// PollSNMPCounter polls SNMP for counter statistics like # bytes and packets transferred.
func (im *InterfaceMetrics) Poll(ctx context.Context, server *gosnmp.GoSNMP, lastDeviceMetrics []*kt.JCHF) ([]*kt.JCHF, error) {
	im.mux.Lock()
	defer im.mux.Unlock()

	deltas := map[string]map[string]uint64{}

	for oid, varName := range im.oidMap {
		if mib, ok := im.oidMibMap[oid[1:]]; ok {
			if !mib.IsPollReady() { // Skip this mib because its time to poll hasn't elapsed yet.
				continue
			}
		}

		results, err := snmp_util.WalkOID(ctx, im.conf, oid, server, im.log, "Counter")
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

			delta, ok := deltas[intId]
			if !ok {
				delta = map[string]uint64{}
				deltas[intId] = delta
			}
			switch variable.Type {
			case gosnmp.Integer:
				// Since its just an int, keep it without computing a delta.
				delta[varName] = value
			default:
				// Treat as a counter
				// Calculate the different of this counter here.
				delta[varName] = im.intValues[intId].SetValueAndReturnDelta(varName, value)
			}
		}
	}

	im.log.Infof("SNMP interface metric poll - found metrics for %d interfaces.", len(deltas))

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
		dst.Timestamp = time.Now().Unix()
		if dst.Provider == kt.ProviderDefault { // Add this to trigger a UI element.
			dst.CustomStr["profile_message"] = kt.DefaultProfileMessage
		}

		metrics := map[string]kt.MetricInfo{}
		for k, v := range cs {
			mib := im.oidMibMap[im.nameOidMap[k][1:]]
			if mib != nil && mib.EnumRev != nil {
				if ev, ok := mib.EnumRev[int64(v)]; ok {
					// If we know what the enum string is, set it here.
					dst.CustomStr[kt.StringPrefix+k] = ev
				} else {
					im.log.Warnf("Missing enum value for interface metric %s %d", k, v)
					dst.CustomStr[kt.StringPrefix+k] = kt.InvalidEnum
				}
			}
			dst.CustomBigInt[k] = int64(v)
			if mib != nil {
				metrics[k] = kt.MetricInfo{Oid: im.nameOidMap[k], Mib: mib.Mib, Profile: im.profileName, Table: mib.Table, Format: kt.CountMetric, PollDur: mib.PollDur}
			} else {
				metrics[k] = kt.MetricInfo{Oid: im.nameOidMap[k], Profile: im.profileName, Format: kt.CountMetric}
			}
		}

		// Drop in Error %s here if appicable.
		if dst.CustomBigInt[SNMP_ifHCInUcastPkts] > 0 {
			if idi, ok := im.nameOidMap[SNMP_ifInErrors]; ok {
				if mib, ok := im.oidMibMap[idi[1:]]; ok {
					dst.CustomBigInt[SNMP_ifInErrorsPercent] = int64(float64(dst.CustomBigInt[SNMP_ifInErrors]) / float64(dst.CustomBigInt[SNMP_ifHCInUcastPkts]) * 100.)
					metrics[SNMP_ifInErrorsPercent] = kt.MetricInfo{Profile: im.profileName, Format: kt.GaugeMetric, Table: mib.Table, PollDur: mib.PollDur, Mib: mib.Mib}
				}
			}
		}
		if dst.CustomBigInt[SNMP_ifHCOutUcastPkts] > 0 {
			if ido, ok := im.nameOidMap[SNMP_ifOutErrors]; ok {
				if mib, ok := im.oidMibMap[ido[1:]]; ok {
					dst.CustomBigInt[SNMP_ifOutErrorsPercent] = int64(float64(dst.CustomBigInt[SNMP_ifOutErrors]) / float64(dst.CustomBigInt[SNMP_ifHCOutUcastPkts]) * 100.)
					metrics[SNMP_ifOutErrorsPercent] = kt.MetricInfo{Profile: im.profileName, Format: kt.GaugeMetric, Table: mib.Table, PollDur: mib.PollDur, Mib: mib.Mib}
				}
			}
		}

		if uptimeDelta > 0 {
			dst.CustomBigInt[Uptime] = int64(uptimeDelta)
		}

		dst.CustomMetrics = metrics // Add this in so that we know what metrics to pull out down the road.
		flows = append(flows, dst)
	}

	return flows
}
