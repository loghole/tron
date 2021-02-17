package parsers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type ProtoDescParser struct {
	project *models.Project
	printer stdout.Printer

	descriptors []*descriptorpb.FileDescriptorProto
}

func NewProtoDescParser(project *models.Project, printer stdout.Printer) *ProtoDescParser {
	return &ProtoDescParser{
		project: project,
		printer: printer,
	}
}

func (p *ProtoDescParser) Parse() error {
	p.printer.VerbosePrintln(color.FgMagenta, "Parse proto files")

	if err := p.scanDescriptors(); err != nil {
		return err
	}

	if err := p.parseDescriptors(); err != nil {
		return err
	}

	if err := p.parsePkgMap(); err != nil {
		return err
	}

	p.printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func (p *ProtoDescParser) scanDescriptors() error {
	dir, err := os.MkdirTemp("", "*_auth_gen")
	if err != nil {
		return fmt.Errorf("create tmp dir: %w", err)
	}

	defer os.RemoveAll(dir)

	for idx, file := range p.project.ProtoFiles {
		output := filepath.Join(dir, strconv.Itoa(idx)+".desc")

		input, err := filepath.Rel(p.project.AbsPath, file)
		if err != nil {
			return fmt.Errorf("relative path: %w", err)
		}

		cmd := exec.Command("protoc",
			"-o", output, fmt.Sprintf("-I.:%s", models.ProjectPathVendorPB), input)
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("'%s': %w", cmd.String(), err)
		}

		data, err := os.ReadFile(output)
		if err != nil {
			return fmt.Errorf("read parsed data: %w", err)
		}

		dest := &descriptorpb.FileDescriptorSet{}

		if err := proto.Unmarshal(data, dest); err != nil {
			return fmt.Errorf("unmarshal proto: %w", err)
		}

		p.descriptors = append(p.descriptors, dest.File...)
	}

	return nil
}

func (p *ProtoDescParser) parseDescriptors() error {
	for _, desc := range p.descriptors {
		result, err := p.parseFileDescriptor(desc)
		if err != nil {
			return fmt.Errorf("%s: %w", desc.GetName(), err)
		}

		p.project.Protos = append(p.project.Protos, result)
	}

	return nil
}

func (p *ProtoDescParser) parseFileDescriptor(desc *descriptorpb.FileDescriptorProto) (*models.Proto, error) {
	if len(desc.GetService()) > 1 {
		return nil, ErrMultipleServices
	}

	if !models.ProtoPkgVersionRegexp.MatchString(desc.GetPackage()) {
		return nil, ErrInvalidProtoPkgName
	}

	result := &models.Proto{
		Name:      strings.TrimSuffix(filepath.Base(desc.GetName()), models.ProtoExt),
		Path:      filepath.Join(p.project.AbsPath, desc.GetName()),
		RelPath:   desc.GetName(),
		Package:   desc.GetPackage(),
		GoPackage: p.goPackage(desc),
		Imports:   desc.GetDependency(),
	}

	if len(desc.GetService()) == 0 {
		return result, nil
	}

	serviceDesc := desc.GetService()[0]

	result.Service = &models.Service{
		Name:    serviceDesc.GetName(),
		Methods: make([]*models.Method, 0, len(serviceDesc.GetMethod())),
	}

	for _, methodDesc := range serviceDesc.GetMethod() {
		result.Service.Methods = append(result.Service.Methods, &models.Method{
			Name:   methodDesc.GetName(),
			Input:  methodDesc.GetInputType(),
			Output: methodDesc.GetOutputType(),
		})
	}

	return result, nil
}

func (p *ProtoDescParser) parsePkgMap() error {
	const pkgParts = 3

	exists := make(map[string]struct{})

	for _, pp := range p.project.Protos {
		for _, val := range pp.Imports {
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
		if len(parts) < pkgParts {
			p.printer.Println(color.FgYellow, "\tcreate import alias failed: %s", val)

			continue
		}

		alias = strings.Join([]string{
			p.project.Module,
			models.ProjectPathPkgClients,
			parts[len(parts)-1],
			parts[len(parts)-2],
		}, "/")

		p.project.ProtoPkgMap = append(p.project.ProtoPkgMap, fmt.Sprintf("M%s=%s", val, alias))
	}

	return nil
}

func (p *ProtoDescParser) goPackage(desc *descriptorpb.FileDescriptorProto) string {
	const importParts = 2

	var (
		result = desc.GetOptions().GetGoPackage()
		parts  = strings.Split(result, ";")
	)

	if result == "" || len(parts) != importParts {
		return strings.ReplaceAll(desc.GetPackage(), ".", "_")
	}

	return parts[1]
}
