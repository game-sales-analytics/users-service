package validate

import (
	"github.com/sirupsen/logrus"

	"github.com/game-sales-analytics/users-service/internal/db/repository"
)

type NormalizedForm struct {
	Email string
}

type Validator interface {
	ValidateRegisterForm(ctx Context, form RegisterForm) (*NormalizedForm, error)
	ValidateLoginForm(ctx Context, form LoginForm) error
	ValidateAuthenticateForm(ctx Context, form AuthenticateForm) error
}

type validator struct {
	logger *logrus.Entry
	repo   *repository.Repo
}

func New(logger *logrus.Entry, repo *repository.Repo) Validator {
	return validator{
		logger,
		repo,
	}
}
