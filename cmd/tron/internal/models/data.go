package models

import (
	"errors"
	"fmt"
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
}

func (p *Proto) NameWithExt() string {
	return p.Name + ProtoExt
}

func (p *Proto) SetService(module, srv, pkg string) error {
	if !ProtoPkgVersionRegexp.MatchString(pkg) {
		return simplerr.Wrap(ErrInvalidProtoPkgName, "use '.v{{ integer }}' at the end of the name")
	}

	p.Version = ProtoPkgVersionRegexp.FindStringSubmatch(pkg)[1]

	p.Service = Service{
		Name: srv,
		Import: strings.Join([]string{
			module,
			ProjectPathImplementation,
			strings.ToLower(strings.ReplaceAll(pkg, ".", "/")),
		}, "/"),
		Alias:   helpers.CamelCase(helpers.GoName(strings.ReplaceAll(pkg, ".", "_"))),
		Package: strings.ReplaceAll(pkg, ".", string(filepath.Separator)),
	}

	return nil
}

func AddProtoExt(name string) string {
	return name + ProtoExt
}

type Service struct {
	Name    string
	Import  string
	Alias   string
	Package string
}

func (s *Service) SnakeCasedName() string {
	return helpers.SnakeCase(s.Name)
}

func (s *Service) Variable() string {
	return fmt.Sprintf("%sHandler", s.Alias)
}
