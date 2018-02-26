package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestRead(t *testing.T) {
	var (
		req *pb.ProfileRequest
		res *pb.ProfileInfo
		err error
	)

	// FIXME: do something better with mock?
	server, db, err := newProfileServerTest(t)
	assert.NotNil(t, server, "Read: newProfileServerTest nil server")
	assert.NotNil(t, db, "Read: newProfileServerTest nil db")
	assert.Nil(t, err, "Read: newProfileServerTest not nil err")

	ctx := context.Background()

	// first create a valid profile
	reqCreation := pb.NewProfileCreationRequestGomeetFaker()
	resCreation, err := server.Create(ctx, reqCreation)
	assert.Nil(t, err, "Read: error on call")
	assert.NotNil(t, resCreation, "Read: error on call")
	assert.True(t, resCreation.Ok, "Read: response Ok")
	assert.NotEmpty(t, resCreation.GetInfo().GetUuid(), "Read: uuid is empty")
	assert.Equal(t, reqCreation.GetGender(), resCreation.GetInfo().GetGender(), "Read: gender doesn't match")
	assert.Equal(t, reqCreation.GetEmail(), resCreation.GetInfo().GetEmail(), "Read: email doesn't match")
	assert.Equal(t, reqCreation.GetName(), resCreation.GetInfo().GetName(), "Read: name doesn't match")
	assert.Equal(t, reqCreation.GetBirthday(), resCreation.GetInfo().GetBirthday(), "Read: birthday doesn't match")
	assert.NotEmpty(t, resCreation.GetInfo().GetCreatedAt(), "Read: created_at is empty")
	assert.NotEmpty(t, resCreation.GetInfo().GetUpdatedAt(), "Read: updated_at is empty")
	assert.Empty(t, resCreation.GetInfo().GetDeletedAt(), "Read: deleted_at is not empty")

	validProfile := resCreation.GetInfo()

	// empty Uuid
	req = &pb.ProfileRequest{
		Uuid: "",
	}
	res, err = server.Read(ctx, req)
	assert.NotNil(t, err, "Read: expected error on call")
	assert.True(t, res.Equal(pb.ProfileInfo{}), "Read: no empty response")

	// invalid Uuid
	req = &pb.ProfileRequest{
		Uuid: "invalid uuid",
	}
	res, err = server.Read(ctx, req)
	assert.NotNil(t, err, "Read: expected error on call")
	assert.True(t, res.Equal(pb.ProfileInfo{Uuid: req.GetUuid()}), "Read: no empty response")

	// valid uuid
	req = &pb.ProfileRequest{
		Uuid: validProfile.GetUuid(),
	}
	res, err = server.Read(ctx, req)
	assert.Nil(t, err, "Read: error on call")
	assert.NotNil(t, res, "Read: error on call")
	assert.Equal(t, validProfile.GetUuid(), res.GetUuid(), "Read: uuid doesn't match")
	assert.Equal(t, validProfile.GetGender(), res.GetGender(), "Read: gender doesn't match")
	assert.Equal(t, validProfile.GetEmail(), res.GetEmail(), "Read: email doesn't match")
	assert.Equal(t, validProfile.GetName(), res.GetName(), "Read: name doesn't match")
	assert.Equal(t, validProfile.GetBirthday(), res.GetBirthday(), "Read: birthday doesn't match")
	assert.Equal(t, validProfile.GetCreatedAt(), res.GetCreatedAt(), "Read: create_at doesn't match")
	assert.Equal(t, validProfile.GetUpdatedAt(), res.GetUpdatedAt(), "Read: update_at doesn't match")
	assert.Equal(t, validProfile.GetDeletedAt(), res.GetDeletedAt(), "Read: deleted_at doesn't match")
}
