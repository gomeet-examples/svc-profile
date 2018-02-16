package remotecli

import (
	"errors"
	"fmt"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (c *remoteCli) cmdEcho(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("Bad arguments : echo <uuid [string]> <content [string]>")
	}

	// request message
	var req *pb.EchoRequest

	// decl req for no nil panic
	req = &pb.EchoRequest{}

	// cast args[0] in req.Uuid - type TYPE_STRING to go type string
	req.Uuid = args[0]

	// cast args[1] in req.Content - type TYPE_STRING to go type string
	req.Content = args[1]

	// message validation - github.com/mwitkow/go-proto-validators
	if reqValidator, ok := interface{}(*req).(interface {
		Validate() error
	}); ok {
		if err := reqValidator.Validate(); err != nil {
			return "", err
		}
	}

	// sending message to server
	r, err := c.c.Echo(c.ctx, req)
	if err != nil {
		return "", fmt.Errorf("Echo service call fail - %v", err)
	}

	return fmt.Sprintf("Echo: %v", r), nil
}
