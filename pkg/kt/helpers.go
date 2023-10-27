package kt

import (
	"encoding/binary"
	"net"
	"os"
	"strconv"
	"strings"
)

func LookupEnvString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if ival, err := strconv.Atoi(val); err == nil {
			return ival
		} else {
			return defaultVal
		}
	}
	return defaultVal
}

func LookupEnvBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		return strings.ToLower(val) == "true"
	}
	return defaultVal
}

func FixupName(name string) string {
	name = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	return name
}

func Int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
