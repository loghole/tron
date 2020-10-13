package models

import (
	"regexp"
)

// nolint:gochecknoglobals //regexp
var (
	VersionRegexp = regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	ImportRegexp  = regexp.MustCompile(`^import "([^"]+)";$`)
	ServiceRegexp = regexp.MustCompile(`^service (\S+) {`)
	PackageRegexp = regexp.MustCompile(`^package ([^;]+);$`)
	ModuleRegexp  = regexp.MustCompile(`^module (.+)$`)

	GoNameRexp      = regexp.MustCompile(`[^a-zA-Z0-9\s_-]+`)
	FirstDigitsRexp = regexp.MustCompile(`^\d+`)
)
