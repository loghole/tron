package app

import (
	"crypto/tls"
	"os"
	"syscall"

	"github.com/loghole/lhw/zap"
	"google.golang.org/grpc"
)

type Option func(opts *Options) error

type Options struct {
	Hostname  string
	PortAdmin uint16
	PortHTTP  uint16
	PortGRPC  uint16
	TLSConfig *tls.Config

	LoggerOptions []zap.Option
	GRPCOptions   []grpc.ServerOption
	ExitSignals   []os.Signal
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
