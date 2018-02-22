package functest

import (
	"fmt"
	"regexp"

	pb "github.com/gomeet-examples/svc-profile/pb"
	"github.com/google/uuid"
)

func testGetCreateRequest(
	config FunctionalTestConfig,
) (reqs []*pb.ProfileCreationRequest, extras map[string]interface{}, err error) {
	// errors cases
	reqs = append(reqs, &pb.ProfileCreationRequest{Gender: 0, Email: "test@example.com", Name: "Profile Test Name", Birthday: "1989-11-09"}) // invalid gender
	reqs = append(reqs, &pb.ProfileCreationRequest{Gender: 1, Email: "test_example.com", Name: "Profile Test Name", Birthday: "1989-11-09"}) // invalid email
	reqs = append(reqs, &pb.ProfileCreationRequest{Gender: 1, Email: "test@example.com", Name: "", Birthday: "1989-11-09"})                  // invalid name
	reqs = append(reqs, &pb.ProfileCreationRequest{Gender: 1, Email: "test@example.com", Name: "Profile Test Name", Birthday: "1906-01-22"}) // invalid birthday
	reqs = append(reqs, &pb.ProfileCreationRequest{Gender: 1, Email: "test@example.com", Name: "Profile Test Name", Birthday: "1976-13-12"}) // invalid birthday

	// valid index 5+
	reqs = append(reqs, pb.NewProfileCreationRequestGomeetFaker())
	return reqs, extras, err
}

func testCreateResponse(
	config FunctionalTestConfig,
	testsType string,
	testCaseResults []*TestCaseResult,
	extras map[string]interface{},
) (failures []TestFailure) {
	validProfiles := []*pb.ProfileInfo{}

	dtRegex := regexp.MustCompile("^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$")

	for i, tr := range testCaseResults {
		var (
			req *pb.ProfileCreationRequest
			res *pb.ProfileResponse
			err error
			ok  bool
		)
		if tr.Request == nil {
			failures = append(failures, TestFailure{Procedure: "Create", Message: "expected request message type pb.ProfileCreationRequest - nil given"})
			continue
		}
		req, ok = tr.Request.(*pb.ProfileCreationRequest)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "Create", Message: "expected request message type pb.ProfileCreationRequest - cast fail"})
			continue
		}

		err = tr.Error
		if i < 5 {
			if err == nil {
				failures = append(failures, TestFailure{Procedure: "Create", Message: "an error is expected"})
			}
			continue
		}

		if err != nil {
			failures = append(failures, TestFailure{Procedure: "Create", Message: "no error expected"})
			continue
		}

		if tr.Response == nil {
			failures = append(failures, TestFailure{Procedure: "Create", Message: "a response is expected"})
			continue
		}

		res, ok = tr.Response.(*pb.ProfileResponse)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "Create", Message: "expected response message type pb.ProfileCreationRequest - cast fail"})
			continue
		}

		if req == nil || res == nil {
			failures = append(failures, TestFailure{Procedure: "Create", Message: "a request and a response are expected"})
			continue
		}

		if _, err := uuid.Parse(res.GetInfo().GetUuid()); err != nil {
			failureMsg := fmt.Sprintf("bad Uuid \"%s\" - %s - for request: %v", res.GetInfo().GetUuid(), err.Error(), req)
			failures = append(failures, TestFailure{Procedure: "Create", Message: failureMsg})
		}
		if res.GetInfo().GetGender() != req.GetGender() {
			failureMsg := fmt.Sprintf("expected Gender \"%s\" but got \"%s\" for request: %v", req.GetGender(), res.GetInfo().GetGender(), req)
			failures = append(failures, TestFailure{Procedure: "Create", Message: failureMsg})
		}
		if res.GetInfo().GetEmail() != req.GetEmail() {
			failureMsg := fmt.Sprintf("expected Email \"%s\" but got \"%s\" for request: %v", req.GetEmail(), res.GetInfo().GetEmail(), req)
			failures = append(failures, TestFailure{Procedure: "Create", Message: failureMsg})
		}
		if res.GetInfo().GetName() != req.GetName() {
			failureMsg := fmt.Sprintf("expected Name \"%s\" but got \"%s\" for request: %v", req.GetName(), res.GetInfo().GetName(), req)
			failures = append(failures, TestFailure{Procedure: "Create", Message: failureMsg})
		}
		if res.GetInfo().GetBirthday() != req.GetBirthday() {
			failureMsg := fmt.Sprintf("expected Birthday \"%s\" but got \"%s\" for request: %v", req.GetBirthday(), res.GetInfo().GetBirthday(), req)
			failures = append(failures, TestFailure{Procedure: "Create", Message: failureMsg})
		}
		if !dtRegex.MatchString(res.GetInfo().GetCreatedAt()) {
			failureMsg := fmt.Sprintf("expected CreatedAt date in good format but got \"%s\" for request: %v", res.GetInfo().GetCreatedAt(), req)
			failures = append(failures, TestFailure{Procedure: "Create", Message: failureMsg})
		}
		if res.GetInfo().GetDeletedAt() != "" {
			failureMsg := fmt.Sprintf("expected DeletedAt \"%s\" must be empty for request: %v", res.GetInfo().GetDeletedAt(), req)
			failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
		}
		validProfiles = append(validProfiles, res.GetInfo())
	}

	if len(validProfiles) > 0 {
		client, ctx, err := grpcClient(config)
		if err != nil {
			failures = append(failures, TestFailure{Procedure: "Create", Message: fmt.Sprintf("gRPC client initialization error (%v)", err)})
			return failures
		}
		defer client.Close()
		for _, validProfile := range validProfiles {
			res, err := client.GetGRPCClient().HardDelete(ctx, &pb.ProfileRequest{Uuid: validProfile.GetUuid()})
			if res == nil || err != nil || res.GetOk() != true {
				failures = append(failures, TestFailure{Procedure: "Create", Message: fmt.Sprintf("deletion of created profile %s fails error (%v) - res (%v)", validProfile.GetUuid(), err, res)})
			}
		}
	}

	return failures
}
