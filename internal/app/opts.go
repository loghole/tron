package app

import (
	"crypto/tls"
	"os"
	"syscall"

	"github.com/go-chi/cors"
	"github.com/loghole/lhw/zap"
	"google.golang.org/grpc"
)

type Option func(opts *Options) error

type RunOption func(opts *Options) error

type Options struct {
	// New options.
	Hostname      string
	PortAdmin     uint16
	PortHTTP      uint16
	PortGRPC      uint16
	LoggerOptions []zap.Option
	ExitSignals   []os.Signal
	CorsOptions   cors.Options

	// Run options.
	TLSConfig   *tls.Config
	GRPCOptions []grpc.ServerOption

	options []RunOption
}

func NewOptions(options ...Option) (*Options, error) {
	parseFlags()

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	opts := &Options{
		Hostname:    hostname,
		ExitSignals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		CorsOptions: cors.Options{},
	}

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

func (o *Options) AddRunOptions(options ...RunOption) {
	o.options = append(o.options, options...)
}

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
