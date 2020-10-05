package tron

import (
	"os"

	"github.com/loghole/tron/internal/app"
)

type Option = app.Option

func WithAdminHTTP(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortAdmin = port

		return nil
	}
}

func WithPublicHTTP(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortHTTP = port

		return nil
	}
}

func WithPublicGRPc(port uint16) Option {
	return func(opts *app.Options) error {
		opts.PortGRPC = port

		return nil
	}
}

func WithExitSignals(sig ...os.Signal) Option {
	return func(opts *app.Options) error {
		opts.ExitSignals = append(opts.ExitSignals, sig...)

		return nil
	}
}
