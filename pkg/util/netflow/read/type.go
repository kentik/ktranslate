package read

import "net"

// LongIPv4 is a 32 bit packed IPv4 address.
type LongIPv4 uint32

func (l LongIPv4) String() string {
	return net.IP{
		uint8(l >> 24),
		uint8(l >> 16),
		uint8(l >> 8),
		uint8(l),
	}.String()
}
