package remotecli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (c *remoteCli) cmdList(args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("Bad arguments : list <page_number [uint32]> <page_size [uint32]> <gender [UNKNOW|MALE|FEMALE]>")
	}

	// request message
	var req *pb.ProfileListRequest

	// decl req for no nil panic
	req = &pb.ProfileListRequest{}

	// cast args[0] in req.PageNumber - type TYPE_UINT32 to go type uint32
	reqPageNumber, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return "", fmt.Errorf("Bad arguments : page_number is not uint32")
	}
	reqPageNumberCast := uint32(reqPageNumber)
	req.PageNumber = reqPageNumberCast

	// cast args[1] in req.PageSize - type TYPE_UINT32 to go type uint32
	reqPageSize, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		return "", fmt.Errorf("Bad arguments : page_size is not uint32")
	}
	reqPageSizeCast := uint32(reqPageSize)
	req.PageSize = reqPageSizeCast

	// cast args[2] in req.Gender - type TYPE_ENUM to go type *grpc.Genders
	reqGender, ok := pb.Genders_value[strings.ToUpper(args[2])]
	if !ok {
		return "", fmt.Errorf("Bad arguments : unknown gender \"%s\"", args[2])
	}
	req.Gender = pb.Genders(reqGender)

	// message validation - github.com/mwitkow/go-proto-validators
	if reqValidator, ok := interface{}(*req).(interface {
		Validate() error
	}); ok {
		if err := reqValidator.Validate(); err != nil {
			return "", err
		}
	}

	// sending message to server
	r, err := c.c.List(c.ctx, req)
	if err != nil {
		return "", fmt.Errorf("List service call fail - %v", err)
	}

	return fmt.Sprintf("List: %v", r), nil
}
