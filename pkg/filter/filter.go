package filter

import (
	"flag"
	"fmt"
	"net"
	"strconv"
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

	OrToken  = " or "
	AndToken = " and "
)

var (
	filters FilterFlag
)

type FilterFlag []string

func (ff *FilterFlag) String() string {
	return strings.Join(*ff, AndToken)
}
func (ff *FilterFlag) Set(val string) error {
	*ff = append(*ff, strings.TrimSpace(val))
	return nil
}

func init() {
	flag.Var(&filters, "filters", "Any filters to use. Format: type dimension operator value")
}

type Operator string

type FilterType string

type Filter interface {
	Filter(*kt.JCHF) bool
}

type FilterWrapper []Filter

type FilterDef struct {
	Dimension string
	Operator  Operator
	Value     string
	FType     FilterType
}

type FilterDefWrapper []FilterDef

func (f *FilterDef) String() string {
	return fmt.Sprintf("%s Filter: %s %s %s", f.FType, f.Dimension, f.Operator, f.Value)
}

func (f FilterDefWrapper) String() string {
	set := make([]string, len(f))
	for i, f := range f {
		set[i] = f.String()
	}
	return strings.Join(set, OrToken)
}

type FilterDefs []FilterDefWrapper

func (ff *FilterDefs) String() string {
	pts := make([]string, len(*ff))
	for i, r := range *ff {
		pts[i] = r.String()
	}
	return strings.Join(pts, "\n")
}

func (i *FilterDefs) Set(value string) error {
	inner := FilterDefWrapper{}
	for _, orSet := range strings.Split(value, OrToken) {
		pts := strings.Split(orSet, ",")
		if len(pts) < 3 {
			return fmt.Errorf("Filter flag is defined by <type> dimension operator value")
		}
		if len(pts) == 3 {
			ftype := String
			// Try to guess the type.
			if _, err := strconv.Atoi(pts[2]); err == nil {
				ftype = Int
			} else if _, _, err := net.ParseCIDR(pts[2]); err == nil {
				ftype = Addr
			}
			pts = append([]string{string(ftype)}, pts...)
		}
		ptn := make([]string, len(pts))
		for i, p := range pts {
			ptn[i] = strings.TrimSpace(p)
		}
		inner = append(inner, FilterDef{
			FType:     FilterType(ptn[0]),
			Dimension: ptn[1],
			Operator:  Operator(ptn[2]),
			Value:     ptn[3],
		})
	}
	*i = append(*i, inner)

	return nil
}

func GetFilters(log logger.Underlying, filters []string) ([]FilterWrapper, error) {
	ff := FilterDefs{}
	for _, f := range filters {
		for _, andSet := range strings.Split(f, AndToken) {
			if err := ff.Set(andSet); err != nil {
				return nil, err
			}
		}
	}
	filterSet := make([]FilterWrapper, 0)
	for _, fSet := range ff {
		orSet := []Filter{}
		for _, fd := range fSet {
			switch fd.FType {
			case String:
				newf, err := newStringFilter(log, fd)
				if err != nil {
					return nil, err
				}
				orSet = append(orSet, newf)
			case Int:
				newf, err := newIntFilter(log, fd)
				if err != nil {
					return nil, err
				}
				orSet = append(orSet, newf)
			case Addr:
				newf, err := newAddrFilter(log, fd)
				if err != nil {
					return nil, err
				}
				orSet = append(orSet, newf)
			default:
				return nil, fmt.Errorf("Invalid type: %s. Valid Types: %s|%s|%s", fd.FType, String, Int, Addr)
			}
		}
		log.Infof("Filter: ", "Added %s", fSet.String())
		filterSet = append(filterSet, orSet)
	}

	return filterSet, nil
}

// This provides and OR wrapper, returning true if any of the filters in this wrapper eval to true.
func (fs FilterWrapper) Filter(chf *kt.JCHF) bool {
	for _, f := range fs {
		if f.Filter(chf) {
			return true
		}
	}
	return false
}
