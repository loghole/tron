package tron

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/internal/grpc"
	"github.com/loghole/tron/internal/http"
)

type servers struct {
	publicGRPC *grpc.Server
	publicHTTP *http.Server
	adminHTTP  *http.Server
}

func (s *servers) init(opts *app.Options) (err error) {
	if opts.PortAdmin == 0 {
		opts.PortAdmin = uint16(viper.GetInt32(app.AdminPortEnv))
	}

	if opts.PortHTTP == 0 {
		opts.PortHTTP = uint16(viper.GetInt32(app.HTTPPortEnv))
	}

	if opts.PortGRPC == 0 {
		opts.PortGRPC = uint16(viper.GetInt32(app.GRPCPortEnv))
	}

	s.publicGRPC = grpc.NewServer(opts.PortGRPC)

	s.publicHTTP = http.NewServer(opts.PortHTTP)

	s.adminHTTP = http.NewServer(opts.PortAdmin)

	return nil
}

func (s *servers) build(opts *app.Options) error {
	if err := s.publicGRPC.BuildServer(opts.TLSConfig, opts.GRPCOptions); err != nil {
		return fmt.Errorf("failed to build public grpc server: %w", err)
	}

	if err := s.publicHTTP.BuildServer(opts.TLSConfig); err != nil {
		return fmt.Errorf("failed to build public http server: %w", err)
	}

	if err := s.adminHTTP.BuildServer(nil); err != nil {
		return fmt.Errorf("failed to build admin http server: %w", err)
	}

	return nil
}
