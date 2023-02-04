package app

import (
	"flag"
)

//nolint:gochecknoglobals //flags
var (
	localConfigEnabled bool
)

func parseFlags() {
	flag.BoolVar(&localConfigEnabled, "local-config-enabled", false, "enable local config YAML parser")

	flag.Parse()
}
