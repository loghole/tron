package tron

import (
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/internal/grpc"
	"github.com/loghole/tron/internal/http"
)

type servers struct {
	logger logger

	publicGRPC *grpc.Server
	publicHTTP *http.Server
	adminHTTP  *http.Server
}

func (s *servers) init(log logger, opts *app.Options) (err error) {
	s.logger = log
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

	s.publicGRPC = grpc.NewServer(opts.PortGRPC)
}

func (s *servers) initPubHTTP(opts *app.Options) {
	if opts.HTTPListener != nil {
		s.publicHTTP = http.NewServerWithListener(opts.GRPCListener)

		return
	}

	s.publicHTTP = http.NewServer(opts.PortHTTP)
}

func (s *servers) initAdmHTTP(opts *app.Options) {
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

func (s *servers) serve(group *errgroup.Group) {
	if s.publicGRPC.IsPresent() {
		group.Go(func() error {
			s.logger.Infof("grpc.public: start server on: %s", s.publicGRPC.Addr())
			defer s.logger.Warn("grpc.public: server stopped")

			return s.publicGRPC.Serve() //nolint:wrapcheck // need clean err
		})
	}

	if s.publicHTTP.IsPresent() {
		group.Go(func() error {
			s.logger.Infof("http.public: start server on: %s", s.publicHTTP.Addr())
			defer s.logger.Warn("http.public: server stopped")

			return s.publicHTTP.Serve() //nolint:wrapcheck // need clean err
		})
	}

	if s.adminHTTP.IsPresent() {
		group.Go(func() error {
			s.logger.Infof("http.admin: start server on: %s", s.adminHTTP.Addr())
			defer s.logger.Warn("http.admin: server stopped")

			return s.adminHTTP.Serve() //nolint:wrapcheck // need clean err
		})
	}
}

func (s *servers) close() {
	if err := s.publicHTTP.Close(); err != nil {
		s.logger.Errorf("error while stopping public http server: %v", err)
	}

	if err := s.publicGRPC.Close(); err != nil {
		s.logger.Errorf("error while stopping public grpc server: %v", err)
	}

	if err := s.adminHTTP.Close(); err != nil {
		s.logger.Errorf("error while stopping admin http server: %v", err)
	}
}
