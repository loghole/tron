// Package errors represents base tron error struct and parsing method.
package errors

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	validation "github.com/gadavy/ozzo-validation/v4"
	"github.com/lissteron/simplerr"
	"github.com/loghole/tracing/tracelog"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error represents base http error struct.
type Error struct {
	rootErr    error
	httpStatus int

	Code    codes.Code `json:"code"`
	Message string     `json:"message,omitempty"`
	Details []*Details `json:"details,omitempty"`
	TraceID string     `json:"trace_id"`
}

// Details represents error details.
type Details struct {
	Code        string `json:"code,omitempty"`
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
}

// ParseError parse error to Error struct.
// Support lissteron/simplerr and ozzo-validation packages.
func ParseError(ctx context.Context, err error) *Error {
	resp := &Error{
		rootErr:    err,
		httpStatus: http.StatusInternalServerError,
		Code:       codes.Unknown,
		TraceID:    tracelog.TraceID(ctx),
	}

	if resp.parseValidationErr(err) {
		resp.httpStatus = http.StatusBadRequest
		resp.Code = codes.InvalidArgument
	} else {
		resp.parseDefaultErr(err)
	}

	resp.Message = http.StatusText(resp.httpStatus)

	return resp
}

// GRPCStatus build status.Status from current error.
func (e *Error) GRPCStatus() *status.Status {
	grpcStatus := status.New(e.Code, e.rootErr.Error())

	if value, err := grpcStatus.WithDetails(e.grpcErrorInfo()); err == nil {
		return value
	}

	return grpcStatus
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

func (e *Error) grpcErrorInfo() *errdetails.ErrorInfo {
	info := &errdetails.ErrorInfo{}
	info.Metadata["trace_id"] = e.TraceID

	for _, detail := range e.Details {
		info.Metadata[detail.Code] = detail.Description
	}

	return info
}

func (e *Error) parseValidationErr(err error) bool {
	if detail, ok := parseValidationErr(err, ""); ok {
		e.Details = append(e.Details, detail)

		return true
	}

	if details, ok := parseValidationErrList(err, ""); ok {
		e.Details = append(e.Details, details...)

		return true
	}

	return false
}

func (e *Error) parseDefaultErr(err error) {
	if errC := simplerr.GetWithCode(err); errC != nil {
		err = errC
	}

	code := simplerr.GetCode(err)

	if gCode := code.GRPC(); gCode > 0 {
		e.Code = codes.Code(gCode)
	}

	if hCode := code.HTTP(); hCode > 0 {
		e.httpStatus = hCode
	}

	details := &Details{
		Code:        strconv.Itoa(code.Int()),
		Description: simplerr.GetText(err), // can by empty.
	}

	if details.Description == "" {
		details.Description = err.Error()
	}

	e.Details = append(e.Details, details)
}

func parseValidationErrList(err error, field string) ([]*Details, bool) {
	var validationErrList validation.Errors

	if !errors.As(err, &validationErrList) {
		return nil, false
	}

	result := make([]*Details, 0, len(validationErrList))

	for secondField, err := range validationErrList {
		field := buildField(field, secondField, ".")

		if detail, ok := parseValidationErr(err, field); ok {
			result = append(result, detail)

			continue
		}

		if details, ok := parseValidationErrList(err, field); ok {
			result = append(result, details...)

			continue
		}

		result = append(result, defaultErr(err, field))
	}

	return result, true
}

func parseValidationErr(err error, field string) (*Details, bool) {
	var validationErr validation.Error

	if !errors.As(err, &validationErr) {
		return nil, false
	}

	result := &Details{
		Field:       field,
		Code:        validationErr.Code(),
		Description: validationErr.Error(),
	}

	return result, true
}

func defaultErr(err error, field string) *Details {
	return &Details{
		Field:       field,
		Description: err.Error(),
	}
}

func buildField(first, second, separator string) string {
	if first == "" {
		return second
	}

	return first + separator + second
}
