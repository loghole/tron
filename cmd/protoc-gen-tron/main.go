package main

import (
	"errors"
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/loghole/tron/cmd/protoc-gen-tron/internal/helpers"
	"github.com/loghole/tron/cmd/protoc-gen-tron/internal/implement"
	"github.com/loghole/tron/cmd/protoc-gen-tron/internal/transport"
	"github.com/loghole/tron/cmd/protoc-gen-tron/internal/version"
)

var (
	ErrMultiplyService = errors.New("files with multiply services aren't supported")
	ErrEmptyPackage    = errors.New("empty package")
)

var (
	showVersion = flag.Bool("version", false, "print the version and exit")
	pkgPath     = flag.String("pkg-path", "", "pkg generated file path")
	implPath    = flag.String("impl-path", "", "pkg generated file path")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("protoc-gen-tron %v\n", version.CliVersion)
		return
	}

	opts := &protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}

	opts.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		module, err := helpers.ModuleFromGoMod()
		if err != nil {
			p.Error(err)
		}

		for _, file := range p.Files {
			if file.Proto.Package == nil {
				return ErrEmptyPackage
			}

			if !file.Generate {
				continue
			}

			if len(file.Proto.Service) > 1 {
				return ErrMultiplyService
			}

			transport.Generate(p, file)

			if err := implement.Generate(file, *implPath, *pkgPath, module); err != nil {
				p.Error(err)
			}
		}

		return nil
	})
}
