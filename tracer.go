package tron

import (
	"fmt"
	"strings"

	"github.com/loghole/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/rtconfig"
)

type tracer struct {
	tracer *tracing.Tracer
}

func (t *tracer) init(info *Info) (err error) {
	t.tracer, err = tracing.NewTracer(t.configuration(info))
	if err != nil {
		return fmt.Errorf("init tracer failed: %w", err)
	}

	if err := tracing.EnablePrometheusMetrics(); err != nil {
		return fmt.Errorf("init prometheus metrics failed: %w", err)
	}

	return nil
}

func (t *tracer) configuration(info *Info) *tracing.Configuration {
	configuration := tracing.DefaultConfiguration(
		info.AppName,
		rtconfig.GetString(app.JaegerAddrEnv),
	)

	switch strings.ToLower(rtconfig.GetString(app.JaegerSamplerType)) {
	case "probabilistic":
		configuration.Sampler = trace.TraceIDRatioBased(rtconfig.GetFloat64(app.JaegerSamplerParam))
	case "const":
		if rtconfig.GetFloat64(app.JaegerSamplerParam) > 0 {
			configuration.Sampler = trace.AlwaysSample()
		} else {
			configuration.Sampler = trace.NeverSample()
		}
	}

	configuration.Attributes = append(configuration.Attributes,
		attribute.String("app.version", info.Version),
		attribute.String("app.namespace", info.Namespace.String()),
		attribute.String("app.name", info.AppName),
		attribute.String("app.git_hash", info.GitHash),
		attribute.String("app.build_at", info.BuildAt),
	)

	return configuration
}

func (t *tracer) Close() {
	_ = t.tracer.Close()
}
