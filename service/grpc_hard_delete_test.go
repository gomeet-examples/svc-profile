package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestHardDelete(t *testing.T) {
	var (
		req *pb.ProfileRequest
		res *pb.ProfileResponseLight
		err error
	)

	// FIXME: do something better with mock?
	server, db, err := newProfileServerTest(t)
	assert.NotNil(t, server, "HardDelete: newProfileServerTest nil server")
	assert.NotNil(t, db, "HardDelete: newProfileServerTest nil db")
	assert.Nil(t, err, "HardDelete: newProfileServerTest not nil err")

	ctx := context.Background()

	// first create a valid profile
	reqCreation := pb.NewProfileCreationRequestGomeetFaker()
	resCreation, err := server.Create(ctx, reqCreation)
	assert.Nil(t, err, "HardDelete: error on call")
	assert.NotNil(t, resCreation, "HardDelete: error on call")
	assert.True(t, resCreation.Ok, "HardDelete: response Ok")
	assert.NotEmpty(t, resCreation.GetInfo().GetUuid(), "HardDelete: uuid is empty")
	assert.Equal(t, reqCreation.GetGender(), resCreation.GetInfo().GetGender(), "HardDelete: gender doesn't match")
	assert.Equal(t, reqCreation.GetEmail(), resCreation.GetInfo().GetEmail(), "HardDelete: email doesn't match")
	assert.Equal(t, reqCreation.GetName(), resCreation.GetInfo().GetName(), "HardDelete: name doesn't match")
	assert.Equal(t, reqCreation.GetBirthday(), resCreation.GetInfo().GetBirthday(), "HardDelete: birthday doesn't match")
	assert.NotEmpty(t, resCreation.GetInfo().GetCreatedAt(), "HardDelete: created_at is empty")
	assert.NotEmpty(t, resCreation.GetInfo().GetUpdatedAt(), "HardDelete: updated_at is empty")
	assert.Empty(t, resCreation.GetInfo().GetDeletedAt(), "HardDelete: deleted_at is not empty")

	validProfile := resCreation.GetInfo()

	// empty Uuid
	req = &pb.ProfileRequest{
		Uuid: "",
	}
	res, err = server.HardDelete(ctx, req)
	assert.NotNil(t, err, "HardDelete: expected error on call")
	assert.NotNil(t, res, "HardDelete: expected error on call")
	assert.False(t, res.GetOk(), "HardDelete: response Ok")

	// invalid Uuid
	req = &pb.ProfileRequest{
		Uuid: "invalid uuid",
	}
	res, err = server.HardDelete(ctx, req)
	assert.NotNil(t, err, "HardDelete: expected error on call")
	assert.NotNil(t, res, "HardDelete: expected error on call")
	assert.False(t, res.GetOk(), "HardDelete: response Ok")

	// not found Uuid - not an error
	req = &pb.ProfileRequest{
		Uuid: uuid.New().String(),
	}
	res, err = server.HardDelete(ctx, req)
	assert.Nil(t, err, "HardDelete: error on call")
	assert.NotNil(t, res, "HardDelete: error on call")
	assert.True(t, res.GetOk(), "HardDelete: response Ok")

	// valid uuid
	req = &pb.ProfileRequest{
		Uuid: validProfile.GetUuid(),
	}
	res, err = server.HardDelete(ctx, req)
	assert.Nil(t, err, "HardDelete: error on call")
	assert.NotNil(t, res, "HardDelete: error on call")
	assert.True(t, res.GetOk(), "HardDelete: response Ok")

	// second round - not an error
	req = &pb.ProfileRequest{
		Uuid: validProfile.GetUuid(),
	}
	res, err = server.HardDelete(ctx, req)
	assert.Nil(t, err, "HardDelete: error on call")
	assert.NotNil(t, res, "HardDelete: error on call")
	assert.True(t, res.GetOk(), "HardDelete: response Ok")
}
