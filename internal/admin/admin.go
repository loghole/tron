package admin

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-openapi/spec"
	jsoniter "github.com/json-iterator/go"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/transport"
)

// Handlers contains http methods with debug service info and swagger docs.
type Handlers struct {
	desc jsoniter.RawMessage
	info *app.Info
	opts *app.Options
}

// NewHandlers create and init handlers object.
func NewHandlers(info *app.Info, opts *app.Options, services ...transport.Service) *Handlers {
	descs := make([]transport.ServiceDesc, 0, len(services))

	for _, service := range services {
		descs = append(descs, service.GetDescription())
	}

	handlers := &Handlers{
		desc: transport.NewCompoundServiceDesc(descs...).SwaggerDef(),
		info: info,
		opts: opts,
	}

	return handlers
}

// InitRoutes init routes for current router with debug service info and swagger docs.
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
	info := struct {
		InstanceUUID string `json:"instance_uuid"`
		ServiceName  string `json:"service_name"`
		Namespace    string `json:"namespace"`
		AppName      string `json:"app_name"`
		GitHash      string `json:"git_hash"`
		Version      string `json:"version"`
		BuildAt      string `json:"build_at"`
		StartTime    string `json:"start_time"`
		UpTime       string `json:"up_time"`
	}{
		InstanceUUID: s.info.InstanceUUID,
		ServiceName:  s.info.ServiceName,
		Namespace:    s.info.Namespace,
		AppName:      s.info.AppName,
		GitHash:      s.info.GitHash,
		Version:      s.info.Version,
		BuildAt:      s.info.BuildAt,
		StartTime:    s.info.StartTime.String(),
		UpTime:       time.Since(s.info.StartTime).String(),
	}

	_ = jsoniter.NewEncoder(w).Encode(info)
}

func (s *Handlers) swaggerDefHandler(w http.ResponseWriter, r *http.Request) {
	if host, _, err := net.SplitHostPort(r.Host); err == nil {
		r.Host = net.JoinHostPort(host, strconv.Itoa(int(s.opts.PortHTTP)))
	}

	var desc spec.Swagger

	_ = jsoniter.Unmarshal(s.desc, &desc)
	desc.Host = r.Host
	desc.Info = &spec.Info{}
	desc.Info.Version = s.info.Version
	desc.Info.Title = s.info.AppName

	_ = jsoniter.NewEncoder(w).Encode(desc)
}
