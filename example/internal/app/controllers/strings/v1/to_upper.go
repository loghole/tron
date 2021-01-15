package stringsV1

import (
	"context"
	"errors"
	"net/http"
	"strings"

	validation "github.com/gadavy/ozzo-validation/v4"
	"github.com/lissteron/simplerr"
	"google.golang.org/grpc/codes"

	typesV1 "example/pkg/types/v1"
)

func (i *Implementation) ToUpper(
	ctx context.Context,
	req *typesV1.String,
) (resp *typesV1.String, err error) {
	if err := validateString(req); err != nil {
		return nil, err
	}

	if err := simpleError(req); err != nil {
		return nil, err
	}

	return &typesV1.String{Str: strings.ToUpper(req.Str)}, nil
}

func validateString(req *typesV1.String) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Str, validation.NotIn("qwerty").ErrorCode("123")),
	)
}

func simpleError(req *typesV1.String) error {
	var someErr = errors.New("internal error")

	switch req.Str {
	case "error-1":
		return simplerr.WrapfWithCode(someErr, Code(1), "some text 1")
	case "error-2":
		return simplerr.WithCode(someErr, Code(2))
	case "error-3":
		return simplerr.Wrap(someErr, "some text 2")
	}

	return nil
}

type Code int

func (c Code) HTTP() int {
	return http.StatusBadGateway
}

func (c Code) GRPC() int {
	return int(codes.Unauthenticated)
}

func (c Code) Int() int {
	return int(c)
}
