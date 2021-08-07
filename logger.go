package tron

import (
	"fmt"

	"github.com/loghole/lhw/zaplog"
	trace "github.com/loghole/tracing/tracelog"
	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

type logger struct {
	*zaplog.Logger
	tracelog trace.Logger
}

func (l *logger) init(info *Info, opts *app.Options) (err error) {
	l.Logger, err = zaplog.NewLogger(&zaplog.Config{
		Level:         viper.GetString(app.LoggerLevelEnv),
		CollectorURL:  viper.GetString(app.LoggerCollectorAddrEnv),
		Hostname:      opts.Hostname,
		Namespace:     info.Namespace,
		Source:        info.ServiceName,
		BuildCommit:   info.GitHash,
		DisableStdout: viper.GetBool(app.LoggerDisableStdoutEnv),
	}, opts.LoggerOptions...)
	if err != nil {
		return fmt.Errorf("init logger failed: %w", err)
	}

	l.tracelog = trace.NewTraceLogger(l.Logger.SugaredLogger)

	return nil
}
