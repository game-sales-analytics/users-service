package auth

import (
	"time"
)

type Auth interface {
	VerifyToken(ctx Context, token string) (*TokenVerificationResult, error)
	LoginWithEmail(ctx Context, creds LoginWithEmailCreds) (*LoginResult, error)
}

type TokenVerificationResultUser struct {
	ID        string
	FirstName string
	LastName  string
}

type TokenVerificationResult struct {
	User TokenVerificationResultUser
}

type LoginResultToken struct {
	ID                 string
	Value              string
	NotBeforeDateTime  time.Time
	ExpirationDateTime time.Time
}

type LoginResult struct {
	Token LoginResultToken
}
