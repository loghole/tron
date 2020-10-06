package generate

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
)

const (
	cmdProtoc = "protoc"

	pkgMap = "Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types," +
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

var ErrProtoc = errors.New("protoc")

type Proto struct {
	project *project.Project
}

func NewProto(pr *project.Project) *Proto {
	return &Proto{project: pr}
}

func (p *Proto) Generate() error {
	fmt.Println("Start generate proto")

	for _, proto := range p.project.Protos {
		fmt.Printf("    generate go-fast %s: ", proto.Service.SnakeCasedName())

		if err := p.generateGoFast(proto); err != nil {
			color.Red("ERROR: %v", err)

			return err
		}

		color.Green("OK")

		// generate go-clay
		fmt.Printf("    generate go-clay %s: ", proto.Service.SnakeCasedName())

		if err := p.generateGoClay(proto); err != nil {
			color.Red("ERROR: %v", err)

			return err
		}

		color.Green("OK")
	}

	color.Green("Success")

	return nil
}

func (p *Proto) generateGoFast(proto *models.Proto) error {
	err := helpers.Mkdir(path.Join(models.ProjectPathPkgClients, proto.Service.PackageName, proto.Name))
	if err != nil {
		return err
	}

	args := []string{
		fmt.Sprintf("--plugin=protoc-gen-gofast=%s", path.Join(p.project.AbsPath, "bin/protoc-gen-gofast")),
		fmt.Sprintf("-I%s:%s", proto.RelativeDir, models.ProjectPathVendorPB),
		fmt.Sprintf("--gofast_out=%s,plugins=grpc:%s", pkgMap, path.Join(models.ProjectPathPkgClients, proto.Service.PackageName)),
		path.Join(proto.RelativeDir, proto.NameWithExt()),
	}

	if err := execProtoc(p.project.AbsPath, args); err != nil {
		return simplerr.Wrap(err, "generate go-fast")
	}

	return nil
}

func (p *Proto) generateGoClay(proto *models.Proto) error {
	err := helpers.Mkdir(filepath.Join(models.ProjectPathImplementation, proto.Service.PackageName, proto.Name))
	if err != nil {
		return simplerr.Wrap(err, "failed to mkdir")
	}

	relToRoot, err := filepath.Rel(filepath.Join(p.project.AbsPath, models.ProjectPathPkgClients, proto.Service.PackageName), p.project.AbsPath)
	if err != nil {
		return simplerr.Wrap(err, "failed to get relative path")
	}

	args := []string{
		fmt.Sprintf("--plugin=protoc-gen-goclay=%s", path.Join(p.project.AbsPath, "bin/protoc-gen-goclay")),
		fmt.Sprintf("-I%s:%s", filepath.Join(relToRoot, proto.RelativeDir), filepath.Join(relToRoot, models.ProjectPathVendorPB)),
		fmt.Sprintf("--goclay_out=%s,impl=true,impl_path=%s,impl_type_name_tmpl=%s:.", pkgMap, path.Join(relToRoot, models.ProjectPathImplementation, proto.Service.PackageName), models.ImplementationName),
		path.Join(relToRoot, proto.RelativeDir, proto.NameWithExt()),
	}

	if err := execProtoc(path.Join(p.project.AbsPath, models.ProjectPathPkgClients, proto.Service.PackageName), args); err != nil {
		return simplerr.Wrap(err, "generate go-clay")
	}

	return nil
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
