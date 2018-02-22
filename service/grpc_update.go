package service

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gomeet/gomeet/utils/log"

	"github.com/gomeet-examples/svc-profile/models"
	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (s *profileServer) Update(ctx context.Context, req *pb.ProfileInfo) (*pb.ProfileResponse, error) {
	log.Debug(ctx, "service call", log.Fields{"req": req})

	// init the response
	res := &pb.ProfileResponse{
		Ok:   false,
		Info: req,
	}

	// validate request
	if err := req.Validate(); err != nil {
		log.Warn(ctx, "invalid request", err, log.Fields{
			"req": req,
			"err": err,
		})

		return res, status.Error(codes.InvalidArgument, err.Error())
	}

	// set uuid
	uuid := req.GetUuid()
	if uuid == "" {
		err := errors.New("Missing Uuid")
		log.Warn(ctx, "invalid request - missing Uuid", err, log.Fields{
			"req": req,
			"err": err,
		})

		return res, status.Error(codes.InvalidArgument, "Uuid is required")
	}

	// init database if not ready yet
	err := s.initDatabaseHandle()
	if err != nil {
		log.Warn(ctx, "Fail to initDatabase", err, log.Fields{
			"err": err,
		})
		return res, status.Errorf(codes.Internal, err.Error())
	}

	dbProfile, err := models.UpdateProfile(
		s.mysqlHandle,
		uuid,
		uint16(req.GetGender()),
		req.GetEmail(),
		req.GetName(),
		req.GetBirthday(),
	)
	if err != nil {
		log.Error(ctx, "database update error", err, log.Fields{
			"req": req,
			"err": err,
		})

		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}

	res.Ok = true
	res.Info = convertProfileFromModelToProtocol(dbProfile)

	return res, nil
}
