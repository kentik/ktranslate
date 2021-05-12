package netflow9

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/kentik/ktranslate/pkg/util/netflow/read"
	"github.com/kentik/ktranslate/pkg/util/netflow/session"
)

const (
	// Version word in the Packet Header
	Version uint16 = 0x0009
)

// Packet consists of a Packet Header followed by one or more FlowSets. The
// FlowSets can be any of the possible three types: Template, Data, or Options
// Template.
//
// The format of the Packet on the wire is:
//
//   +--------+-------------------------------------------+
//   |        | +----------+ +---------+ +----------+     |
//   | Packet | | Template | | Data    | | Options  |     |
//   | Header | | FlowSet  | | FlowSet | | Template | ... |
//   |        | |          | |         | | FlowSet  |     |
//   |        | +----------+ +---------+ +----------+     |
//   +--------+-------------------------------------------+
type Packet struct {
	Header                  PacketHeader
	TemplateFlowSets        []TemplateFlowSet
	OptionsTemplateFlowSets []OptionsTemplateFlowSet
	DataFlowSets            []DataFlowSet
	OptionsDataFlowSets     []OptionsDataFlowSet
}

// PacketHeader is a Packet Header (RFC 3954 section 5.1)
type PacketHeader struct {
	Version        uint16
	Count          uint16
	SysUpTime      uint32
	UnixSecs       uint32
	SequenceNumber uint32
	SourceID       uint32
}

func (p *Packet) UnmarshalFlowSets(r io.Reader, s session.Session, t *Translate) error {
	if debug {
		debugLog.Printf("decoding %d flow sets, sequence number: %d\n", p.Header.Count, p.Header.SequenceNumber)
	}

	// different implementations of netflow v9 treat the header Count field differently; some
	// count the number of sets included in the message, while others count the number of records
	// included in those sets.  Since there's no consistency, we can't depend on the Count field
	// to tell us when to leave the loop; we have to rely on an external source to tell us the
	// size of the v9 message (in practice, this is the size of the UDP datagram that contains
	// the message), and try to read as many sets as are contained in that message.  (Note that
	// IPFIX explicitly trades this field for a Length-in-octets one, so that it can follow the
	// same logic without depending on the transport layer.)

	// Read the rest of the message, containing the sets.
	buffer := new(bytes.Buffer)
	if ndata, err := buffer.ReadFrom(r); err != nil {
		// *IF* r is itself a bytes.Buffer, here, then we'll never hit this case: Buffer.ReadFrom
		// guarantees it won't return io.EOF, and Buffer.Read is supposed to return only io.EOF as an error
		return err
	} else if ndata == 0 {
		return fmt.Errorf("netflow9 message has zero-length body")
	}

	for i := uint16(0); i < p.Header.Count && buffer.Len() > 0; i++ {
		// Read the next set header
		header := FlowSetHeader{}
		if err := header.Unmarshal(buffer); err != nil {
			if debug {
				debugLog.Println("failed to read flow set header:", err)
			}
			return fmt.Errorf("netflow9 FlowSetHeader %d truncated", i)
		}

		// read the body of that set (and return an error if the body has been truncated)
		readSize := int(header.Length) - header.Len()
		if readSize <= 0 {
			if debug {
				debugLog.Printf("short read size of %d\n", readSize)
			}
			return fmt.Errorf("netflow9 FlowSet %d has illegal size %d bytes", i, readSize)
		}
		data := make([]byte, readSize)
		if ndata, err := io.ReadFull(buffer, data); err != nil {
			if debug {
				debugLog.Printf("netflow9 FlowSet %d truncated; header says %d bytes, but message contains only %d", i, readSize, ndata)
			}
			return fmt.Errorf("netflow9 FlowSet %d truncated; header says %d bytes, but message contains only %d", i, readSize, ndata)
		}

		switch header.ID {
		case 0: // Template FlowSet
			tfs := TemplateFlowSet{}
			tfs.Header = header

			if err := tfs.UnmarshalRecords(bytes.NewBuffer(data)); err != nil {
				return err
			}
			if debug {
				debugLog.Printf("unmarshaled %d records: %v\n", len(tfs.Records), tfs)
			}

			for _, tr := range tfs.Records {
				tr.register(p.Header.SourceID, s)
			}

			p.TemplateFlowSets = append(p.TemplateFlowSets, tfs)

		case 1: // Options Template FlowSet
			ots := OptionsTemplateFlowSet{}
			ots.Header = header
			if err := ots.UnmarshalRecords(bytes.NewBuffer(data)); err != nil {
				return err
			}

			for _, otr := range ots.Records {
				otr.register(p.Header.SourceID, s)
			}

			p.OptionsTemplateFlowSets = append(p.OptionsTemplateFlowSets, ots)

		default:
			// If we don't have a session, or no template to resolve the Data
			// Set contained Data Records, we'll store the raw bytes in stead.
			if s == nil {
				if debug {
					debugLog.Printf("no session, storing %d raw bytes in data set\n", len(data))
				}
				continue
			}

			tm, ok := s.GetTemplate(p.Header.SourceID, header.ID)
			if !ok {
				if debug {
					debugLog.Printf("no template for id=%d, storing %d raw bytes in data set\n", header.ID, len(data))
				}
				continue
			}

			switch template := tm.(type) {
			case TemplateRecord:
				dfs := DataFlowSet{}
				dfs.Header = header

				if err := dfs.Unmarshal(bytes.NewBuffer(data), template, t); err != nil {
					return err
				}
				p.DataFlowSets = append(p.DataFlowSets, dfs)

			case OptionsTemplateRecord:
				ods := OptionsDataFlowSet{}
				ods.Header = header

				if err := ods.Unmarshal(bytes.NewBuffer(data), template, t); err != nil {
					return err
				}
				p.OptionsDataFlowSets = append(p.OptionsDataFlowSets, ods)

			default:
				return fmt.Errorf("netflow9 DataFlowSet %d with source-id %d, template-id %d, retrieved template record of type %T; rejected", i, p.Header.SourceID, header.ID, tm)
			}
		}
	}

	return nil
}

func (h PacketHeader) Len() int {
	return 20
}

func (h *PacketHeader) Unmarshal(r io.Reader) error {
	if err := read.Uint16(&h.Version, r); err != nil {
		return fmt.Errorf("netflow9 header truncated while reading version field: %v", err)
	}
	if err := read.Uint16(&h.Count, r); err != nil {
		return fmt.Errorf("netflow9 header truncated while reading count field: %v", err)
	}
	if err := read.Uint32(&h.SysUpTime, r); err != nil {
		return fmt.Errorf("netflow9 header truncated while reading sysUpTime field: %v", err)
	}
	if err := read.Uint32(&h.UnixSecs, r); err != nil {
		return fmt.Errorf("netflow9 header truncated while reading unix-secs field: %v", err)
	}
	if err := read.Uint32(&h.SequenceNumber, r); err != nil {
		return fmt.Errorf("netflow9 header truncated while reading sequence-number field: %v", err)
	}
	if err := read.Uint32(&h.SourceID, r); err != nil {
		return fmt.Errorf("netflow9 header truncated while reading source-id field: %v", err)
	}

	return nil
}

type FlowSetHeader struct {
	ID     uint16
	Length uint16
}

func (h *FlowSetHeader) Len() int {
	return 4
}

func (h *FlowSetHeader) Unmarshal(r io.Reader) error {
	if err := read.Uint16(&h.ID, r); err != nil {
		return err
	}
	if err := read.Uint16(&h.Length, r); err != nil {
		return err
	}

	return nil
}

// TemplateFlowSet enhance the flexibility of the Flow Record format because
// they allow the NetFlow Collector to process Flow Records without necessarily
// knowing the interpretation of all the data in the Flow Record.
//
// The format of the Template FlowSet is as follows:
//
//   0                   1                   2                   3
//   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |       FlowSet ID = 0          |          Length               |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |      Template ID 256          |         Field Count           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |        Field Type 1           |         Field Length 1        |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |        Field Type 2           |         Field Length 2        |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |             ...               |              ...              |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |        Field Type N           |         Field Length N        |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |      Template ID 257          |         Field Count           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |        Field Type 1           |         Field Length 1        |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |        Field Type 2           |         Field Length 2        |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |             ...               |              ...              |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |        Field Type M           |         Field Length M        |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |             ...               |              ...              |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |        Template ID K          |         Field Count           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |             ...               |              ...              |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type TemplateFlowSet struct {
	Header  FlowSetHeader
	Records []TemplateRecord
}

func (tfs *TemplateFlowSet) UnmarshalRecords(r io.Reader) error {
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(r); err != nil {
		return err
	}

	// As long as there are more than 4 bytes in the buffer, we parse the next
	// TemplateRecord, otherwise it's padding.
	tfs.Records = make([]TemplateRecord, 0)
	for buffer.Len() > 4 {
		record := TemplateRecord{}
		if err := record.Unmarshal(buffer); err != nil {
			return err
		}

		tfs.Records = append(tfs.Records, record)
	}

	return nil
}

// TemplateRecord is a Template Record as per RFC3964 section 5.2
type TemplateRecord struct {
	TemplateID uint16
	FieldCount uint16
	Fields     FieldSpecifiers
}

func (tr TemplateRecord) register(sourceID uint32, s session.Session) {
	if s == nil {
		return
	}
	if debug {
		debugLog.Println("register template:", tr)
	}
	s.Lock()
	defer s.Unlock()
	s.AddTemplate(sourceID, tr)
}

func (tr TemplateRecord) ID() uint16 {
	return tr.TemplateID
}

func (tr TemplateRecord) String() string {
	return fmt.Sprintf("id=%d fields=%d (%s)", tr.TemplateID, tr.FieldCount, tr.Fields)
}

func (tr TemplateRecord) Size() int {
	var size int
	for _, f := range tr.Fields {
		size += int(f.Length)
	}
	return size
}

func (tr *TemplateRecord) Unmarshal(r io.Reader) error {
	if err := read.Uint16(&tr.TemplateID, r); err != nil {
		return err
	}
	if err := read.Uint16(&tr.FieldCount, r); err != nil {
		return err
	}

	tr.Fields = make(FieldSpecifiers, tr.FieldCount)
	if err := tr.Fields.Unmarshal(r); err != nil {
		return err
	}

	return nil
}

type FieldSpecifier struct {
	Type   uint16
	Length uint16
}

func (fs *FieldSpecifier) String() string {
	return fmt.Sprintf("type=%d length=%d", fs.Type, fs.Length)
}

func (f *FieldSpecifier) Unmarshal(r io.Reader) error {
	if err := read.Uint16(&f.Type, r); err != nil {
		return err
	}
	if err := read.Uint16(&f.Length, r); err != nil {
		return err
	}

	return nil
}

type FieldSpecifiers []FieldSpecifier

func (fs FieldSpecifiers) String() string {
	v := make([]string, len(fs))
	for i, f := range fs {
		v[i] = f.String()
	}
	return strings.Join(v, ",")
}

func (fs *FieldSpecifiers) Unmarshal(r io.Reader) error {
	for i := 0; i < len(*fs); i++ {
		if err := (*fs)[i].Unmarshal(r); err != nil {
			return err
		}
	}
	return nil
}

// OptionsTemplateRecord (and its corresponding OptionsDataRecord) is used to
// supply information about the NetFlow process configuration or NetFlow
// process specific data, rather than supplying information about IP Flows.
//
// For example, the Options Template FlowSet can report the sample rate
// of a specific interface, if sampling is supported, along with the
// sampling method used.
//
// The format of the Options Template FlowSet follows:
//
//    0                   1                   2                   3
//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |       FlowSet ID = 1          |          Length               |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |         Template ID           |      Option Scope Length      |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |        Option Length          |       Scope 1 Field Type      |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |     Scope 1 Field Length      |               ...             |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |     Scope N Field Length      |      Option 1 Field Type      |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |     Option 1 Field Length     |             ...               |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |     Option M Field Length     |           Padding             |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type OptionsTemplateFlowSet struct {
	Header  FlowSetHeader
	Records []OptionsTemplateRecord
}

func (ots OptionsTemplateFlowSet) String() string {
	return fmt.Sprintf("%v (%d records)", ots.Header, len(ots.Records))
}

func (ots *OptionsTemplateFlowSet) UnmarshalRecords(r io.Reader) error {
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(r); err != nil {
		return err
	}

	// As long as there are more than 4 bytes in the buffer, we parse the next
	// TemplateRecord, otherwise it's padding.
	ots.Records = make([]OptionsTemplateRecord, 0)
	for buffer.Len() > 4 {
		record := OptionsTemplateRecord{}
		if err := record.Unmarshal(buffer); err != nil {
			return err
		}

		ots.Records = append(ots.Records, record)
	}

	return nil
}

type OptionsTemplateRecord struct {
	// Each Options Template Record is given a unique Template ID in the
	// range 256 to 65535.
	TemplateID uint16

	// Number of scope fields in this Options Template Record. The Scope
	// Fields are normal Fields, except that they are interpreted as
	// scope at the Collector. A scope field count of N specifies that
	// the first N Field Specifiers in the Template Record are Scope
	// Fields. The Scope Field Count MUST NOT be zero.
	ScopeFieldCount uint16
	ScopeFields     FieldSpecifiers

	// Number of non-scope fields in this Options Template Record
	FieldCount uint16
	Fields     FieldSpecifiers
}

func (otr OptionsTemplateRecord) register(observationDomainID uint32, s session.Session) {
	if s == nil {
		return
	}
	if debug {
		debugLog.Println("register options template:", otr)
	}
	s.Lock()
	defer s.Unlock()
	s.AddTemplate(observationDomainID, otr)
}

func (otr OptionsTemplateRecord) ID() uint16 {
	return otr.TemplateID
}

func (otr OptionsTemplateRecord) String() string {
	return fmt.Sprintf("id=%d fields=%d (%s) scope fields=%d (%s)",
		otr.TemplateID, otr.FieldCount, otr.Fields, otr.ScopeFieldCount, otr.ScopeFields)
}

func (otr *OptionsTemplateRecord) Unmarshal(r io.Reader) error {
	if err := read.Uint16(&otr.TemplateID, r); err != nil {
		return err
	}
	if err := read.Uint16(&otr.ScopeFieldCount, r); err != nil {
		return err
	}
	if err := read.Uint16(&otr.FieldCount, r); err != nil {
		return err
	}

	// For inexplicable reasons, v9 data templates have counts of the number of fields
	// included, but v9 options templates have counts of the number of bytes included
	// in fields.  But each field always takes up 4 bytes (two for field type, two for
	// field length), so we'll divide by four here and proceed.
	otr.ScopeFieldCount /= 4
	otr.FieldCount /= 4

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r)
	if debug {
		hexdump(buffer.Bytes())
	}

	otr.ScopeFields = make(FieldSpecifiers, otr.ScopeFieldCount)
	if err := otr.ScopeFields.Unmarshal(buffer); err != nil {
		return err
	}

	otr.Fields = make(FieldSpecifiers, otr.FieldCount)
	if err := otr.Fields.Unmarshal(buffer); err != nil {
		return err
	}

	return nil
}

type DataFlowSet struct {
	Header  FlowSetHeader
	Records []DataRecord
	Bytes   []byte
}

func (dfs *DataFlowSet) Unmarshal(r io.Reader, tr TemplateRecord, t *Translate) error {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r)

	dfs.Records = make([]DataRecord, 0)
	for buffer.Len() >= 4 { // Continue until only padding alignment bytes left
		var dr = DataRecord{}
		dr.TemplateID = tr.TemplateID
		if err := dr.Unmarshal(bytes.NewBuffer(buffer.Next(tr.Size())), tr.Fields, t); err != nil {
			return err
		}
		dfs.Records = append(dfs.Records, dr)
	}

	return nil
}

type DataRecord struct {
	TemplateID uint16
	Fields     Fields
}

func (dr *DataRecord) Unmarshal(r io.Reader, fss FieldSpecifiers, t *Translate) error {
	// We don't know how many records there are in a Data Set, so we'll keep
	// reading until we exhausted the buffer.
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(r); err != nil {
		return err
	}

	dr.Fields = make(Fields, 0)
	var err error
	for i := 0; buffer.Len() > 0 && i < len(fss); i++ {
		f := Field{
			Type:   fss[i].Type,
			Length: fss[i].Length,
		}
		if err = f.Unmarshal(buffer); err != nil {
			return err
		}
		dr.Fields = append(dr.Fields, f)
	}

	if t != nil && len(dr.Fields) > 0 {
		if err := t.Record(dr.TemplateID, dr.Fields, fss); err != nil {
			return err
		}
	}

	return nil
}

type OptionsDataFlowSet struct {
	Header  FlowSetHeader
	Bytes   []byte
	Records []OptionsDataRecord
}

func (ods *OptionsDataFlowSet) Unmarshal(r io.Reader, otr OptionsTemplateRecord, t *Translate) error {
	// We don't know how many records there are in a Data Set, so we'll keep
	// reading until we exhausted the buffer.
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r)

	ods.Records = make([]OptionsDataRecord, 0)
	for buffer.Len() >= 4 {
		var odr = OptionsDataRecord{}
		odr.TemplateID = otr.TemplateID
		if err := odr.Unmarshal(buffer, otr.ScopeFields, otr.Fields, t); err != nil {
			// If we hit EOF, we've exhausted the buffer. The current DataRecord is discarded,
			// and we exit normally.
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
		ods.Records = append(ods.Records, odr)
	}

	return nil
}

type OptionsDataRecord struct {
	TemplateID  uint16
	ScopeFields Fields
	Fields      Fields
}

func (odr *OptionsDataRecord) Unmarshal(r io.Reader, scopeFss FieldSpecifiers, fss FieldSpecifiers, t *Translate) error {
	odr.ScopeFields = make(Fields, 0)
	odr.Fields = make(Fields, 0)
	var err error
	for i := 0; i < len(scopeFss); i++ {
		f := Field{
			Type:   scopeFss[i].Type,
			Length: scopeFss[i].Length,
		}
		if err = f.Unmarshal(r); err != nil {
			return err
		}
		odr.ScopeFields = append(odr.ScopeFields, f)
	}
	for i := 0; i < len(fss); i++ {
		f := Field{
			Type:   fss[i].Type,
			Length: fss[i].Length,
		}
		if err = f.Unmarshal(r); err != nil {
			return err
		}
		odr.Fields = append(odr.Fields, f)
	}

	if t != nil {
		if err := t.Record(odr.TemplateID, odr.ScopeFields, scopeFss); err != nil {
			return err
		}
		if err := t.Record(odr.TemplateID, odr.Fields, fss); err != nil {
			return err
		}
	}

	return nil
}

type Field struct {
	Type       uint16
	Length     uint16
	Translated *TranslatedField
	Bytes      []byte
}

func (f *Field) Unmarshal(r io.Reader) error {
	f.Bytes = make([]byte, f.Length)
	if _, err := r.Read(f.Bytes); err != nil {
		return err
	}

	return nil
}

type Fields []Field
