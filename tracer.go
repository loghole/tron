package tron

import (
	"fmt"

	"github.com/loghole/tracing"
	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

type tracer struct {
	tracer *tracing.Tracer
}

func (t *tracer) init(info *Info) (err error) {
	t.tracer, err = tracing.NewTracer(tracing.DefaultConfiguration(
		info.ServiceName,
		viper.GetString(app.JaegerAddrEnv)),
	)
	if err != nil {
		return fmt.Errorf("init tracer failed: %w", err)
	}

	return nil
}

func (t *tracer) Close() {
	_ = t.tracer.Close()
}
