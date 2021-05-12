// +build gofuzz

package ipfix

import "bytes"

func Fuzz(data []byte) int {
	if _, err := Read(bytes.NewBuffer(data)); err != nil {
		return 0
	}
	return 1
}
