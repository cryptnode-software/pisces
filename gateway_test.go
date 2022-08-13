package pisces

import (
	"context"
	"errors"
	"testing"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/utility"
	v1 "go.buf.build/grpc/go/thenewlebowski/pisces/general/v1"
)

var (
	gateway, err = lib.NewGateway(env, utility.Services(env))
	env          = utility.NewEnv(utility.NewLogger())
)

func TestGetSignedURL(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	req := &v1.GetSignedURLRequest{
		FileName: "test.png",
	}

	res, err := gateway.GetSignedURL(ctx, req)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Url == "" {
		t.Error(errors.New("no presigned url returned"))
	}
}
