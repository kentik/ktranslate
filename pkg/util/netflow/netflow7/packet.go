package netflow7

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/kentik/ktranslate/pkg/util/netflow/read"
)

const (
	// Version word in the Packet Header
	Version uint16 = 0x0007
)

// Packet is a NetFlow v7 packet
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
	Version      uint16
	Count        uint16
	SysUptime    time.Duration // 32 bit milliseconds
	Unix         time.Time     // 32 bit seconds + 32 bit nanoseconds
	FlowSequence uint32
	Reserved     uint32
}

func (h PacketHeader) String() string {
	return fmt.Sprintf("v=%d, count=%d, uptime=%s, time=%s, seq=%d",
		h.Version, h.Count, time.Duration(h.SysUptime)*time.Second, h.Unix, h.FlowSequence)
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
	h.SysUptime = time.Duration(u) * time.Millisecond
	var t uint64
	if err := read.Uint64(&t, r); err != nil {
		return err
	}
	h.Unix = time.Unix(int64(t>>32), int64(t&0xffffffff))
	if err := read.Uint32(&h.FlowSequence, r); err != nil {
		return err
	}
	if err := read.Uint32(&h.Reserved, r); err != nil {
		return err
	}
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
	Pad1 uint8 // 36
	// TCP Flags
	TCPFlags uint8 // 37
	// Protocol number
	Protocol uint8 // 38
	// ToS is the IP type of service
	ToS uint8 // 39
	// SrcAS is the source Autonomous System Number
	SrcAS uint16 // 40-41
	// DstAS is the source Autonomous System Number
	DstAS uint16 // 42-43
	// SrcMask is the source network mask
	SrcMask uint8 // 44
	// DstMask is the destination network mask
	DstMask uint8 // 45
	// Flags indicating, among other things, what flows are invalid
	Flags uint16 // 46-47
	// RouterSC is the IP address of the router that is bypassed by the Catalyst 5000 series switch
	RouterSC net.IP // 48-51
}

func (r FlowRecord) String() string {
	return fmt.Sprintf("%s:%d -> %s:%d", r.SrcAddr, r.SrcPort, r.DstAddr, r.DstPort)
}

func (r *FlowRecord) Unmarshal(h io.Reader) error {
	r.SrcAddr = make(net.IP, 4)
	if _, err := h.Read(r.SrcAddr); err != nil { // 0-3
		return err
	}
	r.DstAddr = make(net.IP, 4)
	if _, err := h.Read(r.DstAddr); err != nil { // 4-7
		return err
	}
	r.NextHop = make(net.IP, 4)
	if _, err := h.Read(r.NextHop); err != nil { // 8-11
		return err
	}
	if err := read.Uint16(&r.Input, h); err != nil { // 12-13
		return err
	}
	if err := read.Uint16(&r.Output, h); err != nil { // 14-15
		return err
	}
	if err := read.Uint32(&r.Packets, h); err != nil { // 16-19
		return err
	}
	if err := read.Uint32(&r.Bytes, h); err != nil { // 20-23
		return err
	}
	if err := read.Uint32(&r.First, h); err != nil { // 24-27
		return err
	}
	if err := read.Uint32(&r.Last, h); err != nil { // 28-31
		return err
	}
	if err := read.Uint16(&r.SrcPort, h); err != nil { // 32-33
		return err
	}
	if err := read.Uint16(&r.DstPort, h); err != nil { // 34-35
		return err
	}
	if err := read.Uint8(&r.Pad1, h); err != nil { // 36
		return err
	}
	if err := read.Uint8(&r.TCPFlags, h); err != nil { // 37
		return err
	}
	if err := read.Uint8(&r.Protocol, h); err != nil { // 38
		return err
	}
	if err := read.Uint8(&r.ToS, h); err != nil { // 39
		return err
	}
	if err := read.Uint16(&r.SrcAS, h); err != nil { // 40-41
		return err
	}
	if err := read.Uint16(&r.DstAS, h); err != nil { // 42-43
		return err
	}
	if err := read.Uint8(&r.SrcMask, h); err != nil { // 44
		return err
	}
	if err := read.Uint8(&r.DstMask, h); err != nil { // 45
		return err
	}
	if err := read.Uint16(&r.Flags, h); err != nil { // 46-47
		return err
	}
	r.RouterSC = make(net.IP, 4)
	if _, err := h.Read(r.RouterSC); err != nil { // 48-51
		return err
	}

	return nil
}

func (f FlowRecord) SampleInterval() int {
	return 1
}
