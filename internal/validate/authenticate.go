package validate

import (
	"github.com/getsentry/sentry-go"
)

type AuthenticateForm struct {
	Token string
}

func (v validator) ValidateAuthenticateForm(ctx Context, form AuthenticateForm) error {
	span := ctx.span.StartChild("validate-token")
	span.Status = sentry.SpanStatusOK
	if len(form.Token) == 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return &ValidationError{Field: "auth-token", Message: "cannot be empty"}
	}
	span.Finish()

	return nil
}
