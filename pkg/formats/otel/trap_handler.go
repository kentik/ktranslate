package otel

import (
	"context"
	"log/slog"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/agoda-com/opentelemetry-go/otelslog"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogsgrpc"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// configure common attributes for all logs
func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("ktranslate"),
	)
}

const (
	InstNameSnmpTrapEvent = "snmp-trap-events"
)

type OtelLogger struct {
	ctx    context.Context
	logger *slog.Logger
	log    logger.ContextL
}

func NewLogger(ctx context.Context, log logger.ContextL, cfg *ktranslate.OtelFormatConfig) (*OtelLogger, error) {
	// configure opentelemetry logger provider
	var logExporter sdk.LogRecordExporter

	switch cfg.Protocol {
	case "stdout":
		le, _ := stdoutlogs.NewExporter()
		logExporter = le
	case "http", "https":
		le, err := otlplogs.NewExporter(ctx)
		if err != nil {
			return nil, err
		}
		logExporter = le
	case "grpc":
		le, err := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogsgrpc.NewClient(otlplogsgrpc.WithEndpoint(cfg.Endpoint))))
		if err != nil {
			return nil, err
		}
		logExporter = le
	}

	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(logExporter),
		sdk.WithResource(newResource()),
	)

	otelLogger := slog.New(otelslog.NewOtelHandler(loggerProvider, &otelslog.HandlerOptions{}))

	//configure default logger
	slog.SetDefault(otelLogger)
	ol := &OtelLogger{ctx: ctx, logger: otelLogger, log: log}
	go ol.watchForClose(ctx, loggerProvider)

	return ol, nil
}

// @TODO, actualy call this.
func (ol *OtelLogger) watchForClose(ctx context.Context, loggerProvider *sdk.LoggerProvider) {
	for {
		select {
		case <-ctx.Done():
			// gracefully shutdown logger to flush accumulated signals before program finish
			ol.log.Infof("Done with Trap Logger.")
			loggerProvider.Shutdown(ctx)
			return
		}
	}
}

// For now, just log everything as json
func (ol *OtelLogger) RecordTrap(msg *kt.JCHF) error {
	flat := msg.Flatten()
	strip(flat)

	atrs := make([]slog.Attr, 0)
	for k, v := range flat {
		atrs = append(atrs, slog.Any(k, v))
	}
	slog.LogAttrs(ol.ctx, slog.LevelInfo, "New Trap Event", atrs...)

	return nil
}

func strip(in map[string]interface{}) {
	for k, v := range in {
		switch tv := v.(type) {
		case string:
			if tv == "" || tv == "-" || tv == "--" {
				delete(in, k)
			}
		case int32:
			if tv == 0 {
				delete(in, k)
			}
		case int64:
			if tv == 0 {
				delete(in, k)
			}
		}
	}
	in["instrumentation.provider"] = kt.InstProvider // Let them know who sent this.
	in["instrumentation.name"] = InstNameSnmpTrapEvent
	in["collector.name"] = kt.CollectorName
}
