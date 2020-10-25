// Code generated by protoc-gen-goclay. DO NOT EDIT.
// source: strings.proto

/*
Package stringsV1 is a self-registering gRPC and JSON+Swagger service definition.

It conforms to the github.com/utrack/clay/v2/transport Service interface.
*/
package stringsV1

import (
	"bytes"
	"context"
	"encoding/base64"
	typesV1 "example/pkg/types/v1"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-openapi/spec"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"github.com/pkg/errors"
	"github.com/utrack/clay/v2/transport"
	"github.com/utrack/clay/v2/transport/httpclient"
	"github.com/utrack/clay/v2/transport/httpruntime"
	"github.com/utrack/clay/v2/transport/httpruntime/httpmw"
	"github.com/utrack/clay/v2/transport/httptransport"
	"github.com/utrack/clay/v2/transport/swagger"
	"google.golang.org/grpc"
)

// Update your shared lib or downgrade generator to v1 if there's an error
var _ = transport.IsVersion2

var _ = ioutil.Discard
var _ chi.Router
var _ runtime.Marshaler
var _ bytes.Buffer
var _ context.Context
var _ fmt.Formatter
var _ strings.Reader
var _ errors.Frame
var _ httpruntime.Marshaler
var _ http.Handler
var _ url.Values
var _ base64.Encoding
var _ httptransport.MarshalerError
var _ utilities.DoubleArray

// StringsDesc is a descriptor/registrator for the StringsServer.
type StringsDesc struct {
	svc  StringsServer
	opts httptransport.DescOptions
}

// NewStringsServiceDesc creates new registrator for the StringsServer.
// It implements httptransport.ConfigurableServiceDesc as well.
func NewStringsServiceDesc(svc StringsServer) *StringsDesc {
	return &StringsDesc{
		svc: svc,
	}
}

// RegisterGRPC implements service registrator interface.
func (d *StringsDesc) RegisterGRPC(s *grpc.Server) {
	RegisterStringsServer(s, d.svc)
}

// Apply applies passed options.
func (d *StringsDesc) Apply(oo ...transport.DescOption) {
	for _, o := range oo {
		o.Apply(&d.opts)
	}
}

// SwaggerDef returns this file's Swagger definition.
func (d *StringsDesc) SwaggerDef(options ...swagger.Option) (result []byte) {
	if len(options) > 0 || len(d.opts.SwaggerDefaultOpts) > 0 {
		var err error
		var s = &spec.Swagger{}
		if err = s.UnmarshalJSON(_swaggerDef_strings_proto); err != nil {
			panic("Bad swagger definition: " + err.Error())
		}

		for _, o := range d.opts.SwaggerDefaultOpts {
			o(s)
		}
		for _, o := range options {
			o(s)
		}
		if result, err = s.MarshalJSON(); err != nil {
			panic("Failed marshal spec.Swagger definition: " + err.Error())
		}
	} else {
		result = _swaggerDef_strings_proto
	}
	return result
}

// RegisterHTTP registers this service's HTTP handlers/bindings.
func (d *StringsDesc) RegisterHTTP(mux transport.Router) {
	chiMux, isChi := mux.(chi.Router)

	{
		// Handler for ToUpper, binding: GET /api/v1/strings/upper/{str}
		var h http.HandlerFunc
		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			unmFunc := unmarshaler_goclay_Strings_ToUpper_0(r)
			rsp, err := _Strings_ToUpper_Handler(d.svc, r.Context(), unmFunc, d.opts.UnaryInterceptor)

			if err != nil {
				if err, ok := err.(httptransport.MarshalerError); ok {
					httpruntime.SetError(r.Context(), r, w, errors.Wrap(err.Err, "couldn't parse request"))
					return
				}
				httpruntime.SetError(r.Context(), r, w, err)
				return
			}

			if ctxErr := r.Context().Err(); ctxErr != nil && ctxErr == context.Canceled {
				w.WriteHeader(499) // Client Closed Request
				return
			}

			_, outbound := httpruntime.MarshalerForRequest(r)
			w.Header().Set("Content-Type", outbound.ContentType())
			err = outbound.Marshal(w, rsp)
			if err != nil {
				httpruntime.SetError(r.Context(), r, w, errors.Wrap(err, "couldn't write response"))
				return
			}
		})

		h = httpmw.DefaultChain(h)

		if isChi {
			chiMux.Method("GET", pattern_goclay_Strings_ToUpper_0, h)
		} else {
			panic("query URI params supported only for chi.Router")
		}
	}

	{
		// Handler for GetInfo, binding: GET /api/v1/strings/info
		var h http.HandlerFunc
		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			unmFunc := unmarshaler_goclay_Strings_GetInfo_0(r)
			rsp, err := _Strings_GetInfo_Handler(d.svc, r.Context(), unmFunc, d.opts.UnaryInterceptor)

			if err != nil {
				if err, ok := err.(httptransport.MarshalerError); ok {
					httpruntime.SetError(r.Context(), r, w, errors.Wrap(err.Err, "couldn't parse request"))
					return
				}
				httpruntime.SetError(r.Context(), r, w, err)
				return
			}

			if ctxErr := r.Context().Err(); ctxErr != nil && ctxErr == context.Canceled {
				w.WriteHeader(499) // Client Closed Request
				return
			}

			_, outbound := httpruntime.MarshalerForRequest(r)
			w.Header().Set("Content-Type", outbound.ContentType())
			err = outbound.Marshal(w, rsp)
			if err != nil {
				httpruntime.SetError(r.Context(), r, w, errors.Wrap(err, "couldn't write response"))
				return
			}
		})

		h = httpmw.DefaultChain(h)

		if isChi {
			chiMux.Method("GET", pattern_goclay_Strings_GetInfo_0, h)
		} else {
			mux.Handle(pattern_goclay_Strings_GetInfo_0, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				h(w, r)
			}))
		}
	}

}

type Strings_httpClient struct {
	c    *http.Client
	host string
}

// NewStringsHTTPClient creates new HTTP client for StringsServer.
// Pass addr in format "http://host[:port]".
func NewStringsHTTPClient(c *http.Client, addr string) *Strings_httpClient {
	if strings.HasSuffix(addr, "/") {
		addr = addr[:len(addr)-1]
	}
	return &Strings_httpClient{c: c, host: addr}
}

func (c *Strings_httpClient) ToUpper(ctx context.Context, in *typesV1.String, opts ...grpc.CallOption) (*typesV1.String, error) {
	mw, err := httpclient.NewMiddlewareGRPC(opts)
	if err != nil {
		return nil, err
	}

	path := pattern_goclay_Strings_ToUpper_0_builder(in)

	buf := bytes.NewBuffer(nil)

	m := httpruntime.DefaultMarshaler(nil)

	req, err := http.NewRequest("GET", c.host+path, buf)
	if err != nil {
		return nil, errors.Wrap(err, "can't initiate HTTP request")
	}
	req = req.WithContext(ctx)

	req.Header.Add("Accept", m.ContentType())

	req, err = mw.ProcessRequest(req)
	if err != nil {
		return nil, err
	}
	rsp, err := c.c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error from client")
	}
	defer rsp.Body.Close()

	rsp, err = mw.ProcessResponse(rsp)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(rsp.Body)
		return nil, errors.Errorf("%v %v: server returned HTTP %v: '%v'", req.Method, req.URL.String(), rsp.StatusCode, string(b))
	}

	ret := typesV1.String{}

	err = m.Unmarshal(rsp.Body, &ret)

	return &ret, errors.Wrap(err, "can't unmarshal response")
}

func (c *Strings_httpClient) GetInfo(ctx context.Context, in *typesV1.String, opts ...grpc.CallOption) (*typesV1.String, error) {
	mw, err := httpclient.NewMiddlewareGRPC(opts)
	if err != nil {
		return nil, err
	}

	path := pattern_goclay_Strings_GetInfo_0_builder(in)

	buf := bytes.NewBuffer(nil)

	m := httpruntime.DefaultMarshaler(nil)

	req, err := http.NewRequest("GET", c.host+path, buf)
	if err != nil {
		return nil, errors.Wrap(err, "can't initiate HTTP request")
	}
	req = req.WithContext(ctx)

	req.Header.Add("Accept", m.ContentType())

	req, err = mw.ProcessRequest(req)
	if err != nil {
		return nil, err
	}
	rsp, err := c.c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error from client")
	}
	defer rsp.Body.Close()

	rsp, err = mw.ProcessResponse(rsp)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(rsp.Body)
		return nil, errors.Errorf("%v %v: server returned HTTP %v: '%v'", req.Method, req.URL.String(), rsp.StatusCode, string(b))
	}

	ret := typesV1.String{}

	err = m.Unmarshal(rsp.Body, &ret)

	return &ret, errors.Wrap(err, "can't unmarshal response")
}

// patterns for Strings
var (
	pattern_goclay_Strings_ToUpper_0 = "/api/v1/strings/upper/{str}"

	pattern_goclay_Strings_ToUpper_0_builder = func(in *typesV1.String) string {
		values := url.Values{}

		u := url.URL{
			Path:     fmt.Sprintf("/api/v1/strings/upper/%v", in.Str),
			RawQuery: values.Encode(),
		}
		return u.String()
	}

	unmarshaler_goclay_Strings_ToUpper_0_boundParams = &utilities.DoubleArray{Encoding: map[string]int{"str": 0}, Base: []int{1, 1, 0}, Check: []int{0, 1, 2}}

	pattern_goclay_Strings_GetInfo_0 = "/api/v1/strings/info"

	pattern_goclay_Strings_GetInfo_0_builder = func(in *typesV1.String) string {
		values := url.Values{}
		values.Add("str", fmt.Sprintf("%s", in.Str))

		u := url.URL{
			Path:     fmt.Sprintf("/api/v1/strings/info"),
			RawQuery: values.Encode(),
		}
		return u.String()
	}

	unmarshaler_goclay_Strings_GetInfo_0_boundParams = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

// marshalers for Strings
var (
	unmarshaler_goclay_Strings_ToUpper_0 = func(r *http.Request) func(interface{}) error {
		return func(rif interface{}) error {
			req := rif.(*typesV1.String)

			if err := errors.Wrap(runtime.PopulateQueryParameters(req, r.URL.Query(), unmarshaler_goclay_Strings_ToUpper_0_boundParams), "couldn't populate query parameters"); err != nil {
				return httpruntime.TransformUnmarshalerError(err)
			}

			rctx := chi.RouteContext(r.Context())
			if rctx == nil {
				panic("Only chi router is supported for GETs atm")
			}
			for pos, k := range rctx.URLParams.Keys {
				if err := errors.Wrapf(runtime.PopulateFieldFromPath(req, k, rctx.URLParams.Values[pos]), "can't read '%v' from path", k); err != nil {
					return httptransport.NewMarshalerError(httpruntime.TransformUnmarshalerError(err))
				}
			}

			return nil
		}
	}

	unmarshaler_goclay_Strings_GetInfo_0 = func(r *http.Request) func(interface{}) error {
		return func(rif interface{}) error {
			req := rif.(*typesV1.String)

			if err := errors.Wrap(runtime.PopulateQueryParameters(req, r.URL.Query(), unmarshaler_goclay_Strings_GetInfo_0_boundParams), "couldn't populate query parameters"); err != nil {
				return httpruntime.TransformUnmarshalerError(err)
			}

			return nil
		}
	}
)

var _swaggerDef_strings_proto = []byte(`{
  "swagger": "2.0",
  "info": {
    "title": "strings.proto",
    "version": "version not set"
  },
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/api/v1/strings/info": {
      "get": {
        "operationId": "GetInfo",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/v1String"
            }
          }
        },
        "parameters": [
          {
            "name": "str",
            "in": "query",
            "required": false,
            "type": "string"
          }
        ],
        "tags": [
          "Strings"
        ]
      }
    },
    "/api/v1/strings/upper/{str}": {
      "get": {
        "summary": "Method to upper",
        "operationId": "ToUpper",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/v1String"
            }
          }
        },
        "parameters": [
          {
            "name": "str",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "Strings"
        ]
      }
    }
  },
  "definitions": {
    "v1String": {
      "type": "object",
      "properties": {
        "str": {
          "type": "string"
        }
      }
    }
  }
}

`)