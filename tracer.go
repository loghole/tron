package tron

import (
	"fmt"

	"github.com/loghole/tracing"
	"github.com/opentracing/opentracing-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"

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

func (t *tracer) configuration(info *Info) *jaegerconfig.Configuration {
	configuration := tracing.DefaultConfiguration(
		info.AppName,
		rtconfig.GetString(app.JaegerAddrEnv),
	)

	if v, _ := rtconfig.GetValue(app.JaegerSamplerType); !v.IsNil() {
		configuration.Sampler.Type = v.String()
	}

	if v, _ := rtconfig.GetValue(app.JaegerSamplerParam); !v.IsNil() {
		configuration.Sampler.Param = v.Float64()
	}

	configuration.Tags = append(configuration.Tags,
		opentracing.Tag{Key: "app.version", Value: info.Version},
		opentracing.Tag{Key: "app.namespace", Value: info.Namespace},
		opentracing.Tag{Key: "app.name", Value: info.AppName},
		opentracing.Tag{Key: "app.git_hash", Value: info.GitHash},
		opentracing.Tag{Key: "app.build_at", Value: info.BuildAt},
	)

	return configuration
}

func (t *tracer) Close() {
	_ = t.tracer.Close()
}
