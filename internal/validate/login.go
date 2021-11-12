package validate

import (
	"context"
	"errors"
)

type LoginForm struct {
	Email           string
	Password        string
	DeviceUserAgent string
	IP              string
}

func (v validator) ValidateLoginForm(ctx context.Context, form LoginForm) error {
	if len(form.Password) == 0 {
		return &ValidationError{Field: "password", Message: "cannot be empty"}
	}

	if len(form.Email) == 0 {
		return &ValidationError{Field: "email", Message: "cannot be empty"}
	}

	if len(form.IP) == 0 {
		return &ValidationError{Field: "ip", Message: "cannot be empty"}
	}

	if isValid, err := isIPValid(form.IP); nil != err {
		return errors.New("failed to validate email field")
	} else if !isValid {
		return &ValidationError{Field: "ip", Message: "invalid"}
	}

	if len(form.DeviceUserAgent) == 0 {
		return &ValidationError{Field: "device_user_agent", Message: "cannot be empty"}
	}

	if len(form.Password) < 8 {
		return &ValidationError{Field: "password", Message: "password must have at least 8 characters"}
	}

	if isValid, err := isEmailValid(form.Email); nil != err {
		return errors.New("failed to validate email field")
	} else if !isValid {
		return &ValidationError{Field: "email", Message: "invalid email"}
	}

	return nil
}
