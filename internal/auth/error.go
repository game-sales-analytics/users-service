package auth

import (
	"errors"
)

var (
	ErrTokenNotVerified = errors.New("token is not valid")
	ErrUserNotExists    = errors.New("no user with associated token exists")
)
