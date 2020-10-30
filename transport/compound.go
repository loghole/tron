package transport

import (
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/peterbourgon/mergemap"
	"google.golang.org/grpc"
)

type CompoundServiceDesc struct {
	svc []ServiceDesc
}

func NewCompoundServiceDesc(desc ...ServiceDesc) *CompoundServiceDesc {
	return &CompoundServiceDesc{svc: desc}
}

func (d *CompoundServiceDesc) RegisterGRPC(g *grpc.Server) {
	for _, svc := range d.svc {
		svc.RegisterGRPC(g)
	}
}

func (d *CompoundServiceDesc) RegisterHTTP(mux *runtime.ServeMux) {
	for _, svc := range d.svc {
		svc.RegisterHTTP(mux)
	}
}

func (d *CompoundServiceDesc) SwaggerDef() []byte {
	j := &swagJoiner{}

	for _, svc := range d.svc {
		_ = j.AddDefinition(svc.SwaggerDef())
	}

	return j.SumDefinitions()
}

// swagJoiner glues up several Swagger definitions to one.
// This is one dirty hack...
type swagJoiner struct {
	result map[string]interface{}
}

// AddDefinition adds another definition to the soup.
func (c *swagJoiner) AddDefinition(buf []byte) error {
	def := map[string]interface{}{}

	err := json.Unmarshal(buf, &def)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal JSON def: %w", err)
	}

	if c.result == nil {
		c.result = def
		return nil
	}

	c.result = mergemap.Merge(c.result, def)

	return nil
}

// SumDefinitions returns a (hopefully) valid Swagger definition combined
// from everything that came up .AddDefinition().
func (c *swagJoiner) SumDefinitions() []byte {
	if c.result == nil {
		c.result = map[string]interface{}{}
	}
	ret, err := json.Marshal(c.result)
	if err != nil {
		panic(err)
	}
	return ret
}
