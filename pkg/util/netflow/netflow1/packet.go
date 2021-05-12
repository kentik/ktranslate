package netflow1

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/kentik/ktranslate/pkg/util/netflow/read"
)

const (
	// Version word in the Packet Header
	Version uint16 = 0x0001
)

// Packet is a NetFlow v1 packet
type Packet struct {
	Header  PacketHeader
	Records []*FlowRecord
}

func (p *Packet) Unmarshal(r io.Reader) error {
	if err := p.Header.Unmarshal(r); err != nil {
		return err
	}
	p.Records = make([]*FlowRecord, p.Header.Count)
	for i := range p.Records {
		p.Records[i] = new(FlowRecord)
		if err := p.Records[i].Unmarshal(r); err != nil {
			return err
		}
	}
	return nil
}

// PacketHeader is a NetFlow v1 packet
type PacketHeader struct {
	Version   uint16
	Count     uint16
	SysUptime time.Duration // 32 bit milliseconds
	Unix      time.Time     // 32 bit seconds + 32 bit nanoseconds
}

func (h PacketHeader) String() string {
	return fmt.Sprintf("v=%d, count=%d, uptime=%s, time=%s",
		h.Version, h.Count, time.Duration(h.SysUptime)*time.Second, h.Unix)
}

func (h *PacketHeader) Unmarshal(r io.Reader) error {
	if err := read.Uint16(&h.Version, r); err != nil {
		return err
	}
	if err := read.Uint16(&h.Count, r); err != nil {
		return err
	}
	// The spec says at most 24 flows in one packet, but reality disagrees.
	if h.Count < 1 || h.Count > 32 {
		return fmt.Errorf("protocol error: %d flows out of bounds", h.Count)
	}
	var u uint32
	if err := read.Uint32(&u, r); err != nil {
		return err
	}
	h.SysUptime = time.Duration(u)
	var t uint64
	if err := read.Uint64(&t, r); err != nil {
		return err
	}
	h.Unix = time.Unix(int64(t>>32), int64(t&0xffffffff))
	log.Println(h)
	return nil
}

// FlowRecord is a NetFlow v1 Flow Record
type FlowRecord struct {
	// SrcAddr is the Source IP address
	SrcAddr net.IP // 0-3
	// DstAddr is the Destination IP address
	DstAddr net.IP // 4-7
	// NextHop is the IP address of next hop router
	NextHop net.IP // 8-11
	// Input is the SNMP index of input interface
	Input uint16 // 12-13
	// Output is the SNMP index of output interface
	Output uint16 // 14-15
	// Packets is the number of packets in the flow
	Packets uint32 // 16-19
	// Octets is the number of bytes in the flow
	Bytes uint32 // 20-23
	// First is the SysUptime at start of flow
	First uint32 // 24-27
	// Last is the SysUptime at end of flow
	Last uint32 // 28-31
	// SrcPort is the TCP/UDP source port number or equivalent
	SrcPort uint16 // 32-33
	// DstPort is the TCP/UDP destination port number or equivalent
	DstPort uint16 // 34-35
	// Pad1 are unused bytes
	Pad1 uint16 // 36-37
	// Protocol number (IP)
	Protocol uint8 // 38
	// ToS is the IP type of service
	ToS uint8 // 39
	// Flags are cumulative OR TCP flags
	Flags uint8 // 40
	// Pad2 are unused bytes
	Pad2 uint8
	// Pad3 are unused bytes
	Pad3 uint16
	// Reservered are reserved (unused) bytes
	Reserved uint32
}

func (r FlowRecord) String() string {
	return fmt.Sprintf("%s:%d -> %s:%d", r.SrcAddr, r.SrcPort, r.DstAddr, r.DstPort)
}

func (r *FlowRecord) Unmarshal(h io.Reader) error {
	r.SrcAddr = make(net.IP, 4)
	if _, err := h.Read(r.SrcAddr); err != nil {
		return err
	}
	r.DstAddr = make(net.IP, 4)
	if _, err := h.Read(r.DstAddr); err != nil {
		return err
	}
	r.NextHop = make(net.IP, 4)
	if _, err := h.Read(r.NextHop); err != nil {
		return err
	}
	if err := read.Uint16(&r.Input, h); err != nil {
		return err
	}
	if err := read.Uint16(&r.Output, h); err != nil {
		return err
	}
	if err := read.Uint32(&r.Packets, h); err != nil {
		return err
	}
	if err := read.Uint32(&r.Bytes, h); err != nil {
		return err
	}
	if err := read.Uint32(&r.First, h); err != nil {
		return err
	}
	if err := read.Uint32(&r.Last, h); err != nil {
		return err
	}
	if err := read.Uint16(&r.SrcPort, h); err != nil {
		return err
	}
	if err := read.Uint16(&r.DstPort, h); err != nil {
		return err
	}
	if err := read.Uint16(&r.Pad1, h); err != nil {
		return err
	}
	if err := read.Uint8(&r.Protocol, h); err != nil {
		return err
	}
	if err := read.Uint8(&r.ToS, h); err != nil {
		return err
	}
	if err := read.Uint8(&r.Flags, h); err != nil {
		return err
	}
	if err := read.Uint8(&r.Pad2, h); err != nil {
		return err
	}
	if err := read.Uint16(&r.Pad3, h); err != nil {
		return err
	}
	if err := read.Uint32(&r.Reserved, h); err != nil {
		return err
	}

	return nil
}

func (f FlowRecord) SampleInterval() int {
	return 1
}
