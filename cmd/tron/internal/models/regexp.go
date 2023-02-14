package models

import (
	"errors"
	"regexp"
)

var ErrNotSemanticVersion = errors.New("string is not semantic version")

var (
	Version3Regexp = regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	Version2Regexp = regexp.MustCompile(`(\d+\.\d+)`)
	PackageRegexp  = regexp.MustCompile(`^package ([^;]+);$`)
	ImportRegexp   = regexp.MustCompile(`^import "([^"]+)";$`)
	ModuleRegexp   = regexp.MustCompile(`^module (.+)$`)
)

func ExtractVersion(s string) (string, error) {
	matches := Version3Regexp.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1], nil
	}

	matches = Version2Regexp.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1] + ".0", nil
	}

	return "", ErrNotSemanticVersion
}
