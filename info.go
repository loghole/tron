package tron

import (
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
