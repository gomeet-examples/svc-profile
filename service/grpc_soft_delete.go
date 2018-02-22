package service

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gomeet/gomeet/utils/log"

	"github.com/gomeet-examples/svc-profile/models"
	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (s *profileServer) SoftDelete(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	log.Debug(ctx, "service call", log.Fields{"req": req})

	// init the res
	res := &pb.ProfileResponse{
		Ok:   false,
		Info: &pb.ProfileInfo{Uuid: req.GetUuid()},
	}

	// validate request
	if err := req.Validate(); err != nil {
		log.Warn(ctx, "invalid request", err, log.Fields{
			"req": req,
			"err": err,
		})

		return res, status.Error(codes.InvalidArgument, err.Error())
	}

	// init database if not ready yet
	err := s.initDatabaseHandle()
	if err != nil {
		log.Warn(ctx, "Fail to initDatabase", err, log.Fields{
			"err": err,
		})

		return res, status.Errorf(codes.Internal, err.Error())
	}

	dbProfile, err := models.DeleteProfileLogically(
		s.mysqlHandle,
		req.GetUuid(),
	)
	if err != nil {
		log.Error(ctx, "database delete error", err, log.Fields{
			"req": req,
			"err": err,
		})

		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}

	res.Ok = true
	res.Info = convertProfileFromModelToProtocol(dbProfile)

	return res, nil
}
