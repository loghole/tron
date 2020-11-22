package grpc

import (
	"crypto/tls"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/loghole/tron/transport"
)

type Server struct {
	addr     string
	listener net.Listener
	server   *grpc.Server
}

func NewServer(port uint16) *Server {
	if port == 0 {
		return nil
	}

	return &Server{addr: fmt.Sprintf("0.0.0.0:%d", port)}
}

func NewServerWithListener(listener net.Listener) *Server {
	if listener == nil {
		return nil
	}

	return &Server{addr: listener.Addr().String(), listener: listener}
}

func (s *Server) BuildServer(tlsConfig *tls.Config, opts []grpc.ServerOption) (err error) {
	if s == nil {
		return nil
	}

	s.server = grpc.NewServer(opts...)

	if s.listener != nil {
		return nil
	}

	switch {
	case tlsConfig != nil:
		s.listener, err = tls.Listen("tcp", s.addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("create TLS listener failed: %w", err)
		}
	default:
		s.listener, err = net.Listen("tcp", s.addr)
		if err != nil {
			return fmt.Errorf("create TCP listener failed: %w", err)
		}
	}

	return nil
}

func (s *Server) RegistryDesc(services ...transport.Service) {
	if s == nil {
		return
	}

	for _, service := range services {
		if service != nil {
			service.GetDescription().RegisterGRPC(s.server)
		}
	}
}

func (s *Server) Serve() error {
	if s == nil {
		return nil
	}

	return s.server.Serve(s.listener)
}

func (s *Server) Close() error {
	if s == nil {
		return nil
	}

	if s.server != nil {
		s.server.GracefulStop()
	}

	return nil
}

func (s *Server) Addr() string {
	if s == nil {
		return "-"
	}

	return fmt.Sprintf("grpc://%s", s.addr)
}
