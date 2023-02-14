package models

import (
	"path/filepath"
)

const (
	sep = string(filepath.Separator)
)

const DefaultGoVersion = "1.20"

const (
	GitDir                        = ".git"
	MainFile                      = "main.go"
	GoModFile                     = "go.mod"
	CmdDir                        = "cmd"
	ProtoExt                      = ".proto"
	DockerfileFilepath            = "build" + sep + "Dockerfile"
	DockerfileDevFilepath         = "build" + sep + "Dockerfile.dev"
	DockerComposeFilepath         = "docker-compose.yaml"
	DockerComposeOverrideFilepath = "docker-compose.override.example.yaml"
	GolangciLintFilepath          = ".golangci.yaml"
	GitignoreFilepath             = ".gitignore"
	TronMKFilepath                = "tron.mk"
	MakefileFilepath              = "Makefile"
	ReadmeMDFilepath              = "README.md"
	BufFilepath                   = "buf.yaml"
	BufGenFilepath                = "buf.gen.yaml"
)

const (
	ProjectPathAPI        = "api"
	ProjectPathPkgClients = "pkg"
	ProjectPathVendorPB   = "vendor.pb"
)
