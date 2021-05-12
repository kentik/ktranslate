package translate

//go:generate go run cmd/translate-rfc5102/main.go -output rfc5102.go

import (
	"encoding/binary"
	"errors"
	"math"
	"net"
	"time"

	"github.com/kentik/ktranslate/pkg/util/netflow/session"
)

// Builtin dictionary of information elements
var builtin = make(informationElements)

// Translate knows how to translate the raw bytes from a DataRecord into their actual values.
type Translate struct {
	session.Session
	elements informationElements
}

// NewTranslate creates a new session bound translator.
func NewTranslate(s session.Session) *Translate {
	return &Translate{s, builtin}
}

// Key retrieves the Information Element entry for the given Key.
func (t *Translate) Key(k Key) (InformationElementEntry, bool) {
	i, ok := t.elements[k]
	return i, ok
}

// FieldType is the IPFIX type of an Information Element ("Field").
type FieldType uint8

// The available field types as defined by RFC 5102.
const (
	Unknown FieldType = iota
	Uint8
	Uint16
	Uint32
	Uint64
	Int8
	Int16
	Int32
	Int64
	Float32
	Float64
	Boolean
	MacAddress
	OctetArray
	String
	DateTimeSeconds
	DateTimeMilliseconds
	DateTimeMicroseconds
	DateTimeNanoseconds
	Ipv4Address
	Ipv6Address
)

// FieldTypes are used in the InformationElementEntries map
var FieldTypes = map[string]FieldType{
	"unsigned8":            Uint8,
	"unsigned16":           Uint16,
	"unsigned32":           Uint32,
	"unsigned64":           Uint64,
	"signed8":              Int8,
	"signed16":             Int16,
	"signed32":             Int32,
	"signed64":             Int64,
	"float32":              Float32,
	"float64":              Float64,
	"boolean":              Boolean,
	"macAddress":           MacAddress,
	"octetArray":           OctetArray,
	"string":               String,
	"dateTimeSeconds":      DateTimeSeconds,
	"dateTimeMilliseconds": DateTimeMilliseconds,
	"dateTimeMicroseconds": DateTimeMicroseconds,
	"dateTimeNanoseconds":  DateTimeNanoseconds,
	"ipv4Address":          Ipv4Address,
	"ipv6Address":          Ipv6Address,
}

// minLength is the minimum length of a field of the given type, in bytes.
func (t FieldType) minLength() int {
	switch t {
	case Uint8, Int8, Boolean:
		return 1
	case Uint16, Int16:
		return 1
	case Uint32, Int32:
		return 1
	case Float32, DateTimeSeconds:
		return 4
	case Uint64, Int64:
		return 1
	case Float64:
		return 4 // Float64 can be encoded in 4 bytes though loss of precision may occur
	case DateTimeMilliseconds, DateTimeMicroseconds, DateTimeNanoseconds:
		return 8
	case MacAddress:
		return 6
	case Ipv4Address:
		return 4
	case Ipv6Address:
		return 16
	default:
		return 0
	}
}

// UnmarshalText converts byte slice to FieldType
func (f *FieldType) UnmarshalText(bs []byte) error {
	*f = FieldTypes[string(bs)]
	return nil
}

// InformationElementEntry is an entry in the Information Element map.
type InformationElementEntry struct {
	Name         string
	FieldID      uint16
	EnterpriseID uint32
	Type         FieldType
}

// Key is the key of the Information Element map.
type Key struct {
	EnterpriseID uint32
	FieldID      uint16
}

type informationElements map[Key]InformationElementEntry

var reducedSizeErr error = errors.New("Unable to read reduced size encoding: size not implemented")
var tooManyBitsErr error = errors.New("Unable to read reduced size encoding: too many bits")

// Helper method to read an unsigned reduced size field
func reducedSizeReadUnsigned(bs []byte, maxBits int) (uint64, error) {
	// Exit if `bs` has more bits than we can store
	if len(bs)*8 > maxBits {
		return 0, tooManyBitsErr
	}

	switch len(bs) {
	case 1:
		return uint64(bs[0]), nil
	case 2:
		return uint64(binary.BigEndian.Uint16(bs)), nil
	case 3:
		return uint64(uint32(bs[0])<<16 + uint32(bs[1])<<8 + uint32(bs[2])), nil
	case 4:
		return uint64(binary.BigEndian.Uint32(bs)), nil
	case 8:
		return binary.BigEndian.Uint64(bs), nil
	}
	return 0, reducedSizeErr
}

// Helper method to read a signed reduced size field
func reducedSizeReadSigned(bs []byte, maxBits int) (int64, error) {
	// Exit if `bs` has more bits than we can store
	if len(bs)*8 > maxBits {
		return 0, tooManyBitsErr
	}

	switch len(bs) {
	case 1:
		value := int8(bs[0])
		return int64(value), nil
	case 2:
		value := int16(binary.BigEndian.Uint16(bs))
		return int64(value), nil
	case 4:
		value := int32(binary.BigEndian.Uint32(bs))
		return int64(value), nil
	case 8:
		return int64(binary.BigEndian.Uint64(bs)), nil
	}
	return 0, reducedSizeErr
}

// Read a reduced size field into its full size
func reducedSizeRead(bs []byte, i interface{}) error {
	var unsigned uint64
	var signed int64
	var err error

	switch v := i.(type) {
	case *uint16:
		unsigned, err = reducedSizeReadUnsigned(bs, 16)
		*v = uint16(unsigned)
	case *uint32:
		unsigned, err = reducedSizeReadUnsigned(bs, 32)
		*v = uint32(unsigned)
	case *uint64:
		unsigned, err = reducedSizeReadUnsigned(bs, 64)
		*v = uint64(unsigned)
	case *int16:
		signed, err = reducedSizeReadSigned(bs, 16)
		*v = int16(signed)
	case *int32:
		signed, err = reducedSizeReadSigned(bs, 32)
		*v = int32(signed)
	case *int64:
		signed, err = reducedSizeReadSigned(bs, 64)
		*v = int64(signed)
	default:
		err = reducedSizeErr
	}

	return err
}

// Bytes translates a byte string to a go native type.
func Bytes(bs []byte, t FieldType) interface{} {
	if len(bs) < t.minLength() {
		// Field is too short (corrupt) - return it uninterpreted.
		return bs
	}

	switch t {
	case Uint8:
		return bs[0]
	case Uint16:
		var i uint16
		if err := reducedSizeRead(bs, &i); err == nil {
			return i
		}
	case Uint32:
		var i uint32
		if err := reducedSizeRead(bs, &i); err == nil {
			return i
		}
	case Uint64:
		var i uint64
		if err := reducedSizeRead(bs, &i); err == nil {
			return i
		}
	case Int8:
		return int8(bs[0])
	case Int16:
		var i int16
		if err := reducedSizeRead(bs, &i); err == nil {
			return i
		}
	case Int32:
		var i int32
		if err := reducedSizeRead(bs, &i); err == nil {
			return i
		}
	case Int64:
		var i int64
		if err := reducedSizeRead(bs, &i); err == nil {
			return i
		}
	case Float32:
		return math.Float32frombits(binary.BigEndian.Uint32(bs))
	case Float64:
		switch len(bs) {
		case 4:
			return float64(math.Float32frombits(binary.BigEndian.Uint32(bs)))
		case 8:
			return math.Float64frombits(binary.BigEndian.Uint64(bs))
		}
	case Boolean:
		return bs[0] == 1
	case Unknown, OctetArray:
		return bs
	case String:
		return string(bs)
	case MacAddress:
		return net.HardwareAddr(bs)
	case Ipv4Address, Ipv6Address:
		return net.IP(bs)
	case DateTimeSeconds:
		return time.Unix(int64(binary.BigEndian.Uint32(bs)), 0)
	case DateTimeMilliseconds:
		unixTimeMs := int64(binary.BigEndian.Uint64(bs))
		return time.Unix(0, 0).Add(time.Duration(unixTimeMs) * time.Millisecond)
	case DateTimeMicroseconds:
		unixTimeUs := int64(binary.BigEndian.Uint64(bs))
		return time.Unix(0, 0).Add(time.Duration(unixTimeUs) * time.Microsecond)
	case DateTimeNanoseconds:
		unixTimeNs := int64(binary.BigEndian.Uint64(bs))
		return time.Unix(0, 0).Add(time.Duration(unixTimeNs))
	}
	return bs
}
