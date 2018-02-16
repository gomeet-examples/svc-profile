package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestEcho(t *testing.T) {
	server := newProfileServer()
	ctx := context.Background()

	req := &pb.EchoRequest{}
	// You can generate a fake request see https://github.com/gomeet/go-proto-gomeetfaker
	// req := &pb.EchoRequest{}
	res, err := server.Echo(ctx, req)
	assert.Nil(t, err, "Echo: error on call")
	assert.NotNil(t, res, "Echo: error on call")

	// Do something useful tests with req and res
	// for example :
	// assert.Equal(t, req.GetUuid(), req.GetUuid(), "Echo: Uuid field in response must be the same as that of the request")
}
