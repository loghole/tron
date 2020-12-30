package download

import (
	"errors"
)

var (
	ErrBadImport       = errors.New("bad import")
	ErrReqFailed       = errors.New("request failed")
	ErrVersionNotFound = errors.New("not found")
)
