package templates

import (
	"fmt"
)

type HandlerData struct {
	GoPackage string
	Service   string
	Desc      string
	Imports   map[string]string
}

func NewHandlerData(pkg, srv string) *HandlerData {
	return &HandlerData{
		GoPackage: pkg,
		Service:   srv,
		Imports:   map[string]string{},
	}
}

func (d *HandlerData) AddImport(alias, path string) {
	if _, ok := d.Imports[path]; !ok {
		d.Imports[path] = alias
	}
}

const HandlerTemplate = `package {{ .GoPackage }}

import (
	{{ range $import, $alias := .Imports }}
		{{ $alias }} "{{ $import }}"
	{{- end }}

	"github.com/loghole/tron/transport"
)

type Implementation struct {
	{{ .Desc }}.Unimplemented{{ .Service }}Server
}

func NewImplementation() *Implementation {
	return &Implementation{}
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *Implementation) GetDescription() transport.ServiceDesc {
	return {{ .Desc }}.New{{ .Service }}ServiceDesc(i)
}
`

type MethodData struct {
	GoPackage string
	Name      string
	Input     string
	Output    string
	Imports   map[string]string
}

func NewMethodData() *MethodData {
	return &MethodData{
		Imports: map[string]string{},
	}
}

func (m *MethodData) AddImport(alias, path string) {
	if _, ok := m.Imports[alias]; !ok {
		m.Imports[alias] = path
	}
}

func (m *MethodData) AddInput(alias, importPath, typeName string) {
	m.AddImport(alias, importPath)

	m.Input = fmt.Sprintf("%s.%s", alias, typeName)
}

func (m *MethodData) AddOutput(alias, importPath, typeName string) {
	m.AddImport(alias, importPath)

	m.Output = fmt.Sprintf("%s.%s", alias, typeName)
}

const MethodTemplate = `package {{ .GoPackage }}

import (
	"context"
	"errors"

	{{ range $alias, $import := .Imports }}
		{{ $alias }} "{{ $import }}"
	{{- end }}
)

func (i *Implementation) {{ .Name }} (
	ctx context.Context,
	req *{{ .Input }},
) (resp *{{ .Output }}, err error) {
	return nil, errors.New("method {{ .Name }} unimplemented")
}
`
