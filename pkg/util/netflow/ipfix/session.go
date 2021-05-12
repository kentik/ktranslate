package ipfix

import "sync"

type Session struct {
	buffers   *sync.Pool
	mutex     *sync.Mutex
	templates [][]FieldSpecifier
	minSize   []uint16
}

func NewSession() *Session {
	return &Session{
		buffers: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 65536)
			},
		},
		mutex:     &sync.Mutex{},
		templates: make([][]FieldSpecifier, 65536),
		minSize:   make([]uint16, 65536),
	}
}
