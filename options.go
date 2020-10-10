package tron

import (
	"os"

	"github.com/loghole/lhw/zap"

	"github.com/loghole/tron/internal/app"
)

type Option = app.Option

func WithAdminHTTP(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortAdmin = port

		return nil
	}
}

func WithPublicHTTP(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortHTTP = port

		return nil
	}
}

func WithPublicGRPc(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortGRPC = port

		return nil
	}
}

func WithExitSignals(sig ...os.Signal) Option {
	return func(opts *app.Options) error {
		opts.ExitSignals = append(opts.ExitSignals, sig...)

		return nil
	}
}

// AddLogCaller configures the Logger to annotate each message with the filename
// and line number of zap's caller.  See also WithCaller.
func AddLogCaller() Option {
	return func(opts *app.Options) error {
		opts.LoggerOptions = append(opts.LoggerOptions, zap.AddCaller())

		return nil
	}
}

// AddLogStacktrace configures the Logger to record a stack trace for all messages at
// or above a given level.
func AddLogStacktrace(level string) Option {
	return func(opts *app.Options) error {
		opts.LoggerOptions = append(opts.LoggerOptions, zap.AddStacktrace(level))

		return nil
	}
}

// WithLogField adds field to the Logger.
func WithLogField(key string, value interface{}) Option {
	return func(opts *app.Options) error {
		opts.LoggerOptions = append(opts.LoggerOptions, zap.WithField(key, value))

		return nil
	}
}
