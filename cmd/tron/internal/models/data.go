package models

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
)

const ProtoExt = ".proto"

var ErrInvalidProtoPkgName = errors.New("invalid proto package name")

// Data is a datastruct containing common stuff for the templates
type Data struct {
	Protos []*Proto
	// AppNameUnderscoredUpper is an UPPERCASE underscored project proto
	AppNameUnderscoredUpper string
	// AppName is application name
	AppName string
	// CmdName is a project name
	CmdName string
	// Go Module
	Module string
}

type Proto struct {
	Name        string
	Path        string
	RelativeDir string
	Service     Service
	Version     string
	Imports     []string
}

func (p *Proto) NameWithExt() string {
	return p.Name + ProtoExt
}

func (p *Proto) PkgFilePath() string {
	return filepath.Join(ProjectPathPkgClients, filepath.Join(p.Service.PackageParts...), p.NameWithExt())
}

func (p *Proto) PkgDirPath() string {
	return filepath.Join(ProjectPathPkgClients, filepath.Join(p.Service.PackageParts...))
}

func (p *Proto) SetService(srv, pkg string) error {
	if !ProtoPkgVersionRegexp.MatchString(pkg) {
		return simplerr.Wrap(ErrInvalidProtoPkgName, "use '.v{{ integer }}' at the end of the name")
	}

	p.Version = ProtoPkgVersionRegexp.FindStringSubmatch(pkg)[1]
	p.Service.PackageParts = strings.Split(pkg, ".")

	switch {
	case srv != "":
		p.Service.WithImpl = true
		p.Service.Name = strings.Title(srv)
	default:
		p.Service.Name = strings.Title(p.Service.PackageParts[len(p.Service.PackageParts)-2])
	}

	return nil
}

func AddProtoExt(name string) string {
	return name + ProtoExt
}

type Service struct {
	PackageParts []string
	Name         string
	WithImpl     bool
}

func (s *Service) GoImplImport(module string) string {
	return strings.Join([]string{module, ProjectImportImplementation, strings.Join(s.PackageParts, "/")}, "/")
}

func (s *Service) GoImportAlias() string {
	return helpers.CamelCase(helpers.GoName(strings.Join(s.PackageParts, ".")))
}

func (s *Service) SnakeCasedName() string {
	return helpers.SnakeCase(s.Name)
}
