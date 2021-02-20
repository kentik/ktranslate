package http

import (
	"flag"
)

func init() {
	flag.Var(&headers, "http_header", "Any custom http headers to set on outbound requests")
}
