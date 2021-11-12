package auth

import (
	"context"
)

type Auth interface {
	VerifyToken(ctx context.Context, token string) (*TokenVerificationResult, error)
	LoginWithEmail(ctx context.Context, creds LoginWithEmailCreds) (*LoginResult, error)
}

type TokenVerificationResult struct {
	UserID string
}

type LoginResult struct {
	Token string
}
