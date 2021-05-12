package ipfix

import (
	"github.com/kentik/ktranslate/pkg/util/netflow/session"
	"github.com/kentik/ktranslate/pkg/util/netflow/translate"
)

type TranslatedField struct {
	Name                 string
	InformationElementID uint16
	EnterpriseNumber     uint32
	Value                interface{}
	Bytes                []byte
}

type Translate struct {
	*translate.Translate
}

func NewTranslate(s session.Session) *Translate {
	return &Translate{translate.NewTranslate(s)}
}

func (t *Translate) Record(templateID uint16, fields Fields, fss FieldSpecifiers) error {
	if fss == nil {
		if debug {
			debugLog.Printf("no fields in template id=%d, can't translate\n", templateID)
		}
		return nil
	}

	if debug {
		debugLog.Printf("translating %d/%d fields\n", len(fields), len(fss))
	}

	for i, field := range fss {
		if i > len(fields) {
			break
		}
		f := &(fields[i])
		f.Translated = &TranslatedField{}
		f.Translated.EnterpriseNumber = field.EnterpriseNumber
		f.Translated.InformationElementID = field.InformationElementID

		if element, ok := t.Translate.Key(translate.Key{field.EnterpriseNumber, field.InformationElementID}); ok {
			f.Translated.Name = element.Name
			f.Translated.Value = translate.Bytes(f.Bytes, element.Type)
			if debug {
				debugLog.Printf("translated {%d, %d} to %s, %v\n", field.EnterpriseNumber, field.InformationElementID, f.Translated.Name, f.Translated.Value)
			}
		} else if debug {
			debugLog.Printf("no translator element for {%d, %d}\n", field.EnterpriseNumber, field.InformationElementID)
		}
	}

	return nil
}
