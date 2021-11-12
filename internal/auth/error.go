package auth

import (
	"errors"
)

var (
	ErrTokenNotVerified = errors.New("token is not valid")
	ErrUnauthenticated  = errors.New("invalid credentials provided")
	ErrInternal         = errors.New("internal error occurred")
	ErrUserNotExists    = errors.New("no user with associated token exists")
)
