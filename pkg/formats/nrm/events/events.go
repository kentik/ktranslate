package events

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/kentik/ktranslate/pkg/kt"
)

func SendEvent(msg *kt.JCHF, doGz bool, evts chan []byte) error {
	res, err := toRawJson(msg, doGz)
	if err != nil {
		return err
	}
	select {
	case evts <- res: // Give this guy up to the sender.
	default:
		// noop
	}

	return nil
}

func toRawJson(msg *kt.JCHF, doGz bool) ([]byte, error) {
	msgsNew := []map[string]interface{}{msg.Flatten()}
	t, err := json.Marshal(msgsNew)
	if err != nil {
		return nil, err
	}
	target := t

	if !doGz {
		return target, nil
	}

	buf := &bytes.Buffer{}
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(target)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
