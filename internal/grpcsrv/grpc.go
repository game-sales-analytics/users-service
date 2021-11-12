package grpcsrv

import (
	"github.com/sirupsen/logrus"

	"github.com/game-sales-analytics/users-service/internal/pb"
)

func New(logger *logrus.Logger) GrpcService {
	return server{
		pb.UnimplementedUsersServiceServer{},
		logger,
	}
}
