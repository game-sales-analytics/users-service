package auth

import (
	"errors"

	"github.com/getsentry/sentry-go"

	"github.com/game-sales-analytics/users-service/internal/apm"
	"github.com/game-sales-analytics/users-service/internal/db/repository"
)

func (a authsrv) VerifyToken(ctx Context, token string) (*TokenVerificationResult, error) {
	span := ctx.span.StartChild("verify-raw-token")
	span.Status = sentry.SpanStatusOK
	decodeRes, err := verifyToken(NewContext(ctx, span), token, a.cfg.Secret)
	if nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, ErrTokenNotVerified
	}
	span.Finish()

	span = ctx.span.StartChild("get-user-authentication-info")
	span.Status = sentry.SpanStatusOK
	userAuthInfo, err := a.repo.GetUserAuthenticationInfo(repository.NewDBOperationContext(ctx, span), decodeRes.userID)
	if nil != err {
		defer span.Finish()

		if errors.Is(err, repository.ErrUserNotExists) {
			span.Status = sentry.SpanStatusNotFound
			return nil, ErrUserNotExists
		}

		span.Status = sentry.SpanStatusNotFound
		log := a.logger.WithError(err).WithField("err_code", "E_GET_USER_AUTHENTICATION_INFO")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed retrieving user authentication info")
		return nil, ErrInternal
	}
	span.Finish()

	out := TokenVerificationResult{
		User: TokenVerificationResultUser{
			ID:        decodeRes.userID,
			FirstName: userAuthInfo.FirstName,
			LastName:  userAuthInfo.LastName,
		},
	}

	return &out, nil
}
