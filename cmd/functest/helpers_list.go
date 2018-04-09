package functest

import (
	"fmt"
	"regexp"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func testGetListRequest(
	config FunctionalTestConfig,
) (reqs []*pb.ProfileListRequest, extras map[string]interface{}, err error) {
	extras = make(map[string]interface{}, 6)
	// first creation 3 profiles with gender MALE
	// and 3 profiles with gender FEMALE
	client, ctx, err := grpcClient(config)
	if err != nil {
		return reqs, extras, fmt.Errorf("Read - gRPC client initialization error (%v)", err)
	}
	defer client.Close()

	for i := 0; i < 6; i++ {
		reqCreation := pb.NewProfileCreationRequestGomeetFaker()
		if i%2 == 0 {
			reqCreation.Gender = pb.Genders_MALE
		} else {
			reqCreation.Gender = pb.Genders_FEMALE
		}
		res, err := client.GetGRPCClient().Create(ctx, reqCreation)
		if res == nil || err != nil || res.GetOk() != true {
			return reqs, extras, fmt.Errorf("Read - profile creation error (%v)", err)
		}
		extras[res.GetInfo().GetUuid()] = res.GetInfo()
	}

	// errors case
	reqs = append(reqs, &pb.ProfileListRequest{})
	reqs = append(reqs, &pb.ProfileListRequest{PageNumber: 0})
	reqs = append(reqs, &pb.ProfileListRequest{ExcludeSoftDeleted: true, SoftDeletedOnly: true})

	// valid 3+ - PageSize 1000 - normaly never out of range
	//   idx 4 - only MALE
	//   idx 5 - only FEMALE
	reqs = append(reqs, &pb.ProfileListRequest{PageNumber: 1, PageSize: 1000})
	reqs = append(reqs, &pb.ProfileListRequest{PageNumber: 1, PageSize: 1000, Gender: pb.Genders_MALE})
	reqs = append(reqs, &pb.ProfileListRequest{PageNumber: 1, PageSize: 1000, Gender: pb.Genders_FEMALE})

	return reqs, extras, err
}

func testListResponse(
	config FunctionalTestConfig,
	testsType string,
	testCaseResults []*TestCaseResult,
	extras map[string]interface{},
) (failures []TestFailure) {

	dtRegex := regexp.MustCompile("^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$")

	for i, tr := range testCaseResults {
		var (
			req *pb.ProfileListRequest
			res *pb.ProfileList
			err error
			ok  bool
		)
		if tr.Request == nil {
			failures = append(failures, TestFailure{Procedure: "List", Message: "expected request message type pb.ProfileListRequest - nil given"})
			continue
		}
		req, ok = tr.Request.(*pb.ProfileListRequest)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "List", Message: "expected request message type pb.ProfileListRequest - cast fail"})
			continue
		}

		if tr.Response != nil {
			res, ok = tr.Response.(*pb.ProfileList)
			if !ok {
				failures = append(failures, TestFailure{Procedure: "List", Message: "expected response message type pb.ProfileListRequest - cast fail"})
				continue
			}
		}

		err = tr.Error
		if i < 3 {
			if err == nil {
				failures = append(failures, TestFailure{Procedure: "List", Message: "an error is expected"})
			}
			continue
		}

		if err != nil {
			failures = append(failures, TestFailure{Procedure: "List", Message: "no error expected"})
			continue
		}

		if tr.Response == nil {
			failures = append(failures, TestFailure{Procedure: "List", Message: "a response is expected"})
			continue
		}

		res, ok = tr.Response.(*pb.ProfileList)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "List", Message: "expected response message type pb.ProfileResponse - cast fail"})
			continue
		}

		if req == nil || res == nil {
			failures = append(failures, TestFailure{Procedure: "List", Message: "a request and a response are expected"})
			continue
		}

		profiles := res.GetProfiles()
		for _, profile := range profiles {
			expectedProfile, ok := extras[profile.GetUuid()]
			if !ok {
				// unknown profile in database
				continue
			}
			validProfile, ok := expectedProfile.(*pb.ProfileInfo)
			if !ok {
				failures = append(failures, TestFailure{Procedure: "List", Message: fmt.Sprintf("expected Profile in extras map - cast fail - extras (%v)", extras)})
				continue
			}
			if profile.GetUuid() != validProfile.GetUuid() {
				failureMsg := fmt.Sprintf("expected Uuid \"%s\" but got \"%s\" for request: %v", validProfile.GetUuid(), profile.GetUuid(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if profile.GetGender() != validProfile.GetGender() {
				failureMsg := fmt.Sprintf("expected Gender \"%s\" but got \"%s\" for request: %v", validProfile.GetGender(), profile.GetGender(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if profile.GetEmail() != validProfile.GetEmail() {
				failureMsg := fmt.Sprintf("expected Email \"%s\" but got \"%s\" for request: %v", validProfile.GetEmail(), profile.GetEmail(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if profile.GetName() != validProfile.GetName() {
				failureMsg := fmt.Sprintf("expected Name \"%s\" but got \"%s\" for request: %v", validProfile.GetName(), profile.GetName(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if profile.GetBirthday() != validProfile.GetBirthday() {
				failureMsg := fmt.Sprintf("expected Birthday \"%s\" but got \"%s\" for request: %v", validProfile.GetBirthday(), profile.GetBirthday(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if profile.GetCreatedAt() != validProfile.GetCreatedAt() {
				failureMsg := fmt.Sprintf("expected CreatedAt \"%s\" but got \"%s\" for request: %v", validProfile.GetCreatedAt(), profile.GetCreatedAt(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if !dtRegex.MatchString(profile.GetCreatedAt()) {
				failureMsg := fmt.Sprintf("expected CreatedAt date in good format but got \"%s\" for request: %v", profile.GetCreatedAt(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if profile.GetUpdatedAt() != validProfile.GetUpdatedAt() {
				failureMsg := fmt.Sprintf("expected UpdatedAt \"%s\" but got \"%s\" for request: %v", validProfile.GetUpdatedAt(), profile.GetUpdatedAt(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if !dtRegex.MatchString(profile.GetUpdatedAt()) {
				failureMsg := fmt.Sprintf("expected UpdatedAt date in good format but got \"%s\" for request: %v", profile.GetUpdatedAt(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if profile.GetDeletedAt() != validProfile.GetDeletedAt() {
				failureMsg := fmt.Sprintf("expected DeletedAt \"%s\" but got \"%s\" for request: %v", validProfile.GetDeletedAt(), profile.GetDeletedAt(), req)
				failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
			}
			if profile.GetDeletedAt() != "" {
				failureMsg := fmt.Sprintf("expected DeletedAt \"%s\" must be empty for request: %v", profile.GetDeletedAt(), req)
				failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
			}
			switch i {
			case 4:
				if profile.GetGender() != pb.Genders_MALE {
					failureMsg := fmt.Sprintf("expected Gender \"%s\" but got \"%s\" for request: %v", pb.Genders_MALE, profile.GetGender(), req)
					failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
				}
			case 5:
				if profile.GetGender() != pb.Genders_FEMALE {
					failureMsg := fmt.Sprintf("expected Gender \"%s\" but got \"%s\" for request: %v", pb.Genders_FEMALE, profile.GetGender(), req)
					failures = append(failures, TestFailure{Procedure: "List", Message: failureMsg})
				}
			}
		}
	}

	if len(extras) > 0 {
		client, ctx, err := grpcClient(config)
		if err != nil {
			failures = append(failures, TestFailure{Procedure: "List", Message: fmt.Sprintf("gRPC client initialization error (%v)", err)})
			return failures
		}
		defer client.Close()
		for sUuid := range extras {
			res, err := client.GetGRPCClient().HardDelete(ctx, &pb.ProfileRequest{Uuid: sUuid})
			if res == nil || err != nil || res.GetOk() != true {
				failures = append(failures, TestFailure{Procedure: "List", Message: fmt.Sprintf("deletion of created profile %s fails error (%v) - res (%v)", sUuid, err, res)})
			}
		}
	}

	return failures
}
