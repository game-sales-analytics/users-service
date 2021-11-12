package auth

import (
	"context"
)

func (a authsrv) VerifyToken(ctx context.Context, token string) (*TokenVerificationResult, error) {
	verificationResult, err := verifyToken(token, a.cfg.Secret)
	if nil != err {
		return nil, err
	}

	if exists, err := a.repo.UserWithIDExists(ctx, verificationResult.UserID); nil != err {
		return nil, err
	} else if !exists {
		return nil, ErrUserNotExists
	}

	return verificationResult, nil
}
