package service

import (
	"golang.org/x/net/context"

	"github.com/gomeet/gomeet/utils/log"

	pb "github.com/gomeet-examples/svc-profile/pb"
)

func (s *profileServer) Update(ctx context.Context, req *pb.ProfileInfo) (*pb.ProfileResponse, error) {
	log.Debug(ctx, "service call", log.Fields{"req": req})

	// res := &pb.ProfileResponse{}
	// Do something useful with req and res
	// for now a fake response is returned see https://github.com/gomeet/go-proto-gomeetfaker
	res := pb.NewProfileResponseGomeetFaker()

	return res, nil
}