package http

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/grpclog"

	internalErr "github.com/loghole/tron/internal/errors"
)

const fallback = `{"code": 13, "message": "failed to marshal error message"}`

// ErrorWriter returns runtime.ErrorHandlerFunc to configure error handling.
func ErrorWriter() runtime.ErrorHandlerFunc {
	return func(
		ctx context.Context,
		mux *runtime.ServeMux,
		marshaler runtime.Marshaler,
		w http.ResponseWriter,
		r *http.Request,
		err error) {
		s := internalErr.ParseError(err)

		buf, merr := jsoniter.Marshal(s)
		if merr != nil {
			grpclog.Infof("Failed to marshal error message %q: %v", s, merr)

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, fallback)
		}

		w.Header().Del("Trailer")
		w.Header().Del("Transfer-Encoding")

		contentType := marshaler.ContentType(buf)
		w.Header().Set("Content-Type", contentType)

		w.WriteHeader(s.HTTPStatus())

		if _, err := w.Write(buf); err != nil {
			grpclog.Infof("Failed to marshal error message %q: %v", s, merr)

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, fallback)
		}
	}
}
