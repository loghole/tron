package http

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/utrack/clay/v2/transport/httpruntime"

	internalErr "github.com/loghole/tron/internal/errors"
)

func setClayErrorWriter() {
	httpruntime.SetError = func(context context.Context, r *http.Request, w http.ResponseWriter, err error) {
		w.Header().Set("Content-Type", "application/json")

		resp := internalErr.ParseError(err)

		w.WriteHeader(resp.Status)
		_ = jsoniter.NewEncoder(w).Encode(resp)
	}
}
