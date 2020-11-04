package generate

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("file or dir already exists")
	ErrBadImport     = errors.New("bad import")
	ErrReqFailed     = errors.New("request failed")
	ErrProtoc        = errors.New("protoc")
)
