package app

import (
	"crypto/tls"
	"net/http"
	"os"
	"syscall"

	"github.com/loghole/lhw/zap"
	"google.golang.org/grpc"
)

type Option func(opts *Options) error

type RunOption func(opts *Options)

type Options struct {
	// New options.
	Hostname      string
	PortAdmin     uint16
	PortHTTP      uint16
	PortGRPC      uint16
	TLSConfig     *tls.Config
	LoggerOptions []zap.Option
	ExitSignals   []os.Signal

	// Run options.
	GRPCOptions     []grpc.ServerOption
	HTTPMiddlewares []func(http.Handler) http.Handler
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
	}

	for _, apply := range options {
		if err := apply(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func (o *Options) ApplyRunOptions(options ...RunOption) {
	for _, apply := range options {
		apply(o)
	}
}
