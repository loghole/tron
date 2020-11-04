package tron

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/lissteron/simplerr"
	"github.com/loghole/lhw/zap"
	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracehttp"
	"github.com/loghole/tracing/tracelog"
	"github.com/utrack/clay/v2/transport"
	"golang.org/x/sync/errgroup"

	"github.com/loghole/tron/internal/admin"
	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/internal/config"
	"github.com/loghole/tron/internal/grpc"
)

type App struct {
	info    *Info
	opts    *app.Options
	servers servers
	logger  logger
	tracer  tracer
}

// New init viper config, logger and tracer.
func New(options ...app.Option) (*App, error) {
	opts, err := app.NewOptions(options...)
	if err != nil {
		return nil, simplerr.Wrap(err, "apply opts failed")
	}

	if err := config.Init(); err != nil {
		return nil, simplerr.Wrap(err, "init config failed")
	}

	a := &App{opts: opts, info: initInfo()}

	if err := a.logger.init(a.info, a.opts); err != nil {
		return nil, err
	}

	if err := a.tracer.init(a.info); err != nil {
		return nil, err
	}

	if err := a.servers.init(a.opts); err != nil {
		return nil, simplerr.Wrap(err, "init servers failed")
	}

	// Append tracing and errors middlewares.
	a.opts.AddRunOptions(
		WithUnaryInterceptor(grpc.OpenTracingServerInterceptor(a.Tracer())),
		WithUnaryInterceptor(grpc.SimpleErrorServerInterceptor()),
	)
	a.Router().Use(tracehttp.NewMiddleware(a.Tracer()).Middleware)
	a.Router().Use(cors.New(a.opts.CorsOptions).Handler)

	a.logger.With("app info", a.info).Infof("init app")

	return a, nil
}

func (a *App) Info() *Info {
	return a.info
}

func (a *App) Tracer() *tracing.Tracer {
	return a.tracer.tracer
}

func (a *App) Logger() *zap.Logger {
	return a.logger.Logger
}

func (a *App) TraceLogger() tracelog.Logger {
	return a.logger.tracelog
}

func (a *App) Router() chi.Router {
	return a.servers.publicHTTP.Router()
}

// Close closes tracer and logger.
func (a *App) Close() {
	a.tracer.Close()
	a.logger.Close()
}

// Append some run options.
func (a *App) WithRunOptions(opts ...app.RunOption) *App {
	a.opts.AddRunOptions(opts...)

	return a
}

// Run apply run options if exists and starts servers.
func (a *App) Run(impl ...transport.Service) error { // nolint:funlen // start app
	if err := a.opts.ApplyRunOptions(); err != nil {
		return simplerr.Wrap(err, "apply run options failed")
	}

	if err := a.servers.build(a.opts); err != nil {
		return simplerr.Wrap(err, "build servers failed")
	}

	a.servers.publicGRPC.RegistryDesc(impl...)
	a.servers.publicHTTP.RegistryDesc(impl...)

	admin.NewHandlers(a.info, impl...).InitRoutes(a.servers.adminHTTP.Router())

	a.logger.Info("starting app")

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		a.logger.Infof("start public grpc server on: %s", a.servers.publicGRPC.Addr())

		return a.servers.publicGRPC.Serve()
	})

	eg.Go(func() error {
		a.logger.Infof("start public http server on: %s", a.servers.publicHTTP.Addr())

		return a.servers.publicHTTP.Serve()
	})

	eg.Go(func() error {
		a.logger.Infof("start admin http server on: %s", a.servers.adminHTTP.Addr())

		return a.servers.adminHTTP.Serve()
	})

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, a.opts.ExitSignals...)

	select {
	case <-exit:
		a.logger.Info("stopping application")
	case <-ctx.Done():
		a.logger.Errorf("stopping application with error: %v", ctx.Err())
	}

	signal.Stop(exit)

	if err := a.servers.publicHTTP.Close(); err != nil {
		a.logger.Errorf("error while stopping public http server: %v", err)
	}

	if err := a.servers.publicGRPC.Close(); err != nil {
		a.logger.Errorf("error while stopping public grpc server: %v", err)
	}

	if err := a.servers.publicHTTP.Close(); err != nil {
		a.logger.Errorf("error while stopping admin http server: %v", err)
	}

	_ = a.logger.Sync()

	return nil
}
