package properties

import "regexp"

type PropertyService interface {
	GetString(name string, fallback string) PropertyValue
	GetBool(name string, fallback bool) PropertyValue
	GetUInt64(name string, fallback uint64) PropertyValue
	GetFloat64(name string, fallback float64) PropertyValue
	Refresh()

	GetBoolSub(name string, fallback bool, sub ...interface{}) PropertyValue
	GetUInt64Sub(name string, fallback uint64, sub ...interface{}) PropertyValue
}

var validPropertyName = regexp.MustCompile(`^[a-zA-Z0-9-_.]+$`)

func IsValidPropertyName(propName string) bool {
	return validPropertyName.MatchString(propName)
}

// '_' and '-' not allowed as we are using them for separators
var validPropertySubName = regexp.MustCompile(`^[a-zA-Z0-9.]+$`)

func IsValidPropertySubName(propName string) bool {
	return IsValidPropertyName(propName) &&
		validPropertySubName.MatchString(propName)
}

type PropertyValue interface {
	FromFallback() bool
	Bool() bool
	Uint64() uint64
	String() string
	Float64() float64
}

type propertyValue struct {
	boolValue    bool
	uint64Value  uint64
	stringValue  string
	float64Value float64
	fromFallback bool
}

func (pv *propertyValue) FromFallback() bool {
	return pv.fromFallback
}

func (pv *propertyValue) Bool() bool {
	return pv.boolValue
}

func (pv *propertyValue) Uint64() uint64 {
	return pv.uint64Value
}

func (pv *propertyValue) String() string {
	return pv.stringValue
}

func (pv *propertyValue) Float64() float64 {
	return pv.float64Value
}
