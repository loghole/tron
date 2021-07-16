package models

import (
	"path/filepath"
)

type Dep struct {
	Main    string
	Out     string
	Git     string
	Version string
}

func ProtobufDeps() []*Dep {
	return []*Dep{
		{
			Main:    filepath.Join("grpc-gateway", "protoc-gen-grpc-gateway"),
			Out:     filepath.Join("bin", "protoc-gen-grpc-gateway"),
			Git:     "https://github.com/grpc-ecosystem/grpc-gateway.git",
			Version: "v2.5.0",
		},
		{
			Main:    filepath.Join("grpc-gateway", "protoc-gen-openapiv2"),
			Out:     filepath.Join("bin", "protoc-gen-openapiv2"),
			Git:     "https://github.com/grpc-ecosystem/grpc-gateway.git",
			Version: "v2.5.0",
		},
		{
			Main:    filepath.Join("grpc-go", "cmd", "protoc-gen-go-grpc"),
			Out:     filepath.Join("bin", "protoc-gen-go-grpc"),
			Git:     "https://github.com/grpc/grpc-go.git",
			Version: "cmd/protoc-gen-go-grpc/v1.1.0",
		},
		{
			Main:    filepath.Join("protobuf-go", "cmd", "protoc-gen-go"),
			Out:     filepath.Join("bin", "protoc-gen-go"),
			Git:     "https://github.com/protocolbuffers/protobuf-go.git",
			Version: "v1.26.0",
		},
		{
			Main:    filepath.Join("buf", "cmd", "buf"),
			Out:     filepath.Join("bin", "buf"),
			Git:     "https://github.com/bufbuild/buf.git",
			Version: "v0.43.2",
		},
		{
			Main:    filepath.Join("tron", "cmd", "protoc-gen-tron"),
			Out:     filepath.Join("bin", "protoc-gen-tron"),
			Git:     "https://github.com/loghole/tron.git",
			Version: "v0.17.1-rc1.0",
		},
	}
}
