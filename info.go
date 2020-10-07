package tron

import (
	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

// SetName overrides default the application name.
func SetName(name string) {
	app.ServiceName = name
}

type Info struct {
	InstanceUUID string
	ServiceName  string
	Namespace    string
	AppName      string
	GitHash      string
	Version      string
	BuildAt      string
}

func initInfo() *Info {
	return &Info{
		InstanceUUID: app.InstanceUUID.String(),
		ServiceName:  app.ServiceName,
		AppName:      app.AppName,
		Namespace:    app.ParseNamespace(viper.GetString(app.NamespaceEnv)).String(),
		GitHash:      app.GitHash,
		Version:      app.Version,
		BuildAt:      app.BuildAt,
	}
}
