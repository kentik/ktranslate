package filter

import (
	"flag"
)

func init() {
	flag.Var(&filters, "filters", "Any filters to use. Format: type dimension operator value")
}
