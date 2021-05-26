package kt

import (
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
