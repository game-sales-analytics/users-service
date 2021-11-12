package validate

import (
	"context"
	"errors"

	"github.com/game-sales-analytics/users-service/internal/normalize"
)

type RegisterForm struct {
	Email                string
	Password             string
	PasswordConfirmation string
	FirstName            string
	LastName             string
}

func (v validator) ValidateRegisterForm(ctx context.Context, form RegisterForm) (*NormalizedForm, error) {
	if len(form.Password) == 0 {
		return nil, &ValidationError{Field: "password", Message: "cannot be empty"}
	}

	if len(form.Email) == 0 {
		return nil, &ValidationError{Field: "email", Message: "cannot be empty"}
	}

	if len(form.FirstName) == 0 {
		return nil, &ValidationError{Field: "first_name", Message: "cannot be empty"}
	}

	if len(form.LastName) == 0 {
		return nil, &ValidationError{Field: "last_name", Message: "cannot be empty"}
	}

	if len(form.Password) < 8 {
		return nil, &ValidationError{Field: "password", Message: "password must have at least 8 characters"}
	}

	if form.Password != form.PasswordConfirmation {
		return nil, &ValidationError{
			Field:   "password_confirmation",
			Message: "does not match password",
		}
	}

	if isValid, err := isEmailValid(form.Email); nil != err {
		return nil, errors.New("failed to validate email field")
	} else if !isValid {
		return nil, &ValidationError{Field: "email", Message: "invalid email"}
	}

	normalizedEmail, err := normalize.Email(form.Email)
	if nil != err {
		return nil, errors.New("failed normalizing email address")
	}

	if exists, err := v.repo.NormalizedEmailExists(ctx, normalizedEmail); nil != err {
		v.logger.WithError(err).WithField("err_code", "E_CHECK_NORMALIZED_EMAIL_EXISTENCE").Error("failed checking normalized email existence")
		return nil, err
	} else if exists {
		return nil, &ValidationError{Field: "email", Message: "duplicate email address"}
	}

	return &NormalizedForm{
		Email: normalizedEmail,
	}, nil
}
