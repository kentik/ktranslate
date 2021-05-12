package netflow9

import (
	"bytes"
	"testing"

	"github.com/kentik/ktranslate/pkg/util/netflow/session"
)

func TestSourceIDSeparation(t *testing.T) {
	// rfc3954 section 7 requires that different observation domains within a
	// single device can use the same template ID to refer to different template
	// definitions.  So send the netflow9 decoder two different templates sharing
	// the same template id, but with two different source-ids; then confirm that
	// the appropriate template is used for data from the two source-ids.
	var od1Template = []byte{
		// packet Header
		0x00, 0x09, // Protocol version: 9
		0x00, 0x01, // Count: 1
		0x00, 0x00, 0xff, 0xff, // sysUpTime: 65535 ms
		0x5b, 0x6d, 0x08, 0x1d, // export time: 8/10/2018 03:35:57 UTC
		0x00, 0x00, 0x00, 0x01, // sequence number: 1
		0x00, 0x00, 0x00, 0x01, // source id: 1

		// set header
		0x00, 0x00, // set id: 0
		0x00, 0x10, // set length: 16

		// Template 1
		0x01, 0x00, // template id: 256
		0x00, 0x02, // field count: 2
		0x00, 0x08, 0x00, 0x04, // Field 1: IP_SRC_ADDR(8), 4 bytes
		0x00, 0x0c, 0x00, 0x04, // Field 2: IP_DST_ADDR(12), 4 bytes
	}

	var od2Template = []byte{
		// packet Header
		0x00, 0x09, // Protocol version: 9
		0x00, 0x01, // Count: 1
		0x00, 0x00, 0xff, 0xff, // sysUpTime: 65535 ms
		0x5b, 0x6d, 0x08, 0x1d, // export time: 8/10/2018 03:35:57 UTC
		0x00, 0x00, 0x00, 0x01, // sequence number: 1
		0x00, 0x00, 0x00, 0x02, // source id: 2

		// set header
		0x00, 0x00, // set id: 0
		0x00, 0x10, // set length: 16

		// Template 1
		0x01, 0x00, // template id: 256
		0x00, 0x02, // field count: 2
		0x00, 0x0c, 0x00, 0x04, // Field 1: IP_DST_ADDR(12), 4 bytes
		0x00, 0x08, 0x00, 0x04, // Field 2: IP_SRC_ADDR(8), 4 bytes
	}

	var od1Data = []byte{
		// packet Header
		0x00, 0x09, // Protocol version: 9
		0x00, 0x01, // Count: 1
		0x00, 0x00, 0xff, 0xff, // sysUpTime: 65535 ms
		0x5b, 0x6d, 0x08, 0x1d, // export time: 8/10/2018 03:35:57 UTC
		0x00, 0x00, 0x00, 0x02, // sequence number: 1
		0x00, 0x00, 0x00, 0x01, // source id: 1

		// set header
		0x01, 0x00, // set id: 256
		0x00, 0x0c, // set length: 12

		// flow record 1
		0xbc, 0x41, 0x7e, 0xd5, // src address: 188.65.126.213
		0x2e, 0x63, 0xa4, 0x12, // dst address: 46.99.164.18
	}

	// same bytes as od1Data, except for the observation domain id. But because
	// od2's template 256 has IP_DST_ADDR first, the addresses will be reversed
	// in the message
	var od2Data = []byte{
		// packet Header
		0x00, 0x09, // Protocol version: 9
		0x00, 0x01, // Count: 1
		0x00, 0x00, 0xff, 0xff, // sysUpTime: 65535 ms
		0x5b, 0x6d, 0x08, 0x1d, // export time: 8/10/2018 03:35:57 UTC
		0x00, 0x00, 0x00, 0x02, // sequence number: 1
		0x00, 0x00, 0x00, 0x02, // source id: 1

		// set header
		0x01, 0x00, // set id: 256
		0x00, 0x0c, // set length: 12

		// flow record 1
		0xbc, 0x41, 0x7e, 0xd5, // dst address: 188.65.126.213
		0x2e, 0x63, 0xa4, 0x12, // src address: 46.99.164.18
	}

	s := session.New()

	_, err := Read(bytes.NewBuffer(od1Template), s, nil)
	if err != nil {
		t.Errorf("failed to seed session with od1 template: %v", err)
	}

	_, err = Read(bytes.NewBuffer(od2Template), s, nil)
	if err != nil {
		t.Errorf("failed to seed session with od2 template: %v", err)
	}

	od1Msg, err := Read(bytes.NewBuffer(od1Data), s, nil)
	if err != nil {
		t.Errorf("failed to read od1 data: %v", err)
	}

	if len(od1Msg.DataFlowSets) != 1 || len(od1Msg.DataFlowSets[0].Records) != 1 || len(od1Msg.DataFlowSets[0].Records[0].Fields) != 2 {
		t.Errorf("unexpected data sets from od1 data: %v", od1Msg.DataFlowSets)
	}

	od1Fields := od1Msg.DataFlowSets[0].Records[0].Fields
	if od1Fields[0].Translated == nil || od1Fields[1].Translated == nil {
		t.Errorf("untranslated fields from od1 data: %v", od1Fields)
	}

	if od1Fields[0].Translated.Type != 8 && od1Fields[1].Translated.Type != 12 {
		t.Errorf("fields from od1 data in wrong order: %v %v", od1Fields[0].Translated, od1Fields[1].Translated)
	}

	// Check od2 data in the same way
	od2Msg, err := Read(bytes.NewBuffer(od2Data), s, nil)
	if err != nil {
		t.Errorf("failed to read od2 data: %v", err)
	}

	if len(od2Msg.DataFlowSets) != 1 || len(od2Msg.DataFlowSets[0].Records) != 1 || len(od2Msg.DataFlowSets[0].Records[0].Fields) != 2 {
		t.Errorf("unexpected data sets from od2 data: %v", od2Msg.DataFlowSets)
	}

	od2Fields := od2Msg.DataFlowSets[0].Records[0].Fields
	if od2Fields[0].Translated == nil || od2Fields[1].Translated == nil {
		t.Errorf("untranslated fields from od2 data: %v", od2Fields)
	}

	if od2Fields[0].Translated.Type != 12 && od2Fields[1].Translated.Type != 8 {
		t.Errorf("fields from od2 data in wrong order: %v %v", od2Fields[0].Translated, od2Fields[1].Translated)
	}
}
