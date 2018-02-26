package remotecli

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (c *remoteCli) cmdCreate(args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("Bad arguments : create <gender [UNKNOW|MALE|FEMALE]> <email [string]> <name [string]> <birthday [string]>")
	}

	// request message
	var req *pb.ProfileCreationRequest

	// decl req for no nil panic
	req = &pb.ProfileCreationRequest{}

	// cast args[0] in req.Gender - type TYPE_ENUM to go type *grpc.Genders
	reqGender, ok := pb.Genders_value[strings.ToUpper(args[0])]
	if !ok {
		return "", fmt.Errorf("Bad arguments : unknown gender \"%s\"", args[0])
	}
	req.Gender = pb.Genders(reqGender)

	// cast args[1] in req.Email - type TYPE_STRING to go type string
	req.Email = args[1]

	// cast args[2] in req.Name - type TYPE_STRING to go type string
	req.Name = args[2]

	// cast args[3] in req.Birthday - type TYPE_STRING to go type string
	req.Birthday = args[3]

	// message validation - github.com/mwitkow/go-proto-validators
	if reqValidator, ok := interface{}(*req).(interface {
		Validate() error
	}); ok {
		if err := reqValidator.Validate(); err != nil {
			return "", err
		}
	}

	// sending message to server
	r, err := c.c.Create(c.ctx, req)
	if err != nil {
		return "", fmt.Errorf("Create service call fail - %v", err)
	}

	return fmt.Sprintf("Create: %v", r), nil
}
