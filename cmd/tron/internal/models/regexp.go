package models

import (
	"regexp"
)

// nolint:gochecknoglobals //regexp
var (
	ProtoPkgVersionRegexp = regexp.MustCompile(`\.(v\d+)$`)
	VersionRegexp         = regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	PackageRegexp         = regexp.MustCompile(`^package ([^;]+);$`)
	ImportRegexp          = regexp.MustCompile(`^import "([^"]+)";$`)
	ModuleRegexp          = regexp.MustCompile(`^module (.+)$`)
	TronOptions           = regexp.MustCompile(`tron_option:(\S+)`)
)
