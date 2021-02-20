package properties

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type stdPropertyService struct {
	backing []PropertyBacking
}

/* Create a new property service using the provided "backing" providers, in decreasing order of priority.
 * For example, passing (filesystem, env, map) will give priority to filesystem, then env (for keys that are
 * missing in the filesystem store) then finally map.
 * nil values are acceptable and ignored.
 */
func NewPropertyService(backing ...PropertyBacking) *stdPropertyService {

	// filter out nils
	filteredBacking := make([]PropertyBacking, 0)
	for _, b := range backing {
		if b != nil && !reflect.ValueOf(backing).IsNil() {
			filteredBacking = append(filteredBacking, b)
		}
	}

	return &stdPropertyService{
		backing: filteredBacking,
	}
}

func (propsvc *stdPropertyService) GetProperty(name string) (string, bool) {
	for _, backing := range propsvc.backing {
		if value, present := backing.GetProperty(name); present {
			if present {
				return value, true
			}
		}
	}
	return "", false
}

func (propsvc *stdPropertyService) GetPropertySub(propertyName string, sub ...interface{}) (string, bool) {
	if !IsValidPropertySubName(propertyName) {
		panic(fmt.Errorf("Invalid property sub name '%s'", propertyName))
	}

	if len(sub)%2 != 0 {
		panic(fmt.Errorf("GetPropertySub: sub len must be even"))
	}

	strValue := ""
	gotValue := false

	strValue, gotValue = propsvc.GetProperty(fmt.Sprintf("%s-global", propertyName))

	path := fmt.Sprintf("%s", propertyName)
	for i := 0; i < len(sub); i += 2 {
		name := sub[i]
		value := sub[i+1]
		path = fmt.Sprintf("%s-%v_%v", path, name, value)

		propValue, pOK := propsvc.GetProperty(path)
		if pOK {
			strValue = propValue
			gotValue = true
		}
	}
	return strValue, gotValue
}

func (propsvc *stdPropertyService) GetString(name string, fallback string) PropertyValue {
	ret := &propertyValue{
		stringValue:  fallback,
		fromFallback: true,
	}
	if stringValue, ok := propsvc.GetProperty(name); ok {
		ret.stringValue = stringValue
		ret.fromFallback = false
	}
	return ret
}

func (propsvc *stdPropertyService) GetBool(name string, fallback bool) PropertyValue {
	ret := &propertyValue{
		boolValue:    fallback,
		fromFallback: true,
	}
	if stringValue, ok := propsvc.GetProperty(name); ok {
		if v, err := strconv.ParseBool(strings.TrimSpace(stringValue)); err == nil {
			ret.boolValue = v
			ret.fromFallback = false
		}
	}
	return ret
}

func (propsvc *stdPropertyService) GetBoolSub(name string, fallback bool, sub ...interface{}) PropertyValue {
	ret := &propertyValue{
		boolValue:    fallback,
		fromFallback: true,
	}
	if stringValue, ok := propsvc.GetPropertySub(name, sub...); ok {
		if v, err := strconv.ParseBool(strings.TrimSpace(stringValue)); err == nil {
			ret.boolValue = v
			ret.fromFallback = false
		}
	}
	return ret
}

func (propsvc *stdPropertyService) GetUInt64(name string, fallback uint64) PropertyValue {
	ret := &propertyValue{
		uint64Value:  fallback,
		fromFallback: true,
	}
	if stringValue, ok := propsvc.GetProperty(name); ok {
		if v, err := strconv.ParseUint(strings.TrimSpace(stringValue), 10, 64); err == nil {
			ret.uint64Value = v
			ret.fromFallback = false
		}
	}
	return ret
}

func (propsvc *stdPropertyService) GetUInt64Sub(name string, fallback uint64, sub ...interface{}) PropertyValue {
	ret := &propertyValue{
		uint64Value:  fallback,
		fromFallback: true,
	}
	if stringValue, ok := propsvc.GetPropertySub(name, sub...); ok {
		if v, err := strconv.ParseUint(strings.TrimSpace(stringValue), 10, 64); err == nil {
			ret.uint64Value = v
			ret.fromFallback = false
		}
	}
	return ret
}

func (propsvc *stdPropertyService) GetFloat64(name string, fallback float64) PropertyValue {
	ret := &propertyValue{
		float64Value: fallback,
		fromFallback: true,
	}
	if stringValue, ok := propsvc.GetProperty(name); ok {
		if v, err := strconv.ParseFloat(strings.TrimSpace(stringValue), 64); err == nil {
			ret.float64Value = v
			ret.fromFallback = false
		}
	}
	return ret
}

func (propsvc *stdPropertyService) Refresh() {
	for _, backing := range propsvc.backing {
		backing.Refresh()
	}
}
