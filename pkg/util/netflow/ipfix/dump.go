package ipfix

import (
	"encoding/hex"
	"fmt"
)

func Dump(m *Message) {
	fmt.Println("IPFIX message")
	for _, ds := range m.DataSets {
		fmt.Println("  data set")
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
						fmt.Printf("        %d.%d: %v\n", f.Translated.EnterpriseNumber, f.Translated.InformationElementID, f.Bytes)
					}
				} else {
					fmt.Printf("        %v\n", f.Bytes)
				}
			}
		}
	}
}
