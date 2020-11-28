package tron

import (
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"

	"github.com/loghole/tron/internal/app"
)

// RunOption sets tron run options such as grpc unary interceptors, tls config.
type RunOption = app.RunOption

// WithUnaryInterceptor returns a RunOption that specifies the chained interceptor
// for unary RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All unary interceptors added by this method will be chained.
func WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) RunOption {
	return func(opts *app.Options) error {
		opts.GRPCOptions = append(opts.GRPCOptions, grpc.ChainUnaryInterceptor(interceptor))

		return nil
	}
}

// WithTLSConfig returns a RunOption that set tls configuration for grpc and http servers.
func WithTLSConfig(config *tls.Config) RunOption {
	return func(opts *app.Options) error {
		opts.TLSConfig = config

		return nil
	}
}

// WithTLSKeyPair returns a RunOption that set tls configuration for grpc and http servers from files.
func WithTLSKeyPair(certFile, keyFile string) RunOption {
	if certFile == "" || keyFile == "" {
		return nil
	}

	return func(opts *app.Options) (err error) {
		opts.TLSConfig = &tls.Config{} // nolint:gosec // default by http.ListenAndServeTLS
		opts.TLSConfig.Certificates = make([]tls.Certificate, 1)
		opts.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

		return fmt.Errorf("load X509 key pair failed: %w", err)
	}
}
