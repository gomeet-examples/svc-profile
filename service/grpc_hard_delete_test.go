package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestHardDelete(t *testing.T) {
	server := newProfileServer()
	ctx := context.Background()

	req := &pb.ProfileRequest{}
	// You can generate a fake request see https://github.com/gomeet/go-proto-gomeetfaker
	// req := &pb.ProfileRequest{}
	res, err := server.HardDelete(ctx, req)
	assert.Nil(t, err, "HardDelete: error on call")
	assert.NotNil(t, res, "HardDelete: error on call")

	// Do something useful tests with req and res
	// for example :
	// assert.Equal(t, req.GetUuid(), req.GetUuid(), "HardDelete: Uuid field in response must be the same as that of the request")
}
