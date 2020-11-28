// Package errors represents base tron error struct and parsing method.
package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	validation "github.com/gadavy/ozzo-validation/v4"
	"github.com/lissteron/simplerr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error represents base error struct.
type Error struct {
	rootErr    error
	httpStatus int
	grpcStatus codes.Code

	Message string         `json:"message"`
	Details []ErrorDetails `json:"details"`
}

// ParseError parse error to Error struct.
// Support lissteron/simplerr and ozzo-validation packages.
func ParseError(err error) *Error {
	resp := &Error{
		rootErr:    err,
		httpStatus: http.StatusInternalServerError,
		grpcStatus: codes.Unknown,
		Message:    simplerr.GetText(simplerr.GetWithCode(err)),
	}

	resp.parseErr(err)

	return resp
}

// GRPCStatus build status.Status from current error.
func (e *Error) GRPCStatus() *status.Status {
	st := status.New(e.grpcStatus, e.rootErr.Error())

	for _, val := range e.Details {
		stNew, err := st.WithDetails(&errdetails.DebugInfo{
			Detail: fmt.Sprintf("code: %s, detail: %s", val.Code, val.Detail),
		})
		if err != nil {
			return st
		}

		st = stNew
	}

	return st
}

// HTTPStatus return http status for error.
func (e *Error) HTTPStatus() int {
	return e.httpStatus
}

// Error returns error string.
func (e *Error) Error() string {
	if err := e.GRPCStatus().Err(); err != nil {
		return err.Error()
	}

	return ""
}

// ErrorDetails represents error details.
type ErrorDetails struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (e *Error) parseErr(err error) {
	var validationErrs validation.Errors

	if errors.As(err, &validationErrs) {
		if e.parseValidationErr(validationErrs) {
			return
		}
	}

	if errC := simplerr.GetWithCode(err); errC != nil {
		err = errC
	}

	code := simplerr.GetCode(err)
	e.grpcStatus = codes.Code(code.GRPC())

	if httpCode := code.HTTP(); httpCode > 0 {
		e.httpStatus = httpCode
	}

	e.Details = append(e.Details, ErrorDetails{
		Code:   strconv.Itoa(code.Int()),
		Detail: simplerr.GetText(err),
	})
}

func (e *Error) parseValidationErr(list validation.Errors) bool {
	for field, err := range list {
		var validationErr validation.Error

		if errors.As(err, &validationErr) {
			e.Details = append(e.Details, ErrorDetails{
				Code:   validationErr.Code(),
				Detail: strings.Join([]string{field, validationErr.Error()}, ": "),
			})

			continue
		}

		var validationErrs validation.Errors

		if !errors.As(err, &validationErrs) {
			return false
		}

		if !e.parseValidationErr(validationErrs) {
			return false
		}
	}

	e.httpStatus = http.StatusBadRequest
	e.grpcStatus = codes.InvalidArgument

	return true
}
