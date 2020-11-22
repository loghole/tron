package tron

import (
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"

	"github.com/loghole/tron/internal/app"
)

type RunOption = app.RunOption

func WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) RunOption {
	return func(opts *app.Options) error {
		opts.GRPCOptions = append(opts.GRPCOptions, grpc.ChainUnaryInterceptor(interceptor))

		return nil
	}
}

func WithTLSConfig(config *tls.Config) RunOption {
	return func(opts *app.Options) error {
		opts.TLSConfig = config

		return nil
	}
}

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
