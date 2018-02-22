package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func TestList(t *testing.T) {
	var (
		req            *pb.ProfileListRequest
		res            *pb.ProfileList
		err            error
		expectedLength uint32
		validProfiles  []*pb.ProfileInfo
	)

	// FIXME: do something better with mock?
	server, db, err := newProfileServerTest(t)
	assert.NotNil(t, server, "List: newProfileServerTest nil server")
	assert.NotNil(t, db, "List: newProfileServerTest nil db")
	assert.Nil(t, err, "List: newProfileServerTest not nil err")

	ctx := context.Background()

	// first creation 3 profiles with gender MALE
	// and 3 profiles with gender FEMALE
	validProfilesByGender := make(map[uint32][]*pb.ProfileInfo)
	for i := 0; i < 6; i++ {
		reqCreation := pb.NewProfileCreationRequestGomeetFaker()
		if i%2 == 0 {
			reqCreation.Gender = pb.Genders_MALE
		} else {
			reqCreation.Gender = pb.Genders_FEMALE
		}
		resCreation, err := server.Create(ctx, reqCreation)
		assert.Nil(t, err, "List: error on call")
		assert.NotNil(t, resCreation, "List: error on call")
		assert.True(t, resCreation.Ok, "List:: response Ok")
		assert.NotEmpty(t, resCreation.GetInfo().GetUuid(), "List: uuid is empty")
		assert.Equal(t, reqCreation.GetGender(), resCreation.GetInfo().GetGender(), "List: gender doesn't match")
		assert.Equal(t, reqCreation.GetEmail(), resCreation.GetInfo().GetEmail(), "List: email doesn't match")
		assert.Equal(t, reqCreation.GetName(), resCreation.GetInfo().GetName(), "List: name doesn't match")
		assert.Equal(t, reqCreation.GetBirthday(), resCreation.GetInfo().GetBirthday(), "List: birthday doesn't match")
		assert.NotEmpty(t, resCreation.GetInfo().GetCreatedAt(), "List: created_at is empty")
		assert.NotEmpty(t, resCreation.GetInfo().GetUpdatedAt(), "List: updated_at is empty")
		assert.Empty(t, resCreation.GetInfo().GetDeletedAt(), "List: deleted_at is not empty")
		validProfiles = append(validProfiles, resCreation.GetInfo())
		validProfilesByGender[uint32(reqCreation.GetGender())] = append(validProfilesByGender[uint32(reqCreation.GetGender())], resCreation.GetInfo())
	}

	// invalid request
	req = &pb.ProfileListRequest{}
	res, err = server.List(ctx, req)
	expectedLength = 0
	assert.NotNil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")

	// invalid request ExcludeSoftDeleted == true and SoftDeletedOnly == true
	req = &pb.ProfileListRequest{
		PageNumber:         1,
		PageSize:           50,
		Order:              "created_at asc",
		ExcludeSoftDeleted: true,
		SoftDeletedOnly:    true,
		Gender:             pb.Genders_UNKNOW,
	}
	res, err = server.List(ctx, req)
	expectedLength = 0
	assert.NotNil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")

	// valid request without gender
	req = &pb.ProfileListRequest{
		PageNumber:         1,
		PageSize:           50,
		Order:              "created_at asc",
		ExcludeSoftDeleted: false,
		SoftDeletedOnly:    false,
		Gender:             pb.Genders_UNKNOW,
	}
	res, err = server.List(ctx, req)
	expectedLength = 6
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")

	// valid MALE only
	req.Gender = pb.Genders_MALE
	res, err = server.List(ctx, req)
	expectedLength = 3
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")
	for _, pNfo := range res.GetProfiles() {
		assert.Equal(t, pb.Genders_MALE, pNfo.GetGender(), "List: gender doesn't match")
	}

	// valid MALE only
	req.Gender = pb.Genders_FEMALE
	res, err = server.List(ctx, req)
	expectedLength = 3
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")
	for _, pNfo := range res.GetProfiles() {
		assert.Equal(t, pb.Genders_FEMALE, pNfo.GetGender(), "List: gender doesn't match")
	}

	// valid pagination
	req = &pb.ProfileListRequest{
		PageNumber:         1,
		PageSize:           2,
		Order:              "created_at asc",
		ExcludeSoftDeleted: false,
		SoftDeletedOnly:    false,
		Gender:             pb.Genders_UNKNOW,
	}
	res, err = server.List(ctx, req)
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(6), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.True(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(req.PageSize), uint32(len(res.GetProfiles())), "List: profiles not empty")
	// page 2
	req.PageNumber = 2
	res, err = server.List(ctx, req)
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(6), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.True(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(req.PageSize), uint32(len(res.GetProfiles())), "List: profiles not empty")
	// page 3
	req.PageNumber = 3
	res, err = server.List(ctx, req)
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(6), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(req.PageSize), uint32(len(res.GetProfiles())), "List: profiles not empty")

	// test ExcludeSoftDeleted and SoftDeletedOnly
	profileToDelete := []string{
		validProfilesByGender[uint32(pb.Genders_MALE)][0].GetUuid(),
		validProfilesByGender[uint32(pb.Genders_FEMALE)][0].GetUuid(),
	}
	for _, pUuid := range profileToDelete {
		reqDeletion := &pb.ProfileRequest{
			Uuid: pUuid,
		}
		resDeletion, err := server.SoftDelete(ctx, reqDeletion)
		assert.Nil(t, err, "SoftDelete: error on call")
		assert.NotNil(t, resDeletion, "SoftDelete: error on call")
		assert.True(t, resDeletion.GetOk(), "SoftDelete: response Ok")
	}
	// ExcludeSoftDeleted == true
	req = &pb.ProfileListRequest{
		PageNumber:         1,
		PageSize:           50,
		Order:              "created_at asc",
		ExcludeSoftDeleted: true,
		SoftDeletedOnly:    false,
		Gender:             pb.Genders_UNKNOW,
	}
	res, err = server.List(ctx, req)
	expectedLength = 4
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")

	// valid MALE only
	req.Gender = pb.Genders_MALE
	res, err = server.List(ctx, req)
	expectedLength = 2
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")
	for _, pNfo := range res.GetProfiles() {
		assert.Equal(t, pb.Genders_MALE, pNfo.GetGender(), "List: gender doesn't match")
	}

	// valid MALE only
	req.Gender = pb.Genders_FEMALE
	res, err = server.List(ctx, req)
	expectedLength = 2
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")
	for _, pNfo := range res.GetProfiles() {
		assert.Equal(t, pb.Genders_FEMALE, pNfo.GetGender(), "List: gender doesn't match")
	}

	// SoftDeletedOnly == true
	req = &pb.ProfileListRequest{
		PageNumber:         1,
		PageSize:           50,
		Order:              "created_at asc",
		ExcludeSoftDeleted: false,
		SoftDeletedOnly:    true,
		Gender:             pb.Genders_UNKNOW,
	}
	res, err = server.List(ctx, req)
	expectedLength = 2
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")

	// valid MALE only
	req.Gender = pb.Genders_MALE
	res, err = server.List(ctx, req)
	expectedLength = 1
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")
	for _, pNfo := range res.GetProfiles() {
		assert.Equal(t, pb.Genders_MALE, pNfo.GetGender(), "List: gender doesn't match")
	}

	// valid MALE only
	req.Gender = pb.Genders_FEMALE
	res, err = server.List(ctx, req)
	expectedLength = 1
	assert.Nil(t, err, "List: error on call")
	assert.NotNil(t, res, "List: error on call")
	assert.Equal(t, uint32(expectedLength), uint32(res.GetResultSetSize()), "List: result_set_size not empty")
	assert.False(t, res.GetHasMore(), "List: has_more true false expected")
	assert.Equal(t, uint32(expectedLength), uint32(len(res.GetProfiles())), "List: profiles not empty")
	assert.Equal(t, uint32(res.GetResultSetSize()), uint32(len(res.GetProfiles())), "List: result_set_size not empty")
	for _, pNfo := range res.GetProfiles() {
		assert.Equal(t, pb.Genders_FEMALE, pNfo.GetGender(), "List: gender doesn't match")
	}

}
