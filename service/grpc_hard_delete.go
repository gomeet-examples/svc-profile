package service

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gomeet/gomeet/utils/log"

	"github.com/gomeet-examples/svc-profile/models"
	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (s *profileServer) HardDelete(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponseLight, error) {
	log.Debug(ctx, "service call", log.Fields{"req": req})

	// init the res
	res := &pb.ProfileResponseLight{
		Ok: false,
	}

	// validate request
	if err := req.Validate(); err != nil {
		log.Warn(ctx, "invalid request", err, log.Fields{
			"req": req,
		})

		return res, status.Error(codes.InvalidArgument, err.Error())
	}

	// init database if not ready yet
	err := s.initDatabaseHandle()
	if err != nil {
		log.Warn(ctx, "Fail to initDatabase", err, log.Fields{})

		return res, status.Errorf(codes.Internal, err.Error())
	}

	err = models.DeleteProfilePhysically(
		s.mysqlHandle,
		req.GetUuid(),
	)
	if err != nil {
		log.Error(ctx, "database delete error", err, log.Fields{
			"req": req,
		})

		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}

	res.Ok = true

	return res, nil
}
