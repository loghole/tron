package http

import (
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsoniter "github.com/json-iterator/go"
)

// Implement runtime.Marshaler
type marshaler struct {
	jsoniter.API
}

func newMarshaler() runtime.Marshaler {
	return &marshaler{
		API: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (m *marshaler) ContentType(_ interface{}) string {
	return "application/json"
}

func (m *marshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return m.API.NewDecoder(r)
}

func (m *marshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return m.API.NewEncoder(w)
}
