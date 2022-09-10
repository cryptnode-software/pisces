package pisces

import (
	"context"
	"errors"
	"testing"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/services"
	v1 "go.buf.build/grpc/go/thenewlebowski/pisces/general/v1"
)

var (
	gateway, err = lib.NewGateway(env, services.New(env))
	env          = lib.NewEnv(lib.NewLogger(lib.EnvDev))
)

func TestGetSignedURL(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	req := &v1.StartUploadRequest{
		Key: "test.png",
	}

	res, err := gateway.StartUpload(ctx, req)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Url == "" {
		t.Error(errors.New("no presigned url returned"))
	}
}
