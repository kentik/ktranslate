package enrich

import (
	"errors"
	"fmt"
	"strings"

	"go.starlark.net/starlark"

	"github.com/kentik/ktranslate/pkg/kt"
)

type JCHF struct {
	metric *kt.JCHF
	flat   map[string]interface{}
	frozen bool
}

// Wrap updates the starlark.Metric to wrap a new telegraf.Metric.
func (m *JCHF) Wrap(metric *kt.JCHF) {
	m.metric = metric
	m.flat = metric.Flatten()
	m.frozen = false
}

// Unwrap removes the telegraf.Metric from the startlark.Metric.
func (m *JCHF) Unwrap() *kt.JCHF {
	return m.metric
}

// String returns the starlark representation of the Metric.
//
// The String function is called by both the repr() and str() functions, and so
// it behaves more like the repr function would in Python.
func (m *JCHF) String() string {
	buf := new(strings.Builder)
	buf.WriteString("JCHF(")
	buf.WriteString(fmt.Sprintf("%v", m.flat))
	buf.WriteString(")")
	return buf.String()
}

func (m *JCHF) Type() string {
	return "JCHF"
}

func (m *JCHF) Freeze() {
	m.frozen = true
}

func (m *JCHF) Truth() starlark.Bool {
	return true
}

func (m *JCHF) Hash() (uint32, error) {
	return 0, errors.New("not hashable")
}

// AttrNames implements the starlark.HasAttrs interface.
func (m *JCHF) AttrNames() []string {
	names := []string{"pop"}
	for k, _ := range m.flat {
		names = append(names, k)
	}
	return names
}

// Attr implements the starlark.HasAttrs interface.
func (m *JCHF) Attr(name string) (starlark.Value, error) {
	if name == "pop" {
		return builtinAttr(m, "pop", m.pop)
	}

	if v, ok := m.flat[name]; ok {
		switch nv := v.(type) {
		case kt.Cid:
			return starlark.MakeUint64(uint64(nv)), nil
		case uint64:
			return starlark.MakeUint64(nv), nil
		case int64:
			return starlark.MakeInt64(nv), nil
		case int:
			return starlark.MakeInt(int(nv)), nil
		case int32:
			return starlark.MakeInt(int(nv)), nil
		case string:
			return starlark.String(nv), nil
		}
	}

	// Returning nil, nil indicates "no such field or method"
	return nil, nil
}

// SetField implements the starlark.HasSetField interface.
func (m *JCHF) SetField(name string, value starlark.Value) error {
	if m.frozen {
		return fmt.Errorf("cannot modify frozen metric")
	}

	switch name {
	case "company_id":
		m.metric.CompanyId = kt.Cid(setUint64(value))
	case "timestamp":
		m.metric.Timestamp = setInt64(value)
	case "dst_as":
		m.metric.DstAs = setUint32(value)
	case "dst_geo":
		m.metric.DstGeo = setString(value)
	default:
		switch v := value.(type) {
		case starlark.String:
			m.metric.CustomStr[name] = v.GoString()
		case starlark.Int:
			ns, ok := v.Int64()
			if ok {
				m.metric.CustomBigInt[name] = ns
			}
		}
	}

	return nil
}

// Get implements the starlark.Mapping interface.
func (m *JCHF) Get(key starlark.Value) (v starlark.Value, found bool, err error) {
	if k, ok := key.(starlark.String); ok {
		v, err := m.Attr(k.GoString())
		if err != nil {
			return starlark.None, false, err
		}
		return v, true, nil
	}

	return starlark.None, false, errors.New("key must be of type 'str'")
}

// Delete removes the key and also returns it.
func (m *JCHF) Delete(key starlark.Value) (v starlark.Value, found bool, err error) {
	if k, ok := key.(starlark.String); ok {
		v, err := m.Attr(k.GoString())
		if err != nil {
			return starlark.None, false, err
		}

		// Actually remove the key here.
		delete(m.metric.CustomBigInt, k.GoString())
		delete(m.metric.CustomStr, k.GoString())
		delete(m.metric.CustomInt, k.GoString())

		// And return
		return v, true, nil
	}

	return starlark.None, false, errors.New("key must be of type 'str'")
}

// SetKey implements the starlark.HasSetKey interface to support map update
// using x[k]=v syntax, like a dictionary.
func (m *JCHF) SetKey(k, v starlark.Value) error {
	if m.frozen {
		return fmt.Errorf("cannot modify frozen metric")
	}

	key, ok := k.(starlark.String)
	if !ok {
		return errors.New("field key must be of type 'str'")
	}

	return m.SetField(key.GoString(), v)
}

// Implements the pop method
func (m *JCHF) pop(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var k, d starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &k, &d); err != nil {
		return starlark.None, fmt.Errorf("%s: %v", b.Name(), err)
	}

	if v, found, err := m.Delete(k); err != nil {
		return starlark.None, fmt.Errorf("%s: %v", b.Name(), err)
	} else if found {
		return v, nil
	} else if d != nil {
		return d, nil
	}

	return starlark.None, fmt.Errorf("%s: missing key", b.Name())
}

func setUint64(value starlark.Value) uint64 {
	switch v := value.(type) {
	case starlark.Int:
		ns, ok := v.Uint64()
		if ok {
			return ns
		}
	}
	return 0
}

func setInt64(value starlark.Value) int64 {
	switch v := value.(type) {
	case starlark.Int:
		ns, ok := v.Int64()
		if ok {
			return ns
		}
	}
	return 0
}

func setUint32(value starlark.Value) uint32 {
	switch v := value.(type) {
	case starlark.Int:
		ns, ok := v.Uint64()
		if ok {
			return uint32(ns)
		}
	}
	return 0
}

func setString(value starlark.Value) string {
	switch v := value.(type) {
	case starlark.String:
		return v.GoString()
	}
	return ""
}
