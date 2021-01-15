package tron

import (
	"fmt"

	"github.com/loghole/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

type tracer struct {
	tracer *tracing.Tracer
}

func (t *tracer) init(info *Info) (err error) {
	configuration := tracing.DefaultConfiguration(
		info.AppName,
		viper.GetString(app.JaegerAddrEnv),
	)

	configuration.Tags = append(configuration.Tags,
		opentracing.Tag{Key: "app.version", Value: info.Version},
		opentracing.Tag{Key: "app.namespace", Value: info.Namespace},
		opentracing.Tag{Key: "app.name", Value: info.AppName},
		opentracing.Tag{Key: "app.git_hash", Value: info.GitHash},
		opentracing.Tag{Key: "app.build_at", Value: info.BuildAt},
	)

	t.tracer, err = tracing.NewTracer(configuration)
	if err != nil {
		return fmt.Errorf("init tracer failed: %w", err)
	}

	if err := tracing.EnablePrometheusMetrics(); err != nil {
		return fmt.Errorf("init prometheus metrics failed: %w", err)
	}

	return nil
}

func (t *tracer) Close() {
	_ = t.tracer.Close()
}
