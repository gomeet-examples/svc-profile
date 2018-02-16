package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestList(t *testing.T) {
	server := newProfileServer()
	ctx := context.Background()

	req := &pb.ProfileListRequest{}
	// You can generate a fake request see https://github.com/gomeet/go-proto-gomeetfaker
	// req := &pb.ProfileListRequest{}
	res, err := server.List(ctx, req)
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")

	// Do something useful tests with req and res
	// for example :
	// assert.Equal(t, req.GetUuid(), req.GetUuid(), "List: Uuid field in response must be the same as that of the request")
}
