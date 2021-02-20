package filter

import (
	"fmt"
	"strings"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type StringFilter struct {
	logger.ContextL
	cf        func(map[string]interface{}) bool
	dimension []string
	value     string
}

func newStringFilter(log logger.Underlying, fd FilterDef) (*StringFilter, error) {
	sf := &StringFilter{
		ContextL:  logger.NewContextLFromUnderlying(logger.SContext{S: "stringFilter"}, log),
		dimension: strings.Split(fd.Dimension, "."),
		value:     fd.Value,
	}

	switch fd.Operator {
	case Equal:
		sf.cf = sf.stringEquals
	case NotEqual:
		sf.cf = sf.stringNotEquals
	case Contains:
		sf.cf = sf.stringContains
	default:
		return nil, fmt.Errorf("Invalid operator for string: %s", fd.Operator)
	}

	return sf, nil
}

func (f *StringFilter) Filter(in *kt.JCHF) bool {
	mapr := in.ToMap()
	if !f.cf(mapr) {
		return false
	}
	return true
}

func (f *StringFilter) stringEquals(chf map[string]interface{}) bool {
	if dd, ok := chf[f.dimension[0]]; ok {
		switch dim := dd.(type) {
		case string:
			return dim == f.value
		case map[string]string:
			return dim[f.dimension[1]] == f.value
		}
	}
	return false
}

func (f *StringFilter) stringNotEquals(chf map[string]interface{}) bool {
	return !f.stringEquals(chf)
}

func (f *StringFilter) stringContains(chf map[string]interface{}) bool {
	if dd, ok := chf[f.dimension[0]]; ok {
		switch dim := dd.(type) {
		case string:
			return strings.Contains(dim, f.value)
		case map[string]string:
			return strings.Contains(dim[f.dimension[1]], f.value)
		}
	}
	return false
}
