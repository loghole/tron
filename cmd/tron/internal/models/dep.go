package models

import (
	"fmt"
	"path/filepath"

	"github.com/Masterminds/semver"

	"github.com/loghole/tron/cmd/tron/internal/version"
)

type Dep struct {
	Main    string
	Out     string
	Git     string
	Version string
	LdFlag  string
}

func (d *Dep) IsActual(actual string) (bool, error) {
	targetStr, err := ExtractVersion(d.Version)
	if err != nil {
		return false, fmt.Errorf("extract target version: %w", err)
	}

	actualStr, err := ExtractVersion(actual)
	if err != nil {
		return false, fmt.Errorf("extract actual version: %w", err)
	}

	targetVer, err := semver.NewVersion(targetStr)
	if err != nil {
		return false, fmt.Errorf("parse target version: %w", err)
	}

	actualVer, err := semver.NewVersion(actualStr)
	if err != nil {
		return false, fmt.Errorf("parse actual version: %w", err)
	}

	return targetVer.Equal(actualVer), nil
}

func ProtobufDeps() []*Dep {
	return []*Dep{
		{
			Main:    filepath.Join("grpc-gateway", "protoc-gen-grpc-gateway"),
			Out:     filepath.Join("bin", "protoc-gen-grpc-gateway"),
			Git:     "https://github.com/grpc-ecosystem/grpc-gateway.git",
			Version: "v2.7.3",
			LdFlag:  "main.version=v2.7.3",
		},
		{
			Main:    filepath.Join("grpc-gateway", "protoc-gen-openapiv2"),
			Out:     filepath.Join("bin", "protoc-gen-openapiv2"),
			Git:     "https://github.com/grpc-ecosystem/grpc-gateway.git",
			Version: "v2.7.3",
			LdFlag:  "main.version=v2.7.3",
		},
		{
			Main:    filepath.Join("grpc-go", "cmd", "protoc-gen-go-grpc"),
			Out:     filepath.Join("bin", "protoc-gen-go-grpc"),
			Git:     "https://github.com/grpc/grpc-go.git",
			Version: "cmd/protoc-gen-go-grpc/v1.2.0",
		},
		{
			Main:    filepath.Join("protobuf-go", "cmd", "protoc-gen-go"),
			Out:     filepath.Join("bin", "protoc-gen-go"),
			Git:     "https://github.com/protocolbuffers/protobuf-go.git",
			Version: "v1.27.1",
		},
		{
			Main:    filepath.Join("buf", "cmd", "buf"),
			Out:     filepath.Join("bin", "buf"),
			Git:     "https://github.com/bufbuild/buf.git",
			Version: "v1.0.0",
		},
		{
			Main:    filepath.Join("tron", "cmd", "protoc-gen-tron"),
			Out:     filepath.Join("bin", "protoc-gen-tron"),
			Git:     "https://github.com/loghole/tron.git",
			Version: version.CliVersion,
			LdFlag:  "main.Version=" + version.CliVersion,
		},
	}
}
