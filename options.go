package tron

import (
	"net"
	"os"

	"github.com/loghole/lhw/zap"

	"github.com/loghole/tron/internal/app"
)

// Option sets tron options such as ports, config, etc.
type Option = app.Option

// WithPublicHTTP returns a Option that sets public http port.
// Cannot be changed by config or env.
//
//  Example:
//
//  tron.New(tron.WithPublicGRPC(8080))
func WithPublicHTTP(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortHTTP = port

		return nil
	}
}

// WithAdminHTTP returns a Option that sets admin http port.
// Cannot be changed by config or env.
//
//  Example:
//
//  tron.New(tron.WithAdminHTTP(8081))
func WithAdminHTTP(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortAdmin = port

		return nil
	}
}

// WithPublicGRPC returns a Option that sets public grpc port.
// Cannot be changed by config or env.
//
//  Example:
//
//  tron.New(tron.WithPublicGRPC(8082))
func WithPublicGRPC(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortGRPC = port

		return nil
	}
}

// WithExitSignals returns a Option that sets exit signals for application.
// Default signals is: syscall.SIGTERM, syscall.SIGINT.
//
//  Example:
//
//  tron.New(tron.WithExitSignals(syscall.SIGKILL))
func WithExitSignals(sig ...os.Signal) Option {
	return func(opts *app.Options) error {
		opts.ExitSignals = append(opts.ExitSignals, sig...)

		return nil
	}
}

// WithConfigMap returns a Option that init app config from map and envs.
//
//  Example:
//
//  tron.New(tron.WithConfigMap(map[string]interface{}{
//  	"namespace":         "dev",
//  	"service_port_grpc": 35900,
//  	"cockroach_addr":    "db_addr",
//  	"cockroach_user":    "db_user",
//  	"cockroach_db":      "db_name",
//  }))
func WithConfigMap(cfg map[string]interface{}) Option {
	return func(opts *app.Options) error {
		opts.ConfigMap = cfg

		return nil
	}
}

// WithGRPCListener returns a Option that sets net listener for grpc public server.
// Can be used for create application tests with memory listener.
//
//  Example:
//
//  listener := bufconn.Listen(1024*1024)
//
//  tron.New(tron.WithGRPCListener(listener))
func WithGRPCListener(listener net.Listener) Option {
	return func(opts *app.Options) error {
		opts.GRPCListener = listener

		return nil
	}
}

// WithHTTPListener returns a Option that sets net listener for http public server.
// Can be used for create application tests with memory listener.
//
//  Example:
//
//  listener = bufconn.Listen(1024*1024)
//
//  tron.New(tron.WithHTTPListener(listener))
func WithHTTPListener(listener net.Listener) Option {
	return func(opts *app.Options) error {
		opts.HTTPListener = listener

		return nil
	}
}

// AddLogCaller configures the Logger to annotate each message with the filename
// and line number of zap's caller.
func AddLogCaller() Option {
	return func(opts *app.Options) error {
		opts.LoggerOptions = append(opts.LoggerOptions, zap.AddCaller())

		return nil
	}
}

// AddLogStacktrace configures the Logger to record a stack trace for all messages at
// or above a given level.
//
//  Example:
//
//  tron.New(tron.AddLogStacktrace("error"))
func AddLogStacktrace(level string) Option {
	return func(opts *app.Options) error {
		opts.LoggerOptions = append(opts.LoggerOptions, zap.AddStacktrace(level))

		return nil
	}
}

// WithLogField adds field to the Logger.
//
//  Example:
//
//  tron.New(tron.WithLogField("my_field_key", "my_field_value"))
func WithLogField(key string, value interface{}) Option {
	return func(opts *app.Options) error {
		opts.LoggerOptions = append(opts.LoggerOptions, zap.WithField(key, value))

		return nil
	}
}

// WithRealtimeConfig returns a Option that activate realtime config.
// Realtime config watch only values file for current namespace.
//
//  Example:
//
//  tron.New(tron.WithRealtimeConfig())
func WithRealtimeConfig() Option {
	return func(opts *app.Options) error {
		opts.RealtimeConfig = true

		return nil
	}
}
