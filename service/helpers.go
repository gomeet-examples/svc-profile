package service

import (
	"time"

	"github.com/gomeet-examples/svc-profile/models"
	pb "github.com/gomeet-examples/svc-profile/pb"
)

const (
	defaultPageSize = 50
)

func convertProfileFromModelToProtocol(dbProfile *models.Profile) *pb.ProfileInfo {
	var (
		gender    pb.Genders
		deletedAt string
	)

	if dbProfile.DeletedAt != nil {
		deletedAt = dbProfile.DeletedAt.UTC().Format(time.RFC3339)
	}

	switch dbProfile.Gender {
	case 1:
		gender = pb.Genders_MALE
	case 2:
		gender = pb.Genders_FEMALE
	default:
		gender = pb.Genders_UNKNOW
	}

	return &pb.ProfileInfo{
		Uuid:      dbProfile.UUID.String(),
		Gender:    gender,
		Email:     dbProfile.Email,
		Name:      dbProfile.Name,
		Birthday:  dbProfile.Birthday.UTC().Format("2006-01-02"),
		CreatedAt: dbProfile.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: dbProfile.UpdatedAt.UTC().Format(time.RFC3339),
		DeletedAt: deletedAt,
	}
}
