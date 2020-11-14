package generate

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

const (
	cmdProtoc   = "protoc"
	pkgMapParts = 3
	pkgMap      = "Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/api.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor," +
		"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/source_context.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/type.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types"
)

func Protos(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate proto files")

	generator := &proto{project: p, printer: printer, pkgMap: strings.Split(pkgMap, ",")}

	generator.initPkgMap()

	return generator.run()
}

type proto struct {
	project *project.Project
	printer stdout.Printer
	pkgMap  []string
}

func (p *proto) run() error {
	for _, proto := range p.project.Protos {
		// Generate protoc-gen-go and gateway.
		p.printer.VerbosePrintf(color.Reset, "\tgenerate protoc-gen-go %s: ", proto.Service.SnakeCasedName())

		if err := p.generateProtocGenGo(proto); err != nil {
			p.printer.VerbosePrintf(color.FgRed, "ERROR: %v\n", err)
		}

		p.printer.VerbosePrintln(color.FgGreen, "OK")

		// Generate transport.
		p.printer.VerbosePrintf(color.Reset, "\tgenerate tron transport %s: ", proto.Service.SnakeCasedName())

		if err := p.generateTransport(proto); err != nil {
			p.printer.VerbosePrintf(color.FgRed, "ERROR: %v\n", err)
		}

		p.printer.VerbosePrintln(color.FgGreen, "OK")
	}

	return nil
}

func (p *proto) generateProtocGenGo(proto *models.Proto) error {
	if err := helpers.Mkdir(proto.PkgFilePath()); err != nil {
		return err
	}

	var (
		outputPath = filepath.Join(proto.PkgDirPath())
		pkgMap     = strings.Join(p.pkgMap, ",")
	)

	args := []string{
		fmt.Sprintf("--plugin=protoc-gen-go=%s",
			filepath.Join(p.project.AbsPath, "bin", "protoc-gen-go")),
		fmt.Sprintf("--plugin=protoc-gen-go-grpc=%s",
			filepath.Join(p.project.AbsPath, "bin", "protoc-gen-go-grpc")),
		fmt.Sprintf("--plugin=protoc-gen-grpc-gateway=%s",
			filepath.Join(p.project.AbsPath, "bin", "protoc-gen-grpc-gateway")),
		fmt.Sprintf("--plugin=protoc-gen-openapiv2=%s",
			filepath.Join(p.project.AbsPath, "bin", "protoc-gen-openapiv2")),
		fmt.Sprintf("-I%s:%s", proto.RelativeDir, models.ProjectPathVendorPB),
		fmt.Sprintf("--go_out=%s:%s", pkgMap, outputPath),
		fmt.Sprintf("--go-grpc_out=%s:%s", pkgMap, outputPath),
		fmt.Sprintf("--grpc-gateway_out=%s:%s", pkgMap, outputPath),
		fmt.Sprintf("--openapiv2_out=%s:%s", pkgMap, outputPath),
		"--grpc-gateway_opt", "logtostderr=true",
		"--openapiv2_opt", "fqn_for_openapi_name=true",
		filepath.Join(proto.RelativeDir, proto.NameWithExt()),
	}

	if err := execProtoc(p.project.AbsPath, args); err != nil {
		return simplerr.Wrap(err, "generate protoc-gen-go")
	}

	return nil
}

func (p *proto) generateTransport(proto *models.Proto) error {
	if err := helpers.Mkdir(proto.PkgFilePath()); err != nil {
		return err
	}

	var (
		outputPath = filepath.Join(proto.PkgDirPath())
		pkgMap     = strings.Join(p.pkgMap, ",")
	)

	args := []string{
		fmt.Sprintf("--plugin=protoc-gen-tron=%s", filepath.Join(p.project.AbsPath, "bin", "protoc-gen-tron")),
		fmt.Sprintf("-I%s:%s", proto.RelativeDir, models.ProjectPathVendorPB),
		fmt.Sprintf("--tron_out=%s:%s", pkgMap, outputPath),
		"--tron_opt",
		fmt.Sprintf("impl-path=%s,pkg-path=%s", models.ProjectPathImplementation, models.ProjectPathPkgClients),
		filepath.Join(proto.RelativeDir, proto.NameWithExt()),
	}

	if err := execProtoc(p.project.AbsPath, args); err != nil {
		return simplerr.Wrap(err, "generate protoc-gen-go")
	}

	return nil
}

func (p *proto) initPkgMap() {
	exists := make(map[string]struct{})

	for _, proto := range p.project.Protos {
		for _, val := range proto.Imports {
			if !strings.HasPrefix(val, p.project.Module) {
				continue
			}

			exists[val] = struct{}{}
		}
	}

	for val := range exists {
		alias := strings.TrimPrefix(val, p.project.Module)
		alias = strings.TrimSuffix(alias, models.ProtoExt)

		parts := strings.Split(alias, "/")

		if len(parts) < pkgMapParts {
			p.printer.Println(color.FgYellow, "\tcreate import alias failed: %s", val)

			continue
		}

		alias = strings.Join([]string{
			p.project.Module,
			models.ProjectPathPkgClients,
			parts[len(parts)-1],
			parts[len(parts)-2],
		}, "/")

		p.pkgMap = append(p.pkgMap, fmt.Sprintf("M%s=%s", val, alias))
	}
}

func execProtoc(wd string, args []string) error {
	stderr := bytes.NewBuffer(nil)

	cmd := exec.Command(cmdProtoc, args...)
	cmd.Dir = wd
	cmd.Stderr = stderr

	o, err := cmd.Output()
	if err != nil {
		return simplerr.Wrapf(err, "stderr: %s", stderr.String())
	}

	if len(o) > 0 {
		return simplerr.Wrapf(ErrProtoc, "unexpected output: %s", o)
	}

	return nil
}
