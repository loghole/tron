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
	"github.com/lissteron/simplerr"
	"github.com/utrack/clay/v2/transport"
)

type Server struct {
	addr string
	lis  net.Listener
	rout chi.Router
	*http.Server
}

func NewServer(port uint16, tlsConfig *tls.Config) (server *Server, err error) {
	if port == 0 {
		return nil, nil
	}

	server = &Server{rout: chi.NewRouter(), addr: fmt.Sprintf("0.0.0.0:%d", port)}

	switch {
	case tlsConfig != nil:
		server.lis, err = tls.Listen("tcp", server.addr, tlsConfig)
		if err != nil {
			return nil, simplerr.Wrap(err, "create TLS listener failed")
		}
	default:
		server.lis, err = net.Listen("tcp", server.addr)
		if err != nil {
			return nil, simplerr.Wrap(err, "create TCP listener failed")
		}
	}

	return server, nil
}

func (s *Server) RegistryDesc(services ...transport.Service) {
	if s == nil {
		return
	}

	for _, service := range services {
		if service != nil {
			service.GetDescription().RegisterHTTP(s.rout)
		}
	}
}

func (s *Server) Router() chi.Router {
	if s == nil {
		return nil
	}

	return s.rout
}

func (s *Server) Serve() error {
	if s == nil {
		return nil
	}

	s.Server = &http.Server{Handler: s.rout}

	if len(s.rout.Routes()) == 0 {
		s.rout.HandleFunc("/", http.NotFound)
	}

	if err := s.Server.Serve(s.lis); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Close() error {
	if s == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.Server.SetKeepAlivesEnabled(false)

	return s.Server.Shutdown(ctx)
}

func (s *Server) Addr() string {
	if s == nil {
		return "-"
	}

	return fmt.Sprintf("http://%s", s.addr)
}
