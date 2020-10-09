package grpc

import (
	"fmt"
	"net"

	"github.com/lissteron/simplerr"
	"github.com/utrack/clay/v2/transport"
	"google.golang.org/grpc"
)

type Server struct {
	addr string
	lis  net.Listener
	*grpc.Server
}

func NewServer(port uint16, opts []grpc.ServerOption) (server *Server, err error) {
	if port == 0 {
		return nil, nil
	}

	server = &Server{Server: grpc.NewServer(opts...), addr: fmt.Sprintf("0.0.0.0:%d", port)}

	server.lis, err = net.Listen("tcp", server.addr)
	if err != nil {
		return nil, simplerr.Wrap(err, "create GRPc listener failed")
	}

	return server, nil
}

func (s *Server) RegistryDesc(services ...transport.Service) {
	if s == nil {
		return
	}

	for _, service := range services {
		if service != nil {
			service.GetDescription().RegisterGRPC(s.Server)
		}
	}
}

func (s *Server) Serve() error {
	if s == nil {
		return nil
	}

	return s.Server.Serve(s.lis)
}

func (s *Server) Close() error {
	if s == nil {
		return nil
	}

	if s.Server != nil {
		s.Server.GracefulStop()
	}

	return nil
}

func (s *Server) Addr() string {
	if s == nil {
		return "-"
	}

	return fmt.Sprintf("grpc://%s", s.addr)
}
