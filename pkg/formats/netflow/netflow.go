package netflow

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"
	"github.com/kentik/ktranslate/pkg/util/ic"
	"github.com/kentik/ktranslate/pkg/util/trkdsess"

	"github.com/kentik/ktranslate/pkg/util/netflow"
	"github.com/kentik/ktranslate/pkg/util/netflow/ipfix"
	"github.com/kentik/ktranslate/pkg/util/netflow/netflow9"
	"github.com/kentik/ktranslate/pkg/util/netflow/session"
)

type NetflowFormat struct {
	logger.ContextL
	s               session.Session
	d               *netflow.Decoder
	templateTracker templateTracker
	version         uint16
}

const (
	IN_BYTES      uint16 = 1
	IN_PKTS              = 2
	PROTOCOL             = 4
	TOS                  = 5
	TCP_FLAGS            = 6
	L4_SRC_PORT          = 7
	IPV4_SRC_ADDR        = 8
	INPUT_SNMP           = 10
	L4_DST_PORT          = 11
	IPV4_DST_ADDR        = 12
	OUTPUT_SNMP          = 14
	SRC_AS               = 16
	DST_AS               = 17
	OUT_BYTES            = 23
	OUT_PKTS             = 24
	IPV6_SRC_ADDR        = 27
	IPV6_DST_ADDR        = 28
	SRC_MAC              = 56
	DST_MAC              = 80
)

type Field struct {
	Type   uint16
	Length uint16
}

type EntField struct {
	Field
	EntId uint32
}

var fields = []Field{
	{PROTOCOL, 1},
	{TOS, 1},
	{TCP_FLAGS, 2},
	{SRC_MAC, 6},
	{DST_MAC, 6},
	{IPV4_SRC_ADDR, 4},
	{IPV4_DST_ADDR, 4},
	{L4_SRC_PORT, 2},
	{L4_DST_PORT, 2},
	{IPV6_SRC_ADDR, 16},
	{IPV6_DST_ADDR, 16},
	{SRC_AS, 4},
	{DST_AS, 4},
	{INPUT_SNMP, 4},
	{OUTPUT_SNMP, 4},
	{IN_BYTES, 8},
	{OUT_BYTES, 8},
	{IN_PKTS, 8},
	{OUT_PKTS, 8},
}

const templateSet = uint16(2)
const optionalTemplateSet = uint16(2)
const dataSet = uint16(256)
const templateId = uint16(256)
const templateEntId = uint16(257)

var (
	Version = flag.String("netflow_version", "ipfix", "Version of netflow to produce: (netflow9|ipfix)")
)

func NewFormat(log logger.Underlying, comp kt.Compression) (*NetflowFormat, error) {
	s := trkdsess.New()
	d := netflow.NewDecoder(s)

	ipf := &NetflowFormat{
		ContextL:        logger.NewContextLFromUnderlying(logger.SContext{S: "ipfixFormat"}, log),
		s:               s,
		d:               d,
		templateTracker: newTemplateTracker(),
	}

	switch *Version {
	case "ipfix":
		ipf.version = ipfix.Version
	case "netflow9":
		ipf.version = netflow9.Version
	default:
		return nil, fmt.Errorf("You used an unsupported netflow version: %s. Use netflow9 or ipfix.", *Version)
	}

	if comp != kt.CompressionNone {
		return nil, fmt.Errorf("You cannot use compression on IPFix data.")
	}

	ipf.Infof("Netflow formatter running with version %d", ipf.version)

	return ipf, nil
}

func (f *NetflowFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	switch f.version {
	case netflow9.Version:
		return f.pack9(msgs, serBuf)
	case ipfix.Version:
		return f.packIpfix(msgs, serBuf)
	default:
		return nil, fmt.Errorf("You used an unsupported version.")
	}
}

func (f *NetflowFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	m, err := f.d.Read(bytes.NewBuffer(raw.Body))
	if err != nil {
		return nil, err
	}

	var chfs []*kt.JCHF
	switch p := m.(type) {
	case *netflow9.Packet:
		chfs = f.decodeV9(p, raw.Body)

	case *ipfix.Message:
		chfs = f.decodeIpfix(p, raw.Body)

	default:
		return nil, fmt.Errorf("ipfix parser failed to decode a message: %v: %v", err, raw)
	}

	values := make([]map[string]interface{}, len(chfs))
	for i, m := range chfs {
		values[i] = m.ToMap()
	}
	return values, nil
}

// Not implemented for this format.
func (f *NetflowFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

func (ipf *NetflowFormat) setCHFField(flow *kt.JCHF, id uint16, value interface{}, raw []byte, enterpriseNumber uint32) bool {
	// Lots of boilerplate in this switch expression, but I don't see a better way to
	// get things done without taking a performance hit from reflection.  We're right
	// in the client critical path here, so we'll try the ugly code.
	switch id {

	case PROTOCOL:
		protocol, ok := value.(uint8)
		if !ok {
			return false
		}
		flow.Protocol = ic.PROTO_NAMES[uint16(protocol)]

	case TOS:
		tos, ok := value.(uint8)
		if !ok {
			return false
		}
		flow.Tos = uint32(tos)

	case TCP_FLAGS:
		flags, ok := value.(uint16)
		if !ok {
			return false
		}
		flow.TcpFlags = uint32(flags)

	case L4_SRC_PORT:
		srcPort, ok := value.(uint16)
		if !ok {
			return false
		}
		flow.L4SrcPort = uint32(srcPort)

	case L4_DST_PORT:
		dstPort, ok := value.(uint16)
		if !ok {
			return false
		}
		flow.L4DstPort = uint32(dstPort)

	case IPV4_SRC_ADDR:
		srcV4Addr, ok := value.(net.IP)
		if !ok {
			return false
		}
		if str := srcV4Addr.String(); str != "0.0.0.0" {
			flow.SrcAddr = str
		}

	case IPV4_DST_ADDR:
		dstV4Addr, ok := value.(net.IP)
		if !ok {
			return false
		}
		if str := dstV4Addr.String(); str != "0.0.0.0" {
			flow.DstAddr = str
		}

	case IPV6_SRC_ADDR:
		srcV6Addr, ok := value.(net.IP)
		if !ok {
			return false
		}
		if str := srcV6Addr.String(); str != "::" {
			flow.SrcAddr = str
		}

	case IPV6_DST_ADDR:
		dstV6Addr, ok := value.(net.IP)
		if !ok {
			return false
		}
		if str := dstV6Addr.String(); str != "::" {
			flow.DstAddr = str
		}

	case INPUT_SNMP:
		inputPort, ok := value.(uint32)
		if !ok {
			return false
		}
		flow.InputPort = kt.IfaceID(inputPort)

	case OUTPUT_SNMP:
		outputPort, ok := value.(uint32)
		if !ok {
			return false
		}
		flow.OutputPort = kt.IfaceID(outputPort)

	case SRC_AS:
		srcAS, ok := value.(uint32)
		if !ok {
			return false
		}
		flow.SrcAs = srcAS

	case DST_AS:
		dstAS, ok := value.(uint32)
		if !ok {
			return false
		}
		flow.DstAs = dstAS

	case IN_PKTS:
		packets, ok := value.(uint64)
		if !ok {
			return false
		}
		flow.InPkts = packets

	case IN_BYTES:
		bytes, ok := value.(uint64)
		if !ok {
			return false
		}
		flow.InBytes = bytes

	case OUT_PKTS:
		packets, ok := value.(uint64)
		if !ok {
			return false
		}
		flow.OutPkts = packets

	case OUT_BYTES:
		bytes, ok := value.(uint64)
		if !ok {
			return false
		}
		flow.OutBytes = bytes

	case SRC_MAC:
		srcMac, ok := value.(net.HardwareAddr)
		if !ok {
			return false
		}
		flow.SrcEthMac = srcMac.String()

	case DST_MAC:
		dstMac, ok := value.(net.HardwareAddr)
		if !ok {
			return false
		}
		flow.DstEthMac = dstMac.String()
	}
	return true
}

func (ipf *NetflowFormat) writeWithHeader(buf *bytes.Buffer, id uint16, body *bytes.Buffer) error {
	var err error
	switch ipf.version {
	case ipfix.Version:
		head := &ipfix.SetHeader{
			ID:     id,
			Length: uint16(body.Len() + 4),
		}
		err = write(buf, head, body.Bytes())
	case netflow9.Version:
		head := &netflow9.FlowSetHeader{
			ID:     id,
			Length: uint16(body.Len() + 4),
		}
		err = write(buf, head, body.Bytes())
	}
	return err
}

func encodeTemplate(flow *kt.JCHF, tempId uint16) (*bytes.Buffer, error) {
	fieldCount := uint16(len(fields))

	data := []interface{}{tempId, fieldCount}
	for i := range fields {
		data = append(data, fields[i])
	}

	buf := &bytes.Buffer{}
	err := write(buf, data...)

	return buf, err
}

func encodeFlow(flow *kt.JCHF) (*bytes.Buffer, error) {
	empty := [16]byte{}

	src4, src6 := ip2int(net.ParseIP(flow.SrcAddr))
	dst4, dst6 := ip2int(net.ParseIP(flow.DstAddr))

	if src6 == nil {
		src6 = empty[:16]
	}

	if dst6 == nil {
		dst6 = empty[:16]
	}

	smac, _ := net.ParseMAC(flow.SrcEthMac)
	if smac == nil {
		smac = empty[:6]
	}

	dmac, _ := net.ParseMAC(flow.DstEthMac)
	if dmac == nil {
		dmac = empty[:6]
	}

	buf := &bytes.Buffer{}
	for _, field := range fields {
		var value interface{}

		switch field.Type {
		case PROTOCOL:
			value = uint8(ic.PROTO_NUMS[flow.Protocol])
		case TOS:
			value = uint8(flow.Tos)
		case TCP_FLAGS:
			value = uint16(flow.TcpFlags)
		case SRC_MAC:
			value = smac[:6]
		case DST_MAC:
			value = dmac[:6]
		case IPV4_SRC_ADDR:
			value = src4
		case IPV4_DST_ADDR:
			value = dst4
		case L4_SRC_PORT:
			value = uint16(flow.L4SrcPort)
		case L4_DST_PORT:
			value = uint16(flow.L4DstPort)
		case IPV6_SRC_ADDR:
			value = src6
		case IPV6_DST_ADDR:
			value = dst6
		case SRC_AS:
			value = uint32(flow.SrcAs)
		case DST_AS:
			value = uint32(flow.DstAs)
		case INPUT_SNMP:
			value = uint32(flow.InputPort)
		case OUTPUT_SNMP:
			value = uint32(flow.OutputPort)
		case IN_BYTES:
			value = uint64(flow.InBytes)
		case OUT_BYTES:
			value = uint64(flow.OutBytes)
		case IN_PKTS:
			value = uint64(flow.InPkts)
		case OUT_PKTS:
			value = uint64(flow.OutPkts)
		}

		err := write(buf, value)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func write(buf *bytes.Buffer, values ...interface{}) error {
	for _, v := range values {
		switch v.(type) {
		case string:
			v = []byte(v.(string))
		}
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func ip2int(ip net.IP) (uint32, []byte) {
	if ip == nil {
		return 0, nil
	}
	if ip.To4() != nil {
		return binary.BigEndian.Uint32(ip.To4()), nil
	}
	return 0, ip.To16()
}
