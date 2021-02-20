package properties

import (
	"os"
	"sync"
)

type PropertyBacking interface {
	GetProperty(name string) (value string, present bool)
	Refresh()
}

type PropertyBackingGetter func(name string) (value string, present bool)

// map backing

type staticMapPropertyBacking struct {
	propertyMap map[string]string
}

func NewStaticMapPropertyBacking(propertyMap map[string]string) *staticMapPropertyBacking {
	mapCopy := make(map[string]string)
	for k, v := range propertyMap {
		mapCopy[k] = v
	}
	return &staticMapPropertyBacking{
		propertyMap: mapCopy,
	}
}

func (back *staticMapPropertyBacking) GetProperty(name string) (string, bool) {
	value, present := back.propertyMap[name]
	return value, present
}

func (back *staticMapPropertyBacking) Refresh() {
}

// environment variables backing

type envPropertyBacking struct{}

func NewEnvPropertyBacking() *envPropertyBacking {
	return &envPropertyBacking{}
}

func (back *envPropertyBacking) GetProperty(name string) (string, bool) {
	value, present := os.LookupEnv(name)
	return value, present
}

func (back *envPropertyBacking) Refresh() {
}

// filesystem backing

type filesystemPropertyBacking struct {
	root        string
	propertyMap map[string]string
	sync.RWMutex
}

func NewFileSystemPropertyBacking(root string) *filesystemPropertyBacking {
	back := &filesystemPropertyBacking{
		root: root,
	}
	back.load()
	return back
}

func (back *filesystemPropertyBacking) load() {
	if newMap, err := loadPropertiesFromFilesystem(back.root); err == nil {
		back.Lock()
		defer back.Unlock()
		back.propertyMap = newMap
	}
}

func (back *filesystemPropertyBacking) GetProperty(name string) (string, bool) {
	back.RLock()
	defer back.RUnlock()
	value, present := back.propertyMap[name]
	return value, present
}

func (back *filesystemPropertyBacking) Refresh() {
	back.load()
}
