package models

import (
	"github.com/loghole/tron/cmd/tron/internal/helpers"
)

const ProtoExt = ".proto"

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
}

func (p *Proto) NameWithExt() string {
	return p.Name + ProtoExt
}

func AddProtoExt(name string) string {
	return name + ProtoExt
}

type Service struct {
	Name           string
	KebabCasedName string
	Package        string
	PackageName    string
}

func (s *Service) SnakeCasedName() string {
	return helpers.SnakeCase(s.Name)
}
