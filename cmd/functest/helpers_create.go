package functest

import (
	pb "github.com/gomeet-examples/svc-profile/pb"
)

func testGetCreateRequest() (reqs []*pb.ProfileCreationRequest) {
	// return an array of pb.ProfileCreationRequest struct pointers,
	// each of them will be passed as an argument to the grpc Create method

	reqs = append(reqs, &pb.ProfileCreationRequest{})
	return reqs
}

func testCreateResponse(
	testsType string,
	testCaseResults []*TestCaseResult,
) (failures []TestFailure) {
	// Do something useful functional test with
	// testCaseResults[n].Request, testCaseResults[n].Response and testCaseResults[n].Error
	// then return a array of TestFailure struct
	// testsType value is value of FUNCTEST_HTTP (HTTP) and FUNCTEST_GRPC (GRPC) constants cf. types.go
	for _, tr := range testCaseResults {
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

		if tr.Response != nil {
			res, ok = tr.Response.(*pb.ProfileResponse)
			if !ok {
				failures = append(failures, TestFailure{Procedure: "Create", Message: "expected response message type pb.ProfileCreationRequest - cast fail"})
				continue
			}
		}

		// Do something useful functional test with req, res and err
		err = tr.Error
		if err != nil {
			// if no error are expected do something like this
			// failures = append(failures, TestFailure{Procedure: "Create", Message: "no error expected"})
			// continue
		}

		if req != nil && res != nil {
			// for example :
			// if res.GetId() != req.GetId() {
			//     failureMsg := fmt.Sprintf("expected ID \"%s\" but got \"%s\" for request: %v", req.GetId(), res.GetId(), req)
			//     failures = append(failures, TestFailure{Procedure: "Create", Message: failureMsg})
			// }
		}
	}

	return failures
}
