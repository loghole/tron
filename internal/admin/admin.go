package admin

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/loghole/tron/internal/admin/swagger"
	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/transport"
)

// Health contains handlers to check app is healthy.
type Health interface {
	LivenessHandler(w http.ResponseWriter, r *http.Request)
	ReadinessHandler(w http.ResponseWriter, r *http.Request)
}

// Handlers contains http methods with debug service info and swagger docs.
type Handlers struct {
	desc   json.RawMessage
	info   *app.Info
	opts   *app.Options
	health Health
}

// NewHandlers create and init handlers object.
func NewHandlers(info *app.Info, opts *app.Options, health Health, services ...transport.Service) *Handlers {
	descs := make([]transport.ServiceDesc, 0, len(services))

	for _, service := range services {
		descs = append(descs, service.GetDescription())
	}

	handlers := &Handlers{
		desc:   transport.NewCompoundServiceDesc(descs...).SwaggerDef(),
		info:   info,
		opts:   opts,
		health: health,
	}

	return handlers
}

// InitRoutes init routes for current router with debug service info and swagger docs.
func (s *Handlers) InitRoutes(router chi.Router) {
	if router == nil {
		return
	}

	router.HandleFunc("/", s.index)

	router.Handle("/metrics", promhttp.Handler())

	router.Get("/info", s.serviceInfoHandler)

	router.Mount("/debug", middleware.Profiler())

	router.Mount("/docs", http.StripPrefix("/docs", http.FileServer(http.FS(swagger.Content))))

	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	router.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger.json", http.StatusMovedPermanently)
	})

	router.HandleFunc("/swagger.json", s.swaggerDefHandler)

	router.Get("/heath/live", s.health.LivenessHandler)
	router.Get("/heath/ready", s.health.ReadinessHandler)
}

func (s *Handlers) serviceInfoHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	//nolint:errchkjson // not need check error.
	_ = encoder.Encode(info{
		InstanceUUID: s.info.InstanceUUID,
		ServiceName:  s.info.ServiceName,
		AppName:      s.info.AppName,
		GitHash:      s.info.GitHash,
		Version:      s.info.Version,
		BuildAt:      s.info.BuildAt,
		StartTime:    s.info.StartTime.String(),
		UpTime:       time.Since(s.info.StartTime).String(),
	})
}

func (s *Handlers) swaggerDefHandler(w http.ResponseWriter, r *http.Request) {
	if host, _, err := net.SplitHostPort(r.Host); err == nil {
		r.Host = net.JoinHostPort(host, strconv.Itoa(int(s.opts.PortHTTP)))
	}

	var desc spec.Swagger

	_ = json.Unmarshal(s.desc, &desc)
	desc.Host = r.Host
	desc.Info = &spec.Info{}
	desc.Info.Version = s.info.Version
	desc.Info.Title = s.info.AppName

	//nolint:errchkjson // not need check error.
	_ = json.NewEncoder(w).Encode(desc)
}

func (s *Handlers) index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, _ = w.Write([]byte(`<html>
		<head>
			<title>Admin console</title>
		</head>
		<body>
		  	<p><a href="/docs">Swagger docs</a></p>
		  	<p><a href="/metrics">Metrics</a></p>
		  	<p><a href="/info">Info</a></p>
		  	<p><a href="/debug">Debug</a></p>
		  	<p><a href="/swagger.json">Swagger json</a></p>
		</body>
	</html>`))
}
