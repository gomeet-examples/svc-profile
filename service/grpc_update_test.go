package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestUpdate(t *testing.T) {
	var (
		req *pb.ProfileInfo
		res *pb.ProfileResponse
		err error
	)

	// FIXME: do something better with mock?
	server, db, err := newProfileServerTest(t)
	assert.NotNil(t, server, "Update: newProfileServerTest nil server")
	assert.NotNil(t, db, "Update: newProfileServerTest nil db")
	assert.Nil(t, err, "Update: newProfileServerTest not nil err")

	ctx := context.Background()

	// first create a valid profile
	reqCreation := pb.NewProfileCreationRequestGomeetFaker()
	res, err = server.Create(ctx, reqCreation)
	assert.Nil(t, err, "Update: error on call")
	assert.NotNil(t, res, "Update: error on call")
	assert.True(t, res.Ok, "Update: error on call")
	assert.NotEmpty(t, res.GetInfo().GetUuid(), "Update: uuid is empty")
	assert.Equal(t, reqCreation.GetGender(), res.GetInfo().GetGender(), "Update: gender doesn't match")
	assert.Equal(t, reqCreation.GetEmail(), res.GetInfo().GetEmail(), "Update: email doesn't match")
	assert.Equal(t, reqCreation.GetName(), res.GetInfo().GetName(), "Update: name doesn't match")
	assert.Equal(t, reqCreation.GetBirthday(), res.GetInfo().GetBirthday(), "Update: birthday doesn't match")
	assert.NotEmpty(t, res.GetInfo().GetCreatedAt(), "Update: created_at is empty")
	assert.NotEmpty(t, res.GetInfo().GetUpdatedAt(), "Update: updated_at is empty")
	assert.Empty(t, res.GetInfo().GetDeletedAt(), "Update: deleted_at is not empty")

	validProfile := res.GetInfo()

	// second: invalid update

	// empty Uuid
	req = &pb.ProfileInfo{
		Uuid:     "",
		Gender:   validProfile.GetGender(),
		Email:    validProfile.GetEmail(),
		Name:     validProfile.GetName(),
		Birthday: validProfile.GetBirthday(),
	}
	res, err = server.Update(ctx, req)
	assert.NotNil(t, err, "Update: expected error on call")
	assert.NotNil(t, res, "Update: expected response on call")
	assert.False(t, res.Ok, "Update: expected Ok false response")

	// invalid Uuid
	req = &pb.ProfileInfo{
		Uuid:     "invalid uuid",
		Gender:   validProfile.GetGender(),
		Email:    validProfile.GetEmail(),
		Name:     validProfile.GetName(),
		Birthday: validProfile.GetBirthday(),
	}
	res, err = server.Update(ctx, req)
	assert.NotNil(t, err, "Update: expected error on call")
	assert.NotNil(t, res, "Update: expected response on call")
	assert.False(t, res.Ok, "Update: expected Ok false response")

	// not found Uuid
	req = &pb.ProfileInfo{
		Uuid:     uuid.New().String(),
		Gender:   validProfile.GetGender(),
		Email:    validProfile.GetEmail(),
		Name:     validProfile.GetName(),
		Birthday: validProfile.GetBirthday(),
	}
	res, err = server.Update(ctx, req)
	assert.NotNil(t, err, "Update: expected error on call")
	assert.NotNil(t, res, "Update: expected response on call")
	assert.False(t, res.Ok, "Update: expected Ok false response")

	// invalid gender
	req = &pb.ProfileInfo{
		Uuid:     validProfile.GetUuid(),
		Gender:   pb.Genders_UNKNOW,
		Email:    validProfile.GetEmail(),
		Name:     validProfile.GetName(),
		Birthday: validProfile.GetBirthday(),
	}
	res, err = server.Update(ctx, req)
	assert.NotNil(t, err, "Update: expected error on call")
	assert.NotNil(t, res, "Update: expected response on call")
	assert.False(t, res.Ok, "Update: expected Ok false response")

	// invalid email
	req = &pb.ProfileInfo{
		Uuid:     validProfile.GetUuid(),
		Gender:   validProfile.GetGender(),
		Email:    "test_example.com",
		Name:     validProfile.GetName(),
		Birthday: validProfile.GetBirthday(),
	}
	res, err = server.Update(ctx, req)
	assert.NotNil(t, err, "Update: expected error on call")
	assert.NotNil(t, res, "Update: expected response on call")
	assert.False(t, res.Ok, "Update: expected Ok false response")

	// invalid name
	req = &pb.ProfileInfo{
		Uuid:     validProfile.GetUuid(),
		Gender:   validProfile.GetGender(),
		Email:    validProfile.GetEmail(),
		Name:     "",
		Birthday: validProfile.GetBirthday(),
	}
	res, err = server.Update(ctx, req)
	assert.NotNil(t, err, "Update: expected error on call")
	assert.NotNil(t, res, "Update: expected response on call")
	assert.False(t, res.Ok, "Update: expected Ok false response")

	// invalid birthday
	req = &pb.ProfileInfo{
		Uuid:     validProfile.GetUuid(),
		Gender:   validProfile.GetGender(),
		Email:    validProfile.GetEmail(),
		Name:     validProfile.GetName(),
		Birthday: "1906-01-22",
	}
	res, err = server.Update(ctx, req)
	assert.NotNil(t, err, "Update: expected error on call")
	assert.NotNil(t, res, "Update: expected response on call")
	assert.False(t, res.Ok, "Update: expected Ok false response")

	// valid changes
	req = &pb.ProfileInfo{
		Uuid:     validProfile.GetUuid(),
		Gender:   pb.Genders_MALE,
		Email:    "test@example.com",
		Name:     "Profile Test Name",
		Birthday: "1976-12-13",
	}
	res, err = server.Update(ctx, req)
	assert.Nil(t, err, "Update: error on call")
	assert.NotNil(t, res, "Update: error on call")
	assert.True(t, res.Ok, "Update: error on call")
	assert.NotEmpty(t, res.GetInfo().GetUuid(), "Update: uuid is empty")
	assert.Equal(t, req.GetGender(), res.GetInfo().GetGender(), "Update: gender doesn't match")
	assert.Equal(t, req.GetEmail(), res.GetInfo().GetEmail(), "Update: email doesn't match")
	assert.Equal(t, req.GetName(), res.GetInfo().GetName(), "Update: name doesn't match")
	assert.Equal(t, req.GetBirthday(), res.GetInfo().GetBirthday(), "Update: birthday doesn't match")
	assert.NotEmpty(t, res.GetInfo().GetUpdatedAt(), "Update: created_at is empty")
	assert.NotEmpty(t, res.GetInfo().GetUpdatedAt(), "Update: updated_at is empty")
	assert.Empty(t, res.GetInfo().GetDeletedAt(), "Update: deleted_at is not empty")

	// invalid duplicate email
	// insert a second profile
	reqCreation = pb.NewProfileCreationRequestGomeetFaker()
	res, err = server.Create(ctx, reqCreation)
	assert.Nil(t, err, "Update: error on call")
	assert.NotNil(t, res, "Update: error on call")
	assert.True(t, res.Ok, "Update: error on call")
	assert.NotEmpty(t, res.GetInfo().GetUuid(), "Update: uuid is empty")
	assert.Equal(t, reqCreation.GetGender(), res.GetInfo().GetGender(), "Update: gender doesn't match")
	assert.Equal(t, reqCreation.GetEmail(), res.GetInfo().GetEmail(), "Update: email doesn't match")
	assert.Equal(t, reqCreation.GetName(), res.GetInfo().GetName(), "Update: name doesn't match")
	assert.Equal(t, reqCreation.GetBirthday(), res.GetInfo().GetBirthday(), "Update: birthday doesn't match")
	assert.NotEmpty(t, res.GetInfo().GetCreatedAt(), "Update: created_at is empty")
	assert.NotEmpty(t, res.GetInfo().GetUpdatedAt(), "Update: updated_at is empty")
	assert.Empty(t, res.GetInfo().GetDeletedAt(), "Update: deleted_at is not empty")

	validProfile2 := res.GetInfo()

	// invalid duplicate email
	req = &pb.ProfileInfo{
		Uuid:     validProfile.GetUuid(),
		Gender:   validProfile.GetGender(),
		Email:    validProfile2.GetEmail(),
		Name:     validProfile.GetName(),
		Birthday: validProfile.GetBirthday(),
	}
	res, err = server.Update(ctx, req)
	assert.NotNil(t, err, "Update: expected error on call")
	assert.NotNil(t, res, "Update: expected response on call")
	assert.False(t, res.Ok, "Update: expected Ok false response")
}
