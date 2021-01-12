package models

import (
	"bytes"
	"strings"
)

type Project struct {
	Version string

	// parsed data.
	AbsPath     string
	Module      string
	Name        string
	ProtoFiles  []string
	ValuesEnv   []*ConfigValue
	Protos      []*Proto
	ProtoPkgMap []string

	// from flags.
	ProtoDirs []string
}

func (p *Project) WithProtos() bool {
	return len(p.ProtoFiles) > 0
}

func (p *Project) Dockerfile() string {
	return DockerfileFilepath
}

func (p *Project) DockerImage() string {
	if parts := strings.Split(p.Module, "/"); len(parts) > 1 {
		return strings.Join(parts[1:], "/")
	}

	return p.Module
}

func (p *Project) ServiceName() string {
	return p.Name
}

func (p *Project) AppName() string {
	return p.Module
}

func (p *Project) Mainfile() string {
	return strings.Join([]string{CmdDir, p.Name, MainFile}, "/")
}

func (p *Project) GenerateCmd() string {
	buf := bytes.NewBufferString("tron generate -v")

	for _, val := range p.ProtoDirs {
		if val == "" {
			continue
		}

		buf.WriteString(" --proto=")
		buf.WriteString(val)
	}

	return buf.String()
}

type ConfigValue struct {
	Name string
	Key  string
}
