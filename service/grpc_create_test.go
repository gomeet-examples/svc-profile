package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestCreate(t *testing.T) {
	server := newProfileServer()
	ctx := context.Background()

	req := &pb.ProfileCreationRequest{}
	// You can generate a fake request see https://github.com/gomeet/go-proto-gomeetfaker
	// req := &pb.ProfileCreationRequest{}
	res, err := server.Create(ctx, req)
	assert.Nil(t, err, "Create: error on call")
	assert.NotNil(t, res, "Create: error on call")

	// Do something useful tests with req and res
	// for example :
	// assert.Equal(t, req.GetUuid(), req.GetUuid(), "Create: Uuid field in response must be the same as that of the request")
}
