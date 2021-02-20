package traps

import (
	"fmt"
	"net"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/gosnmp"
	"github.com/kentik/ktranslate/pkg/kt"
	snmp_util "github.com/kentik/ktranslate/pkg/snmp/util"
)

type SnmpTrap struct {
	log      logger.ContextL
	jchfChan chan []*kt.JCHF
	listen   string
	tl       *gosnmp.TrapListener
	metrics  *kt.SnmpMetricSet
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

func NewSnmpTrapListener(conf *kt.SnmpTrapConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, log logger.ContextL) (*SnmpTrap, error) {
	st := &SnmpTrap{
		jchfChan: jchfChan,
		log:      log,
		listen:   conf.Listen,
		metrics:  metrics,
	}

	// Some quick defaults.
	if conf.Transport == "" {
		conf.Transport = "udp"
	}
	if conf.Community == "" {
		conf.Community = "hello"
	}

	// Now set things up.
	tl := gosnmp.NewTrapListener()
	tl.OnNewTrap = st.handle
	tl.Params = &gosnmp.GoSNMP{
		Port:               161,
		Transport:          conf.Transport,
		Community:          conf.Community,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}
	switch conf.Version {
	case "v2c", "":
		tl.Params.Version = gosnmp.Version2c
	case "v3":
		tl.Params.Version = gosnmp.Version3
	default:
		return nil, fmt.Errorf("Invalid trap version: %s", conf.Version)
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

	return st, nil
}

func (s *SnmpTrap) Listen() {
	err := s.tl.Listen(s.listen)
	if err != nil {
		s.log.Errorf("error in Trap listen: %s", err)
	}
}

func (s *SnmpTrap) handle(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	s.log.Debugf("got trapdata from %s", addr.IP)
	s.metrics.Traps.Mark(1)

	dst := kt.NewJCHF()
	dst.CustomStr = make(map[string]string)
	dst.CustomInt = make(map[string]int32)
	dst.CustomBigInt = make(map[string]int64)
	dst.EventType = kt.KENTIK_EVENT_SNMP_TRAP
	dst.SrcAddr = addr.IP.String()

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
