package app

import (
	"flag"
)

var (
	localConfigEnabled bool
)

func parseFlags() {
	flag.BoolVar(&localConfigEnabled, "local-config-enabled", false, "enable local config YAML parser")

	flag.Parse()
}
