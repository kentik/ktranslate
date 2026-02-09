package maps

import (
	"context"

	kkapi "github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/maps/api"
	"github.com/kentik/ktranslate/pkg/maps/file"
)

type Mapper string

const (
	FileMapper Mapper = "file"
	ApiMapper  Mapper = "api"
	NullMapper        = "null"
)

type TagMapper interface {
	LookupTagValue(kt.Cid, uint32, string) (string, string, bool)
	LookupTagValueBig(kt.Cid, int64, string) (string, string, bool)
	LookupKV(uint32) string
	Run(ctx context.Context)
}

func LoadMapper(mtype Mapper, log logger.Underlying, tagMapFilePath string, apic *kkapi.KentikApi) (TagMapper, error) {
	switch mtype {
	case FileMapper:
		if tagMapFilePath != "" {
			return file.NewFileTagMapper(log, tagMapFilePath)
		} else {
			if apic != nil { // Try the api version if the client is present.
				return api.NewApiTagMapper(log, apic)
			} else {
				return &NullType{ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "nullMapper"}, log)}, nil
			}
		}
	case ApiMapper:
		return api.NewApiTagMapper(log, apic)
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

func (ntm *NullType) LookupTagValueBig(cid kt.Cid, tagval int64, colname string) (string, string, bool) {
	return "", "", false
}

func (ntm *NullType) LookupKV(k uint32) string {
	return ""
}

func (ntm *NullType) Run(ctx context.Context) {}
