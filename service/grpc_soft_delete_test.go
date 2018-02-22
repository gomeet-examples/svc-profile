package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestSoftDelete(t *testing.T) {
	var (
		req *pb.ProfileRequest
		res *pb.ProfileResponse
		err error
	)

	// FIXME: do something better with mock?
	server, db, err := newProfileServerTest(t)
	assert.NotNil(t, server, "SoftDelete: newProfileServerTest nil server")
	assert.NotNil(t, db, "SoftDelete: newProfileServerTest nil db")
	assert.Nil(t, err, "SoftDelete: newProfileServerTest not nil err")

	ctx := context.Background()

	// first create a valid profile
	reqCreation := pb.NewProfileCreationRequestGomeetFaker()
	resCreation, err := server.Create(ctx, reqCreation)
	assert.Nil(t, err, "SoftDelete: error on call")
	assert.NotNil(t, resCreation, "SoftDelete: error on call")
	assert.True(t, resCreation.Ok, "SoftDelete: response Ok")
	assert.NotEmpty(t, resCreation.GetInfo().GetUuid(), "SoftDelete: uuid is empty")
	assert.Equal(t, reqCreation.GetGender(), resCreation.GetInfo().GetGender(), "SoftDelete: gender doesn't match")
	assert.Equal(t, reqCreation.GetEmail(), resCreation.GetInfo().GetEmail(), "SoftDelete: email doesn't match")
	assert.Equal(t, reqCreation.GetName(), resCreation.GetInfo().GetName(), "SoftDelete: name doesn't match")
	assert.Equal(t, reqCreation.GetBirthday(), resCreation.GetInfo().GetBirthday(), "SoftDelete: birthday doesn't match")
	assert.NotEmpty(t, resCreation.GetInfo().GetCreatedAt(), "SoftDelete: created_at is empty")
	assert.NotEmpty(t, resCreation.GetInfo().GetUpdatedAt(), "SoftDelete: updated_at is empty")
	assert.Empty(t, resCreation.GetInfo().GetDeletedAt(), "SoftDelete: deleted_at is not empty")

	validProfile := resCreation.GetInfo()

	// empty Uuid
	req = &pb.ProfileRequest{
		Uuid: "",
	}
	res, err = server.SoftDelete(ctx, req)
	assert.NotNil(t, err, "SoftDelete: expected error on call")
	assert.NotNil(t, res, "SoftDelete: expected error on call")
	assert.False(t, res.GetOk(), "SoftDelete: response Ok")
	assert.True(t, res.Equal(pb.ProfileResponse{Ok: false, Info: &pb.ProfileInfo{}}), "SoftDelete: no empty response")

	// invalid Uuid
	req = &pb.ProfileRequest{
		Uuid: "invalid uuid",
	}
	res, err = server.SoftDelete(ctx, req)
	assert.NotNil(t, err, "SoftDelete: expected error on call")
	assert.NotNil(t, res, "SoftDelete: expected error on call")
	assert.False(t, res.GetOk(), "SoftDelete: response Ok")
	assert.True(t, res.Equal(pb.ProfileResponse{Ok: false, Info: &pb.ProfileInfo{Uuid: req.GetUuid()}}), "SoftDelete: no empty response")

	// not found Uuid
	req = &pb.ProfileRequest{
		Uuid: uuid.New().String(),
	}
	res, err = server.SoftDelete(ctx, req)
	assert.NotNil(t, err, "SoftDelete: expected error on call")
	assert.NotNil(t, res, "SoftDelete: expected error on call")
	assert.False(t, res.GetOk(), "SoftDelete: response Ok")
	assert.True(t, res.Equal(pb.ProfileResponse{Ok: false, Info: &pb.ProfileInfo{Uuid: req.GetUuid()}}), "SoftDelete: no empty response")

	// valid uuid
	req = &pb.ProfileRequest{
		Uuid: validProfile.GetUuid(),
	}
	res, err = server.SoftDelete(ctx, req)
	assert.Nil(t, err, "SoftDelete: error on call")
	assert.NotNil(t, res, "SoftDelete: error on call")
	assert.True(t, res.GetOk(), "SoftDelete: response Ok")
	assert.Equal(t, validProfile.GetUuid(), res.GetInfo().GetUuid(), "SoftDelete: uuid doesn't match")
	assert.Equal(t, validProfile.GetGender(), res.GetInfo().GetGender(), "SoftDelete: gender doesn't match")
	assert.Equal(t, validProfile.GetEmail(), res.GetInfo().GetEmail(), "SoftDelete: email doesn't match")
	assert.Equal(t, validProfile.GetName(), res.GetInfo().GetName(), "SoftDelete: name doesn't match")
	assert.Equal(t, validProfile.GetBirthday(), res.GetInfo().GetBirthday(), "SoftDelete: birthday doesn't match")
	assert.Equal(t, validProfile.GetCreatedAt(), res.GetInfo().GetCreatedAt(), "SoftDelete: create_at doesn't match")
	assert.Equal(t, validProfile.GetUpdatedAt(), res.GetInfo().GetUpdatedAt(), "SoftDelete: update_at doesn't match")
	assert.NotEmpty(t, res.GetInfo().GetDeletedAt(), "SoftDelete: deleted_at is empty")
	assert.NotEqual(t, validProfile.GetDeletedAt(), res.GetInfo().GetDeletedAt(), "SoftDelete: deleted_at match")

	// second round - return a not found
	req = &pb.ProfileRequest{
		Uuid: validProfile.GetUuid(),
	}
	res, err = server.SoftDelete(ctx, req)
	assert.NotNil(t, err, "SoftDelete: expected error on call")
	assert.NotNil(t, res, "SoftDelete: expected error on call")
	assert.False(t, res.GetOk(), "SoftDelete: response Ok")
	assert.True(t, res.Equal(pb.ProfileResponse{Ok: false, Info: &pb.ProfileInfo{Uuid: req.GetUuid()}}), "SoftDelete: no empty response")
}
