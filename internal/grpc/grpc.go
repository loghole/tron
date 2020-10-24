package grpc

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/lissteron/simplerr"
	"github.com/utrack/clay/v2/transport"
	"google.golang.org/grpc"
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

func (s *Server) BuildServer(tlsConfig *tls.Config, opts []grpc.ServerOption) (err error) {
	if s == nil {
		return nil
	}

	s.server = grpc.NewServer(opts...)

	switch {
	case tlsConfig != nil:
		s.listener, err = tls.Listen("tcp", s.addr, tlsConfig)
		if err != nil {
			return simplerr.Wrap(err, "create TLS listener failed")
		}
	default:
		s.listener, err = net.Listen("tcp", s.addr)
		if err != nil {
			return simplerr.Wrap(err, "create TCP listener failed")
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
