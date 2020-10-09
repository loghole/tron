package tron

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-chi/chi"
	"github.com/lissteron/simplerr"
	"github.com/loghole/lhw/zap"
	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracelog"
	"github.com/spf13/viper"
	"github.com/utrack/clay/v2/transport"
	"golang.org/x/sync/errgroup"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/internal/config"
	"github.com/loghole/tron/internal/grpc"
	"github.com/loghole/tron/internal/http"
	"github.com/loghole/tron/internal/swagger"
)

type servers struct {
	publicGRPC *grpc.Server
	publicHTTP *http.Server
	adminHTTP  *http.Server
}

func (s *servers) init(opts *app.Options) (err error) {
	if opts.PortAdmin == 0 {
		opts.PortAdmin = uint16(viper.GetInt32(app.AdminPortEnv))
	}

	if opts.PortHTTP == 0 {
		opts.PortHTTP = uint16(viper.GetInt32(app.HTTPPortEnv))
	}

	if opts.PortGRPC == 0 {
		opts.PortGRPC = uint16(viper.GetInt32(app.GRPCPortEnv))
	}

	s.publicGRPC, err = grpc.NewServer(opts.PortGRPC, opts.GRPCOptions...)
	if err != nil {
		return err
	}

	s.publicHTTP, err = http.NewServer(opts.PortHTTP, opts.TLSConfig)
	if err != nil {
		return err
	}

	s.adminHTTP, err = http.NewServer(opts.PortAdmin, nil)
	if err != nil {
		return err
	}

	return nil
}

type logger struct {
	*zap.Logger
	traceLogger tracelog.Logger
}

func (l *logger) init(info *Info) (err error) {
	l.Logger, err = zap.NewLogger(&zap.Config{
		Level:         viper.GetString(app.LoggerLevelEnv),
		CollectorURL:  viper.GetString(app.LoggerCollectorAddrEnv),
		Hostname:      info.ServiceName,
		Namespace:     info.Namespace,
		Source:        info.ServiceName,
		BuildCommit:   info.GitHash,
		DisableStdout: viper.GetBool(app.LoggerDisableStdoutEnv),
	})
	if err != nil {
		return err
	}

	l.traceLogger = tracelog.NewTraceLogger(l.Logger.SugaredLogger)

	return nil
}

type tracer struct {
	tracer *tracing.Tracer
}

func (t *tracer) init(info *Info) (err error) {
	t.tracer, err = tracing.NewTracer(tracing.DefaultConfiguration(info.ServiceName, viper.GetString(app.JaegerAddrEnv)))
	if err != nil {
		return err
	}

	return nil
}

type App struct {
	info *Info
	opts *app.Options
	servers
	logger
	tracer
}

func New(options ...Option) (*App, error) {
	opts, err := app.NewOptions(options...)
	if err != nil {
		return nil, err
	}

	if err := config.Init(); err != nil {
		return nil, simplerr.Wrap(err, "init config failed")
	}

	a := &App{opts: opts, info: initInfo()}

	if err := a.logger.init(a.info); err != nil {
		return nil, simplerr.Wrap(err, "init logger failed")
	}

	if err := a.tracer.init(a.info); err != nil {
		return nil, simplerr.Wrap(err, "init tracer failed")
	}

	if err := a.servers.init(opts); err != nil {
		return nil, simplerr.Wrap(err, "init servers failed")
	}

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
	return a.traceLogger
}

func (a *App) Router() chi.Router {
	return a.publicHTTP.Router()
}

func (a *App) Run(impl ...transport.Service) {
	descs := make([]transport.ServiceDesc, 0, len(impl))

	for _, service := range impl {
		descs = append(descs, service.GetDescription())
	}

	desc := transport.NewCompoundServiceDesc(descs...)

	a.logger.Info("starting app")

	a.publicGRPC.RegistryDesc(desc)
	a.publicHTTP.RegistryDesc(desc)

	swagger.New(desc, a.info.Version).InitRoutes(a.adminHTTP.Router())

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		a.logger.Infof("start public grpc server on: %s", a.publicGRPC.Addr())

		return a.publicGRPC.Serve()
	})

	eg.Go(func() error {
		a.logger.Infof("start public http server on: %s", a.publicHTTP.Addr())

		return a.publicHTTP.Serve()
	})

	eg.Go(func() error {
		a.logger.Infof("start admin http server on: %s", a.adminHTTP.Addr())

		return a.adminHTTP.Serve()
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

	if err := a.publicHTTP.Close(); err != nil {
		a.logger.Errorf("error while stopping public http server: %v", err)
	}

	if err := a.publicGRPC.Close(); err != nil {
		a.logger.Errorf("error while stopping public grpc server: %v", err)
	}

	if err := a.publicHTTP.Close(); err != nil {
		a.logger.Errorf("error while stopping admin http server: %v", err)
	}

	_ = a.logger.Sync()
}
