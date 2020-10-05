package models

import (
	"regexp"
)

var (
	VersionRegexp = regexp.MustCompile(`([0-9]+)\.([0-9]+)(\.([0-9]+))?(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?`)
	ImportRegexp  = regexp.MustCompile(`^import "(.*?)";$`)
	ServiceRegexp = regexp.MustCompile(`^service (.*?) {`)
	PackageRegexp = regexp.MustCompile(`^package[\s]*?(\w*);$`)
	ModuleRegexp  = regexp.MustCompile(`^module (.*)$`)

	GoNameRexp      = regexp.MustCompile(`[^a-zA-Z0-9\\s_-]+`)
	FirstDigitsRexp = regexp.MustCompile(`^\\d+`)
)
