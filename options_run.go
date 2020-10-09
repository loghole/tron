package tron

import (
	"net/http"

	"google.golang.org/grpc"

	"github.com/loghole/tron/internal/app"
)

type RunOption = app.RunOption

func WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) RunOption {
	return func(opts *app.Options) {
		opts.GRPCOptions = append(opts.GRPCOptions, grpc.ChainUnaryInterceptor(interceptor))
	}
}

func WithHTTPMiddleware(middleware func(http.Handler) http.Handler) RunOption {
	return func(opts *app.Options) {
		opts.HTTPMiddlewares = append(opts.HTTPMiddlewares, middleware)
	}
}
