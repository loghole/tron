package project

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
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
	projectPathApi = "api"
	gomodFile      = "go.mod"
	protoExt       = ".proto"
)

var (
	ErrEmptyModule    = errors.New("can't create project with empty module")
	ErrModuleNotFound = errors.New("project module does not exists")
)

type Project struct {
	AbsPath string
	Module  string
	Name    string
	Protos  []*models.Proto

	serviceRgx *regexp.Regexp
	packageRgx *regexp.Regexp
}

func NewProject(module string) (project *Project, err error) {
	if module == "" {
		module, err = moduleFromGoMod()
		if err != nil {
			return nil, ErrEmptyModule
		}
	}

	parts := strings.Split(module, "/")

	project = &Project{
		Module:     module,
		Name:       parts[len(parts)-1],
		Protos:     make([]*models.Proto, 0),
		serviceRgx: regexp.MustCompile(`^service (.*?) {`),
		packageRgx: regexp.MustCompile(`^package[\s]*?(\w*);$`),
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
	files := []struct {
		name string
		tmpl string
	}{
		{name: "tron.mk", tmpl: templates.TronMK},
		{name: "Makefile", tmpl: templates.Makefile},
	}

	for _, file := range files {
		if err := helpers.WriteToFile(file.name, []byte(file.tmpl)); err != nil {
			return err
		}
	}

	return nil
}

func (p *Project) InitMainFile() error {
	fpath := filepath.Join(p.AbsPath, "cmd", p.Name, "main.go")

	data := templates.MainData{
		Data: models.Data{
			Protos: p.Protos,
		},
		ConfigPackage: filepath.Join(p.Module, "./config"),
	}
	data.AddImport("context")
	data.AddImport("log")
	data.AddImport("github.com/loghole/tron")
	data.AddImport(data.ConfigPackage, "_")

	for _, p := range p.Protos {
		data.AddImport(p.Service.Package)
	}

	mainScript, err := helpers.ExecTemplate(templates.MainTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	formattedBytes, err := imports.Process("", []byte(mainScript), nil)
	if err != nil {
		return simplerr.Wrap(err, "failed to imports process")
	}

	return helpers.WriteToFile(fpath, formattedBytes)
}

func (p *Project) InitGitignore() error {
	if !helpers.ConfirmOverwrite(".gitignore") {
		return nil
	}

	return helpers.WriteToFile(".gitignore", []byte(templates.GitignoreTemplate))
}

func (p *Project) InitDockerfile() error {
	dockerfile, err := helpers.ExecTemplate(templates.DefaultDockerfileTemplate, p)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteToFile(".deploy/docker/default/Dockerfile", []byte(dockerfile)); err != nil {
		return err
	}

	return nil
}

func (p *Project) InitLinter() error {
	dockerfile, err := helpers.ExecTemplate(templates.GolangCILintTemplate, p)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteToFile(".golangci.yaml", []byte(dockerfile)); err != nil {
		return err
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

	if err := helpers.WriteToFile(".deploy/config/values.yaml", []byte(values)); err != nil {
		return err
	}

	if err := helpers.WriteToFile(".deploy/config/values_dev.yaml", []byte(templates.ValuesDevTemplate)); err != nil {
		return err
	}

	if err := helpers.WriteToFile(".deploy/config/values_stg.yaml", []byte(templates.ValuesStgTemplate)); err != nil {
		return err
	}

	if err := helpers.WriteToFile(".deploy/config/values_prod.yaml", []byte(templates.ValuesProdTemplate)); err != nil {
		return err
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
		var (
			newDir  = path.Join(projectPathApi, proto.Service.PackageName)
			newName = proto.Service.SnakeCasedName()
			oldPath = path.Join(proto.RelativeDir, proto.Name+protoExt)
			newPath = path.Join(newDir, newName+protoExt)
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

	if filepath.Ext(path) != protoExt {
		return nil
	}

	proto := &models.Proto{Name: strings.TrimSuffix(info.Name(), protoExt), Path: path}

	switch {
	case strings.HasPrefix(path, p.AbsPath):
		proto.RelativeDir, err = filepath.Rel(p.AbsPath, filepath.Dir(path))
	default:
		proto.RelativeDir = filepath.Join("api", filepath.Base(filepath.Dir(path)))
	}

	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer helpers.Close(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		if m := p.packageRgx.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			if proto.Service.PackageName != "" {
				return fmt.Errorf("package '%s/%s.proto' has multiple package entries", proto.RelativeDir, proto.Name)
			}

			proto.Service.PackageName = m[1]
			proto.Service.Package = strings.Join([]string{p.Module, "internal/app/controllers", m[1]}, "/")
		}

		if m := p.serviceRgx.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			if proto.Service.Name != "" {
				return fmt.Errorf("package '%s/%s.proto' has multiple service entries", proto.RelativeDir, proto.Name)
			}

			proto.Service.Name = m[1]
		}

		if proto.Service.PackageName != "" && proto.Service.Name != "" {
			break
		}
	}

	color.Yellow("\tcollected proto '%s%s.proto'", proto.Path, proto.Name)

	p.Protos = append(p.Protos, proto)

	return nil
}

func moduleFromGoMod() (string, error) {
	file, err := os.Open(gomodFile)
	if err != nil {
		return "", err
	}

	defer helpers.Close(file)

	var reg = regexp.MustCompile(`^module (.*)$`)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}

		if m := reg.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			return m[1], nil
		}
	}

	return "", ErrModuleNotFound
}
