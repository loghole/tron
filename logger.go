package tron

import (
	"fmt"

	"github.com/loghole/lhw/zaplog"
	"github.com/loghole/tracing/tracelog"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/rtconfig"
)

type logger struct {
	*zaplog.Logger
	tracelog tracelog.Logger
}

func (l *logger) init(info *Info, opts *app.Options) (err error) {
	l.Logger, err = zaplog.NewLogger(&zaplog.Config{
		Level:         rtconfig.GetString(app.LoggerLevelEnv),
		CollectorURL:  rtconfig.GetString(app.LoggerCollectorAddrEnv),
		Hostname:      opts.Hostname,
		Namespace:     info.Namespace,
		Source:        info.ServiceName,
		BuildCommit:   info.GitHash,
		DisableStdout: rtconfig.GetBool(app.LoggerDisableStdoutEnv),
	}, opts.LoggerOptions...)
	if err != nil {
		return fmt.Errorf("init logger failed: %w", err)
	}

	l.tracelog = tracelog.NewTraceLogger(l.Logger.SugaredLogger)

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

	l.Logger.SetLevel(newValue.String())
}
