package generator

import (
	"google.golang.org/protobuf/compiler/protogen"
)

// nolint:funlen // generation can be big
func (gen *Generator) generateTransport(p *protogen.Plugin, f *protogen.File) {
	var (
		service     = f.Services[0]
		descName    = service.GoName + "ServiceDesc"
		protoSource = f.Proto.GetName()
		fileName    = fileName(protoSource)
	)

	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.tron.go", f.GoImportPath)
	g.P("// Code generated by protoc-gen-tron. DO NOT EDIT.")

	if protoSource != "" {
		g.P("// source: ", protoSource)
	}

	g.P()
	g.P("package ", f.GoPackageName)
	g.P()

	g.P("import (")
	g.P(`"context"`)
	g.P(`"embed"`)
	g.P()
	g.P(`"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"`)
	g.P(`"github.com/loghole/tron/transport"`)
	g.P(`"google.golang.org/grpc"`)
	g.P(")")
	g.P()

	if protoSource != "" {
		g.P("//go:embed ", fileName, ".swagger.json")
	}

	g.P("var swagger embed.FS")
	g.P()

	g.P("// ", descName, " is description for the ", service.GoName, "Server.")
	g.P("type ", descName, " struct {")
	g.P("svc ", service.GoName, "Server")
	g.P("}")
	g.P()

	g.P("func New", descName, "(s ", service.GoName, "Server) ", "transport.ServiceDesc {")
	g.P("return &", descName, "{svc: s}")
	g.P("}")
	g.P()

	g.P("func(d *", descName, ") RegisterGRPC(s *grpc.Server) {")
	g.P("Register", service.GoName, "Server(s, d.svc)")
	g.P("}")
	g.P()

	g.P("func(d *", descName, ") RegisterHTTP(mux *runtime.ServeMux) {")
	g.P("Register", service.GoName, "HandlerServer(context.Background(), mux, d.svc)")
	g.P("}")
	g.P()

	g.P("func(d *", descName, ") SwaggerDef() []byte {")
	g.P(`b, _ := swagger.ReadFile("`, fileName, `.swagger.json")`)
	g.P()
	g.P("return b")
	g.P("}")
	g.P()
}
