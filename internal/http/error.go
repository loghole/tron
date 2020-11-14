package http

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func ErrorWriter() runtime.ErrorHandlerFunc {
	return runtime.DefaultHTTPErrorHandler // TODO: add custom err writer
}
