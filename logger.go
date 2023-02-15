package tron

import (
	"fmt"

	"github.com/loghole/tracing/tracelog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/loghole/tron/internal/app"
)

type logger struct {
	*zap.SugaredLogger
	tracelog tracelog.Logger
}

func (l *logger) init(opts *app.Options) (err error) {
	zlog, err := opts.LoggerConfig.Build(opts.LoggerOptions...)
	if err != nil {
		return fmt.Errorf("build zap logger: %w", err)
	}

	l.SugaredLogger = zlog.Sugar()
	l.tracelog = tracelog.NewTraceLogger(l.SugaredLogger)

	return nil
}

func parseZapLevel(s string) zapcore.Level {
	var lvl zapcore.Level

	if err := lvl.Set(s); err != nil {
		return zapcore.InfoLevel
	}

	return lvl
}
