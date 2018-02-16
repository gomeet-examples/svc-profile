package functest

import (
	pb "github.com/gomeet-examples/svc-profile/pb"
)

func testGetUpdateRequest() (reqs []*pb.ProfileInfo) {
	// return an array of pb.ProfileInfo struct pointers,
	// each of them will be passed as an argument to the grpc Update method

	reqs = append(reqs, &pb.ProfileInfo{})
	return reqs
}

func testUpdateResponse(
	testsType string,
	testCaseResults []*TestCaseResult,
) (failures []TestFailure) {
	// Do something useful functional test with
	// testCaseResults[n].Request, testCaseResults[n].Response and testCaseResults[n].Error
	// then return a array of TestFailure struct
	// testsType value is value of FUNCTEST_HTTP (HTTP) and FUNCTEST_GRPC (GRPC) constants cf. types.go
	for _, tr := range testCaseResults {
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

		if tr.Response != nil {
			res, ok = tr.Response.(*pb.ProfileResponse)
			if !ok {
				failures = append(failures, TestFailure{Procedure: "Update", Message: "expected response message type pb.ProfileInfo - cast fail"})
				continue
			}
		}

		// Do something useful functional test with req, res and err
		err = tr.Error
		if err != nil {
			// if no error are expected do something like this
			// failures = append(failures, TestFailure{Procedure: "Update", Message: "no error expected"})
			// continue
		}

		if req != nil && res != nil {
			// for example :
			// if res.GetId() != req.GetId() {
			//     failureMsg := fmt.Sprintf("expected ID \"%s\" but got \"%s\" for request: %v", req.GetId(), res.GetId(), req)
			//     failures = append(failures, TestFailure{Procedure: "Update", Message: failureMsg})
			// }
		}
	}

	return failures
}
