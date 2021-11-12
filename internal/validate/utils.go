package validate

import (
	"net/mail"
)

func isEmailValid(email string) (bool, error) {
	_, err := mail.ParseAddress(email)

	return err == nil, nil
}
