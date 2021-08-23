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
	logTee chan string
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, logTee chan string) (*StdoutSink, error) {
	return &StdoutSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "stdoutSink"}, log),
		logTee:   logTee,
	}, nil
}

func (s *StdoutSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	if s.logTee != nil {
		go s.watchLogs(ctx)
	}

	return nil
}

func (s *StdoutSink) Send(ctx context.Context, payload *kt.Output) {
	fmt.Printf("%s\n", string(payload.Body))
}

func (s *StdoutSink) Close() {}

func (s *StdoutSink) HttpInfo() map[string]float64 {
	return map[string]float64{}
}

func (s *StdoutSink) watchLogs(ctx context.Context) {
	s.Infof("Receiving logs...")
	for {
		select {
		case log := <-s.logTee:
			s.Send(ctx, kt.NewOutput([]byte(log)))
		case <-ctx.Done():
			s.Infof("Logs received")
			return
		}
	}
}
