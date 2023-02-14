// Package tron contains base application object, some options and functions to create base app.
package tron

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracegrpc"
	"github.com/loghole/tracing/tracehttp"
	"github.com/loghole/tracing/tracelog"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/loghole/tron/healthcheck"
	"github.com/loghole/tron/internal/admin"
	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/internal/grpc"
	"github.com/loghole/tron/transport"
)

// App is base application object with grpc and http servers, logger and tracer.
type App struct {
	info    *Info
	opts    *app.Options
	servers servers
	logger  logger
	tracer  tracer
	health  health

	errGroupCtx context.Context //nolint:containedctx // need internal context
	errGroup    *errgroup.Group
}

// New init viper config, logger and tracer.
func New(options ...app.Option) (*App, error) {
	info := initInfo()

	opts, err := app.NewOptions(info, options...)
	if err != nil {
		return nil, fmt.Errorf("apply opts failed: %w", err)
	}

	application := &App{opts: opts, info: info}
	application.health.init()

	if err := application.logger.init(application.opts); err != nil {
		return nil, err
	}

	if err := application.tracer.init(application.opts); err != nil {
		return nil, err
	}

	if err := application.servers.init(application.logger, application.opts); err != nil {
		return nil, fmt.Errorf("init servers failed: %w", err)
	}

	// Append recover, tracing and errors middlewares.
	application.opts.AddRunOptions(
		WithUnaryInterceptor(grpc.RecoverServerInterceptor(application.logger.tracelog)),
		WithUnaryInterceptor(tracegrpc.UnaryServerInterceptor(application.Tracer())),
		WithStreamInterceptor(tracegrpc.StreamServerInterceptor(application.Tracer())),
		WithUnaryInterceptor(grpc.SimpleErrorServerInterceptor()),
	)

	application.servers.publicHTTP.UseMiddleware(tracehttp.Handler(application.Tracer()))

	application.logger.With("app info", application.info).Infof("init app")

	application.errGroup, application.errGroupCtx = errgroup.WithContext(context.Background())

	return application, nil
}

// Info returns application info.
func (a *App) Info() *Info {
	return a.info
}

// Tracer returns wrapped jaeger tracer.
func (a *App) Tracer() *tracing.Tracer {
	return a.tracer.tracer
}

// Logger returns default zap logger.
func (a *App) Logger() *zap.SugaredLogger {
	return a.logger.SugaredLogger
}

// TraceLogger returns wrapped zap logger with opentracing metadata injection to log records.
func (a *App) TraceLogger() tracelog.Logger {
	return a.logger.tracelog
}

// Router returns http router that runs on public port.
func (a *App) Router() chi.Router {
	return a.servers.publicHTTP.Router()
}

// Health returns application health checker.
func (a *App) Health() healthcheck.Checker {
	return a.health
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (a *App) Go(f func() error) {
	a.errGroup.Go(f)
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (a *App) Wait() error {
	return a.errGroup.Wait() //nolint:wrapcheck // need clean err.
}

// Close closes tracer and logger.
func (a *App) Close() {
	_ = a.logger.Sync()

	a.tracer.Close()
}

// WithRunOptions appends some run options.
func (a *App) WithRunOptions(opts ...app.RunOption) *App {
	a.opts.AddRunOptions(opts...)

	return a
}

// Run apply run options if exists and starts servers.
func (a *App) Run(impl ...transport.Service) error {
	if err := a.opts.ApplyRunOptions(); err != nil {
		return fmt.Errorf("apply run options failed: %w", err)
	}

	if err := a.servers.build(a.opts); err != nil {
		return fmt.Errorf("build servers failed: %w", err)
	}

	a.servers.publicGRPC.RegistryDesc(impl...)
	a.servers.publicHTTP.RegistryDesc(impl...)

	admin.NewHandlers(a.info, a.opts, a.health, impl...).InitRoutes(a.servers.adminHTTP.Router())

	a.logger.Info("starting app")

	a.servers.serve(a.errGroup)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, a.opts.ExitSignals...)

	a.health.setReady()

	select {
	case <-exit:
		a.logger.Info("stopping application")
	case <-a.errGroupCtx.Done():
		a.logger.Errorf("stopping application with error: %v", a.errGroupCtx.Err())
	}

	a.health.setUnready()

	signal.Stop(exit)

	a.servers.close()

	return nil
}
