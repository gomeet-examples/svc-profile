// Code generated by protoc-gen-gomeet-service. DO NOT EDIT.
// source: pb/profile.proto
package functest

import (
	"fmt"
)

func TestCreate(config FunctionalTestConfig) (failures []TestFailure) {
	client, ctx, err := grpcClient(config)
	if err != nil {
		failures = append(failures, TestFailure{Procedure: "Create", Message: fmt.Sprintf("gRPC client initialization error (%v)", err)})
		return failures
	}
	defer client.Close()

	var testCaseResults []*TestCaseResult
	for _, req := range testGetCreateRequest() {
		res, err := client.GetGRPCClient().Create(ctx, req)
		testCaseResults = append(testCaseResults, &TestCaseResult{req, res, err})
	}

	return testCreateResponse(FUNCTEST_GRPC, testCaseResults)
}
