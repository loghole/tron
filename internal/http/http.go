package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/transport"
)

type Server struct {
	addr     string
	listener net.Listener
	server   *http.Server
	router   *chi.Mux
}

func NewServer(port uint16) *Server {
	if port == 0 {
		return nil
	}

	return &Server{router: chi.NewRouter(), addr: fmt.Sprintf("0.0.0.0:%d", port)}
}

func (s *Server) BuildServer(tlsConfig *tls.Config) (err error) {
	if s == nil {
		return nil
	}

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

	mux := runtime.NewServeMux(runtime.WithErrorHandler(ErrorWriter()))

	for _, service := range services {
		if service != nil {
			service.GetDescription().RegisterHTTP(mux)
		}
	}

	s.router.Handle("/*", mux)
}

func (s *Server) Router() chi.Router {
	if s == nil {
		return nil
	}

	return s.router
}

func (s *Server) Serve() error {
	if s == nil {
		return nil
	}

	s.server = &http.Server{Handler: s.router}

	if len(s.router.Routes()) == 0 {
		s.router.HandleFunc("/", http.NotFound)
	}

	if err := s.server.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return simplerr.Wrap(err, "serve http failed")
	}

	return nil
}

func (s *Server) Close() error {
	if s == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)

	return s.server.Shutdown(ctx)
}

func (s *Server) Addr() string {
	if s == nil {
		return "-"
	}

	return fmt.Sprintf("http://%s", s.addr)
}
