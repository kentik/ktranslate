package kaping

import (
	"fmt"
	"net"
	"net/netip"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Socket4 struct {
	conn *icmp.PacketConn
	mode int
}

func (s *Socket4) Send(addr netip.Addr, echo icmp.Echo) error {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &echo,
	}

	bytes, err := msg.Marshal(nil)
	if err != nil {
		return fmt.Errorf("packet marshal: %w", err)
	}

	_, err = s.conn.WriteTo(bytes, address(addr, s.mode))
	if err != nil {
		return fmt.Errorf("packet transmit: %w", err)
	}

	return nil
}

func (s *Socket4) Recv(buffer []byte) (*icmp.Message, netip.Addr, error) {
	n, addr, err := s.conn.ReadFrom(buffer)
	if err != nil {
		return nil, netip.Addr{}, fmt.Errorf("packet receive: %w", err)
	}

	msg, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), buffer[:n])
	if err != nil {
		return nil, netip.Addr{}, fmt.Errorf("ICMP message: %w", err)
	}

	var peer netip.Addr
	switch addr := addr.(type) {
	case *net.IPAddr:
		peer, _ = netip.AddrFromSlice(addr.IP)
	case *net.UDPAddr:
		peer, _ = netip.AddrFromSlice(addr.IP)
	default:
		return nil, netip.Addr{}, fmt.Errorf("unreachable peer: %w", err)
	}

	return msg, peer, nil
}

func (s *Socket4) Close() {
	s.conn.Close()
}
