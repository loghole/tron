package tron

import (
	"github.com/loghole/lhw/zap"
	trace "github.com/loghole/tracing/tracelog"
	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

type logger struct {
	*zap.Logger
	tracelog trace.Logger
}

func (l *logger) init(info *Info, opts *app.Options) (err error) {
	l.Logger, err = zap.NewLogger(&zap.Config{
		Level:         viper.GetString(app.LoggerLevelEnv),
		CollectorURL:  viper.GetString(app.LoggerCollectorAddrEnv),
		Hostname:      opts.Hostname,
		Namespace:     info.Namespace,
		Source:        info.ServiceName,
		BuildCommit:   info.GitHash,
		DisableStdout: viper.GetBool(app.LoggerDisableStdoutEnv),
	})
	if err != nil {
		return err
	}

	l.tracelog = trace.NewTraceLogger(l.Logger.SugaredLogger)

	return nil
}
