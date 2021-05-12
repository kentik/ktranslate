package netflow9

import (
	"encoding/hex"
	"fmt"
)

func Dump(p *Packet) {
	fmt.Println("NetFlow version 9 packet")
	for _, ds := range p.DataFlowSets {
		fmt.Printf("  data set template %d, length: %d\n", ds.Header.ID, ds.Header.Length)
		if ds.Records == nil {
			fmt.Printf("    %d raw bytes:\n", len(ds.Bytes))
			fmt.Println(hex.Dump(ds.Bytes))
			continue
		}
		fmt.Printf("    %d records:\n", len(ds.Records))
		for i, dr := range ds.Records {
			fmt.Printf("      record %d:\n", i)
			for _, f := range dr.Fields {
				if f.Translated != nil {
					if f.Translated.Name != "" {
						fmt.Printf("        %s: %v\n", f.Translated.Name, f.Translated.Value)
					} else {
						fmt.Printf("        %d: %v\n", f.Translated.Type, f.Bytes)
					}
				} else {
					fmt.Printf("        %d: %v (raw)\n", f.Type, f.Bytes)
				}
			}
		}
	}
}
