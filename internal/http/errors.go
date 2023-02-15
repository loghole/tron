package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	internalErr "github.com/loghole/tron/internal/errors"
)

const fallback = `{"code": 500, "message": "failed to marshal error message"}`

// ErrorWriter returns runtime.ErrorHandlerFunc to configure error handling.
func ErrorWriter() runtime.ErrorHandlerFunc {
	return func(
		ctx context.Context,
		mux *runtime.ServeMux,
		marshaler runtime.Marshaler,
		w http.ResponseWriter,
		r *http.Request,
		err error,
	) {
		parsed := internalErr.ParseError(ctx, err)

		buf, merr := json.Marshal(parsed)
		if merr != nil {
			w.WriteHeader(http.StatusInternalServerError)

			_, _ = io.WriteString(w, fallback)
		}

		w.Header().Del("Trailer")
		w.Header().Del("Transfer-Encoding")
		w.Header().Set("Content-Type", marshaler.ContentType(buf))

		w.WriteHeader(parsed.HTTPStatus())

		if _, err := w.Write(buf); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			_, _ = io.WriteString(w, fallback)
		}
	}
}
