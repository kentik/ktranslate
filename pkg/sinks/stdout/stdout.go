package stdout

import (
	"context"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

type StdoutSink struct {
	logger.ContextL
}

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*StdoutSink, error) {
	return &StdoutSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "stdoutSink"}, log),
	}, nil
}

func (s *StdoutSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	return nil
}

func (s *StdoutSink) Send(ctx context.Context, payload *kt.Output) {
	fmt.Printf("%s\n", string(payload.Body))
}

func (s *StdoutSink) Close() {}

func (s *StdoutSink) HttpInfo() map[string]float64 {
	return map[string]float64{}
}
