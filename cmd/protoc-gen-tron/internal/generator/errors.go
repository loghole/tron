package generator

import (
	"errors"
)

var (
	ErrMultiplyService  = errors.New("files with multiply services aren't supported")
	ErrInvalidPackage   = errors.New("invalid proto package")
	ErrInvalidGoPackage = errors.New("invalid proto go_package")
)
