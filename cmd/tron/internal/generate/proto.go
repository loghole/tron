package generate

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
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
		p.printer.VerbosePrintf(color.Reset, "\tgenerate go-fast %s: ", proto.Service.SnakeCasedName())

		if err := p.generateGoFast(proto); err != nil {
			p.printer.VerbosePrintln(color.FgRed, "ERROR: %v", err)

			return err
		}

		p.printer.VerbosePrintln(color.FgGreen, "OK")

		if !proto.Service.WithImpl {
			continue
		}

		// generate go-clay
		p.printer.VerbosePrintf(color.Reset, "\tgenerate go-clay %s: ", proto.Service.SnakeCasedName())

		if err := p.generateGoClay(proto); err != nil {
			p.printer.VerbosePrintln(color.FgRed, "ERROR: %v", err)

			return err
		}

		p.printer.VerbosePrintln(color.FgGreen, "OK")
	}

	return nil
}

func (p *proto) generateGoFast(proto *models.Proto) error {
	err := helpers.Mkdir(path.Join(models.ProjectPathPkgClients, proto.Service.Package, proto.Name))
	if err != nil {
		return err
	}

	args := []string{
		fmt.Sprintf("--plugin=protoc-gen-gofast=%s", path.Join(p.project.AbsPath, "bin/protoc-gen-gofast")),
		fmt.Sprintf("-I%s:%s", proto.RelativeDir, models.ProjectPathVendorPB),
		fmt.Sprintf("--gofast_out=%s,plugins=grpc:%s",
			strings.Join(p.pkgMap, ","),
			path.Join(models.ProjectPathPkgClients, proto.Service.Package),
		),
		path.Join(proto.RelativeDir, proto.NameWithExt()),
	}

	if err := execProtoc(p.project.AbsPath, args); err != nil {
		return simplerr.Wrap(err, "generate go-fast")
	}

	return nil
}

func (p *proto) generateGoClay(proto *models.Proto) error {
	err := helpers.Mkdir(filepath.Join(models.ProjectPathImplementation, proto.Service.Package, proto.Name))
	if err != nil {
		return simplerr.Wrap(err, "failed to mkdir")
	}

	wd := filepath.Join(p.project.AbsPath, models.ProjectPathPkgClients, proto.Service.Package)

	relToRoot, err := filepath.Rel(wd, p.project.AbsPath)
	if err != nil {
		return simplerr.Wrap(err, "failed to get relative path")
	}

	args := []string{
		fmt.Sprintf("--plugin=protoc-gen-goclay=%s",
			path.Join(p.project.AbsPath, "bin/protoc-gen-goclay"),
		),
		fmt.Sprintf("-I%s:%s",
			filepath.Join(relToRoot, proto.RelativeDir),
			filepath.Join(relToRoot, models.ProjectPathVendorPB),
		),
		fmt.Sprintf("--goclay_out=%s,impl=true,impl_path=%s,impl_type_name_tmpl=%s:.",
			strings.Join(p.pkgMap, ","),
			path.Join(relToRoot, models.ProjectPathImplementation, proto.Service.Package),
			models.ImplementationName,
		),
		path.Join(relToRoot, proto.RelativeDir, proto.NameWithExt()),
	}

	if err := execProtoc(wd, args); err != nil {
		return simplerr.Wrap(err, "generate go-clay")
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
