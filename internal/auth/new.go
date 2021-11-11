package auth

import (
	"github.com/sirupsen/logrus"

	"github.com/game-sales-analytics/users-service/internal/config"
	"github.com/game-sales-analytics/users-service/internal/db/repository"
)

type authsrv struct {
	repo   *repository.Repo
	logger *logrus.Logger
	cfg    *config.JwtConfig
}

func New(
	repo *repository.Repo,
	logger *logrus.Logger,
	cfg *config.JwtConfig,
) Auth {
	return authsrv{
		repo,
		logger,
		cfg,
	}
}
