package templates

import (
	"strconv"
	"strings"

	"github.com/loghole/tron/cmd/tron/internal/models"
)

// MainData is a datastruct for main.go
type MainData struct {
	*models.Project
	Imports map[string]Import
}

func NewMainData(project *models.Project) *MainData {
	return &MainData{
		Project: project,
		Imports: make(map[string]Import),
	}
}

type Import struct {
	Alias string
	Pkg   string
}

// Pkg returns pkg prefix with dot
func (m *MainData) Pkg(pkg string) string {
	i := m.Imports[pkg]

	if i.Alias != "" {
		return i.Alias + `.`
	}

	return ""
}

// AddImport adds pkg to Imports map, and guarantees a unique alias
// if alias[0] passed it will be used as alias for this pkg, alias[1..N] are ignored
func (m *MainData) AddImport(pkg string, alias ...string) {
	if m.Imports == nil {
		m.Imports = make(map[string]Import)
	}

	if len(alias) != 0 {
		m.Imports[pkg] = Import{Pkg: pkg, Alias: alias[0]}

		return
	}

	path := strings.Split(pkg, "/")

	if len(path) == 0 {
		return
	}

	lastPart := strings.ReplaceAll(path[len(path)-1], "-", "_")

	a := func(i int) string {
		if i == 1 {
			return lastPart
		}

		return lastPart + strconv.Itoa(i)
	}

	n := 1
	found := false

	for !found {
		for _, i := range m.Imports {
			if i.Alias == a(n) {
				n++

				break
			}
		}

		found = true
	}

	m.Imports[pkg] = Import{Pkg: pkg, Alias: a(n)}
}

const MainTemplate = `package main

import (
	"log"	

	"{{ .Module }}/config"
	{{ range $import := .Imports -}}
		{{ $import.Alias }} "{{ $import.Pkg }}"
	{{ end }}

	"github.com/loghole/tron"
)

func main() {
	app, err := tron.New(tron.AddLogCaller())
	if err != nil {
		log.Fatalf("can't create app: %s", err)
	}

	defer app.Close()

	app.Logger().Info(config.GetExampleValue())

	// Init all ..

	{{ if .WithProtos }}
	var (
		{{ range $proto := .Protos -}}
			{{- if $proto.WithImpl -}}
			{{ $proto.GoPackage }}Handler = {{ $proto.GoPackage }}.NewImplementation()
			{{ end -}}
		{{ end }}
	)
	{{ end }}

	if err := app.WithRunOptions().Run(
			{{- range $proto := .Protos -}} 
				{{- if $proto.WithImpl -}}
				{{ $proto.GoPackage }}Handler,
				{{- end -}}
			{{- end -}}
		); err != nil {
		app.Logger().Fatalf("can't run app: %v", err)
	}

	// Stop all...
}
`
