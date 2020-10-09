package swagger

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utrack/clay/v2/transport"
	"github.com/utrack/clay/v2/transport/swagger"
)

type Docs struct {
	desc []byte
}

func New(version string, services ...transport.Service) *Docs {
	descs := make([]transport.ServiceDesc, 0, len(services))

	for _, service := range services {
		descs = append(descs, service.GetDescription())
	}

	return &Docs{desc: transport.NewCompoundServiceDesc(descs...).SwaggerDef(swagger.WithVersion(version))}
}

func (s *Docs) InitRoutes(r chi.Router) {
	if r == nil {
		return
	}

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
}
