package tron

import (
	"fmt"

	"github.com/loghole/tracing/tracelog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/rtconfig"
)

type logger struct {
	level zap.AtomicLevel

	*zap.SugaredLogger
	tracelog tracelog.Logger
}

func (l *logger) init(info *Info, opts *app.Options) (err error) {
	l.level = zap.NewAtomicLevelAt(parseZapLevel(rtconfig.GetString(app.LoggerLevelEnv)))

	var cfg zap.Config

	if info.Namespace == app.NamespaceLocal {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig = l.encoderConfig()
	}

	cfg.Level = l.level
	cfg.DisableStacktrace = true
	cfg.InitialFields = map[string]interface{}{
		"host":         info.InstanceUUID,
		"namespace":    info.Namespace,
		"source":       info.ServiceName,
		"version":      info.Version,
		"build_commit": info.GitHash,
	}

	logger, err := cfg.Build(opts.LoggerOptions...)
	if err != nil {
		return fmt.Errorf("build zap logger: %w", err)
	}

	l.SugaredLogger = logger.Sugar()
	l.tracelog = tracelog.NewTraceLogger(l.SugaredLogger)

	if err := rtconfig.WatchVariable(app.LoggerLevelEnv, l.levelWatcher); err != nil {
		return fmt.Errorf("start watch log level: %w", err)
	}

	return nil
}

func (l *logger) levelWatcher(oldValue, newValue rtconfig.Value) {
	if newValue.IsNil() {
		return
	}

	if oldValue.String() == newValue.String() {
		return
	}

	var lvl zapcore.Level

	if err := lvl.Set(newValue.String()); err != nil {
		l.Errorf("invalid level value '%s'", newValue.String())

		return
	}

	l.level.SetLevel(parseZapLevel(newValue.String()))
	l.Errorf("update log level to: '%s'", newValue.String())
}

func (l *logger) encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func parseZapLevel(s string) zapcore.Level {
	var lvl zapcore.Level

	if err := lvl.Set(s); err != nil {
		return zapcore.InfoLevel
	}

	return lvl
}
