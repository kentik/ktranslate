package rollup

import (
	"fmt"

	gohll "github.com/sasha-s/go-hll"
)

func init() {
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
