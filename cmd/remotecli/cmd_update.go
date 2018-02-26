package remotecli

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (c *remoteCli) cmdUpdate(args []string) (string, error) {
	if len(args) < 5 {
		return "", errors.New("Bad arguments : update <uuid [string]> <gender [UNKNOW|MALE|FEMALE]> <email [string]> <name [string]> <birthday [string]>")
	}

	// request message
	var req *pb.ProfileInfo

	// decl req for no nil panic
	req = &pb.ProfileInfo{}

	// cast args[0] in req.Uuid - type TYPE_STRING to go type string
	req.Uuid = args[0]

	// cast args[1] in req.Gender - type TYPE_ENUM to go type *grpc.Genders
	reqGender, ok := pb.Genders_value[strings.ToUpper(args[1])]
	if !ok {
		return "", fmt.Errorf("Bad arguments : unknown gender \"%s\"", args[1])
	}
	req.Gender = pb.Genders(reqGender)

	// cast args[2] in req.Email - type TYPE_STRING to go type string
	req.Email = args[2]

	// cast args[3] in req.Name - type TYPE_STRING to go type string
	req.Name = args[3]

	// cast args[4] in req.Birthday - type TYPE_STRING to go type string
	req.Birthday = args[4]

	// message validation - github.com/mwitkow/go-proto-validators
	if reqValidator, ok := interface{}(*req).(interface {
		Validate() error
	}); ok {
		if err := reqValidator.Validate(); err != nil {
			return "", err
		}
	}

	// sending message to server
	r, err := c.c.Update(c.ctx, req)
	if err != nil {
		return "", fmt.Errorf("Update service call fail - %v", err)
	}

	return fmt.Sprintf("Update: %v", r), nil
}
