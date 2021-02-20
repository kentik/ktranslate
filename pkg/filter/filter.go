package filter

import (
	"fmt"
	"strings"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

const (
	Equal       Operator = "=="
	NotEqual             = "!="
	LessThan             = "<"
	GreaterThan          = ">"
	Contains             = "%"

	String FilterType = "string"
	Int               = "int"
	Addr              = "addr"
)

var (
	filters FilterFlag
)

type Operator string

type FilterType string

type Filter interface {
	Filter(*kt.JCHF) bool
}

type FilterDef struct {
	Dimension string
	Operator  Operator
	Value     string
	FType     FilterType
}

func (f *FilterDef) String() string {
	return fmt.Sprintf("%s Filter: %s %s %s", f.FType, f.Dimension, f.Operator, f.Value)
}

type FilterFlag []FilterDef

func (ff *FilterFlag) String() string {
	pts := make([]string, len(*ff))
	for i, r := range *ff {
		pts[i] = r.String()
	}
	return strings.Join(pts, "\n")
}

func (i *FilterFlag) Set(value string) error {
	pts := strings.Split(value, ",")
	if len(pts) < 3 {
		return fmt.Errorf("Filter flag is defined by type dimension operator value")
	}
	ptn := make([]string, len(pts))
	for i, p := range pts {
		ptn[i] = strings.TrimSpace(p)
	}
	*i = append(*i, FilterDef{
		FType:     FilterType(ptn[0]),
		Dimension: ptn[1],
		Operator:  Operator(ptn[2]),
		Value:     ptn[3],
	})
	return nil
}

func GetFilters(log logger.Underlying) ([]Filter, error) {
	filterSet := make([]Filter, 0)
	for _, fd := range filters {
		switch fd.FType {
		case String:
			newf, err := newStringFilter(log, fd)
			if err != nil {
				return nil, err
			}
			filterSet = append(filterSet, newf)
		case Int:
			newf, err := newIntFilter(log, fd)
			if err != nil {
				return nil, err
			}
			filterSet = append(filterSet, newf)
		case Addr:
			newf, err := newAddrFilter(log, fd)
			if err != nil {
				return nil, err
			}
			filterSet = append(filterSet, newf)
		}
	}

	return filterSet, nil
}
