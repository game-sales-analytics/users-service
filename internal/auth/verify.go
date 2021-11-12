package auth

import (
	"context"
	"errors"

	"github.com/game-sales-analytics/users-service/internal/db/repository"
)

func (a authsrv) VerifyToken(ctx context.Context, token string) (*TokenVerificationResult, error) {
	decodeRes, err := verifyToken(token, a.cfg.Secret)
	if nil != err {
		return nil, ErrTokenNotVerified
	}

	userAuthInfo, err := a.repo.GetUserAuthenticationInfo(ctx, decodeRes.userID)
	if nil != err {
		if errors.Is(err, repository.ErrUserNotExists) {
			return nil, ErrUserNotExists
		}
		return nil, ErrInternal
	}

	out := TokenVerificationResult{
		User: TokenVerificationResultUser{
			ID:        decodeRes.userID,
			FirstName: userAuthInfo.FirstName,
			LastName:  userAuthInfo.LastName,
		},
	}

	return &out, nil
}
