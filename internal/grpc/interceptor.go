package grpc

import (
	"context"
	"runtime/debug"

	"github.com/loghole/tracing/tracelog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	internalErr "github.com/loghole/tron/internal/errors"
)

// SimpleErrorServerInterceptor returns error parser grpc interceptor.
func SimpleErrorServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			return resp, internalErr.ParseError(ctx, err)
		}

		return resp, err
	}
}

// RecoverServerInterceptor returns recovery grpc interceptor.
func RecoverServerInterceptor(logger tracelog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "recovered from panic: %v", r)

				logger.With("stack_trace", string(debug.Stack())).
					Errorf(ctx, "panic in '%s', text %v", info.FullMethod, r)
			}
		}()

		return handler(ctx, req)
	}
}
