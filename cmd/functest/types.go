package functest

import (
	"github.com/golang/protobuf/proto"
)

const FUNCTEST_HTTP = "HTTP"
const FUNCTEST_GRPC = "GRPC"

// FunctionalTestConfig encodes client configuration for functional tests.
type FunctionalTestConfig struct {
	ServerAddress       string
	GrpcServerAddress   string
	HttpServerAddress   string
	CaCertificate       string
	ClientCertificate   string
	ClientPrivateKey    string
	MySQLDataSourceName string
	JsonWebToken        string
	TimeoutSeconds      int
}

// TestFailure encodes a functional test failure description.
type TestFailure struct {
	Message   string
	Procedure string
}

// TestCaseResult encodes a functional test case result description.
type TestCaseResult struct {
	Request  proto.Message
	Response proto.Message
	Error    error
}
