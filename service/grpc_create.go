package service

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gomeet/gomeet/utils/log"

	"github.com/gomeet-examples/svc-profile/models"
	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (s *profileServer) Create(ctx context.Context, req *pb.ProfileCreationRequest) (*pb.ProfileResponse, error) {
	log.Debug(ctx, "service call", log.Fields{"req": req})

	// init the response
	res := &pb.ProfileResponse{
		Ok:   false,
		Info: &pb.ProfileInfo{},
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

	// create profile in database
	dbProfile, err := models.CreateProfile(
		s.mysqlHandle,
		uint16(req.GetGender()),
		req.GetEmail(),
		req.GetName(),
		req.GetBirthday(),
	)
	if err != nil {
		log.Error(ctx, "database insert error", err, log.Fields{
			"req": req,
		})

		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// cast profile from model to protocol
	res.Ok = true
	res.Info = convertProfileFromModelToProtocol(dbProfile)

	return res, nil
}
