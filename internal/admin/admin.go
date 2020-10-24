package admin

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	jsoniter "github.com/json-iterator/go"
	"github.com/utrack/clay/v2/transport"
	"github.com/utrack/clay/v2/transport/swagger"

	"github.com/loghole/tron/internal/app"
)

const adminToPublicPort = 2

type Handlers struct {
	desc transport.ServiceDesc
	info *app.Info
}

func NewHandlers(info *app.Info, services ...transport.Service) *Handlers {
	descs := make([]transport.ServiceDesc, 0, len(services))

	for _, service := range services {
		descs = append(descs, service.GetDescription())
	}

	return &Handlers{desc: transport.NewCompoundServiceDesc(descs...), info: info}
}

func (s *Handlers) InitRoutes(r chi.Router) {
	if r == nil {
		return
	}

	r.Get("/info", s.serviceInfoHandler)

	r.Mount("/debug", middleware.Profiler())

	r.Mount("/docs", http.StripPrefix("/docs", http.FileServer(AssetFile())))

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	r.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger.json", http.StatusMovedPermanently)
	})

	r.HandleFunc("/swagger.json", s.swaggerDefHandler)
}

func (s *Handlers) serviceInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"InstanceUUID": s.info.InstanceUUID,
		"ServiceName":  s.info.ServiceName,
		"Namespace":    s.info.Namespace,
		"AppName":      s.info.AppName,
		"GitHash":      s.info.GitHash,
		"Version":      s.info.Version,
		"BuildAt":      s.info.BuildAt,
		"StartTime":    s.info.StartTime,
		"UpTime":       time.Since(s.info.StartTime),
	}

	_ = jsoniter.NewEncoder(w).Encode(info)
}

func (s *Handlers) swaggerDefHandler(w http.ResponseWriter, r *http.Request) {
	if host, port, err := net.SplitHostPort(r.Host); err == nil {
		if port, err := strconv.Atoi(port); err == nil {
			r.Host = net.JoinHostPort(host, strconv.Itoa(port-adminToPublicPort))
		}
	}

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(s.desc.SwaggerDef(
		swagger.WithVersion(s.info.Version),
		swagger.WithTitle(s.info.AppName),
		swagger.WithHost(r.Host),
	))
}
