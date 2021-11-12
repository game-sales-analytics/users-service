package validate

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/game-sales-analytics/users-service/internal/db/repository"
)

type NormalizedForm struct {
	Email string
}

type Validator interface {
	ValidateRegisterForm(ctx context.Context, form RegisterForm) (*NormalizedForm, error)
	ValidateLoginForm(ctx context.Context, form LoginForm) error
	ValidateAuthenticateForm(ctx context.Context, form AuthenticateForm) error
}

type validator struct {
	logger *logrus.Logger
	repo   *repository.Repo
}

func New(logger *logrus.Logger, repo *repository.Repo) Validator {
	return validator{
		logger,
		repo,
	}
}
