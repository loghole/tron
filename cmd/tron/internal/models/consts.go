package models

import (
	"path/filepath"

	"github.com/loghole/tron/internal/app"
)

const (
	sep = string(filepath.Separator)

	dockerfilePath = app.DeploymentsDir + sep + "docker" + sep
	valuesDirPath  = app.DeploymentsDir + sep + app.ValuesDir + sep
)

const (
	GitDir    = ".git"
	MainFile  = "main.go"
	GoModFile = "go.mod"
	CmdDir    = "cmd"
	ProtoExt  = ".proto"

	ValuesBaseFilepath   = valuesDirPath + app.ValuesBaseName + "." + app.ValuesExt
	ValuesDevFilepath    = valuesDirPath + app.ValuesDevName + "." + app.ValuesExt
	ValuesLocalFilepath  = valuesDirPath + app.ValuesLocalName + "." + app.ValuesExt
	ValuesStgFilepath    = valuesDirPath + app.ValuesStgName + "." + app.ValuesExt
	ValuesProdFilepath   = valuesDirPath + app.ValuesProdName + "." + app.ValuesExt
	DockerfileFilepath   = dockerfilePath + "Dockerfile"
	GolangciLintFilepath = ".golangci.yaml"
	GitignoreFilepath    = ".gitignore"
	TronMKFilepath       = "tron.mk"
	MakefileFilepath     = "Makefile"
	ReadmeMDFilepath     = "README.md"
	ConfigConstFilepath  = "config" + sep + "constants.go"
	ConfigFilepath       = "config" + sep + "config.go"
	BufFilepath          = "buf.yaml"
	BufGenFilepath       = "buf.gen.yaml"
)

const (
	ProjectPathAPI        = "api"
	ProjectPathPkgClients = "pkg"
	ProjectPathVendorPB   = "vendor.pb"
)

const (
	TronOptionsSep = ","
	TronOptionsTag = "tron_option:"
	TronOptionJSON = "json"
)
