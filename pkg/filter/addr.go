package filter

import (
	"fmt"
	"net"
	"strings"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type AddrFilter struct {
	logger.ContextL
	cf        func(map[string]interface{}) bool
	dimension []string
	value     *net.IPNet
}

func newAddrFilter(log logger.Underlying, fd FilterDef) (*AddrFilter, error) {
	sf := &AddrFilter{
		ContextL:  logger.NewContextLFromUnderlying(logger.SContext{S: "addrFilter"}, log),
		dimension: strings.Split(fd.Dimension, "."),
	}

	_, val, err := net.ParseCIDR(fd.Value)
	if err != nil {
		return nil, err
	}
	sf.value = val

	switch fd.Operator {
	case Equal:
		sf.cf = sf.addrEquals
	case NotEqual:
		sf.cf = sf.addrNotEquals
	case Contains:
		sf.cf = sf.addrContains
	default:
		return nil, fmt.Errorf("Invalid operator for addr: %s", fd.Operator)
	}

	return sf, nil
}

func (f *AddrFilter) Filter(in *kt.JCHF) bool {
	mapr := in.ToMap()
	return f.FilterMap(mapr)
}

func (f *AddrFilter) FilterMap(mapr map[string]interface{}) bool {
	if !f.cf(mapr) {
		return false
	}
	return true
}

func (f *AddrFilter) addrEquals(chf map[string]interface{}) bool {
	if dd, ok := chf[f.dimension[0]]; ok {
		switch dim := dd.(type) {
		case string:
			vv := net.ParseIP(dim)
			return f.value.Contains(vv)
		case map[string]string:
			vv := net.ParseIP(dim[f.dimension[1]])
			return f.value.Contains(vv)
		}
	} else if dd, ok := chf["custom_str"]; ok { // Fall back and try all strings here.
		switch dim := dd.(type) {
		case map[string]string:
			if _, ok := dim[f.dimension[0]]; ok {
				vv := net.ParseIP(dim[f.dimension[0]])
				return f.value.Contains(vv)
			}
		}
	}
	return false
}

func (f *AddrFilter) addrNotEquals(chf map[string]interface{}) bool {
	return !f.addrEquals(chf)
}

func (f *AddrFilter) addrContains(chf map[string]interface{}) bool {
	return f.addrEquals(chf)
}
