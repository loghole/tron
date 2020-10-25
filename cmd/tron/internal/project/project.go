package project

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

var (
	ErrEmptyModule      = errors.New("can't create project with empty module")
	ErrModuleNotFound   = errors.New("project module does not exists")
	ErrMultiplePackages = errors.New("multiple package entries")
	ErrMultipleServices = errors.New("multiple service entries")
)

type Project struct {
	AbsPath string
	Module  string
	Name    string
	Protos  []*models.Proto
	printer stdout.Printer
}

func NewProject(module string, printer stdout.Printer) (project *Project, err error) {
	if module == "" {
		module, err = moduleFromGoMod()
		if err != nil {
			return nil, ErrEmptyModule
		}
	}

	parts := strings.Split(module, "/")

	project = &Project{
		Module:  module,
		Name:    parts[len(parts)-1],
		Protos:  make([]*models.Proto, 0),
		printer: printer,
	}

	project.AbsPath, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (p *Project) FindProtoFiles(dirs ...string) error {
	p.printer.VerbosePrintln(color.FgMagenta, "Find proto files")

	for _, dir := range dirs {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			return err
		}

		if _, err := os.Stat(absPath); err != nil {
			return err
		}

		if err := filepath.Walk(absPath, p.getProtoFileInfo); err != nil {
			return err
		}
	}

	return nil
}

func (p *Project) MoveProtoFiles() error {
	p.printer.VerbosePrintln(color.FgMagenta, "Move proto files")

	for _, proto := range p.Protos {
		if strings.Contains(proto.RelativeDir, models.ProjectPathVendorPB) {
			continue
		}

		var (
			newDir  = filepath.Join(models.ProjectPathAPI, proto.Version)
			newName = proto.Service.SnakeCasedName()
			oldPath = filepath.Join(proto.RelativeDir, proto.NameWithExt())
			newPath = filepath.Join(newDir, models.AddProtoExt(newName))
		)

		if oldPath == newPath {
			continue
		}

		p.printer.VerbosePrintf(color.Reset, "\tmove proto %s >> %s", oldPath, newPath)

		data, err := ioutil.ReadFile(proto.Path)
		if err != nil {
			return err
		}

		err = helpers.WriteToFile(filepath.Join(p.AbsPath, newPath), data)
		if err != nil {
			return err
		}

		err = os.Remove(oldPath)
		if err != nil {
			return err
		}

		proto.Name = newName
		proto.RelativeDir = newDir
	}

	return nil
}

func (p *Project) getProtoFileInfo(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if filepath.Ext(path) != models.ProtoExt {
		return nil
	}

	proto := &models.Proto{Name: strings.TrimSuffix(info.Name(), models.ProtoExt), Path: path}

	switch {
	case strings.HasPrefix(path, p.AbsPath):
		proto.RelativeDir, err = filepath.Rel(p.AbsPath, filepath.Dir(path))
	default:
		proto.RelativeDir = filepath.Join(models.ProjectPathAPI, filepath.Base(filepath.Dir(path)))
	}

	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer helpers.Close(file)

	proto, err = p.scanProtoFile(file, proto)
	if err != nil {
		return err
	}

	p.printer.VerbosePrintf(color.Reset, "\tcollected proto '%s%s.proto'\n", proto.Path, proto.Name)

	p.Protos = append(p.Protos, proto)

	return nil
}

func (p *Project) scanProtoFile(file io.Reader, proto *models.Proto) (*models.Proto, error) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var srv, pkg string

	for scanner.Scan() {
		if m := models.PackageRegexp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			if pkg != "" {
				return nil, simplerr.Wrapf(ErrMultiplePackages, "'%s/%s.proto'", proto.RelativeDir, proto.Name)
			}

			pkg = m[1]
		}

		if m := models.ServiceRegexp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			if srv != "" {
				return nil, simplerr.Wrapf(ErrMultipleServices, "'%s/%s.proto'", proto.RelativeDir, proto.Name)
			}

			srv = m[1]
		}

		if m := models.ImportRegexp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			proto.Imports = append(proto.Imports, m[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := proto.SetService(p.Module, srv, pkg); err != nil {
		return nil, err
	}

	return proto, nil
}

func moduleFromGoMod() (string, error) {
	file, err := os.Open(models.GoModFile)
	if err != nil {
		return "", err
	}

	defer helpers.Close(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}

		if m := models.ModuleRegexp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			return m[1], nil
		}
	}

	return "", ErrModuleNotFound
}
