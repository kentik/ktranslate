package enrich

import (
	"errors"
	"fmt"
	"strings"

	"go.starlark.net/starlark"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type MibMetric struct {
	logger.ContextL
	idx    string
	key    string
	strs   map[string]string
	ints   map[string]int64
	frozen bool
}

// Wrap updates the starlark.Metric to wrap a new telegraf.Metric.
func (m *MibMetric) Wrap(idx string, key string, ints map[string]int64, strs map[string]string) {
	m.idx = idx
	m.key = key
	m.strs = strs
	m.ints = ints
	m.frozen = false
}

// String returns the starlark representation of the Metric.
//
// The String function is called by both the repr() and str() functions, and so
// it behaves more like the repr function would in Python.
func (m *MibMetric) String() string {
	buf := new(strings.Builder)
	buf.WriteString("MibMetric(")
	buf.WriteString(fmt.Sprintf("Index: %s, ", m.idx))
	buf.WriteString(fmt.Sprintf("Key: %s, ", m.key))
	buf.WriteString(fmt.Sprintf("Ints: %v", m.ints))
	buf.WriteString(fmt.Sprintf("Strings: %v", m.strs))
	buf.WriteString(")")
	return buf.String()
}

func (m *MibMetric) Type() string {
	return "MibMetric"
}

func (m *MibMetric) Freeze() {
	m.frozen = true
}

func (m *MibMetric) Truth() starlark.Bool {
	return true
}

func (m *MibMetric) Hash() (uint32, error) {
	return 0, errors.New("not hashable")
}

// AttrNames implements the starlark.HasAttrs interface.
func (m *MibMetric) AttrNames() []string {
	names := []string{"this.idx", "this.key"}
	for k, _ := range m.strs {
		names = append(names, k)
	}
	for k, _ := range m.ints {
		names = append(names, k)
	}
	return names
}

// Attr implements the starlark.HasAttrs interface.
func (m *MibMetric) Attr(name string) (starlark.Value, error) {
	switch name {
	case "this.idx":
		return starlark.String(m.idx), nil
	case "this.key":
		return starlark.String(m.key), nil
	default:
		if v, ok := m.ints[name]; ok {
			return starlark.MakeInt64(v), nil
		}
		if v, ok := m.strs[name]; ok {
			return starlark.String(v), nil
		}
	}

	// By default, return empty string for all other variables.
	return starlark.String(""), nil
}

// SetField implements the starlark.HasSetField interface.
func (m *MibMetric) SetField(name string, value starlark.Value) error {
	if m.frozen {
		return fmt.Errorf("cannot modify frozen metric")
	}

	switch name {
	case "this.idx":
		return starlark.NoSuchAttrError(
			fmt.Sprintf("cannot assign to field '%s'", name))
	case "this.key":
		return starlark.NoSuchAttrError(
			fmt.Sprintf("cannot assign to field '%s'", name))
	default:
		// Copy over the info about this new key from the key which started this process.
		switch v := value.(type) {
		case starlark.String:
			m.strs[name] = v.GoString()
		case starlark.Int:
			ns, ok := v.Int64()
			if ok {
				m.ints[name] = ns
			}
		}
	}

	return nil
}

// Get implements the starlark.Mapping interface.
func (m *MibMetric) Get(key starlark.Value) (v starlark.Value, found bool, err error) {
	if k, ok := key.(starlark.String); ok {
		v, err := m.Attr(k.GoString())
		if err != nil {
			return starlark.None, false, err
		}
		return v, true, nil
	}

	return starlark.None, false, errors.New("key must be of type 'str'")
}

// SetKey implements the starlark.HasSetKey interface to support map update
// using x[k]=v syntax, like a dictionary.
func (m *MibMetric) SetKey(k, v starlark.Value) error {
	if m.frozen {
		return fmt.Errorf("cannot modify frozen metric")
	}

	key, ok := k.(starlark.String)
	if !ok {
		return errors.New("field key must be of type 'str'")
	}

	return m.SetField(key.GoString(), v)
}
