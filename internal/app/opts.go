package app

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/loghole/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// Option sets tron options such as ports, config, etc.
type Option func(opts *Options) error

// RunOption sets tron run options such as grpc unary interceptors, tls config.
type RunOption func(opts *Options) error

// Options is base tron options.
type Options struct {
	// New options.
	Hostname      string
	PortAdmin     uint16
	PortHTTP      uint16
	PortGRPC      uint16
	LoggerOptions []zap.Option
	LoggerConfig  zap.Config
	TracerConfig  *tracing.Configuration
	ExitSignals   []os.Signal
	ConfigMap     map[string]interface{}
	GRPCListener  net.Listener
	HTTPListener  net.Listener

	// Run options.
	TLSConfig   *tls.Config
	GRPCOptions []grpc.ServerOption

	options []RunOption
}

// NewOptions returns Options with applied Option list.
func NewOptions(info *Info, options ...Option) (*Options, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname failed: %w", err)
	}

	opts := &Options{
		Hostname:     hostname,
		ExitSignals:  []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		LoggerConfig: zap.NewProductionConfig(),
		TracerConfig: tracing.DefaultConfiguration(info.AppName, ""),
	}

	opts.LoggerConfig.InitialFields = map[string]interface{}{
		"host":         info.InstanceUUID,
		"source":       info.ServiceName,
		"version":      info.Version,
		"build_commit": info.GitHash,
	}

	opts.LoggerConfig.EncoderConfig = defaultZapEncoderConfig()

	opts.TracerConfig.Attributes = append(opts.TracerConfig.Attributes,
		attribute.String("app.version", info.Version),
		attribute.String("app.name", info.AppName),
		attribute.String("app.git_hash", info.GitHash),
		attribute.String("app.build_at", info.BuildAt),
	)

	for _, apply := range options {
		if apply == nil {
			continue
		}

		if err := apply(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

// AddRunOptions append run options to others options.
func (o *Options) AddRunOptions(options ...RunOption) {
	o.options = append(o.options, options...)
}

// ApplyRunOptions sets run options.
func (o *Options) ApplyRunOptions() error {
	for _, apply := range o.options {
		if apply == nil {
			continue
		}

		if err := apply(o); err != nil {
			return err
		}
	}

	o.options = nil

	return nil
}

func defaultZapEncoderConfig() zapcore.EncoderConfig {
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
