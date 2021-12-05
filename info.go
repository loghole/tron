package tron

import (
	"github.com/loghole/tron/config"
	"github.com/loghole/tron/internal/app"
)

// Info contains service information.
type Info = app.Info

// SetName overrides default the application name.
func SetName(name string) {
	app.ServiceName = name
}

// GetInfo returns base service information.
func GetInfo() *Info {
	return initInfo()
}

func initInfo() *Info {
	return &Info{
		InstanceUUID: app.InstanceUUID.String(),
		ServiceName:  app.ServiceName,
		AppName:      app.AppName,
		Namespace:    app.ParseNamespace(config.GetString(app.NamespaceEnv)).String(),
		GitHash:      app.GitHash,
		Version:      app.Version,
		BuildAt:      app.BuildAt,
		StartTime:    app.StartTime,
	}
}
