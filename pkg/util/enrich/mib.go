package enrich

import (
	"errors"
	"fmt"
	"strings"

	"go.starlark.net/starlark"

	"github.com/kentik/ktranslate/pkg/kt"
)

type Mib struct {
	attr   map[string]interface{}
	idx    string
	frozen bool
}

// Wrap updates the starlark.Metric to wrap a new telegraf.Metric.
func (m *Mib) Wrap(idx string, attr map[string]interface{}) {
	m.attr = attr
	m.idx = idx
	m.frozen = false
}

// String returns the starlark representation of the Metric.
//
// The String function is called by both the repr() and str() functions, and so
// it behaves more like the repr function would in Python.
func (m *Mib) String() string {
	buf := new(strings.Builder)
	buf.WriteString("Mib(")
	buf.WriteString(fmt.Sprintf("%s", m.idx))
	buf.WriteString(fmt.Sprintf("%v", m.attr))
	buf.WriteString(")")
	return buf.String()
}

func (m *Mib) Type() string {
	return "Mib"
}

func (m *Mib) Freeze() {
	m.frozen = true
}

func (m *Mib) Truth() starlark.Bool {
	return true
}

func (m *Mib) Hash() (uint32, error) {
	return 0, errors.New("not hashable")
}

// AttrNames implements the starlark.HasAttrs interface.
func (m *Mib) AttrNames() []string {
	names := []string{"idx"}
	for k, _ := range m.attr {
		names = append(names, k)
	}
	return names
}

// Attr implements the starlark.HasAttrs interface.
func (m *Mib) Attr(name string) (starlark.Value, error) {
	if name == "idx" {
		return starlark.String(m.idx), nil
	}
	if v, ok := m.attr[name]; ok {
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

	return nil, nil
}

// SetField implements the starlark.HasSetField interface.
func (m *Mib) SetField(name string, value starlark.Value) error {
	if m.frozen {
		return fmt.Errorf("cannot modify frozen metric")
	}

	switch name {
	case "idx":
		m.idx = setString(value)
	default:
		switch v := value.(type) {
		case starlark.String:
			m.attr[name] = v.GoString()
		case starlark.Int:
			ns, ok := v.Int64()
			if ok {
				m.attr[name] = ns
			}
		}
	}

	return nil
}
