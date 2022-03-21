package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/util/vendor"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	CONV_HWADDR    = "hwaddr"
	CONV_POWERSET  = "powerset_status"
	CONV_HEXTOINT  = "hextoint"
	CONV_HEXTOIP   = "hextoip"
	CONV_ENGINE_ID = "engine_id"
	CONV_REGEXP    = "regexp"
	CONV_ONE       = "to_one"
)

var (
	SNMP_POLL_SLEEP_TIME = 10 * time.Second

	// Device Manufacturer aka sysDescr
	SNMP_DEVICE_MANUFACTURER_OID = "1.3.6.1.2.1.1.1"

	// Keep a cache of seen regexps seen to speed things up.
	reCache = map[string]*regexp.Regexp{}
)

func ContainsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

var (
	TRUNCATE     = true
	NO_TRUNCATE  = false
	MAX_SNMP_LEN = 128
)

func ReadOctetString(variable gosnmp.SnmpPDU, truncate bool) (string, bool) {
	if variable.Type != gosnmp.OctetString {
		return "", false
	}

	value := variable.Value.([]byte)
	value = bytes.Trim(value, "\x00")
	if truncate && len(value) > MAX_SNMP_LEN {
		value = value[0:MAX_SNMP_LEN]
	}
	return string(value), true
}

// toInt64 converts most of the numeric types in Go (and *all* the numeric
// types I expect to see from snmp) to int64.
func ToInt64(val interface{}) int64 {
	switch v := val.(type) {
	case uint:
		return int64(v)
	case uint8: // same as byte
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32: // same as rune
		return int64(v)
	case int64:
		return int64(v)
	}
	return 0 // Probably not reached, in context.
}

// getIndex returns the rest of value after prefix, e.g. for 1.2.3.4.5.6 and
// 2.3.4, returns .5.6  Modeled on getIndex in
// topology-demo/devicemetrics/main.go.  Prefix must occur in value, otherwise
// results are undefined (aka "wrong").
func GetIndex(value, prefix string) string {
	posit := strings.Index(value, prefix) + len(prefix)
	if len(value) > posit {
		return value[posit:]
	}
	return ""
}

// walk the OID subtree under a root, returning a slice of varbinds
func WalkOID(ctx context.Context, device *kt.SnmpDeviceConfig, oid string, server *gosnmp.GoSNMP, log logger.ContextL, logName string) ([]gosnmp.SnmpPDU, error) {

	// New strategy -- for each varbind, we'll try three times:
	//   first, with GetBulk
	//   if that fails, try GetNext without sleeping in-between
	//   if that fails, sleep for a while, then try GetNext again
	// The first retry is really a fallback -- see if the device is just unresponsive
	// to some GetBulk requests.  Turns out that happens reasonably frequently.
	// The second one is to see if there's a temporary load issue on the device,
	// and if we wait a while, things will get better.
	type pollTry struct {
		walk  func(string) ([]gosnmp.SnmpPDU, error)
		sleep time.Duration
	}

	tries := []pollTry{
		pollTry{walk: server.BulkWalkAll, sleep: time.Duration(0)},
		pollTry{walk: server.WalkAll, sleep: time.Duration(0)},
		pollTry{walk: server.WalkAll, sleep: SNMP_POLL_SLEEP_TIME},
	}

	// If we are overriding with a test set.
	if walker := device.GetTestWalker(); walker != nil {
		tries = []pollTry{pollTry{walk: walker.WalkAll, sleep: time.Duration(0)}}
	} else if device.NoUseBulkWalkAll { // If the device says to not use bulkwalkall, trim this out now.
		tries = tries[1:]
	}

	var err error
	var results []gosnmp.SnmpPDU
	for i, try := range tries {
		time.Sleep(try.sleep)

		results, err = try.walk(oid)
		if err == nil {
			if i > 0 {
				log.Infof("%s SNMP retry %d on OID %s succeeded", logName, i, oid)
			}
			return results, nil
		}

		log.Warnf("There was an SNMP polling error with the %s walking OID %s after %d retries: %v.", logName, oid, i, err)
	}

	log.Warnf("There was an error with the %s OID in the SNMP retry: %s.", logName, oid)
	return nil, err
}

type snmpWalker interface {
	WalkAll(string) ([]gosnmp.SnmpPDU, error)
}

func GetDeviceManufacturer(server snmpWalker, log logger.ContextL) string {
	results, err := server.WalkAll(SNMP_DEVICE_MANUFACTURER_OID)
	if err != nil {
		log.Debugf("Error retrieving SNMP device manufacturer; ignoring it: %v", err)
		return ""
	}
	if len(results) == 0 {
		return ""
	}
	deviceManufacturerEnc, ok := results[0].Value.([]byte)
	// Don't know why it wouldn't be a []byte, but just in case
	if !ok {
		log.Debugf("getDeviceManufacturer: received a non-[]byte: %v", results[0].Value)
		return ""
	}
	deviceManufacturerBytes, err := base64.StdEncoding.DecodeString(string(deviceManufacturerEnc))
	// An error (probably) just means it's not actually base64 encoded; assume plain text.
	if err != nil {
		deviceManufacturerBytes = deviceManufacturerEnc
	}
	deviceManufacturer :=
		strings.TrimSpace(
			strings.Replace(
				strings.Replace(string(deviceManufacturerBytes), "\n", "/", -1),
				"\r", "", -1))
	deviceManufacturerRunes := []rune(deviceManufacturer)
	if len(deviceManufacturerRunes) > 128 {
		deviceManufacturerRunes = deviceManufacturerRunes[:128]
	}
	deviceManufacturer = strconv.QuoteToASCII(string(deviceManufacturerRunes))
	// Strip the leading & trailing quotes that QuoteToASCII adds.
	deviceManufacturer = deviceManufacturer[1 : len(deviceManufacturer)-1]
	return deviceManufacturer
}

func PrettyPrint(pdu gosnmp.SnmpPDU, format string, log logger.ContextL) string {
	switch pdu.Type {
	case gosnmp.OctetString:
		_, s := GetFromConv(pdu, format, log)
		return s
	case gosnmp.IPAddress:
		return pdu.Value.(string)
	case gosnmp.ObjectIdentifier:
		return pdu.Value.(string)
	default:
		v := ToInt64(pdu.Value)
		return strconv.Itoa(int(v))
	}
}

// Does a walk of the targeted device and exits.
func DoWalk(device string, baseOid string, format string, conf *kt.SnmpConfig, connectTimeout time.Duration, retries int, log logger.ContextL) error {
	dconf := conf.Devices[device]
	if dconf == nil {
		return fmt.Errorf("The %s device was not found in the SNMP configuration file.", device)
	}

	server, err := InitSNMP(dconf, connectTimeout, retries, "", log)
	if err != nil {
		return err
	}

	res, err := WalkOID(context.Background(), dconf, baseOid, server, log, "DoWalk")
	if err != nil {
		return err
	}

	for _, variable := range res {
		log.Infof("%s snmpwalk result: %s = %v: %s", device, variable.Name, variable.Type, PrettyPrint(variable, format, log))
	}

	time.Sleep(200 * time.Millisecond)
	return fmt.Errorf("ok")
}

// Handle the case of wierd ints encoded as byte arrays.
func GetFromConv(pdu gosnmp.SnmpPDU, conv string, log logger.ContextL) (int64, string) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Invalid Conversion: %s %v %v", pdu.Name, pdu.Value, r)
		}
	}()

	bv, ok := pdu.Value.([]byte)
	if !ok || len(bv) == 0 {
		return 0, ""
	}

	switch conv {
	case CONV_HWADDR: // If there's an encoded mac addr here.
		return 0, net.HardwareAddr(bv).String()
	case CONV_POWERSET:
		return vendor.HandlePowersetStatus(bv)
	case CONV_HEXTOIP:
		return hexToIP(bv)
	case CONV_ENGINE_ID:
		return engineID(bv)
	case CONV_ONE:
		return toOne(bv)
	default:
		// Otherwise, try out some custom conversions.
		split := strings.Split(conv, ":")
		if split[0] == CONV_HEXTOINT && len(split) == 3 {
			endian := split[1]
			bit := split[2]

			if endian == "LittleEndian" {
				switch bit {
				case "uint64":
					return 0, fmt.Sprintf("%d", binary.LittleEndian.Uint64(bv))
				case "uint32":
					return 0, fmt.Sprintf("%d", binary.LittleEndian.Uint32(bv))
				case "uint16":
					return 0, fmt.Sprintf("%d", binary.LittleEndian.Uint16(bv))
				default:
					log.Errorf("invalid bit value (%s) for hex to int conversion", bit)
					return 0, ""
				}
			} else if endian == "BigEndian" {
				switch bit {
				case "uint64":
					return 0, fmt.Sprintf("%d", binary.BigEndian.Uint64(bv))
				case "uint32":
					return 0, fmt.Sprintf("%d", binary.BigEndian.Uint32(bv))
				case "uint16":
					return 0, fmt.Sprintf("%d", binary.BigEndian.Uint16(bv))
				default:
					log.Errorf("invalid bit value (%s) for hex to int conversion", bit)
					return 0, ""
				}
			} else {
				log.Errorf("invalid Endian value (%s) for hex to int conversion", endian)
				return 0, ""
			}
		} else if split[0] == CONV_REGEXP && len(split) >= 2 {
			return fromRegexp(bv, strings.Join(split[1:], ":")) // Put back together just in case RE has a : in it.
		}
	}

	return 0, string(bv) // Default down to here.
}

/**
Some OID's don't store IP as a string, they store it as a hex value that we are going to want to translate.
I need to take this:
.1.3.6.1.4.1.9.9.42.1.2.2.1.2.1 = Hex-String: 0A00640A
and display it as a string 10.0.100.10
*/
func hexToIP(bv []byte) (int64, string) {
	switch len(bv) {
	case 8:
		return 0, fmt.Sprintf("%d.%d.%d.%d",
			binary.BigEndian.Uint16(bv[0:2]),
			binary.BigEndian.Uint16(bv[2:4]),
			binary.BigEndian.Uint16(bv[4:6]),
			binary.BigEndian.Uint16(bv[6:8]),
		)
	case 4:
		nv := []byte{0x00, bv[0], 0x00, bv[1], 0x00, bv[2], 0x00, bv[3]}
		return 0, fmt.Sprintf("%d.%d.%d.%d",
			binary.BigEndian.Uint16(nv[0:2]),
			binary.BigEndian.Uint16(nv[2:4]),
			binary.BigEndian.Uint16(nv[4:6]),
			binary.BigEndian.Uint16(nv[6:8]),
		)
	default:
		return 0, ""
	}
}

func engineID(bv []byte) (int64, string) {
	buf := make([]byte, 0, 3*len(bv))
	x := buf[1*len(bv) : 3*len(bv)]
	hex.Encode(x, bv)
	for i := 0; i < len(x); i += 2 {
		buf = append(buf, x[i], x[i+1], ':')
	}
	return 0, string(buf[:len(buf)-1])
}

/**
Ubiquity and maybe others can be annoying in returning a string version of CPU.
This lets people parse it out.

	agentSwitchCpuProcessTotalUtilization
	1.3.6.1.4.1.4413.1.1.1.1.4.9.0 = STRING: "    5 Secs ( 96.3762%)   60 Secs ( 62.8549%)  300 Secs ( 25.2877%)"
*/
func fromRegexp(bv []byte, reg string) (int64, string) {
	r := reCache[reg]
	if r == nil {
		rn, err := regexp.Compile(reg)
		if err != nil {
			return 0, ""
		}
		reCache[reg] = rn
		r = rn
	}

	// Now, lets run some regexps.
	res := r.FindSubmatch(bv)
	if len(res) < 2 {
		return 0, ""
	}
	ival, err := strconv.Atoi(string(res[1]))
	if err != nil { // If we can't parse as int, return as string.
		return 0, string(res[1])
	}
	return int64(ival), string(res[1]) // Parsed as int but return both just in case.
}

// This one is used for certain string valued oids which we want to poll as metrics. Just converts to 1.
func toOne(bv []byte) (int64, string) {
	return 1, string(bv)
}
