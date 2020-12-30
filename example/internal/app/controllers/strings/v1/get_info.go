package stringsV1

import (
	"context"

	"example/config"
	typesV1 "example/pkg/types/v1"
)

func (i *Implementation) GetInfo(
	ctx context.Context,
	req *typesV1.String,
) (resp *typesV1.String, err error) {
	resp = &typesV1.String{
		Str: config.GetExampleValue(),
	}

	return resp, nil
}
