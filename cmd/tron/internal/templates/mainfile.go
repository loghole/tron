package templates

import (
	"strconv"
	"strings"

	"github.com/loghole/tron/cmd/tron/internal/models"
)

// MainData is a datastruct for main.go
type MainData struct {
	models.Data

	Imports map[string]Import

	ConfigPackage string
}

type Import struct {
	Alias string
	Pkg   string
}

// Pkg returns pkg prefix with dot
func (d MainData) Pkg(pkg string) string {
	i := d.Imports[pkg]
	if i.Alias != "" {
		return i.Alias + `.`
	}
	return ""
}

// AddImport adds pkg to Imports map, and guarantees a unique alias
// if alias[0] passed it will be used as alias for this pkg, alias[1..N] are ignored
func (d *MainData) AddImport(pkg string, alias ...string) *MainData {
	if d.Imports == nil {
		d.Imports = make(map[string]Import)
	}
	if len(alias) != 0 {
		d.Imports[pkg] = Import{Pkg: pkg, Alias: alias[0]}
		return d
	}

	path := strings.Split(pkg, "/")
	if len(path) == 0 {
		return d
	}

	lastPart := strings.Replace(path[len(path)-1], "-", "_", -1)

	a := func(i int) string {
		if i == 1 {
			return lastPart
		}
		return lastPart + strconv.Itoa(i)
	}

	n := 1
	found := false
	for !found {
		for _, i := range d.Imports {
			if i.Alias == a(n) {
				n++
				break
			}
		}
		found = true
	}
	d.Imports[pkg] = Import{Pkg: pkg, Alias: a(n)}

	return d
}

const MainTemplate = `package main

import (
	{{ range $import := .Imports -}}
		{{ $import.Alias }} "{{ $import.Pkg }}"
	{{ end }}
)

func main() {
	app, err := {{ pkg "github.com/loghole/tron" }}New()
	if err != nil {
		{{ pkg "log" }}Fatalf("can't create app: %s", err)
	}

	// Init all ..

	var (
		{{ range $proto := .Protos -}}
			{{ $proto.Service.SnakeCasedName }} = {{ pkg $proto.Service.Package }}New{{ $proto.Service.Name }}()
		{{- end }}
	)

	app.Run({{- range $proto := .Protos -}} {{ $proto.Service.SnakeCasedName }}, {{- end -}})

	// Stop all...
}
`
