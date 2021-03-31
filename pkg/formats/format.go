package formats

import (
	"fmt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/kentik/ktranslate/pkg/formats/avro"
	"github.com/kentik/ktranslate/pkg/formats/influx"
	"github.com/kentik/ktranslate/pkg/formats/json"
	"github.com/kentik/ktranslate/pkg/formats/netflow"
	"github.com/kentik/ktranslate/pkg/formats/nrm"
	"github.com/kentik/ktranslate/pkg/formats/prom"
	"github.com/kentik/ktranslate/pkg/formats/splunk"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"
)

type Formatter interface {
	To([]*kt.JCHF, []byte) ([]byte, error)
	From([]byte) ([]map[string]interface{}, error)
	Rollup([]rollup.Rollup) ([]byte, error)
}

type Format string

const (
	FORMAT_AVRO    Format = "avro"
	FORMAT_JSON           = "json"
	FORMAT_NETFLOW        = "netflow"
	FORMAT_INFLUX         = "influx"
	FORMAT_PROM           = "prometheus"
	FORMAT_NR             = "new_relic"
	FORMAT_NRM            = "new_relic_metric"
	FORMAT_SPLUNK         = "splunk"
)

func NewFormat(format Format, log logger.Underlying, compression kt.Compression) (Formatter, error) {
	switch format {
	case FORMAT_AVRO:
		return avro.NewFormat(log, compression)
	case FORMAT_JSON:
		return json.NewFormat(log, compression, false)
	case FORMAT_NETFLOW:
		return netflow.NewFormat(log, compression)
	case FORMAT_INFLUX:
		return influx.NewFormat(log, compression)
	case FORMAT_PROM:
		return prom.NewFormat(log, compression)
	case FORMAT_NR:
		return json.NewFormat(log, compression, true)
	case FORMAT_NRM:
		return nrm.NewFormat(log, compression)
	case FORMAT_SPLUNK:
		return splunk.NewFormat(log, compression)
	default:
		return nil, fmt.Errorf("Unknown format %v", format)
	}
}
