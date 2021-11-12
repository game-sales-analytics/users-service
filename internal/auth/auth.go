package auth

import (
	"context"
	"time"
)

type Auth interface {
	VerifyToken(ctx context.Context, token string) (*TokenVerificationResult, error)
	LoginWithEmail(ctx context.Context, creds LoginWithEmailCreds) (*LoginResult, error)
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
