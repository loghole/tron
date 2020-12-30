package parsers

import (
	"errors"
)

var (
	ErrPackageNotPound     = errors.New("package not found")
	ErrModuleNotFound      = errors.New("project module does not exists")
	ErrMultipleServices    = errors.New("multiple service entries")
	ErrInvalidProtoPkgName = errors.New("invalid proto package name, use '.v{{ integer }}' at the end of pkg name")
)
