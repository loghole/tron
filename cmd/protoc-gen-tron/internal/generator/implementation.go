package generator

import (
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
)

//nolint:funlen // generation can be big
func (gen *Generator) generateImpl(p *protogen.Plugin, f *protogen.File) {
	var (
		service         = f.Services[0]
		implementDir    = gen.implDir(f.Proto.GetPackage())
		implementImport = gen.implImport(f.Proto.GetPackage())
	)

	for _, m := range service.Methods {
		implPath := filepath.Join(implementDir, snakeCase(m.GoName)+".go")

		if st, err := os.Stat(implPath); err == nil && !st.IsDir() {
			continue
		}

		g := p.NewGeneratedFile(implPath, implementImport)
		g.P("// Generated by protoc-gen-tron.")

		g.P()
		g.P("package ", f.GoPackageName)
		g.P()

		var (
			inputImportAlias  = p.FilesByPath[m.Input.Location.SourceFile].GoPackageName
			outputImportAlias = p.FilesByPath[m.Output.Location.SourceFile].GoPackageName
		)

		imports := make(map[protogen.GoPackageName]protogen.GoImportPath)
		imports[inputImportAlias] = m.Input.GoIdent.GoImportPath
		imports[outputImportAlias] = m.Output.GoIdent.GoImportPath

		g.P("import (")
		g.P(`"context"`)
		g.P(`"errors"`)
		g.P()

		for alias, importPath := range imports {
			g.P(alias, ` `, importPath)
		}

		g.P(")")
		g.P()

		g.P("func (i *Implementation) ", m.GoName, "(")
		g.P("ctx context.Context,")
		g.P("req *", inputImportAlias, ".", m.Input.GoIdent.GoName, ",")
		g.P(") (*", outputImportAlias, ".", m.Output.GoIdent.GoName, ", error) {")
		g.P(`return nil, errors.New("unimplemented")`)
		g.P("}")
		g.P()
	}

	implPath := filepath.Join(implementDir, "handler.go")

	if st, err := os.Stat(implPath); err == nil && !st.IsDir() {
		return
	}

	g := p.NewGeneratedFile(implPath, implementImport)
	g.P("// Generated by protoc-gen-tron.")

	g.P()
	g.P("package ", f.GoPackageName)
	g.P()

	g.P("import (")
	g.P(`"github.com/loghole/tron/transport"`)
	g.P()
	g.P(f.GoPackageName, ` `, f.GoImportPath)
	g.P(")")
	g.P()

	g.P("type Implementation struct {")
	g.P(f.GoPackageName, ".Unimplemented", service.GoName, "Server")
	g.P("}")
	g.P()

	g.P("func NewImplementation() *Implementation {")
	g.P("return &Implementation{}")
	g.P("}")
	g.P()

	g.P("// GetDescription is a simple alias to the ServiceDesc constructor.")
	g.P("// It makes it possible to register the service implementation @ the server.")
	g.P("func (i *Implementation) GetDescription() transport.ServiceDesc {")
	g.P("return ", f.GoPackageName, ".New", service.GoName, "ServiceDesc(i)")
	g.P("}")
	g.P()
}
