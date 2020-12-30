package models

import (
	"path/filepath"
	"strings"
)

const (
	ProtoExt     = ".proto"
	PbGoExt      = ".pb.go"
	PbSwaggerExt = ".swagger.json"
	PbTronExt    = ".pb.tron.go"
)

type Proto struct {
	Name      string
	Path      string
	RelPath   string
	Package   string
	GoPackage string
	Imports   []string
	Service   *Service
}

func (p *Proto) WithImpl() bool {
	return p.Service != nil
}

func (p *Proto) PkgDir() string {
	return filepath.Join(ProjectPathPkgClients, strings.ReplaceAll(p.Package, ".", sep))
}

func (p *Proto) PkgTypesFile() string {
	return filepath.Join(p.PkgDir(), p.Name+PbGoExt)
}

func (p *Proto) PkgSwaggerFile() string {
	return filepath.Join(p.PkgDir(), p.Name+PbSwaggerExt)
}

func (p *Proto) PkgTronFile() string {
	return filepath.Join(p.PkgDir(), p.Name+PbTronExt)
}

type Service struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Name   string
	Input  string
	Output string
}
