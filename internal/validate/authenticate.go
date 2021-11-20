package validate

import (
	"context"
)

type AuthenticateForm struct {
	Token string
}

func (v validator) ValidateAuthenticateForm(ctx context.Context, form AuthenticateForm) error {
	if len(form.Token) == 0 {
		return &ValidationError{Field: "auth-token", Message: "cannot be empty"}
	}

	return nil
}
