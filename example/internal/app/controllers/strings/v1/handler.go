package stringsV1

import (
	stringsV1 "example/pkg/strings/v1"

	"github.com/loghole/tron/transport"
)

type Implementation struct {
	stringsV1.UnimplementedStringsServer
}

func NewImplementation() *Implementation {
	return &Implementation{}
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *Implementation) GetDescription() transport.ServiceDesc {
	return stringsV1.NewStringsServiceDesc(i)
}
