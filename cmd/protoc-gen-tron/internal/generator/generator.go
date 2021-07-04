package generator

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

type Generator struct {
	module         string
	moduleName     string
	protoPkgPrefix string
}

func NewGenerator(module string) *Generator {
	parts := strings.Split(module, "/")

	gen := &Generator{
		module:         module,
		moduleName:     parts[len(parts)-1],
		protoPkgPrefix: strings.ReplaceAll(parts[len(parts)-1], "-", "_"),
	}

	return gen
}

func (gen *Generator) Generate(p *protogen.Plugin) error {
	p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	for _, f := range p.Files {
		if !f.Generate {
			continue
		}

		if !strings.HasPrefix(f.Proto.GetPackage(), gen.protoPkgPrefix) {
			p.Error(fmt.Errorf(
				"ERROR: file %s has %w, package can be %s",
				f.Proto.GetName(),
				ErrInvalidPackage,
				gen.protoPkgPrefix,
			))

			continue
		}

		if !strings.HasPrefix(f.Proto.GetOptions().GetGoPackage(), gen.module) {
			p.Error(fmt.Errorf(
				"ERROR: file %s has %w, go_package can be %s/pkg/name/version",
				f.Proto.GetName(),
				ErrInvalidGoPackage,
				gen.module,
			))
		}

		if len(f.Services) == 0 {
			continue
		}

		if len(f.Services) > 1 {
			p.Error(fmt.Errorf("ERROR: file %s has %w", f.Proto.GetName(), ErrMultiplyService))

			continue
		}

		gen.generateTransport(p, f)
		gen.generateImpl(p, f)
	}

	gen.generateMain(p)

	return nil
}
