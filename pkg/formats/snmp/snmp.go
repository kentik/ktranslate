package snmp

import (
	"flag"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/gosnmp/gosnmp"
)

const (
	snmpTrapOID_0   = ".1.3.6.1.6.3.1.1.4.1.0"
	KTransTrapIdent = ".1.3.6.1.4.1.41263.6169"
)

type SnmpFormat struct {
	logger.ContextL
	conf *kt.SnmpConfig
	ts   *gosnmp.GoSNMP
}

var (
	confFile string
)

func init() {
	flag.StringVar(&confFile, "snmp.format.conf", "", "Parse this file for the snmp format option. Same format as -snmp flag.")
}

func NewFormat(log logger.Underlying, cfg *ktranslate.SnmpFormatConfig) (*SnmpFormat, error) {
	sf := &SnmpFormat{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "snmpFormat"}, log),
	}

	if cfg.ConfigFile == "" {
		return nil, fmt.Errorf("-snmp.format.conf or ConfigFile required for snmp format.")
	}
	snmpConf, err := snmp.ParseConfig(cfg.ConfigFile, sf)
	if err != nil {
		return nil, err
	}
	sf.conf = snmpConf

	if sf.conf.Trap.Endpoint == "" {
		return nil, fmt.Errorf("endpoint in trap required for snmp format.")
	}

	port := sf.conf.Trap.EndpointPort
	if port == 0 {
		port = 162
	}
	ts := &gosnmp.GoSNMP{
		Transport:          sf.conf.Trap.Transport,
		Community:          sf.conf.Trap.Community,
		Timeout:            time.Duration(sf.conf.Global.TimeoutMS) * time.Millisecond,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
		Target:             sf.conf.Trap.Endpoint,
		Port:               port,
	}
	switch sf.conf.Trap.Version {
	case "v1":
		ts.Version = gosnmp.Version1
	case "v2c", "":
		ts.Version = gosnmp.Version2c
	case "v3":
		params, flags, contextEngineID, contextName, err := snmp_util.ParseV3Config(sf.conf.Trap.V3)
		if err != nil {
			return nil, err
		}
		ts.Version = gosnmp.Version3
		ts.SecurityModel = gosnmp.UserSecurityModel
		ts.MsgFlags = flags
		ts.SecurityParameters = params
		ts.ContextEngineID = contextEngineID
		ts.ContextName = contextName
	default:
		return nil, fmt.Errorf("Invalid trap version: %s", sf.conf.Trap.Version)
	}

	ts.Logger = gosnmp.NewLogger(logWrapper{
		print: func(v ...interface{}) {
			sf.Debugf("GoSNMP Trap Send:" + fmt.Sprint(v...))
		},
		printf: func(format string, v ...interface{}) {
			sf.Debugf("GoSNMP Trap Send:  "+format, v...)
		},
	})
	sf.ts = ts

	err = sf.ts.Connect()
	if err != nil {
		return nil, err
	}

	sf.Infof("Online. Sending to %s:%d via %s", sf.conf.Trap.Endpoint, port, sf.conf.Trap.Version)

	return sf, nil
}

// Outputs as a snmp trap.
// Turn each message into the contents of a trap, send it on.
func (f *SnmpFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	for _, m := range msgs {
		trap := gosnmp.SnmpTrap{
			Variables: []gosnmp.SnmpPDU{
				gosnmp.SnmpPDU{
					Name:  snmpTrapOID_0,
					Type:  gosnmp.ObjectIdentifier,
					Value: KTransTrapIdent,
				},
			},
		}

		flat := m.Flatten()
		strip(flat)
		for k, v := range flat {
			switch tv := v.(type) {
			case string:
				trap.Variables = append(trap.Variables, gosnmp.SnmpPDU{
					Name:  getKeyOid(k),
					Type:  gosnmp.OctetString,
					Value: tv,
				})
			case int32:
				trap.Variables = append(trap.Variables, gosnmp.SnmpPDU{
					Name:  getKeyOid(k),
					Type:  gosnmp.Gauge32,
					Value: uint32(tv),
				})
			case int64:
				trap.Variables = append(trap.Variables, gosnmp.SnmpPDU{
					Name:  getKeyOid(k),
					Type:  gosnmp.Counter64,
					Value: uint64(tv),
				})
			}
		}

		// And send this on.
		_, err := f.ts.SendTrap(trap)
		if err != nil {
			f.Errorf("SendTrap() err: %v", err)
		}

	}

	return nil, nil
}

func (f *SnmpFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	return nil, nil
}

func (f *SnmpFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

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

func strip(in map[string]interface{}) {
	for k, v := range in {
		switch tv := v.(type) {
		case string:
			if tv == "" || tv == "-" || tv == "--" {
				delete(in, k)
			}
		case int32:
			if tv == 0 {
				delete(in, k)
			}
		case int64:
			if tv == 0 {
				delete(in, k)
			}
		}
	}
	in["instrumentation.provider"] = kt.InstProvider // Let them know who sent this.
	in["collector.name"] = kt.CollectorName
}

func getKeyOid(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%s.%s", KTransTrapIdent, formatKeyOid(h.Sum32()))
}

func formatKeyOid(n uint32) string {
	in := strconv.FormatUint(uint64(n), 10)
	numOfDigits := len(in)
	numOfDots := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfDots)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = '.'
		}
	}
}
