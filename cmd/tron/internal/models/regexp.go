package models

import (
	"regexp"
)

// nolint:gochecknoglobals //regexp
var (
	ProtoPkgVersionRegexp = regexp.MustCompile(`\.(v\d+)$`)
	VersionRegexp         = regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	ImportRegexp          = regexp.MustCompile(`^import "([^"]+)";$`)
	ServiceRegexp         = regexp.MustCompile(`^service (\S+) {`)
	PackageRegexp         = regexp.MustCompile(`^package ([^;]+);$`)
	GoPackageRegexp       = regexp.MustCompile(`^option go_package "([^"]+)"$`)
	ModuleRegexp          = regexp.MustCompile(`^module (.+)$`)
)
