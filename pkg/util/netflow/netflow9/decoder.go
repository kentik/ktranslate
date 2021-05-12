package netflow9

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kentik/ktranslate/pkg/util/netflow/session"
)

func errInvalidVersion(v uint16) error {
	return fmt.Errorf("version %d is not a valid NetFlow packet version", v)
}

func errProtocol(f string, v ...interface{}) error {
	return fmt.Errorf("protocol error: "+f, v...)
}

func errTemplateNotFound(t uint16) error {
	return fmt.Errorf("template with id=%d not found", t)
}

// Decoder can decode multiple IPFIX messages from a stream.
type Decoder struct {
	io.Reader
	session.Session
	*Translate
}

func NewDecoder(r io.Reader, s session.Session) *Decoder {
	return &Decoder{r, s, NewTranslate(s)}
}

// Decode decodes a single message from a buffer of bytes.
func (d *Decoder) Decode(data []byte) (*Packet, error) {
	return Read(bytes.NewBuffer(data), d.Session, d.Translate)
}

// Next decodes the next message from the stream. Note that if there is an
// exception, depending on where the exception originated from, the decoder
// results can no longer be trusted and the stream should be reset.
func (d *Decoder) Next() (*Packet, error) {
	return Read(d.Reader, d.Session, d.Translate)
}

// Read a single Netflow packet from the provided reader and decode all the sets.
func Read(r io.Reader, s session.Session, t *Translate) (*Packet, error) {
	p := new(Packet)

	if t == nil && s != nil {
		t = NewTranslate(s)
	}

	if err := p.Header.Unmarshal(r); err != nil {
		return nil, err
	}
	if p.Header.Version != Version {
		return nil, errInvalidVersion(p.Header.Version)
	}
	if p.Header.Len() < 4 {
		return nil, io.ErrShortBuffer
	}
	if p.Header.Count == 0 {
		return p, nil
	}
	return p, p.UnmarshalFlowSets(r, s, t)
}
