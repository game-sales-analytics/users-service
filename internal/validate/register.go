package validate

import (
	"errors"

	"github.com/getsentry/sentry-go"

	"github.com/game-sales-analytics/users-service/internal/apm"
	"github.com/game-sales-analytics/users-service/internal/db/repository"
	"github.com/game-sales-analytics/users-service/internal/normalize"
)

type RegisterForm struct {
	Email                string
	Password             string
	PasswordConfirmation string
	FirstName            string
	LastName             string
}

func (v validator) ValidateRegisterForm(ctx Context, form RegisterForm) (*NormalizedForm, error) {
	span := ctx.span.StartChild("validate-password")
	span.Status = sentry.SpanStatusOK
	if len(form.Password) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, &ValidationError{Field: "password", Message: "cannot be empty"}
	}
	if len(form.Password) < 8 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, &ValidationError{Field: "password", Message: "password must have at least 8 characters"}
	}
	span.Finish()

	span = ctx.span.StartChild("validate-password-confirmation")
	span.Status = sentry.SpanStatusOK
	if form.Password != form.PasswordConfirmation {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, &ValidationError{
			Field:   "password_confirmation",
			Message: "does not match password",
		}
	}
	span.Finish()

	span = ctx.span.StartChild("validate-email")
	span.Status = sentry.SpanStatusOK
	if len(form.Email) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, &ValidationError{Field: "email", Message: "cannot be empty"}
	}
	if isValid, err := isEmailValid(form.Email); nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := v.logger.WithError(err).WithField("err_code", "E_VALIDATE_EMAIL_FORMAT")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed validating email format")
		return nil, errors.New("failed to validate email field")
	} else if !isValid {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, &ValidationError{Field: "email", Message: "invalid email"}
	}
	normalizedEmail, err := normalize.Email(form.Email)
	if nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := v.logger.WithError(err).WithField("err_code", "E_NORMALIZE_EMAIL")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed normalizing email address")
		return nil, errors.New("failed normalizing email address")
	}
	if exists, err := v.repo.NormalizedEmailExists(repository.NewDBOperationContext(ctx, span), normalizedEmail); nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := v.logger.WithError(err).WithField("err_code", "E_CHECK_NORMALIZED_EMAIL_EXISTENCE")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed checking normalized email existence")
		return nil, err
	} else if exists {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, &ValidationError{Field: "email", Message: "duplicate email address"}
	}
	span.Finish()

	span = ctx.span.StartChild("validate-first-name")
	span.Status = sentry.SpanStatusOK
	if len(form.FirstName) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, &ValidationError{Field: "first_name", Message: "cannot be empty"}
	}
	span.Finish()

	span = ctx.span.StartChild("validate-last-name")
	span.Status = sentry.SpanStatusOK
	if len(form.LastName) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, &ValidationError{Field: "last_name", Message: "cannot be empty"}
	}
	span.Finish()

	return &NormalizedForm{
		Email: normalizedEmail,
	}, nil
}
