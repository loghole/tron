package tron

import (
	"fmt"

	"github.com/loghole/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"

	"github.com/loghole/tron/internal/app"
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

func (t *tracer) configuration(info *Info) *config.Configuration {
	configuration := tracing.DefaultConfiguration(
		info.AppName,
		viper.GetString(app.JaegerAddrEnv),
	)

	if v := viper.GetString(app.JaegerSamplerType); v != "" {
		configuration.Sampler.Type = v
	}

	if v := viper.GetFloat64(app.JaegerSamplerParam); v != 0 {
		configuration.Sampler.Param = v
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
