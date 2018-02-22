package service

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gomeet/gomeet/utils/log"

	"github.com/gomeet-examples/svc-profile/models"
	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (s *profileServer) List(ctx context.Context, req *pb.ProfileListRequest) (*pb.ProfileList, error) {
	log.Debug(ctx, "service call", log.Fields{"req": req})

	// init the response
	res := &pb.ProfileList{}

	// validate request
	if err := req.Validate(); err != nil {
		log.Warn(ctx, "invalid request", err, log.Fields{
			"req": req,
			"err": err,
		})

		return res, status.Error(codes.InvalidArgument, err.Error())
	}

	// page size result set
	pageSize := uint(defaultPageSize)
	if req.GetPageSize() > 0 {
		pageSize = uint(req.GetPageSize())
	}

	// page number result set
	pageNumber := uint(1)
	if req.GetPageNumber() > 1 {
		pageNumber = uint(req.GetPageNumber())
	}

	order := "created_at asc"
	if req.GetOrder() != "" {
		order = req.GetOrder() // FIXME: validation
	}

	criteria := make(map[string]interface{})

	gender := req.GetGender()
	switch gender {
	case pb.Genders_MALE:
		criteria["gender"] = 1
	case pb.Genders_FEMALE:
		criteria["gender"] = 2
	}

	// init database if not ready yet
	err := s.initDatabaseHandle()
	if err != nil {
		log.Warn(ctx, "Fail to initDatabase", err, log.Fields{
			"err": err,
		})

		return res, status.Errorf(codes.Internal, err.Error())
	}

	dbProfileList, resultSetSize, hasMore, err := models.ListProfiles(
		s.mysqlHandle,
		(pageNumber-1)*pageSize,
		pageSize,
		order,
		criteria,
		req.GetExcludeSoftDeleted(),
		req.GetSoftDeletedOnly(),
	)
	if err != nil {
		log.Error(ctx, "database select error", err, log.Fields{
			"req": req,
			"err": err,
		})

		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}

	res.ResultSetSize = uint32(resultSetSize)
	res.HasMore = hasMore
	for _, dbProfile := range dbProfileList {
		res.Profiles = append(res.Profiles, convertProfileFromModelToProtocol(&dbProfile))
	}

	return res, nil
}
