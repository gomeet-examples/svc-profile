package functest

import (
	"fmt"
	"regexp"

	pb "github.com/gomeet-examples/svc-profile/pb"
	"github.com/google/uuid"
)

func testGetUpdateRequest(
	config FunctionalTestConfig,
) (reqs []*pb.ProfileInfo, extras map[string]interface{}, err error) {
	// create a valid profile
	client, ctx, err := grpcClient(config)
	if err != nil {
		return reqs, extras, fmt.Errorf("Update - gRPC client initialization error (%v)", err)
	}
	defer client.Close()

	res, err := client.GetGRPCClient().Create(ctx, pb.NewProfileCreationRequestGomeetFaker())
	if res == nil || err != nil || res.GetOk() != true {
		return reqs, extras, fmt.Errorf("Update - profile creation error (%v)", err)
	}
	// set valid profile into extras
	validProfile := res.GetInfo()

	// create a second valid Profile
	res, err = client.GetGRPCClient().Create(ctx, pb.NewProfileCreationRequestGomeetFaker())
	if res == nil || err != nil || res.GetOk() != true {
		return reqs, extras, fmt.Errorf("Update - profile creation error (%v)", err)
	}
	// set valid profile into extras
	validProfile2 := res.GetInfo()

	// init a valid changes
	validChanges := *validProfile
	switch validChanges.GetGender() {
	case pb.Genders_MALE:
		validChanges.Gender = pb.Genders_FEMALE
	case pb.Genders_FEMALE:
		validChanges.Gender = pb.Genders_MALE
	}
	validChanges.Name = fmt.Sprintf("%s (updated)", validProfile.GetName())
	validChanges.Email = fmt.Sprintf("%s@%s.com", uuid.New().String(), uuid.New().String())
	if validProfile.GetBirthday() != "1976-12-13" {
		validChanges.Birthday = "1976-12-13"
	} else {
		validChanges.Birthday = "1980-01-13"
	}

	extras = make(map[string]interface{}, 2)
	extras[validProfile.GetUuid()] = []*pb.ProfileInfo{validProfile, &validChanges}
	extras[validProfile2.GetUuid()] = []*pb.ProfileInfo{validProfile2, nil}

	// errors cases
	//   - empty Uuid
	//   - invalid Uuid
	//   - not found Uuid
	//   - invalid gender
	//   - invalid email
	//   - invalid email duplicate
	//   - empty name
	//   - invalid birthday 1976-13-12
	//   - invalid birthday > 99
	reqs = append(reqs, &pb.ProfileInfo{Uuid: "", Gender: validChanges.GetGender(), Email: validChanges.GetEmail(), Name: validChanges.GetName(), Birthday: validChanges.GetBirthday()})
	reqs = append(reqs, &pb.ProfileInfo{Uuid: "invalid uuid", Gender: validChanges.GetGender(), Email: validChanges.GetEmail(), Name: validChanges.GetName(), Birthday: validChanges.GetBirthday()})
	reqs = append(reqs, &pb.ProfileInfo{Uuid: uuid.New().String(), Gender: validChanges.GetGender(), Email: validChanges.GetEmail(), Name: validChanges.GetName(), Birthday: validChanges.GetBirthday()})
	reqs = append(reqs, &pb.ProfileInfo{Uuid: validChanges.GetUuid(), Gender: pb.Genders_UNKNOW, Email: validChanges.GetEmail(), Name: validChanges.GetName(), Birthday: validChanges.GetBirthday()})
	reqs = append(reqs, &pb.ProfileInfo{Uuid: validChanges.GetUuid(), Gender: validChanges.GetGender(), Email: "test_example.com", Name: validChanges.GetName(), Birthday: validChanges.GetBirthday()})
	reqs = append(reqs, &pb.ProfileInfo{Uuid: validChanges.GetUuid(), Gender: validChanges.GetGender(), Email: validProfile2.GetEmail(), Name: validChanges.GetName(), Birthday: validChanges.GetBirthday()})
	reqs = append(reqs, &pb.ProfileInfo{Uuid: validChanges.GetUuid(), Gender: validChanges.GetGender(), Email: validChanges.GetEmail(), Name: "", Birthday: validChanges.GetBirthday()})
	reqs = append(reqs, &pb.ProfileInfo{Uuid: validChanges.GetUuid(), Gender: validChanges.GetGender(), Email: validChanges.GetEmail(), Name: validChanges.GetName(), Birthday: "1976-13-12"})
	reqs = append(reqs, &pb.ProfileInfo{Uuid: validChanges.GetUuid(), Gender: validChanges.GetGender(), Email: validChanges.GetEmail(), Name: validChanges.GetName(), Birthday: "1906-01-22"})

	// valid index 9+
	reqs = append(reqs, &pb.ProfileInfo{Uuid: validChanges.GetUuid(), Gender: validChanges.GetGender(), Email: validChanges.GetEmail(), Name: validChanges.GetName(), Birthday: validChanges.GetBirthday()})
	return reqs, extras, err
}

func testUpdateResponse(
	config FunctionalTestConfig,
	testsType string,
	testCaseResults []*TestCaseResult,
	extras map[string]interface{},
) (failures []TestFailure) {

	dtRegex := regexp.MustCompile("^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$")

	for i, tr := range testCaseResults {
		var (
			req *pb.ProfileInfo
			res *pb.ProfileResponse
			err error
			ok  bool
		)
		if tr.Request == nil {
			failures = append(failures, TestFailure{Procedure: "Update", Message: "expected request message type pb.ProfileInfo - nil given"})
			continue
		}
		req, ok = tr.Request.(*pb.ProfileInfo)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "Update", Message: "expected request message type pb.ProfileInfo - cast fail"})
			continue
		}

		err = tr.Error
		if i < 9 {
			if err == nil {
				failures = append(failures, TestFailure{Procedure: "Update", Message: "an error is expected"})
			}
			continue
		}

		if err != nil {
			failures = append(failures, TestFailure{Procedure: "Update", Message: "no error expected"})
			continue
		}

		if tr.Response == nil {
			failures = append(failures, TestFailure{Procedure: "Update", Message: "a response is expected"})
			continue
		}
		res, ok = tr.Response.(*pb.ProfileResponse)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "Update", Message: "expected response message type pb.ProfileInfo - cast fail"})
			continue
		}

		if req == nil || res == nil {
			failures = append(failures, TestFailure{Procedure: "Update", Message: "a request and a response are expected"})
			continue
		}

		if !res.GetOk() {
			failureMsg := fmt.Sprintf("expected response Ok \"true\" but got \"false\" for request: %v", req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
			continue
		}

		resProfile := res.GetInfo()
		expectedInfo, ok := extras[resProfile.GetUuid()]
		if !ok {
			failures = append(failures, TestFailure{Procedure: "Update", Message: fmt.Sprintf("profile %s is not in expected extras map - res (%v) extras (%v)", resProfile.GetUuid(), res, extras)})
			continue
		}
		expectedProfiles, ok := expectedInfo.([]*pb.ProfileInfo)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "Update", Message: fmt.Sprintf("expected Profile in extras map - cast fail - extras (%v)", extras)})
			continue
		}
		validProfile, validChanges := expectedProfiles[0], expectedProfiles[1]

		if resProfile.GetUuid() != req.GetUuid() {
			failureMsg := fmt.Sprintf("expected Uuid \"%s\" but got \"%s\" for request: %v", req.GetUuid(), resProfile.GetUuid(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if resProfile.GetUuid() != validChanges.GetUuid() {
			failureMsg := fmt.Sprintf("expected Uuid \"%s\" but got \"%s\" for request: %v", validChanges.GetUuid(), resProfile.GetUuid(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if resProfile.GetGender() != validChanges.GetGender() {
			failureMsg := fmt.Sprintf("expected Gender \"%s\" but got \"%s\" for request: %v", validChanges.GetGender(), resProfile.GetGender(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if resProfile.GetEmail() != validChanges.GetEmail() {
			failureMsg := fmt.Sprintf("expected Email \"%s\" but got \"%s\" for request: %v", validChanges.GetEmail(), resProfile.GetEmail(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if resProfile.GetName() != validChanges.GetName() {
			failureMsg := fmt.Sprintf("expected Name \"%s\" but got \"%s\" for request: %v", validChanges.GetName(), resProfile.GetName(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if resProfile.GetBirthday() != validChanges.GetBirthday() {
			failureMsg := fmt.Sprintf("expected Birthday \"%s\" but got \"%s\" for request: %v", validChanges.GetBirthday(), resProfile.GetBirthday(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if resProfile.GetCreatedAt() != validChanges.GetCreatedAt() {
			failureMsg := fmt.Sprintf("expected CreatedAt \"%s\" but got \"%s\" for request: %v", validChanges.GetCreatedAt(), resProfile.GetCreatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if resProfile.GetCreatedAt() != validProfile.GetCreatedAt() {
			failureMsg := fmt.Sprintf("expected CreatedAt \"%s\" but got \"%s\" for request: %v", validProfile.GetCreatedAt(), resProfile.GetCreatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if !dtRegex.MatchString(resProfile.GetCreatedAt()) {
			failureMsg := fmt.Sprintf("expected CreatedAt date in good format but got \"%s\" for request: %v", resProfile.GetCreatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		if !dtRegex.MatchString(resProfile.GetUpdatedAt()) {
			failureMsg := fmt.Sprintf("expected UpdatedAt date in good format but got \"%s\" for request: %v", resProfile.GetUpdatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		// the update is too fast for that :
		/* if resProfile.GetUpdatedAt() == validProfile.GetUpdatedAt() {
			failureMsg := fmt.Sprintf("expected UpdatedAt \"%s\" doesn't change for request: %v", resProfile.GetUpdatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		} */
		if resProfile.GetDeletedAt() != "" {
			failureMsg := fmt.Sprintf("expected DeletedAt \"%s\" must be empty for request: %v", resProfile.GetDeletedAt(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
	}

	if len(extras) > 0 {
		client, ctx, err := grpcClient(config)
		if err != nil {
			failures = append(failures, TestFailure{Procedure: "Update", Message: fmt.Sprintf("gRPC client initialization error (%v)", err)})
			return failures
		}
		defer client.Close()
		for sUuid, _ := range extras {
			res, err := client.GetGRPCClient().HardDelete(ctx, &pb.ProfileRequest{Uuid: sUuid})
			if res == nil || err != nil || res.GetOk() != true {
				failures = append(failures, TestFailure{Procedure: "Update", Message: fmt.Sprintf("deletion of created profile %s fails error (%v) - res (%v)", sUuid, err, res)})
			}
		}
	}

	return failures
}
