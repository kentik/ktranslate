// Package read provides convenience read functions to extract values from a reader.
package read

import (
	"encoding/binary"
	"io"
)

// Uint8 reads a single byte
func Uint8(v *uint8, r io.Reader) error {
	var b [1]byte
	_, err := io.ReadFull(r, b[:])
	*v = b[0]
	return err
}

// Uint16 reads an unsigned word
func Uint16(v *uint16, r io.Reader) error {
	var b [2]byte
	_, err := io.ReadFull(r, b[:])
	*v = binary.BigEndian.Uint16(b[:])
	return err
}

// Uint32 reads an unsigned dword
func Uint32(v *uint32, r io.Reader) error {
	var b [4]byte
	_, err := io.ReadFull(r, b[:])
	*v = binary.BigEndian.Uint32(b[:])
	return err
}

// Uint32IPv4 reads an unsigned dword as IP address
func Uint32IPv4(v *LongIPv4, r io.Reader) error {
	var u uint32
	if err := Uint32(&u, r); err != nil {
		return err
	}
	*v = LongIPv4(u)
	return nil
}

// Uint64 reads an unsigned quad word
func Uint64(v *uint64, r io.Reader) error {
	var b [8]byte
	_, err := io.ReadFull(r, b[:])
	*v = binary.BigEndian.Uint64(b[:])
	return err
}

// VariableLength reads a variable length byte stream as per RFC 7011 section 7.
func VariableLength(p []byte, r io.Reader) ([]byte, error) {
	var l0 uint8
	if err := Uint8(&l0, r); err != nil {
		return nil, err
	}

	var l int
	if l0 < 0xff {
		l = int(l0)
	} else {
		var l1 uint16
		if err := Uint16(&l1, r); err != nil {
			return nil, err
		}
		l = int(l1)
	}

	var b = p
	if cap(b) < l {
		// Allocate new slice for p if there it not enough capacity
		b = make([]byte, l)
	} else if len(b) < l {
		// Grow the passed slice using current capacity
		for i := len(b); i < l; i++ {
			b = append(b, 0x00)
		}
	}

	if _, err := r.Read(b); err != nil {
		return b, err
	}

	return b, nil
}
