package enrich

import (
	"errors"
	"fmt"
	"strings"

	"go.starlark.net/starlark"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type Mib struct {
	logger.ContextL
	idx    string
	key    string
	attr   map[string]interface{}
	lm     *kt.LastMetadata
	frozen bool
}

// Wrap updates the starlark.Metric to wrap a new telegraf.Metric.
func (m *Mib) Wrap(idx string, key string, attr map[string]interface{}, lm *kt.LastMetadata) {
	m.idx = idx
	m.key = key
	m.attr = attr
	m.lm = lm
	m.frozen = false
}

// String returns the starlark representation of the Metric.
//
// The String function is called by both the repr() and str() functions, and so
// it behaves more like the repr function would in Python.
func (m *Mib) String() string {
	buf := new(strings.Builder)
	buf.WriteString("Mib(")
	buf.WriteString(fmt.Sprintf("Index: %s, ", m.idx))
	buf.WriteString(fmt.Sprintf("Key: %s, ", m.key))
	buf.WriteString(fmt.Sprintf("Attr: %v", m.attr))
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
	names := []string{"this.idx", "this.key"}
	for k, _ := range m.attr {
		names = append(names, k)
	}
	return names
}

// Attr implements the starlark.HasAttrs interface.
func (m *Mib) Attr(name string) (starlark.Value, error) {
	switch name {
	case "this.idx":
		return starlark.String(m.idx), nil
	case "this.key":
		return starlark.String(m.key), nil
	default:
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
	}

	// By default, return empty string for all other variables.
	return starlark.String(""), nil
}

// SetField implements the starlark.HasSetField interface.
func (m *Mib) SetField(name string, value starlark.Value) error {
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
		if m.lm != nil {
			if _, ok := m.lm.XtraInfo[name]; !ok {
				m.lm.XtraInfo[name] = m.lm.XtraInfo[m.key]
			}
		}

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
