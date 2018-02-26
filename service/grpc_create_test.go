package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestCreate(t *testing.T) {
	var (
		req *pb.ProfileCreationRequest
		res *pb.ProfileResponse
		err error
	)

	// FIXME: do something better
	server, db, err := newProfileServerTest(t)
	assert.NotNil(t, server, "Create: newProfileServerTest nil server")
	assert.NotNil(t, db, "Create: newProfileServerTest nil db")
	assert.Nil(t, err, "Create: newProfileServerTest not nil err")

	ctx := context.Background()

	// invalid gender
	req = &pb.ProfileCreationRequest{
		Gender:   pb.Genders_UNKNOW,
		Email:    "test@example.com",
		Name:     "Profile Test Name",
		Birthday: "1989-11-09",
	}
	res, err = server.Create(ctx, req)
	assert.NotNil(t, err, "Create: expected error on call")
	assert.NotNil(t, res, "Create: expected response on call")
	assert.False(t, res.Ok, "Create: expected Ok false response")

	// invalid email
	req = &pb.ProfileCreationRequest{
		Gender:   pb.Genders_MALE,
		Email:    "test_example.com",
		Name:     "Profile Test Name",
		Birthday: "1989-11-09",
	}
	res, err = server.Create(ctx, req)
	assert.NotNil(t, err, "Create: expected error on call")
	assert.NotNil(t, res, "Create: expected response on call")
	assert.False(t, res.Ok, "Create: expected Ok false response")

	// invalid name
	req = &pb.ProfileCreationRequest{
		Gender:   pb.Genders_MALE,
		Email:    "test@example.com",
		Name:     "",
		Birthday: "1989-11-09",
	}
	res, err = server.Create(ctx, req)
	assert.NotNil(t, err, "Create: expected error on call")
	assert.NotNil(t, res, "Create: expected response on call")
	assert.False(t, res.Ok, "Create: expected Ok false response")

	// invalid birthday
	req = &pb.ProfileCreationRequest{
		Gender:   pb.Genders_MALE,
		Email:    "test@example.com",
		Name:     "Profile Test Name",
		Birthday: "1906-01-22",
	}
	res, err = server.Create(ctx, req)
	assert.NotNil(t, err, "Create: expected error on call")
	assert.NotNil(t, res, "Create: expected response on call")
	assert.False(t, res.Ok, "Create: expected Ok false response")

	// valid
	req = pb.NewProfileCreationRequestGomeetFaker()
	res, err = server.Create(ctx, req)
	assert.Nil(t, err, "Create: error on call")
	assert.NotNil(t, res, "Create: error on call")
	assert.True(t, res.Ok, "Create: error on call")
	assert.NotEmpty(t, res.GetInfo().GetUuid(), "Create: uuid is empty")
	assert.Equal(t, req.GetGender(), res.GetInfo().GetGender(), "Create: gender doesn't match")
	assert.Equal(t, req.GetEmail(), res.GetInfo().GetEmail(), "Create: email doesn't match")
	assert.Equal(t, req.GetName(), res.GetInfo().GetName(), "Create: name doesn't match")
	assert.Equal(t, req.GetBirthday(), res.GetInfo().GetBirthday(), "Create: birthday doesn't match")
	assert.NotEmpty(t, res.GetInfo().GetCreatedAt(), "Create: created_at is empty")
	assert.NotEmpty(t, res.GetInfo().GetUpdatedAt(), "Create: updated_at is empty")
	assert.Empty(t, res.GetInfo().GetDeletedAt(), "Create: deleted_at is not empty")

	// no duplicate email
	res, err = server.Create(ctx, req)
	assert.NotNil(t, err, "Create: expected error on call")
	assert.NotNil(t, res, "Create: expected response on call")
	assert.False(t, res.Ok, "Create: expected Ok false response")
}
