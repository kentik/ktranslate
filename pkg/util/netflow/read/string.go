package read

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// TCP Flags
const (
	tcpFIN = 1 << iota
	tcpSYN
	tcpRST
	tcpPSH
	tcpACK
	tcpURG
	tcpECE
	tcpCWR
	tcpNS

	tcpFlags = "NCEUAPRSF"
)

var protocol = map[uint8]string{}

func init() {
	if f, err := os.Open("/etc/protocols"); err == nil {
		defer f.Close()

		s := bufio.NewScanner(f)
		for s.Scan() {
			line := s.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			if n, err := strconv.Atoi(fields[1]); err == nil {
				protocol[uint8(n)] = fields[0]
			}
		}
	}
}

// Protocol returns the protocol name
func Protocol(p uint8) string {
	return protocol[p]
}

// TCPFlags returns the TCP flags
func TCPFlags(f uint8) string {
	flags := []byte{}
	for i := uint8(0); i < 8; i++ {
		if f&0x01 > 0 {
			flags = append(flags, tcpFlags[8-i])
		} else {
			flags = append(flags, '.')
		}
		f >>= 1
	}
	return "[" + string(flags) + "]"
}
