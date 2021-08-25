package netflow9

import (
	"fmt"

	"github.com/kentik/ktranslate/pkg/util/netflow/session"
	"github.com/kentik/ktranslate/pkg/util/netflow/translate"
)

type TranslatedField struct {
	Name  string
	Type  uint16
	Value interface{}
	Bytes []byte
}

func (tf TranslatedField) String() string {
	return fmt.Sprintf("%s=%v", tf.Name, tf.Value)
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
		if i >= len(fields) {
			break
		}
		f := &fields[i]
		f.Translated = &TranslatedField{}
		f.Translated.Type = field.Type

		if element, ok := t.Translate.Key(translate.Key{EnterpriseID: 0, FieldID: field.Type}); ok {
			f.Translated.Name = element.Name
			f.Translated.Value = translate.Bytes(fields[i].Bytes, element.Type)
		} else if debug {
			debugLog.Printf("no translator element for {0, %d}\n", field.Type)
		}
	}

	return nil
}
