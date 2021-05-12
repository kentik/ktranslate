package netflow1

import (
	"fmt"

	"github.com/kentik/ktranslate/pkg/util/netflow/read"
)

func Dump(p *Packet) {
	fmt.Println("NetFlow version 1 packet", p.Header)
	fmt.Printf("  %d flow records:\n", len(p.Records))
	for i, r := range p.Records {
		fmt.Printf("  record %d:\n", i)
		fmt.Println("    srcAddr: ", r.SrcAddr)
		fmt.Println("    srcPort: ", r.SrcPort)
		fmt.Println("    dstAddr: ", r.DstAddr)
		fmt.Println("    dstPort: ", r.DstPort)
		fmt.Println("    nextHop: ", r.NextHop)
		fmt.Println("    bytes:   ", r.Bytes)
		fmt.Println("    packets: ", r.Packets)
		fmt.Println("    first:   ", r.First)
		fmt.Println("    last:    ", r.Last)
		fmt.Println("    protocol:", r.Protocol, read.Protocol(r.Protocol))
		fmt.Println("    tos:     ", r.ToS)
		fmt.Println("    flags:   ", r.Flags, read.TCPFlags(r.Flags))
	}
}
