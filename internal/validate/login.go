package validate

import (
	"errors"

	"github.com/getsentry/sentry-go"
)

type LoginForm struct {
	Email           string
	Password        string
	DeviceUserAgent string
	IP              string
}

func (v validator) ValidateLoginForm(ctx Context, form LoginForm) error {
	span := ctx.span.StartChild("validate-password")
	span.Status = sentry.SpanStatusOK
	if len(form.Password) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return &ValidationError{Field: "password", Message: "cannot be empty"}
	}
	if len(form.Password) < 8 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return &ValidationError{Field: "password", Message: "password must have at least 8 characters"}
	}
	span.Finish()

	span = ctx.span.StartChild("validate-email")
	span.Status = sentry.SpanStatusOK
	if len(form.Email) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return &ValidationError{Field: "email", Message: "cannot be empty"}
	}
	if isValid, err := isEmailValid(form.Email); nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return errors.New("failed to validate email field")
	} else if !isValid {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return &ValidationError{Field: "email", Message: "invalid email"}
	}
	span.Finish()

	span = ctx.span.StartChild("validate-ip")
	span.Status = sentry.SpanStatusOK
	if len(form.IP) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return &ValidationError{Field: "ip", Message: "cannot be empty"}
	}
	if isValid, err := isIPValid(form.IP); nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return errors.New("failed to validate email field")
	} else if !isValid {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return &ValidationError{Field: "ip", Message: "invalid"}
	}
	span.Finish()

	span = ctx.span.StartChild("validate-device-user-agent")
	span.Status = sentry.SpanStatusOK
	if len(form.DeviceUserAgent) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return &ValidationError{Field: "device_user_agent", Message: "cannot be empty"}
	}
	span.Finish()

	return nil
}
