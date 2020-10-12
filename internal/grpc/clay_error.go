package grpc

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	jsoniter "github.com/json-iterator/go"
	"github.com/utrack/clay/v2/transport/httpruntime"
	grpcStatus "google.golang.org/grpc/status"
)

func setClayErrorWriter() {
	httpruntime.SetError = func(context context.Context, r *http.Request, w http.ResponseWriter, err error) {
		w.Header().Set("Content-Type", "application/json")

		status, ok := grpcStatus.FromError(err)
		switch {
		case ok:
			w.WriteHeader(runtime.HTTPStatusFromCode(status.Code()))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = jsoniter.NewEncoder(w).Encode(status.Proto())
	}
}
