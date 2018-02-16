package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestUpdate(t *testing.T) {
	server := newProfileServer()
	ctx := context.Background()

	req := &pb.ProfileInfo{}
	// You can generate a fake request see https://github.com/gomeet/go-proto-gomeetfaker
	// req := &pb.ProfileInfo{}
	res, err := server.Update(ctx, req)
	assert.Nil(t, err, "Update: error on call")
	assert.NotNil(t, res, "Update: error on call")

	// Do something useful tests with req and res
	// for example :
	// assert.Equal(t, req.GetUuid(), req.GetUuid(), "Update: Uuid field in response must be the same as that of the request")
}
