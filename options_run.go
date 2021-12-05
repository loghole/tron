package tron

import (
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"

	"github.com/loghole/tron/internal/app"
)

// WithUnaryInterceptor returns a RunOption that specifies the chained interceptor
// for unary RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All unary interceptors added by this method will be chained.
func WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) app.RunOption {
	return func(opts *app.Options) error {
		opts.GRPCOptions = append(opts.GRPCOptions, grpc.ChainUnaryInterceptor(interceptor))

		return nil
	}
}

// WithStreamInterceptor returns a RunOption that specifies the chained interceptor
// for streaming RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All stream interceptors added by this method will be chained.
func WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) app.RunOption {
	return func(opts *app.Options) error {
		opts.GRPCOptions = append(opts.GRPCOptions, grpc.ChainStreamInterceptor(interceptor))

		return nil
	}
}

// WithTLSConfig returns a RunOption that set tls configuration for grpc and http servers.
func WithTLSConfig(config *tls.Config) app.RunOption {
	return func(opts *app.Options) error {
		opts.TLSConfig = config

		return nil
	}
}

// WithTLSKeyPair returns a RunOption that set tls configuration for grpc and http servers from files.
func WithTLSKeyPair(certFile, keyFile string) app.RunOption {
	if certFile == "" || keyFile == "" {
		return nil
	}

	return func(opts *app.Options) (err error) {
		opts.TLSConfig = &tls.Config{} // nolint:gosec // default by http.ListenAndServeTLS
		opts.TLSConfig.Certificates = make([]tls.Certificate, 1)

		opts.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("load X509 key pair failed: %w", err)
		}

		return nil
	}
}
