// Package http implements an HTTP server.
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

	"github.com/loghole/tron/transport"
)

// A Server defines configuration for running an HTTP server.
type Server struct {
	addr     string
	listener net.Listener
	server   *http.Server
	router   *chi.Mux
}

// NewServer returns http server with port.
func NewServer(port uint16) *Server {
	if port == 0 {
		return nil
	}

	return &Server{router: chi.NewRouter(), addr: fmt.Sprintf("0.0.0.0:%d", port)}
}

// NewServerWithListener returns http server with listener.
func NewServerWithListener(listener net.Listener) *Server {
	if listener == nil {
		return nil
	}

	return &Server{router: chi.NewRouter(), addr: listener.Addr().String(), listener: listener}
}

// BuildServer init http server.
func (s *Server) BuildServer(tlsConfig *tls.Config) (err error) {
	if s == nil {
		return nil
	}

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

// RegistryDesc in http server.
func (s *Server) RegistryDesc(services ...transport.Service) {
	if s == nil {
		return
	}

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(ErrorWriter()),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, newMarshaler()),
	)

	for _, service := range services {
		if service != nil {
			service.GetDescription().RegisterHTTP(mux)
		}
	}

	s.router.Handle("/*", mux)
}

// UseMiddleware appends a middleware handler to the Mux middleware stack.
func (s *Server) UseMiddleware(middlewares ...func(http.Handler) http.Handler) {
	if s == nil {
		return
	}

	s.router.Use(middlewares...)
}

// Router returns http router.
func (s *Server) Router() chi.Router {
	if s == nil {
		return nil
	}

	return s.router
}

// Serve starts serving incoming connections.
func (s *Server) Serve() error {
	if s == nil {
		return nil
	}

	s.server = &http.Server{Handler: s.router}

	if len(s.router.Routes()) == 0 {
		s.router.HandleFunc("/", http.NotFound)
	}

	if err := s.server.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http failed: %w", err)
	}

	return nil
}

// Close closes the server.
func (s *Server) Close() error {
	if s == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)

	return s.server.Shutdown(ctx)
}

// Addr returns the server address.
func (s *Server) Addr() string {
	if s == nil {
		return "-"
	}

	return fmt.Sprintf("http://%s", s.addr)
}
