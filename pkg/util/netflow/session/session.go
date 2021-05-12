// Package session provides sessions for the Netflow version 9 and IPFIX
// decoders that need to track templates bound to a session.
package session

import "sync"

type Template interface {
	ID() uint16
}

type Session interface {
	Lock()
	Unlock()

	// To keep track of templates
	AddTemplate(uint32, Template)
	GetTemplate(uint32, uint16) (t Template, found bool)
}

// templates ids are only unique within a sub-device scope that the netflow v9
// spec calls a "source id" and ipfix calls an "observation domain id". Use the
// combination of one of these ids and the template id to look up templates.
type templateKey struct {
	sourceID   uint32
	templateID uint16
}

type basicSession struct {
	mutex     *sync.Mutex
	templates map[templateKey]Template
}

func New() *basicSession {
	return &basicSession{
		mutex:     &sync.Mutex{},
		templates: make(map[templateKey]Template, 65536),
	}
}

func (s *basicSession) Lock() {
	s.mutex.Lock()
}

func (s *basicSession) Unlock() {
	s.mutex.Unlock()
}

func (s *basicSession) AddTemplate(sourceID uint32, t Template) {
	s.templates[templateKey{sourceID: sourceID, templateID: t.ID()}] = t
}

func (s *basicSession) GetTemplate(sourceID uint32, templateID uint16) (t Template, found bool) {
	t, found = s.templates[templateKey{sourceID: sourceID, templateID: templateID}]
	return
}

// Test if basicSession is compliant
var _ Session = (*basicSession)(nil)
