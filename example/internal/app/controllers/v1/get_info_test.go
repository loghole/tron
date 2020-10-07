// Code generated by protoc-gen-goclay, but your can (must) modify it.
// source: strings.proto

package v1

import (
	"context"
	"testing"

	desc "github.com/loghole/example/pkg/v1"
	"github.com/stretchr/testify/require"
)

func TestImplementation_GetInfo(t *testing.T) {
	api := NewStrings()
	_, err := api.GetInfo(context.Background(), &desc.String{})

	require.NotNil(t, err)
	require.Equal(t, "GetInfo not implemented", err.Error())
}
