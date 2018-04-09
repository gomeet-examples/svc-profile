package functest

import (
	"fmt"

	pb "github.com/gomeet-examples/svc-profile/pb"
	"github.com/google/uuid"
)

func testGetHardDeleteRequest(
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
	reqs = append(reqs, &pb.ProfileRequest{Uuid: ""})             // empty Uuid
	reqs = append(reqs, &pb.ProfileRequest{Uuid: "invalid uuid"}) // invalid Uuid

	// valid index 2+
	reqs = append(reqs, &pb.ProfileRequest{Uuid: uuid.New().String()})
	reqs = append(reqs, &pb.ProfileRequest{Uuid: validProfile.GetUuid()})
	return reqs, extras, err
}

func testHardDeleteResponse(
	config FunctionalTestConfig,
	testsType string,
	testCaseResults []*TestCaseResult,
	extras map[string]interface{},
) (failures []TestFailure) {
	client, ctx, err := grpcClient(config)
	if err != nil {
		failures = append(failures, TestFailure{Procedure: "HardDelete", Message: fmt.Sprintf("gRPC client initialization error (%v)", err)})
		return failures
	}
	defer client.Close()

	for i, tr := range testCaseResults {
		var (
			req *pb.ProfileRequest
			res *pb.ProfileResponseLight
			err error
			ok  bool
		)

		if tr.Request == nil {
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: "expected request message type pb.ProfileRequest - nil given"})
			continue
		}
		req, ok = tr.Request.(*pb.ProfileRequest)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: "expected request message type pb.ProfileRequest - cast fail"})
			continue
		}

		err = tr.Error
		if i < 2 {
			if err == nil {
				failures = append(failures, TestFailure{Procedure: "HardDelete", Message: "an error is expected"})
			}
			continue
		}

		if err != nil {
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: "no error expected"})
			continue
		}

		if tr.Response == nil {
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: "a response is expected"})
			continue
		}
		res, ok = tr.Response.(*pb.ProfileResponseLight)
		if !ok {
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: "expected response message type pb.ProfileInfo - cast fail"})
			continue
		}

		if req == nil || res == nil {
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: "a request and a response are expected"})
			continue
		}

		if !res.GetOk() {
			failureMsg := fmt.Sprintf("expected response Ok \"true\" but got \"false\" for request: %v", req)
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: failureMsg})
			continue
		}

		// When reading the profile the not found response is expected
		resRead, err := client.GetGRPCClient().Read(ctx, &pb.ProfileRequest{Uuid: req.GetUuid()})
		if resRead != nil {
			failureMsg := fmt.Sprintf("expected nil response on read but got \"%v\" for request: %v", resRead, req)
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: failureMsg})
			continue
		}
		if err == nil {
			failureMsg := fmt.Sprintf("expected error on read but got \"nil\" for request: %v", req)
			failures = append(failures, TestFailure{Procedure: "HardDelete", Message: failureMsg})
			continue
		}
	}

	if len(extras) > 0 {
		for sUuid := range extras {
			res, err := client.GetGRPCClient().HardDelete(ctx, &pb.ProfileRequest{Uuid: sUuid})
			if res == nil || err != nil || res.GetOk() != true {
				failures = append(failures, TestFailure{Procedure: "HardDelete", Message: fmt.Sprintf("deletion of created profile %s fails error (%v) - res (%v)", sUuid, err, res)})
			}
		}
	}

	return failures
}
