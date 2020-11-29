// Package transport contains base service desc interfaces.
package transport

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Service is a registerable collection of endpoints.
// These functions should be autogenerated by protoc-gen-tron.
type Service interface {
	GetDescription() ServiceDesc
}

// ServiceDesc is a description of an endpoints' collection.
// These functions should be autogenerated by protoc-gen-tron.
type ServiceDesc interface {
	RegisterGRPC(server *grpc.Server)
	RegisterHTTP(mux *runtime.ServeMux)
	SwaggerDef() []byte
}
