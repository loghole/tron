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
	Status     int        `json:"-"`
	GRPCStatus codes.Code `json:"-"`
	Errors     []Error    `json:"errors"`
}

func ParseError(err error) *ErrResponse {
	resp := &ErrResponse{Status: http.StatusInternalServerError}
	resp.parseErr(err)

	return resp
}

func (r *ErrResponse) ToStatus() *status.Status {
	st := status.New(r.GRPCStatus, r.GRPCStatus.String())

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

type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (r *ErrResponse) parseErr(err error) {
	code := simplerr.GetCode(err)
	r.GRPCStatus = codes.Code(code.GRPC())

	if code.HTTP() == http.StatusBadRequest {
		validationErr, ok := errors.Unwrap(err).(validation.Errors)
		if ok && r.parseValidationErr(validationErr) {
			return
		}
	}

	if httpCode := code.HTTP(); httpCode > 0 {
		r.Status = httpCode
	}

	r.Errors = append(r.Errors, Error{
		Code:   strconv.Itoa(code.Int()),
		Detail: simplerr.GetText(err),
	})
}

func (r *ErrResponse) parseValidationErr(err error) bool {
	e1, ok := err.(validation.Errors)
	if !ok {
		return false
	}

	for field, e2 := range e1 {
		if e3, ok := e2.(validation.Error); ok {
			r.Errors = append(r.Errors, Error{
				Code:   e3.Code(),
				Detail: strings.Join([]string{field, e3.Error()}, ": "),
			})

			continue
		}

		if !r.parseValidationErr(e2) {
			return false
		}
	}

	r.Status = http.StatusBadRequest

	return true
}
