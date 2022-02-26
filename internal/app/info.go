// Package app contains base application constants, flags and options.
package app

import (
	"time"

	"github.com/google/uuid"
)

// nolint:gochecknoglobals // build args
var (
	InstanceUUID = uuid.New()
	StartTime    = time.Now()
	ServiceName  = "-"
	AppName      = "-"
	GitHash      = "-"
	Version      = "-"
	BuildAt      = "0001-01-01T00:00:00"
)

// Info contains service information.
type Info struct {
	InstanceUUID string
	ServiceName  string
	Namespace    Namespace
	AppName      string
	GitHash      string
	Version      string
	BuildAt      string
	StartTime    time.Time
}
