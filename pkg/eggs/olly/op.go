package olly

import (
	"fmt"
	"strings"
)

type ollyOp struct {
	string
}

// Op() and the associated type are only here to help find Op keys in the code
func Op(s string) ollyOp {
	if !strings.Contains(s, ".") {
		panic(fmt.Errorf("olly: Op must use some.dotted.notation"))
	}
	return ollyOp{
		string: s,
	}
}

