package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type IntFilter struct {
	logger.ContextL
	cf        func(map[string]interface{}) bool
	dimension []string
	value     int64
}

func newIntFilter(log logger.Underlying, fd FilterDef) (*IntFilter, error) {
	sf := &IntFilter{
		ContextL:  logger.NewContextLFromUnderlying(logger.SContext{S: "intFilter"}, log),
		dimension: strings.Split(fd.Dimension, "."),
	}

	val, err := strconv.Atoi(fd.Value)
	if err != nil {
		return nil, fmt.Errorf("Invalid int value: %s", fd.Value)
	}
	sf.value = int64(val)

	switch fd.Operator {
	case Equal:
		sf.cf = sf.intEquals
	case NotEqual:
		sf.cf = sf.intNotEquals
	case LessThan:
		sf.cf = sf.intLessThan
	case GreaterThan:
		sf.cf = sf.intGreaterThan
	default:
		return nil, fmt.Errorf("Invalid operator for int: %s", fd.Operator)
	}

	return sf, nil
}

func (f *IntFilter) Filter(in *kt.JCHF) bool {
	mapr := in.ToMap()
	if !f.cf(mapr) {
		return false
	}
	return true
}

func (f *IntFilter) intEquals(chf map[string]interface{}) bool {
	if dd, ok := chf[f.dimension[0]]; ok {
		switch dim := dd.(type) {
		case int64:
			return dim == f.value
		case map[string]int32:
			return dim[f.dimension[1]] == int32(f.value)
		case map[string]int64:
			return dim[f.dimension[1]] == f.value
		}
	} else if dd, ok := chf["custom_bigint"]; ok { // Fall back and try all ints here.
		switch dim := dd.(type) {
		case map[string]int64:
			if _, ok := dim[f.dimension[0]]; ok {
				return dim[f.dimension[0]] == f.value
			}
		}
	}
	return false
}

func (f *IntFilter) intNotEquals(chf map[string]interface{}) bool {
	return !f.intEquals(chf)
}

func (f *IntFilter) intLessThan(chf map[string]interface{}) bool {
	if dd, ok := chf[f.dimension[0]]; ok {
		switch dim := dd.(type) {
		case int64:
			return dim < f.value
		case map[string]int32:
			return dim[f.dimension[1]] < int32(f.value)
		case map[string]int64:
			return dim[f.dimension[1]] < f.value
		}
	} else if dd, ok := chf["custom_bigint"]; ok { // Fall back and try all ints here.
		switch dim := dd.(type) {
		case map[string]int64:
			if _, ok := dim[f.dimension[0]]; ok {
				return dim[f.dimension[0]] < f.value
			}
		}
	}
	return false
}

func (f *IntFilter) intGreaterThan(chf map[string]interface{}) bool {
	if dd, ok := chf[f.dimension[0]]; ok {
		switch dim := dd.(type) {
		case int64:
			return dim > f.value
		case map[string]int32:
			return dim[f.dimension[1]] > int32(f.value)
		case map[string]int64:
			return dim[f.dimension[1]] > f.value
		}
	} else if dd, ok := chf["custom_bigint"]; ok { // Fall back and try all ints here.
		switch dim := dd.(type) {
		case map[string]int64:
			if _, ok := dim[f.dimension[0]]; ok {
				return dim[f.dimension[0]] > f.value
			}
		}
	}
	return false
}
