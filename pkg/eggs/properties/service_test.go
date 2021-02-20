package properties

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint: errcheck
func TestPropertyServiceWithMapBacking(t *testing.T) {

	os.Setenv("float.property.from.env", "1.99")
	os.Setenv("float.property", "1.99") // should be shadowed by map value

	propertyMap := map[string]string{
		"bool.property":   "true",
		"float.property":  "0.17",
		"int.property":    "12987",
		"string.property": "cool",
	}

	props := NewPropertyService(
		NewStaticMapPropertyBacking(propertyMap),
		NewEnvPropertyBacking(),
	)

	// string
	assert.Equal(t, props.GetString("string.property", "meh").String(), "cool")
	assert.Equal(t, props.GetString("nosuch", "meh").String(), "meh")
	assert.True(t, props.GetString("nosuch", "meh").FromFallback())

	// bool
	assert.True(t, props.GetBool("bool.property", false).Bool())
	assert.False(t, props.GetBool("nosuch", false).Bool())
	assert.True(t, props.GetBool("nosuch", false).FromFallback())

	// uint64
	assert.Equal(t, props.GetUInt64("int.property", 4242).Uint64(), uint64(12987))
	assert.Equal(t, props.GetUInt64("nosuch", 4242).Uint64(), uint64(4242))
	assert.True(t, props.GetUInt64("nosuch", 4242).FromFallback())

	// float64
	assert.Equal(t, props.GetFloat64("float.property.from.env", 0.42).Float64(), 1.99)
	assert.Equal(t, props.GetFloat64("float.property", 0.42).Float64(), 0.17)
	assert.Equal(t, props.GetFloat64("nosuch", 0.42).Float64(), 0.42)
	assert.True(t, props.GetFloat64("nosuch", 0.0).FromFallback())

}
