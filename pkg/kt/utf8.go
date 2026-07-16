package kt

import (
	"encoding/hex"
	"unicode/utf8"
)

// SanitizeUTF8 returns s unchanged when it is valid UTF-8. Otherwise it
// returns a hex encoding of the raw bytes so callers can still export the
// value (OTLP/gRPC and JSON reject invalid UTF-8 and will drop whole batches).
// This mirrors the SNMP trap handling of non-UTF-8 octet strings.
func SanitizeUTF8(s string) string {
	if s == "" || utf8.ValidString(s) {
		return s
	}
	return hex.EncodeToString([]byte(s))
}

// SanitizeUTF8WithHint is like SanitizeUTF8 but also reports whether the
// value was rewritten.
func SanitizeUTF8WithHint(s string) (string, bool) {
	if s == "" || utf8.ValidString(s) {
		return s, false
	}
	return hex.EncodeToString([]byte(s)), true
}
