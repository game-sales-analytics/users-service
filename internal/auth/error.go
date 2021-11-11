package auth

import (
	"errors"
)

var (
	ErrTokenNotVerified = errors.New("token is not valid")
)
