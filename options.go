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

func WithPublicGRPC(port uint16) Option {
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

// WithConfigMap init app with config from map and envs.
func WithConfigMap(cfg map[string]interface{}) Option {
	return func(opts *app.Options) error {
		opts.ConfigMap = cfg

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

// WithCORSAllowedOrigins sets allowed origin domains for cross-domain requests.
// If the special "*" value is present in the list, all origins will be allowed.
// An origin may contain a wildcard (*) to replace 0 or more characters
// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality.
// Only one wildcard can be used per origin.
// Default option enable CORS for requests on admin port.
func WithCORSAllowedOrigins(origins []string) Option {
	return func(opts *app.Options) error {
		opts.CorsOptions.AllowedOrigins = origins

		return nil
	}
}

// WithCORSAllowedAuthentication sets if HTTP client allows authentication like
// cookies, SSL certs and HTTP basic auth.
func WithCORSAllowedAuthentication(allow bool) Option {
	return func(opts *app.Options) error {
		opts.CorsOptions.AllowCredentials = allow

		return nil
	}
}

// WithCORSAllowedHeaders sets a list of non simple headers the HTTP
// client is allowed to use with cross-domain requests.
// If the special "*" value is present in the list, all headers will be allowed.
// Default value is [] but "Origin" is always appended to the list.
func WithCORSAllowedHeaders(headers []string) Option {
	return func(opts *app.Options) error {
		opts.CorsOptions.AllowedHeaders = headers

		return nil
	}
}
