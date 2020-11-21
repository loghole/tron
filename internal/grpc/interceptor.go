package grpc

import (
	"context"
	"runtime/debug"
	"strings"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/loghole/tracing/tracelog"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	optlog "github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	internalErr "github.com/loghole/tron/internal/errors"
)

const ComponentName = "net/grpc"

func OpenTracingServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		spanContext, _ := extractSpanContext(ctx, tracer)

		span := tracer.StartSpan(defaultNameFunc(info), ext.RPCServerOption(spanContext))
		defer span.Finish()

		ext.Component.Set(span, ComponentName)

		ctx = opentracing.ContextWithSpan(ctx, span)

		resp, err = handler(ctx, req)
		if err != nil {
			otgrpc.SetSpanTags(span, err, false)
			span.LogFields(optlog.String("event", "error"), optlog.String("message", err.Error()))
		}

		return resp, err
	}
}

func SimpleErrorServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			return resp, internalErr.ParseError(err)
		}

		return resp, err
	}
}

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

func extractSpanContext(ctx context.Context, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	return tracer.Extract(opentracing.HTTPHeaders, metadataReaderWriter{md})
}

type metadataReaderWriter struct {
	metadata.MD
}

func (w metadataReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)

	w.MD[key] = append(w.MD[key], val)
}

func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func defaultNameFunc(r *grpc.UnaryServerInfo) string {
	return strings.Join([]string{"GRPC", r.FullMethod}, " ")
}
