package traps

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kentik/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	snmp_util "github.com/kentik/ktranslate/pkg/snmp/util"
)

type SnmpTrap struct {
	log       logger.ContextL
	jchfChan  chan []*kt.JCHF
	listen    string
	tl        *gosnmp.TrapListener
	metrics   *kt.SnmpMetricSet
	conf      *kt.SnmpConfig
	mux       sync.RWMutex
	deviceMap map[string]*kt.SnmpDeviceConfig
}

//Move to util?
type logWrapper struct {
	print  func(v ...interface{})
	printf func(format string, v ...interface{})
}

func (l logWrapper) Print(v ...interface{}) {
	l.print(v...)
}

func (l logWrapper) Printf(format string, v ...interface{}) {
	l.printf(format, v...)
}

func NewSnmpTrapListener(conf *kt.SnmpConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, log logger.ContextL) (*SnmpTrap, error) {
	st := &SnmpTrap{
		jchfChan:  jchfChan,
		log:       log,
		listen:    conf.Trap.Listen,
		metrics:   metrics,
		deviceMap: map[string]*kt.SnmpDeviceConfig{},
	}

	// Some quick defaults.
	if conf.Trap.Transport == "" {
		conf.Trap.Transport = "udp"
	}
	if conf.Trap.Community == "" {
		conf.Trap.Community = "hello"
	}
	if conf.Global.TimeoutMS == 0 {
		conf.Global.TimeoutMS = 5000
	}

	pts := strings.Split(conf.Trap.Listen, ":")
	port := 161
	if len(pts) > 1 {
		port, _ = strconv.Atoi(pts[1])
	}

	// Now set things up.
	tl := gosnmp.NewTrapListener()
	tl.OnNewTrap = st.handle
	tl.Params = &gosnmp.GoSNMP{
		Port:               uint16(port),
		Transport:          conf.Trap.Transport,
		Community:          conf.Trap.Community,
		Timeout:            time.Duration(conf.Global.TimeoutMS) * time.Millisecond,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}
	switch conf.Trap.Version {
	case "v2c", "":
		tl.Params.Version = gosnmp.Version2c
	case "v3":
		tl.Params.Version = gosnmp.Version3
	default:
		return nil, fmt.Errorf("Invalid trap version: %s", conf.Trap.Version)
	}

	tl.Params.Logger = logWrapper{
		print: func(v ...interface{}) {
			log.Debugf("GoSNMP Trap:" + fmt.Sprint(v...))
		},
		printf: func(format string, v ...interface{}) {
			log.Debugf("GoSNMP Trap:  "+format, v...)
		},
	}
	st.tl = tl

	for _, device := range conf.Devices {
		st.deviceMap[device.DeviceIP] = device
	}

	return st, nil
}

func (s *SnmpTrap) Listen() {
	err := s.tl.Listen(s.listen)
	if err != nil {
		s.log.Errorf("error in Trap listen: %s", err)
	}
}

func (s *SnmpTrap) handle(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	s.log.Infof("got trapdata from %s", addr.IP)
	s.metrics.Traps.Mark(1)
	s.mux.RLock()
	defer s.mux.RUnlock()

	dev := s.deviceMap[addr.IP.String()] // See if we know which device this is coming from.
	dst := kt.NewJCHF()
	dst.CustomStr = make(map[string]string)
	dst.CustomInt = make(map[string]int32)
	dst.CustomBigInt = make(map[string]int64)
	dst.EventType = kt.KENTIK_EVENT_SNMP_TRAP
	dst.SrcAddr = addr.IP.String()
	if dev != nil {
		dst.DeviceName = dev.DeviceName
		dst.Provider = dev.Provider
	}

	for _, v := range packet.Variables {
		switch v.Type {
		case gosnmp.OctetString:
			if value, ok := snmp_util.ReadOctetString(v, snmp_util.NO_TRUNCATE); ok {
				dst.CustomStr[v.Name] = value
			}
		case gosnmp.Counter64, gosnmp.Counter32, gosnmp.Gauge32, gosnmp.TimeTicks, gosnmp.Uinteger32:
			dst.CustomBigInt[v.Name] = gosnmp.ToBigInt(v.Value).Int64()
		default:
			s.log.Infof("trap variable with unknown type handling, skipping: %+v", v)
		}
	}

	s.jchfChan <- []*kt.JCHF{dst}
}
