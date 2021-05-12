package netflow

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/kentik/ktranslate/pkg/util/netflow/ipfix"
	"github.com/kentik/ktranslate/pkg/util/netflow/netflow1"
	"github.com/kentik/ktranslate/pkg/util/netflow/netflow5"
	"github.com/kentik/ktranslate/pkg/util/netflow/netflow6"
	"github.com/kentik/ktranslate/pkg/util/netflow/netflow7"
	"github.com/kentik/ktranslate/pkg/util/netflow/netflow9"
	"github.com/kentik/ktranslate/pkg/util/netflow/session"
)

// Decoder for NetFlow messages.
type Decoder struct {
	session.Session
}

// Message generlized interface.
type Message interface {
}

// NewDecoder sets up a decoder suitable for reading NetFlow packets.
func NewDecoder(s session.Session) *Decoder {
	return &Decoder{s}
}

// Read a single Netflow message from the network. If an error is returned,
// there is no guarantee the following reads will be succesful.
func (d *Decoder) Read(r io.Reader) (Message, error) {
	data := [2]byte{}
	if _, err := io.ReadFull(r, data[:]); err != nil {
		return nil, fmt.Errorf("failed to read netflow version number: %v", err)
	}

	version := binary.BigEndian.Uint16(data[:])
	buffer := bytes.NewBuffer(data[:])
	mr := io.MultiReader(buffer, r)

	switch version {
	case netflow1.Version:
		return netflow1.Read(mr)

	case netflow5.Version:
		return netflow5.Read(mr)

	case netflow6.Version:
		return netflow6.Read(mr)

	case netflow7.Version:
		return netflow7.Read(mr)

	case netflow9.Version:
		return netflow9.Read(mr, d.Session, nil)

	case ipfix.Version:
		return ipfix.Read(mr, d.Session, nil)

	default:
		return nil, fmt.Errorf("netflow: unsupported version %d", version)
	}
}
