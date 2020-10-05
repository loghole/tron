package app

import (
	"crypto/tls"
	"os"
	"syscall"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Option func(opts *Options) error

type Options struct {
	ConfigName string

	Hostname  string
	PortAdmin uint16
	PortHTTP  uint16
	PortGRPC  uint16

	TLSConfig *tls.Config

	GRPCOptions []grpc.ServerOption

	ExitSignals []os.Signal
}

func NewOptions(options ...Option) (*Options, error) {
	parseFlags()

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	opts := &Options{
		ConfigName:  "dev",
		Hostname:    hostname,
		PortAdmin:   uint16(viper.GetInt32(AdminPortEnv)),
		PortHTTP:    uint16(viper.GetInt32(HTTPPortEnv)),
		PortGRPC:    uint16(viper.GetInt32(GRPCPortEnv)),
		ExitSignals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
	}

	for _, apply := range options {
		if err := apply(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}
