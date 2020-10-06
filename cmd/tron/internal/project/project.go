package project

import (
	"bufio"
	"errors"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"
	"golang.org/x/tools/imports"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/templates"
	"github.com/loghole/tron/internal/app"
)

const (
	gomodFile = "go.mod"
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
	IsNew   bool
}

func NewProject(module string) (project *Project, err error) {
	var isNew bool

	if module == "" {
		module, err = moduleFromGoMod()
		if err != nil {
			return nil, ErrEmptyModule
		}

		isNew = false
	}

	parts := strings.Split(module, "/")

	project = &Project{
		Module: module,
		Name:   parts[len(parts)-1],
		Protos: make([]*models.Proto, 0),
		IsNew:  isNew,
	}

	project.AbsPath, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (p *Project) InitGoMod() error {
	if _, err := os.Stat("go.mod"); err != nil {
		if os.IsNotExist(err) {
			return exec.Command("go", "mod", "init", p.Module).Run()
		}

		return err
	}

	return nil
}

func (p *Project) InitMakeFile() error {
	data := templates.NewTronMKData(
		strings.Join([]string{models.CmdDir, p.Name, models.MainFile}, "/"),
		models.DockerfileFilepath,
		p.Module,
	)

	tronMK, err := helpers.ExecTemplate(templates.TronMK, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteToFile(models.TronMKFilepath, []byte(tronMK)); err != nil {
		return err
	}

	if err := helpers.WriteToFile(models.MakefileFilepath, []byte(templates.Makefile)); err != nil {
		return err
	}

	return nil
}

func (p *Project) InitMainFile() error {
	fpath := filepath.Join(p.AbsPath, models.CmdDir, p.Name, models.MainFile)

	data := templates.NewMainData(models.Data{Protos: p.Protos})

	data.AddImport("log")
	data.AddImport("github.com/loghole/tron")
	data.AddImport(strings.Join([]string{p.Module, "config"}, "/"))

	for _, p := range p.Protos {
		data.AddImport(p.Service.Package)
	}

	mainScript, err := helpers.ExecTemplate(templates.MainTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	formattedBytes, err := format.Source([]byte(mainScript))
	if err != nil {
		return simplerr.Wrap(err, "failed to format process")
	}

	formattedBytes, err = imports.Process("", formattedBytes, nil)
	if err != nil {
		return simplerr.Wrap(err, "failed to imports process")
	}

	return helpers.WriteToFile(fpath, formattedBytes)
}

func (p *Project) InitGitignore() error {
	if !helpers.ConfirmOverwrite(models.GitignoreFilepath) {
		return nil
	}

	return helpers.WriteToFile(models.GitignoreFilepath, []byte(templates.GitignoreTemplate))
}

func (p *Project) InitDockerfile() error {
	dockerfile, err := helpers.ExecTemplate(templates.DefaultDockerfileTemplate, p)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteToFile(models.DockerfileFilepath, []byte(dockerfile)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	return nil
}

func (p *Project) InitLinter() error {
	dockerfile, err := helpers.ExecTemplate(templates.GolangCILintTemplate, p)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteToFile(models.GolangciLintFilepath, []byte(dockerfile)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	return nil
}

func (p *Project) InitValues() error {
	data := templates.ValuesData{
		List: []templates.Env{
			{Key: strings.ToUpper(app.LoggerLevelEnv), Val: "info"},
			{Key: strings.ToUpper(app.LoggerCollectorAddrEnv), Val: ""},
			{Key: strings.ToUpper(app.LoggerDisableStdoutEnv), Val: "false"},
			{Key: strings.ToUpper(app.JaegerAddrEnv), Val: "127.0.0.1:6831"},
			{Key: strings.ToUpper(app.HTTPPortEnv), Val: "8080"},
			{Key: strings.ToUpper(app.GRPCPortEnv), Val: "8081"},
			{Key: strings.ToUpper(app.AdminPortEnv), Val: "8082"},
		},
	}

	values, err := helpers.ExecTemplate(templates.ValuesTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteToFile(models.ValuesBaseFilepath, []byte(values)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	if err := helpers.WriteToFile(models.ValuesDevFilepath, []byte(templates.ValuesDevTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	if err := helpers.WriteToFile(models.ValuesLocalFilepath, []byte(templates.ValuesLocalTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	if err := helpers.WriteToFile(models.ValuesStgFilepath, []byte(templates.ValuesStgTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	if err := helpers.WriteToFile(models.ValuesProdFilepath, []byte(templates.ValuesProdTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	return nil
}

func (p *Project) FindProtoFiles(dirs ...string) error {
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
	for _, proto := range p.Protos {
		if strings.Contains(proto.RelativeDir, models.ProjectPathVendorPB) {
			continue
		}

		var (
			newDir  = filepath.Join(models.ProjectPathAPI, proto.Service.PackageName)
			newName = proto.Service.SnakeCasedName()
			oldPath = filepath.Join(proto.RelativeDir, proto.NameWithExt())
			newPath = filepath.Join(newDir, models.AddProtoExt(newName))
		)

		if oldPath == newPath {
			continue
		}

		color.Yellow("\tmove proto %s >> %s", oldPath, newPath)

		data, err := ioutil.ReadFile(proto.Path)
		if err != nil {
			return err
		}

		err = helpers.WriteToFile(path.Join(p.AbsPath, newPath), data)
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

	color.Yellow("\tcollected proto '%s%s.proto'", proto.Path, proto.Name)

	p.Protos = append(p.Protos, proto)

	return nil
}

func (p *Project) scanProtoFile(file io.Reader, proto *models.Proto) (*models.Proto, error) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		if m := models.PackageRegexp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			if proto.Service.PackageName != "" {
				return nil, simplerr.Wrapf(ErrMultiplePackages, "'%s/%s.proto'", proto.RelativeDir, proto.Name)
			}

			proto.Service.PackageName = m[1]
			proto.Service.Package = strings.Join([]string{p.Module, models.ProjectPathImplementation, m[1]}, "/")
		}

		if m := models.ServiceRegexp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			if proto.Service.Name != "" {
				return nil, simplerr.Wrapf(ErrMultipleServices, "'%s/%s.proto'", proto.RelativeDir, proto.Name)
			}

			proto.Service.Name = m[1]
		}

		if proto.Service.PackageName != "" && proto.Service.Name != "" {
			break
		}
	}

	return proto, nil
}

func moduleFromGoMod() (string, error) {
	file, err := os.Open(gomodFile)
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
