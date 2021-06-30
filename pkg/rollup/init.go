package rollup

import (
	"flag"
	"fmt"

	gohll "github.com/sasha-s/go-hll"
)

func init() {
	flag.Var(&rollups, "rollups", "Any rollups to use. Format: type, name, metric, dimension 1, dimension 2, ..., dimension n: sum,bytes,in_bytes,dst_addr")
	verifyHLLConstantsOrPanic()
}

func verifyHLLConstantsOrPanic() {
	// This is a constant calculation, so run at package init to double check
	// our constants.
	s, err := gohll.SizeByError(hllErrRateHigh)
	if err != nil {
		panic(err)
	}
	if s != hllSizeForErrRateHigh {
		panic(fmt.Sprintf("hllSizeForErrRateHigh=%d, but gohll.SizeByError(%v)=%d",
			hllSizeForErrRateHigh, hllErrRateHigh, s))
	}
}
