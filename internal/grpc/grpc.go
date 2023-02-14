// Package grpc implements an gRPC server.
package grpc

import (
	"crypto/tls"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/loghole/tron/transport"
)

// A Server defines configuration for running an gRPC server.
type Server struct {
	addr     string
	listener net.Listener
	server   *grpc.Server
}

// NewServer returns gRPC server with port.
func NewServer(port uint16) *Server {
	if port == 0 {
		return nil
	}

	return &Server{addr: fmt.Sprintf("0.0.0.0:%d", port)}
}

// NewServerWithListener returns gRPC server with listener.
func NewServerWithListener(listener net.Listener) *Server {
	if listener == nil {
		return nil
	}

	return &Server{addr: listener.Addr().String(), listener: listener}
}

// BuildServer init gRPC server.
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

// RegistryDesc in gRPC server.
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

func (s *Server) IsPresent() bool {
	return s != nil
}

// Serve starts serving incoming connections.
func (s *Server) Serve() error {
	if s == nil {
		return nil
	}

	if err := s.server.Serve(s.listener); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

// Close closes the server.
func (s *Server) Close() error {
	if s == nil {
		return nil
	}

	if s.server != nil {
		s.server.GracefulStop()
	}

	return nil
}

// Addr returns the server address.
func (s *Server) Addr() string {
	if s == nil {
		return "-"
	}

	return fmt.Sprintf("grpc://%s", s.addr)
}
