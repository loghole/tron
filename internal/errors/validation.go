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

type ErrResponse struct {
	httpStatus int
	grpcStatus codes.Code
	Errors     []Error `json:"errors"`
}

func ParseError(err error) *ErrResponse {
	resp := &ErrResponse{httpStatus: http.StatusInternalServerError}
	resp.parseErr(err)

	return resp
}

func (r *ErrResponse) GRPCStatus() *status.Status {
	st := status.New(r.grpcStatus, r.grpcStatus.String())

	for _, val := range r.Errors {
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

func (r *ErrResponse) HTTPStatus() int {
	return r.httpStatus
}

func (r *ErrResponse) Error() string {
	return r.GRPCStatus().Err().Error()
}

type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (r *ErrResponse) parseErr(err error) {
	var validationErrs validation.Errors

	if errors.As(err, &validationErrs) {
		if r.parseValidationErr(validationErrs) {
			return
		}
	}

	if errC := simplerr.GetWithCode(err); errC != nil {
		err = errC
	}

	code := simplerr.GetCode(err)
	r.grpcStatus = codes.Code(code.GRPC())

	if httpCode := code.HTTP(); httpCode > 0 {
		r.httpStatus = httpCode
	}

	r.Errors = append(r.Errors, Error{
		Code:   strconv.Itoa(code.Int()),
		Detail: simplerr.GetText(err),
	})
}

func (r *ErrResponse) parseValidationErr(list validation.Errors) bool {
	for field, err := range list {
		var validationErr validation.Error

		if errors.As(err, &validationErr) {
			r.Errors = append(r.Errors, Error{
				Code:   validationErr.Code(),
				Detail: strings.Join([]string{field, validationErr.Error()}, ": "),
			})

			continue
		}

		var validationErrs validation.Errors

		if !errors.As(err, &validationErrs) {
			return false
		}

		if !r.parseValidationErr(validationErrs) {
			return false
		}
	}

	r.httpStatus = http.StatusBadRequest

	return true
}
