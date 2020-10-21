package admin

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	jsoniter "github.com/json-iterator/go"
	"github.com/utrack/clay/v2/transport"
	"github.com/utrack/clay/v2/transport/swagger"

	"github.com/loghole/tron/internal/app"
)

type Handlers struct {
	desc []byte
	info *app.Info
}

func NewHandlers(info *app.Info, services ...transport.Service) *Handlers {
	descs := make([]transport.ServiceDesc, 0, len(services))

	for _, service := range services {
		descs = append(descs, service.GetDescription())
	}

	return &Handlers{desc: transport.NewCompoundServiceDesc(descs...).SwaggerDef(
		swagger.WithVersion(info.Version), swagger.WithTitle(info.AppName),
	)}
}

func (s *Handlers) InitRoutes(r chi.Router) {
	if r == nil {
		return
	}

	r.Mount("/debug", middleware.Profiler())

	r.Mount("/docs", http.StripPrefix("/docs", http.FileServer(AssetFile())))

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(s.desc)
	})

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	r.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger.json", http.StatusMovedPermanently)
	})

	r.Get("/info", s.writeServiceInfo)
}

func (s *Handlers) writeServiceInfo(w http.ResponseWriter, r *http.Request) {
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
