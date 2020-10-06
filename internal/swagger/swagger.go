package swagger

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utrack/clay/v2/transport/swagger"
)

type ServiceDesc interface {
	SwaggerDef(options ...swagger.Option) []byte
}

type Docs struct {
	desc []byte
}

func New(desc ServiceDesc, version string) *Docs {
	return &Docs{desc: desc.SwaggerDef(swagger.WithVersion(version))}
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
