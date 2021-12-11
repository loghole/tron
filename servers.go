package tron

import (
	"fmt"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/internal/grpc"
	"github.com/loghole/tron/internal/http"
	"github.com/loghole/tron/rtconfig"
)

type servers struct {
	publicGRPC *grpc.Server
	publicHTTP *http.Server
	adminHTTP  *http.Server
}

func (s *servers) init(opts *app.Options) (err error) {
	s.initPubGRPC(opts)
	s.initPubHTTP(opts)
	s.initAdmHTTP(opts)

	return nil
}

func (s *servers) initPubGRPC(opts *app.Options) {
	if opts.GRPCListener != nil {
		s.publicGRPC = grpc.NewServerWithListener(opts.GRPCListener)

		return
	}

	if opts.PortGRPC == 0 {
		opts.PortGRPC = uint16(rtconfig.GetInt32(app.GRPCPortEnv))
	}

	s.publicGRPC = grpc.NewServer(opts.PortGRPC)
}

func (s *servers) initPubHTTP(opts *app.Options) {
	if opts.HTTPListener != nil {
		s.publicHTTP = http.NewServerWithListener(opts.GRPCListener)

		return
	}

	if opts.PortHTTP == 0 {
		opts.PortHTTP = uint16(rtconfig.GetInt32(app.HTTPPortEnv))
	}

	s.publicHTTP = http.NewServer(opts.PortHTTP)
}

func (s *servers) initAdmHTTP(opts *app.Options) {
	if opts.PortAdmin == 0 {
		opts.PortAdmin = uint16(rtconfig.GetInt32(app.AdminPortEnv))
	}

	s.adminHTTP = http.NewServer(opts.PortAdmin)
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
