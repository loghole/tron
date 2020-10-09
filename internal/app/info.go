package app

import (
	"github.com/google/uuid"
)

// nolint:gochecknoglobals // build args
var (
	InstanceUUID = uuid.New()
	ServiceName  = "-"
	AppName      = "-"
	GitHash      = "-"
	Version      = "-"
	BuildAt      = "0001-01-01T00:00:00"
)
