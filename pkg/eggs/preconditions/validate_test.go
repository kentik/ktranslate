package preconditions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateBasics(t *testing.T) {
	type testStruct struct {
		someField1 int `validate:"unknown_validator"`
	}

	ValidateStruct(new(int))      // non structs are ignored
	ValidateStruct(new(struct{})) // structs without validate tags pass validation
}

func TestValidateNilChecks(t *testing.T) {

	type nilChecksStruct struct {
		mustNotBeNil *int `validate:"not_nil"`
		mustBeNil    *int `validate:"nil"`
	}

	goodObj := nilChecksStruct{
		mustBeNil:    nil,
		mustNotBeNil: new(int),
	}
	ValidateStruct(goodObj)
	ValidateStruct(&goodObj)

	badObj := nilChecksStruct{
		mustBeNil:    nil,
		mustNotBeNil: nil,
	}
	assert.Panics(t, func() { ValidateStruct(badObj) })
	assert.Panics(t, func() { ValidateStruct(&badObj) })
}

func TestValidateNoNilPointersOption(t *testing.T) {
	type nilChecksStruct struct {
		mustNotBeNil *int
	}

	goodObj := nilChecksStruct{
		mustNotBeNil: new(int),
	}
	ValidateStruct(&goodObj, NoNilPointers)

	badObj := nilChecksStruct{
		mustNotBeNil: nil,
	}
	assert.Panics(t, func() { ValidateStruct(badObj, NoNilPointers) })
	assert.Panics(t, func() { ValidateStruct(&badObj, NoNilPointers) })
}

func TestValidateNoNilMapsOption(t *testing.T) {
	type nilChecksStruct struct{ mustNotBeNil map[string]int }

	goodObj := nilChecksStruct{mustNotBeNil: make(map[string]int)}
	ValidateStruct(&goodObj, NoNilMaps)

	badObj := nilChecksStruct{mustNotBeNil: nil}
	assert.Panics(t, func() { ValidateStruct(badObj, NoNilMaps) })
	assert.Panics(t, func() { ValidateStruct(&badObj, NoNilMaps) })
}

func TestValidateNotEmptyStringsOption(t *testing.T) {

	type nilChecksStruct struct {
		mustNotBeEmpty string
	}

	goodObj := nilChecksStruct{
		mustNotBeEmpty: "ohai",
	}
	ValidateStruct(&goodObj, NoEmptyStrings)

	badObj := nilChecksStruct{
		mustNotBeEmpty: "",
	}
	assert.Panics(t, func() { ValidateStruct(badObj, NoEmptyStrings) })
	assert.Panics(t, func() { ValidateStruct(&badObj, NoEmptyStrings) })
}
