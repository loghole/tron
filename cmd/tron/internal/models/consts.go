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
	GitDir         = ".git"
	MainFile       = "main.go"
	GoModFile      = "go.mod"
	CmdDir         = "cmd"
	InternalDir    = "internal"
	AppDir         = "app"
	ControllersDir = "controllers"

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
)

const (
	ProjectPathAPI            = "api"
	ProjectPathImplementation = InternalDir + sep + AppDir + sep + ControllersDir
	ProjectPathPkgClients     = "pkg"
	ProjectPathVendorPB       = "vendor.pb"
)

const (
	ProjectImportImplementation = InternalDir + "/" + AppDir + "/" + ControllersDir
)
