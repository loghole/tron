package app

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/go-chi/cors"
	"github.com/loghole/lhw/zap"
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
	ExitSignals   []os.Signal
	CorsOptions   cors.Options
	ConfigMap     map[string]interface{}
	GRPCListener  net.Listener
	HTTPListener  net.Listener

	// Run options.
	TLSConfig   *tls.Config
	GRPCOptions []grpc.ServerOption

	options []RunOption
}

// NewOptions returns Options with applied Option list.
func NewOptions(options ...Option) (*Options, error) {
	parseFlags()

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname failed: %w", err)
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
