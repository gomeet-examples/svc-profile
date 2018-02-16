package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestRead(t *testing.T) {
	server := newProfileServer()
	ctx := context.Background()

	req := &pb.ProfileRequest{}
	// You can generate a fake request see https://github.com/gomeet/go-proto-gomeetfaker
	// req := &pb.ProfileRequest{}
	res, err := server.Read(ctx, req)
	assert.Nil(t, err, "Read: error on call")
	assert.NotNil(t, res, "Read: error on call")

	// Do something useful tests with req and res
	// for example :
	// assert.Equal(t, req.GetUuid(), req.GetUuid(), "Read: Uuid field in response must be the same as that of the request")
}