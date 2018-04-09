package functest

import (
	"fmt"
	"regexp"

	pb "github.com/gomeet-examples/svc-profile/pb"
	"github.com/google/uuid"
)

func testGetSoftDeleteRequest(
	config FunctionalTestConfig,
) (reqs []*pb.ProfileRequest, extras map[string]interface{}, err error) {
	// create a valid profile
	client, ctx, err := grpcClient(config)
	if err != nil {
		return reqs, extras, fmt.Errorf("Read - gRPC client initialization error (%v)", err)
	}
	defer client.Close()
	res, err := client.GetGRPCClient().Create(ctx, pb.NewProfileCreationRequestGomeetFaker())
	if res == nil || err != nil || res.GetOk() != true {
		return reqs, extras, fmt.Errorf("Read - profile creation error (%v)", err)
	}
	// set valid profile into extras
	validProfile := res.GetInfo()
	extras = make(map[string]interface{}, 1)
	extras[validProfile.GetUuid()] = validProfile

	// errors cases
	reqs = append(reqs, &pb.ProfileRequest{Uuid: ""})                  // empty Uuid
	reqs = append(reqs, &pb.ProfileRequest{Uuid: "invalid uuid"})      // invalid Uuid
	reqs = append(reqs, &pb.ProfileRequest{Uuid: uuid.New().String()}) // valid uuid not found

	// valid index 3+
	reqs = append(reqs, &pb.ProfileRequest{Uuid: validProfile.GetUuid()})
	return reqs, extras, err
}

func testSoftDeleteResponse(
	config FunctionalTestConfig,
	testsType string,
	testCaseResults []*TestCaseResult,
	extras map[string]interface{},
) (failures []TestFailure) {
	client, ctx, err := grpcClient(config)
	if err != nil {
		failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: fmt.Sprintf("gRPC client initialization error (%v)", err)})
		return failures
	}
	defer client.Close()

	dtRegex := regexp.MustCompile("^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$")

	for i, tr := range testCaseResults {
		var (
			req *pb.ProfileRequest
			res *pb.ProfileResponse
			err error
			ok  bool
		)

		if tr.Request == nil {
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: "expected request message type pb.ProfileRequest - nil given"})
			continue
		}
		req, ok = tr.Request.(*pb.ProfileRequest)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: "expected request message type pb.ProfileRequest - cast fail"})
			continue
		}

		err = tr.Error
		if i < 3 {
			if err == nil {
				failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: "an error is expected"})
			}
			continue
		}

		if err != nil {
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: "no error expected"})
			continue
		}

		if tr.Response == nil {
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: "a response is expected"})
			continue
		}
		res, ok = tr.Response.(*pb.ProfileResponse)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: "expected response message type pb.ProfileInfo - cast fail"})
			continue
		}

		if req == nil || res == nil {
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: "a request and a response are expected"})
			continue
		}

		expectedProfile, ok := extras[res.GetInfo().GetUuid()]
		if !ok {
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: fmt.Sprintf("profile %s is not in expected extras map - res (%v) extras (%v)", res.GetInfo().GetUuid(), res, extras)})
			continue
		}
		validProfile, ok := expectedProfile.(*pb.ProfileInfo)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: fmt.Sprintf("expected Profile in extras map - cast fail - extras (%v)", extras)})
			continue
		}
		if res.GetInfo().GetUuid() != req.GetUuid() {
			failureMsg := fmt.Sprintf("expected Uuid \"%s\" but got \"%s\" for request: %v", req.GetUuid(), res.GetInfo().GetUuid(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if res.GetInfo().GetUuid() != validProfile.GetUuid() {
			failureMsg := fmt.Sprintf("expected Uuid \"%s\" but got \"%s\" for request: %v", validProfile.GetUuid(), res.GetInfo().GetUuid(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if res.GetInfo().GetGender() != validProfile.GetGender() {
			failureMsg := fmt.Sprintf("expected Gender \"%s\" but got \"%s\" for request: %v", validProfile.GetGender(), res.GetInfo().GetGender(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if res.GetInfo().GetEmail() != validProfile.GetEmail() {
			failureMsg := fmt.Sprintf("expected Email \"%s\" but got \"%s\" for request: %v", validProfile.GetEmail(), res.GetInfo().GetEmail(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if res.GetInfo().GetName() != validProfile.GetName() {
			failureMsg := fmt.Sprintf("expected Name \"%s\" but got \"%s\" for request: %v", validProfile.GetName(), res.GetInfo().GetName(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if res.GetInfo().GetBirthday() != validProfile.GetBirthday() {
			failureMsg := fmt.Sprintf("expected Birthday \"%s\" but got \"%s\" for request: %v", validProfile.GetBirthday(), res.GetInfo().GetBirthday(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if res.GetInfo().GetCreatedAt() != validProfile.GetCreatedAt() {
			failureMsg := fmt.Sprintf("expected CreatedAt \"%s\" but got \"%s\" for request: %v", validProfile.GetCreatedAt(), res.GetInfo().GetCreatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if !dtRegex.MatchString(res.GetInfo().GetCreatedAt()) {
			failureMsg := fmt.Sprintf("expected CreatedAt date in good format but got \"%s\" for request: %v", res.GetInfo().GetCreatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if res.GetInfo().GetUpdatedAt() != validProfile.GetUpdatedAt() {
			failureMsg := fmt.Sprintf("expected UpdatedAt \"%s\" but got \"%s\" for request: %v", validProfile.GetUpdatedAt(), res.GetInfo().GetUpdatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if !dtRegex.MatchString(res.GetInfo().GetUpdatedAt()) {
			failureMsg := fmt.Sprintf("expected UpdatedAt date in good format but got \"%s\" for request: %v", res.GetInfo().GetUpdatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if res.GetInfo().GetDeletedAt() == "" {
			failureMsg := fmt.Sprintf("expected DeletedAt \"%s\" doesn't be empty for request: %v", res.GetInfo().GetDeletedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if !dtRegex.MatchString(res.GetInfo().GetDeletedAt()) {
			failureMsg := fmt.Sprintf("expected DeletedAt date in good format but got \"%s\" for request: %v", res.GetInfo().GetDeletedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		// When reading the profile the field "deletedAt" must be present
		resRead, err := client.GetGRPCClient().Read(ctx, &pb.ProfileRequest{Uuid: res.GetInfo().GetUuid()})
		if res.GetInfo().GetDeletedAt() != resRead.GetDeletedAt() {
			failureMsg := fmt.Sprintf("expected DeletedAt \"%s\" but got \"%s\" for request: %v", res.GetInfo().GetDeletedAt(), resRead.GetDeletedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if resRead.GetDeletedAt() == "" {
			failureMsg := fmt.Sprintf("expected DeletedAt \"%s\" doesn't be empty for request: %v", resRead.GetDeletedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
		if !dtRegex.MatchString(resRead.GetDeletedAt()) {
			failureMsg := fmt.Sprintf("expected DeletedAt date in good format but got \"%s\" for request: %v", resRead.GetDeletedAt(), req)
			failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: failureMsg})
		}
	}

	if len(extras) > 0 {
		for sUuid := range extras {
			res, err := client.GetGRPCClient().HardDelete(ctx, &pb.ProfileRequest{Uuid: sUuid})
			if res == nil || err != nil || res.GetOk() != true {
				failures = append(failures, TestFailure{Procedure: "SoftDelete", Message: fmt.Sprintf("deletion of created profile %s fails error (%v) - res (%v)", sUuid, err, res)})
			}
		}
	}

	return failures
}
