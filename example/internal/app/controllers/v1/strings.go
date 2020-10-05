// Code generated by protoc-gen-goclay, but your can (must) modify it.
// source: strings.proto

package v1

import (
	desc "example/pkg/v1"

	"github.com/utrack/clay/v2/transport"
)

type Implementation struct{}

// NewStrings create new Implementation
func NewStrings() *Implementation {
	return &Implementation{}
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *Implementation) GetDescription() transport.ServiceDesc {
	return desc.NewStringsServiceDesc(i)
}
