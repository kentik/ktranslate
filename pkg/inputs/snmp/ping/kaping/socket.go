package kaping

import (
	"net"
	"net/netip"

	"golang.org/x/net/icmp"
)

const (
	ICMP = 0
	RAW  = 1
)

type Socket interface {
	Send(addr netip.Addr, echo icmp.Echo) error
	Recv(buffer []byte) (*icmp.Message, netip.Addr, error)
	Close()
}

func NewSocket(addr netip.Addr, mode int) (Socket, error) {
	if addr.Is4() {
		return socket4(addr, mode)
	} else {
		return socket6(addr, mode)
	}
}

func socket4(addr netip.Addr, mode int) (*Socket4, error) {
	var network = "udp4"

	if mode == RAW {
		network = "ip4:icmp"
	}

	conn, err := icmp.ListenPacket(network, addr.String())
	if err != nil {
		return nil, err
	}

	return &Socket4{conn, mode}, nil
}

func socket6(addr netip.Addr, mode int) (*Socket6, error) {
	var network = "udp6"

	if mode == RAW {
		network = "ip6:ipv6-icmp"
	}

	conn, err := icmp.ListenPacket(network, addr.String())
	if err != nil {
		return nil, err
	}

	return &Socket6{conn, mode}, nil
}

func address(addr netip.Addr, mode int) net.Addr {
	switch mode {
	case ICMP:
		return &net.UDPAddr{
			IP:   net.IP(addr.AsSlice()),
			Zone: addr.Zone(),
		}
	default:
		return &net.IPAddr{
			IP:   net.IP(addr.AsSlice()),
			Zone: addr.Zone(),
		}
	}
}
