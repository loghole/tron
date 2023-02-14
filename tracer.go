package tron

import (
	"fmt"

	"github.com/loghole/tracing"

	"github.com/loghole/tron/internal/app"
)

type tracer struct {
	tracer *tracing.Tracer
}

func (t *tracer) init(options *app.Options) (err error) {
	t.tracer, err = tracing.NewTracer(options.TracerConfig)
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
