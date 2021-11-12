package grpcsrv

import (
	"github.com/sirupsen/logrus"

	"github.com/game-sales-analytics/users-service/internal/auth"
	"github.com/game-sales-analytics/users-service/internal/db/repository"
	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func New(
	logger *logrus.Logger,
	repo *repository.Repo,
	validator validate.Validator,
	auth auth.Auth,
) GrpcService {
	return server{
		pb.UnimplementedUsersServiceServer{},
		logger,
		repo,
		validator,
		auth,
	}
}
