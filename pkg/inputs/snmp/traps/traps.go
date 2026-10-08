package traps

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/resolv"
)

const (
	snmpTrapOID   = ".1.3.6.1.6.3.1.1.4.1"
	snmpTrapOID_0 = ".1.3.6.1.6.3.1.1.4.1.0"

	// snmpTrapAddress (RFC 3584) carries the original Agent Address when a
	// proxy translates a v1 trap to v2c/v3. Checked as part of UseStdSender.
	snmpTrapAddressOID = ".1.3.6.1.6.3.18.1.3.0"

	senderSourceStdAddr = "snmpTrapAddress"
	senderSourceAgent   = "agentAddr"
	senderSourceUDP     = "udp"
)

var (
	v1TrapEnum = map[int]string{
		0: "coldStart",
		1: "warmStart",
		2: "linkDown",
		3: "linkUp",
		4: "authenticationFailure",
		5: "egpNeighborLoss",
		6: "enterpriseSpeciﬁc",
	}
)

type SnmpTrap struct {
	log           logger.ContextL
	jchfChan      chan []*kt.JCHF
	listen        string
	tl            *gosnmp.TrapListener
	metrics       *kt.SnmpMetricSet
	conf          *kt.SnmpConfig
	mux           sync.RWMutex
	mibdb         *mibs.MibDB
	resolv        *resolv.Resolver
	ctx           context.Context
	deviceMap     map[string]*kt.SnmpDeviceConfig
	deviceMapName map[string]*kt.SnmpDeviceConfig
	trustedRelays map[string]bool // Empty/nil means any UDP source may have its sender overridden.
	baseTags      map[string]string
}

// Move to util?
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

func NewSnmpTrapListener(ctx context.Context, conf *kt.SnmpConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, mibdb *mibs.MibDB,
	log logger.ContextL, resolv *resolv.Resolver, serviceName string, globalTags map[string]string) (*SnmpTrap, error) {
	st := &SnmpTrap{
		jchfChan:      jchfChan,
		log:           log,
		mibdb:         mibdb,
		listen:        conf.Trap.Listen,
		metrics:       metrics,
		deviceMap:     map[string]*kt.SnmpDeviceConfig{},
		deviceMapName: map[string]*kt.SnmpDeviceConfig{},
		resolv:        resolv,
		ctx:           ctx,
		conf:          conf,
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

	// Now set things up.
	tl := gosnmp.NewTrapListener()
	tl.OnNewTrap = st.handle
	tl.Params = &gosnmp.GoSNMP{
		Transport:          conf.Trap.Transport,
		Community:          conf.Trap.Community,
		Timeout:            time.Duration(conf.Global.TimeoutMS) * time.Millisecond,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}
	switch conf.Trap.Version {
	case "v1":
		tl.Params.Version = gosnmp.Version1
	case "v2c", "":
		tl.Params.Version = gosnmp.Version2c
	case "v3":
		params, flags, contextEngineID, contextName, err := snmp_util.ParseV3Config(conf.Trap.V3)
		if err != nil {
			return nil, err
		}
		tl.Params.Version = gosnmp.Version3
		tl.Params.SecurityModel = gosnmp.UserSecurityModel
		tl.Params.MsgFlags = flags
		tl.Params.SecurityParameters = params
		tl.Params.ContextEngineID = contextEngineID
		tl.Params.ContextName = contextName
	default:
		return nil, fmt.Errorf("Invalid trap version: %s", conf.Trap.Version)
	}

	tl.Params.Logger = gosnmp.NewLogger(logWrapper{
		print: func(v ...interface{}) {
			log.Debugf("GoSNMP Trap:" + fmt.Sprint(v...))
		},
		printf: func(format string, v ...interface{}) {
			log.Debugf("GoSNMP Trap:  "+format, v...)
		},
	})
	st.tl = tl
	log.Infof("Trap listener setup with version %s on %s. DropUndefined: %v", conf.Trap.Version, conf.Trap.Listen, conf.Trap.DropUndefined)

	for _, device := range conf.Devices {
		st.deviceMap[device.DeviceIP] = device
		if device.DeviceName != "" {
			st.deviceMapName[device.DeviceName] = device
		}
	}

	if len(conf.Trap.TrustedRelays) > 0 {
		st.trustedRelays = map[string]bool{}
		for _, relay := range conf.Trap.TrustedRelays {
			st.trustedRelays[relay] = true
		}
	}

	// Set up some default tags if the device sending isn't found.
	baseTags := map[string]string{}
	if serviceName != "" {
		baseTags[kt.UserTagPrefix+"container_service"] = serviceName
	}
	for k, v := range globalTags {
		key := k
		if !strings.HasPrefix(key, kt.UserTagPrefix) {
			key = kt.UserTagPrefix + k
		}
		baseTags[key] = v
	}
	st.baseTags = baseTags

	return st, nil
}

func (s *SnmpTrap) Listen() {
	err := s.tl.Listen(s.listen)
	if err != nil {
		s.log.Errorf("error in Trap listen: %s", err)
	}
}

// varValueAsString pulls a string out of a varbind value, for the types that can
// plausibly carry an IP address or hostname (OctetString, IPAddress, ObjectIdentifier).
// An empty value is reported as absent, so resolveSender keeps looking rather than
// settling on a varbind that names nothing.
func varValueAsString(v gosnmp.SnmpPDU) (string, bool) {
	switch v.Type {
	case gosnmp.OctetString:
		// ReadOctetString reports ok for a value that trims down to "", hence the recheck.
		// Sanitize as the trap variable loop below does: a device returning binary here
		// should not reach the output as invalid UTF-8.
		if s, ok := snmp_util.ReadOctetString(v, snmp_util.NO_TRUNCATE); ok && s != "" {
			return kt.SanitizeUTF8(s), true
		}
	case gosnmp.IPAddress, gosnmp.ObjectIdentifier:
		if s, ok := v.Value.(string); ok && s != "" {
			return s, true
		}
	}
	return "", false
}

// findVarValue returns the value of the first varbind matching oid, if any.
func findVarValue(vars []gosnmp.SnmpPDU, oid string) (string, bool) {
	for _, v := range vars {
		if v.Name == oid {
			return varValueAsString(v)
		}
	}
	return "", false
}

// resolveSender looks for the trap's original source device inside the packet itself,
// so traps relayed through a proxy or NAT can still be mapped to the right device. See
// SnmpTrapConfig.SenderVarOids / UseStdSender. Returns ("", senderSourceUDP) when nothing
// in the packet should override the UDP source address.
func resolveSender(packet *gosnmp.SnmpPacket, udpAddr string, cfg *kt.SnmpTrapConfig, trustedRelays map[string]bool) (string, string) {
	if len(trustedRelays) > 0 && !trustedRelays[udpAddr] { // Only trust the payload from a known relay, when configured.
		return "", senderSourceUDP
	}

	for _, oid := range cfg.SenderVarOids {
		if v, ok := findVarValue(packet.Variables, oid); ok {
			return v, "varoid:" + oid
		}
	}

	if cfg.UseStdSender {
		// A v1 trap names its originator in the PDU's own AgentAddress field, so that
		// wins there. snmpTrapAddress is what RFC 3584 has a proxy add when it translates
		// such a trap to v2c/v3, so it only applies to those versions.
		if packet.Version == gosnmp.Version1 {
			if packet.AgentAddress != "" {
				return packet.AgentAddress, senderSourceAgent
			}
		} else if v, ok := findVarValue(packet.Variables, snmpTrapAddressOID); ok {
			return v, senderSourceStdAddr
		}
	}

	return "", senderSourceUDP
}

// lookupDevice finds a configured device by the resolved sender value, trying it as an
// IP address first and then as a device name.
func (s *SnmpTrap) lookupDevice(sender string) *kt.SnmpDeviceConfig {
	if sender == "" {
		return nil
	}
	if dev, ok := s.deviceMap[sender]; ok {
		return dev
	}
	return s.deviceMapName[sender]
}

func (s *SnmpTrap) handle(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	engineID := "" // Decode the context engine id if sent here.
	if packet.Version == gosnmp.Version3 && packet.ContextEngineID != "" {
		_, eid, _ := snmp_util.EngineID([]byte(packet.ContextEngineID))
		engineID = eid
	}

	s.log.Infof("got trapdata from %s, EngineID %s", addr.IP, engineID)
	s.metrics.Traps.Mark(1)
	s.mux.RLock()
	defer s.mux.RUnlock()

	udpAddr := addr.IP.String()
	sender, senderSource := resolveSender(packet, udpAddr, s.conf.Trap, s.trustedRelays)

	dev := s.lookupDevice(sender) // See if the trap itself tells us who really sent this.
	senderResolved := dev != nil
	if dev == nil { // Fall back to the UDP source address.
		dev = s.deviceMap[udpAddr]
	}
	if dev == nil && engineID != "" { // If the device is still nil, try looking up via the ContextEngineID.
		for _, d := range s.deviceMap {
			if d.EngineID == engineID {
				dev = d
				break
			}
		}
	}
	dst := kt.NewJCHF()
	dst.CustomStr = make(map[string]string)
	dst.CustomInt = make(map[string]int32)
	dst.CustomBigInt = make(map[string]int64)
	dst.EventType = kt.KENTIK_EVENT_SNMP_TRAP
	dst.SrcAddr = udpAddr
	if senderResolved { // Only report the override when it actually resolved a device.
		// Formats parse SrcAddr as an IP (see pkg/formats/kflow), so a sender matched by
		// DeviceName contributes its configured address here and its declared name below.
		if dev.DeviceIP != "" {
			dst.SrcAddr = dev.DeviceIP
		}
		dst.CustomStr["SenderAddr"] = sender
		dst.CustomStr["SenderSource"] = senderSource
		dst.CustomStr["RelayAddr"] = udpAddr
	}
	if dev != nil {
		dst.DeviceName = dev.DeviceName
		if s.conf.Trap.TrapOnly {
			dst.Provider = kt.ProviderTrapUnknown
		} else {
			dst.Provider = dev.Provider
		}
		dev.SetUserTags(dst.CustomStr)
	} else {
		dst.DeviceName = addr.IP.String()
		if s.resolv != nil { // Try pulling device name from reverse IP if possible.
			dm := s.resolv.Resolve(s.ctx, dst.DeviceName, true)
			if dm != "" {
				dst.DeviceName = dm
			}
		}
		dst.Provider = kt.ProviderTrapUnknown
		for k, v := range s.baseTags {
			dst.CustomStr[k] = v
		}
	}

	// What trap is this from?
	found := false // Set to true ONLY if the TrapName is set.
	trapOid := ""
	var trap *mibs.Trap
	for _, v := range packet.Variables {
		if v.Name == snmpTrapOID || v.Name == snmpTrapOID_0 {
			if v.Type == gosnmp.ObjectIdentifier {
				toid := v.Value.(string)
				trap = s.mibdb.GetTrap(toid)
				dst.CustomStr["TrapOID"] = toid
				trapOid = toid
				if trap != nil {
					dst.CustomStr["TrapName"] = trap.Name
					idx := snmp_util.GetIndex(toid, trap.Oid)
					if idx != "" {
						dst.CustomStr["Index"] = idx
					}
					found = true
				}
			}
		}
	}

	// Handle v1 traps here.
	if packet.Version == gosnmp.Version1 {
		dst.CustomStr["GenericTrap"] = fmt.Sprintf("%s (%d)", v1TrapEnum[packet.GenericTrap], packet.GenericTrap)
		dst.CustomInt["SpecificTrap"] = int32(packet.SpecificTrap)
		dst.CustomStr["Enterprise"] = packet.Enterprise
		dst.CustomBigInt["Timestamp"] = int64(packet.Timestamp)
	}

	for _, v := range packet.Variables {
		if v.Name == snmpTrapOID || v.Name == snmpTrapOID_0 {
			continue
		}

		// Do we know this guy?
		res, vars, err := s.mibdb.GetForKey(v.Name)
		if err != nil {
			s.log.Errorf("Cannot look up OID in trap: %v", err)
		}

		// If we don't want undefined vars, pass along here.
		if res == nil && (s.conf.Trap.DropUndefined || trap.DropUndefinedVars()) {
			s.log.Infof("Trap variable dropped: %s", v.Name)
			continue
		}

		// Load any variables defined in the name here.
		for n, value := range vars {
			dst.CustomStr[n] = value
		}

		switch v.Type {
		case gosnmp.OctetString:
			if value, ok := snmp_util.ReadOctetString(v, snmp_util.NO_TRUNCATE); ok {
				if res != nil && res.Conversion != "" { // Adjust for any hard coded values here.
					_, sval, mval := snmp_util.GetFromConv(v, res.Conversion, s.log)
					if len(mval) > 0 { // List the regex matches.
						for k, v := range mval {
							dst.CustomStr[k] = v
						}
					} else {
						dst.CustomStr[res.GetName()] = sval
					}
				} else { // No conversion.
					if res != nil {
						dst.CustomStr[res.GetName()] = kt.SanitizeUTF8(value)
					} else {
						dst.CustomStr[v.Name] = kt.SanitizeUTF8(value)
					}
				}
			}
		case gosnmp.Counter64, gosnmp.Counter32, gosnmp.Gauge32, gosnmp.TimeTicks, gosnmp.Uinteger32, gosnmp.Integer:
			if res != nil {
				dst.CustomBigInt[res.GetName()] = gosnmp.ToBigInt(v.Value).Int64()
			} else {
				dst.CustomBigInt[v.Name] = gosnmp.ToBigInt(v.Value).Int64()
			}
		case gosnmp.ObjectIdentifier, gosnmp.IPAddress: // Both decode to a plain string value in gosnmp.
			value := v.Value.(string)
			if res != nil && res.Conversion != "" { // Adjust for any hard coded values here.
				_, value, _ = snmp_util.GetFromConv(v, res.Conversion, s.log)
			}
			value = kt.SanitizeUTF8(value)
			if res != nil {
				dst.CustomStr[res.GetName()] = value
			} else {
				dst.CustomStr[v.Name] = value
			}
		default:
			s.log.Infof("trap variable with unknown type (%v) handling, skipping: %+v", v.Type, v)
		}
	}

	// If dropping, only send on if found.
	if s.conf.Trap.DropUndefined || trap.DropUndefinedVars() {
		if found {
			s.jchfChan <- []*kt.JCHF{dst}
		} else {
			s.log.Infof("Whole trap packet dropped: %s", trapOid)
		}
	} else { // Else, keep to current behavor.
		s.jchfChan <- []*kt.JCHF{dst}
	}
}
