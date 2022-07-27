package maps

import (
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/maps/file"
)

type Mapper string

const (
	FileMapper Mapper = "file"
	NullMapper        = "null"
)

type TagMapper interface {
	LookupTagValue(kt.Cid, uint32, string) (string, string, bool)
}

func LoadMapper(mtype Mapper, log logger.Underlying, tagMapFilePath string) (TagMapper, error) {
	switch mtype {
	case FileMapper:
		return file.NewFileTagMapper(log, tagMapFilePath)
	default:
		return &NullType{ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "nullMapper"}, log)}, nil
	}
}

type NullType struct {
	logger.ContextL
}

func (ntm *NullType) LookupTagValue(cid kt.Cid, tagval uint32, colname string) (string, string, bool) {
	return "", "", false
}
