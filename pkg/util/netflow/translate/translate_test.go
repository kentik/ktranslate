package translate

import (
	"math"
	"reflect"
	"testing"
)

// If nicer test failure output like line numbers is desired, one can stub in
// the commented assertMatch()/assertNotMatch() below and import
// github.com/stretchr/testify/assert
//func assertMatch(t *testing.T, ft FieldType, buf []byte, expectedValue interface{}) {
//	actualValue := Bytes(buf, ft)
//	assert.Equal(t, expectedValue, actualValue)
//}
//func assertNotMatch(t *testing.T, ft FieldType, buf []byte, expectedValue interface{}) {
//	actualValue := Bytes(buf, ft)
//	assert.NotEqual(t, expectedValue, actualValue)
//}

func assertMatch(t *testing.T, ft FieldType, buf []byte, expectedValue interface{}) {
	actualValue := Bytes(buf, ft)
	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Fatalf("expectedValue was %v with type %T, received %v with %T instead", expectedValue, expectedValue, actualValue, actualValue)
	}
}

func assertNotMatch(t *testing.T, ft FieldType, buf []byte, expectedValue interface{}) {
	actualValue := Bytes(buf, ft)
	if reflect.DeepEqual(expectedValue, actualValue) {
		t.Fatalf("expectedValue should not match actualValue: %v with type %T", expectedValue, expectedValue)
	}
}

//////////////////////////////////////////////////////////////////////////////
// Full-length field tests
//////////////////////////////////////////////////////////////////////////////

func TestFieldTypeUnsigned8(t *testing.T) {
	buf := []byte{0xff}
	assertMatch(t, Uint8, buf, uint8(0xff))
}

func TestFieldTypeUnsigned16(t *testing.T) {
	buf := []byte{0xff, 0}
	assertMatch(t, Uint16, buf, uint16(0xff00))
}

func TestFieldTypeUnsigned32(t *testing.T) {
	buf := []byte{0xff, 0, 0, 0}
	assertMatch(t, Uint32, buf, uint32(0xff000000))
}

func TestFieldTypeUnsigned64(t *testing.T) {
	buf := []byte{0xff, 0, 0, 0, 0, 0, 0, 0}
	assertMatch(t, Uint64, buf, uint64(0xff00000000000000))
}

func TestFieldTypeSigned8(t *testing.T) {
	buf := []byte{0xff}
	assertMatch(t, Int8, buf, int8(-1))
}

func TestFieldTypeSigned16(t *testing.T) {
	buf := []byte{0xff, 0xff}
	assertMatch(t, Int16, buf, int16(-1))
}

func TestFieldTypeSigned32(t *testing.T) {
	buf := []byte{0xff, 0xff, 0xff, 0xff}
	assertMatch(t, Int32, buf, int32(-1))
}

func TestFieldTypeSigned64(t *testing.T) {
	buf := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assertMatch(t, Int64, buf, int64(-1))
}

func TestFieldTypeFloat64(t *testing.T) {
	buf := []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}
	assertMatch(t, Float64, buf, math.Pi)
}

func TestFieldTypeFloat32(t *testing.T) {
	buf := []byte{0x3f, 0x99, 0x99, 0x9a}
	assertMatch(t, Float32, buf, float32(1.2))
}

//////////////////////////////////////////////////////////////////////////////
// Reduced-size encoding tests
//////////////////////////////////////////////////////////////////////////////

func TestZeroByteUnsigned16(t *testing.T) {
	buf := []byte{}
	assertMatch(t, Uint16, buf, buf)
}

func TestZeroByteSigned16(t *testing.T) {
	buf := []byte{}
	assertMatch(t, Uint16, buf, buf)
}

func TestOneByteUnsigned16(t *testing.T) {
	buf := []byte{0xff}
	assertMatch(t, Uint16, buf, uint16(0xff))
}

func TestOneByteUnsigned32(t *testing.T) {
	buf := []byte{0xff}
	assertMatch(t, Uint32, buf, uint32(0xff))
}

func TestOneByteUnsigned64(t *testing.T) {
	buf := []byte{0xff}
	assertMatch(t, Uint64, buf, uint64(0xff))
}

func TestOneByteSigned16(t *testing.T) {
	buf := []byte{0x01}
	assertMatch(t, Int16, buf, int16(1))

	buf = []byte{0x80}
	assertMatch(t, Int16, buf, int16(math.MinInt8))
}

func TestOneByteSigned32(t *testing.T) {
	buf := []byte{0x01}
	assertMatch(t, Int32, buf, int32(1))

	buf = []byte{0x80}
	assertMatch(t, Int32, buf, int32(math.MinInt8))
}

func TestOneByteSigned64(t *testing.T) {
	buf := []byte{0x01}
	assertMatch(t, Int64, buf, int64(1))

	buf = []byte{0x80}
	assertMatch(t, Int64, buf, int64(math.MinInt8))
}

func TestTwoByteSigned16(t *testing.T) {
	buf := []byte{0x00, 0x01}
	assertMatch(t, Int16, buf, int16(1))

	buf = []byte{0x80, 0x00}
	assertMatch(t, Int16, buf, int16(math.MinInt16))
}

func TestTwoByteSigned32(t *testing.T) {
	buf := []byte{0x00, 0x01}
	assertMatch(t, Int32, buf, int32(1))

	buf = []byte{0x80, 0x00}
	assertMatch(t, Int32, buf, int32(math.MinInt16))
}

func TestTwoByteSigned64(t *testing.T) {
	buf := []byte{0x00, 0x01}
	assertMatch(t, Int64, buf, int64(1))

	buf = []byte{0x80, 0x00}
	assertMatch(t, Int64, buf, int64(math.MinInt16))
}

func TestThreeByteUnsigned32(t *testing.T) {
	buf := []byte{1, 2, 3}
	assertMatch(t, Uint32, buf, uint32(0x010203))
}

func TestFourByteSigned32(t *testing.T) {
	buf := []byte{0, 0, 0, 1}
	assertMatch(t, Int32, buf, int32(1))

	buf = []byte{0x80, 0, 0, 0}
	assertMatch(t, Int32, buf, int32(math.MinInt32))
}

func TestFourByteSigned64(t *testing.T) {
	buf := []byte{0, 0, 0, 1}
	assertMatch(t, Int64, buf, int64(1))

	buf = []byte{0x80, 0, 0, 0}
	assertMatch(t, Int64, buf, int64(math.MinInt32))
}

func TestFourByteFloat64(t *testing.T) {
	buf := []byte{0x41, 0x46, 0x00, 0x00}
	// 12.375 can be represented in float32 without loss of precision
	assertMatch(t, Float64, buf, float64(12.375))
}

//////////////////////////////////////////////////////////////////////////////
// Reduced-size overflow tests
//////////////////////////////////////////////////////////////////////////////

func TestTwoByteSigned8(t *testing.T) {
	buf := []byte{0, 1}
	assertNotMatch(t, Int8, buf, int8(1))
}

func TestTwoByteUnsigned8(t *testing.T) {
	buf := []byte{0, 1}
	assertNotMatch(t, Uint8, buf, uint8(1))
}

func TestFourByteSigned16(t *testing.T) {
	buf := []byte{0, 0, 0, 1}
	assertNotMatch(t, Int16, buf, int16(1))
}

func TestFourByteUnsigned16(t *testing.T) {
	buf := []byte{0, 0, 0, 1}
	assertNotMatch(t, Int16, buf, uint16(1))
}

func TestEightByteSigned32(t *testing.T) {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 1}
	assertNotMatch(t, Int32, buf, int32(1))
}

func TestEightByteUnsigned32(t *testing.T) {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 1}
	assertNotMatch(t, Uint32, buf, uint32(1))
}

//////////////////////////////////////////////////////////////////////////////
// Direct tests of reducedSizeRead()
//////////////////////////////////////////////////////////////////////////////

func TestReducedSizeUnknownType(t *testing.T) {
	err := reducedSizeRead([]byte{0}, []byte{})
	if err == nil {
		t.Fatal("Expected unknown type to fail")
	}
}

func TestReducedSizeReadSigned(t *testing.T) {
	buf := make([]byte, 9)
	_, err := reducedSizeReadSigned(buf, 128)
	if err == nil {
		t.Fatal("Expected reducedSizeReadSigned() to fail with large byte slice")
	}
}

func TestReducedSizeReadUnsigned(t *testing.T) {
	buf := make([]byte, 9)
	_, err := reducedSizeReadUnsigned(buf, 128)
	if err == nil {
		t.Fatal("Expected reducedSizeReadUnsigned() to fail with large byte slice")
	}
}
