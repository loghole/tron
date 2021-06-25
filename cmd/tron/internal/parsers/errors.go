package parsers

import (
	"errors"
)

var (
	ErrPackageNotPound     = errors.New("package not found")
	ErrModuleNotFound      = errors.New("project module does not exists")
	ErrInvalidProtoPkgName = errors.New("invalid proto package name")
)
