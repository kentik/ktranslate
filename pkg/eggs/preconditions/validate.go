package preconditions

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	validateTagName      = "validate"
	validateTagSeparator = ","
)

type validatorFunc func(reflect.StructField, reflect.Value)

var validators = map[string]validatorFunc{
	"nil": func(field reflect.StructField, value reflect.Value) {
		if !value.IsNil() {
			panic(fmt.Errorf("The %s field must be null.", field.Name))
		}
	},
	"not_nil": func(field reflect.StructField, value reflect.Value) {
		if value.IsNil() {
			panic(fmt.Errorf("The %s field must not be null.", field.Name))
		}
	},
	"not_empty": func(field reflect.StructField, value reflect.Value) {
		if value.Len() == 0 {
			panic(fmt.Errorf("You need to fill in the %s field.", field.Name))
		}
	},
}

type validationConfig struct {
	byKind map[reflect.Kind][]validatorFunc
}

func newConfig(options ...validationOption) (config *validationConfig) {
	config = &validationConfig{
		byKind: make(map[reflect.Kind][]validatorFunc, 0),
	}
	for _, opt := range options {
		opt(config)
	}
	return
}

func (config *validationConfig) addValidatorForKind(kind reflect.Kind, validator validatorFunc) {
	config.byKind[kind] = append(config.byKind[kind], validator)
}

type validationOption func(*validationConfig)

func NoNilPointers(config *validationConfig) {
	config.addValidatorForKind(reflect.Ptr, validators["not_nil"])
}

func NoNilMaps(config *validationConfig) {
	config.addValidatorForKind(reflect.Map, validators["not_nil"])
}

func NoEmptyStrings(config *validationConfig) {
	config.addValidatorForKind(reflect.String, validators["not_empty"])
}

func ValidateStruct(obj interface{}, options ...validationOption) {

	config := newConfig(options...)

	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
		objType = objValue.Type()
	}

	if objType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)
		for _, validatorName := range strings.Split(field.Tag.Get(validateTagName), validateTagSeparator) {
			if validatorFunc, ok := validators[validatorName]; ok {
				validatorFunc(field, fieldValue)
			}
		}
		for _, validatorFunc := range config.byKind[fieldValue.Kind()] {
			validatorFunc(field, fieldValue)
		}
	}
}
