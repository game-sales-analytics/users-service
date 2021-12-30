package auth

import (
	"errors"
	"math/rand"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/game-sales-analytics/users-service/internal/apm"
	"github.com/game-sales-analytics/users-service/internal/db/repository"
	"github.com/game-sales-analytics/users-service/internal/id"
	"github.com/game-sales-analytics/users-service/internal/passhash"
)

type LoginDefaultCreds struct {
	UserIPAddress       string
	UserDeviceUserAgent string
}

type LoginWithEmailCreds struct {
	LoginDefaultCreds
	Email    string
	Password string
}

func (a authsrv) LoginWithEmail(ctx Context, creds LoginWithEmailCreds) (*LoginResult, error) {
	span := ctx.span.StartChild("get-user-login-info")
	span.Status = sentry.SpanStatusOK
	user, err := a.repo.GetUserLoginInfo(repository.NewDBOperationContext(ctx, span), creds.Email)
	if nil != err {
		defer span.Finish()

		if errors.Is(err, repository.ErrUserNotExists) {
			span.Status = sentry.SpanStatusUnauthenticated

			randSpan := rand.Int63n(251)
			time.Sleep(time.Millisecond * time.Duration((6136 + randSpan)))
			return nil, ErrUnauthenticated
		}

		span.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_RETRIEVE_USER_LOGIN_INFO")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed retrieving user login information")
		return nil, ErrInternal
	}
	span.Finish()

	span = ctx.span.StartChild("verify-user-password")
	span.Status = sentry.SpanStatusOK
	matched, err := passhash.Verify(creds.Password, user.Password)
	if nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_VERIFY_PASSWORD")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed verifying user password for login")
		return nil, ErrInternal
	}
	if !matched {
		defer span.Finish()

		span.Status = sentry.SpanStatusUnauthenticated
		return nil, ErrUnauthenticated
	}
	span.Finish()

	span = ctx.span.StartChild("generate-auth-token")
	span.Status = sentry.SpanStatusOK
	token, err := a.generateToken(NewContext(ctx, span), user.ID, a.cfg.Secret)
	if nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_GENERATE_AUTH_TOKEN")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed generating authentication token for user")
		return nil, ErrInternal
	}
	span.Finish()

	span = ctx.span.StartChild("generate-user-login-attempt")
	span.Status = sentry.SpanStatusOK
	loginRecordID, err := id.GenerateUserLoginID()
	if nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_GENERATE_NEW_LOGIN_RECORD_ID")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed generating id for new user login record")
		return nil, ErrInternal
	}
	span.Finish()

	loginRecord := repository.NewUserLoginToSave{
		ID:                  loginRecordID,
		UserID:              user.ID,
		LoggedInAt:          time.Now(),
		UserIPAddress:       creds.UserIPAddress,
		UserDeviceUserAgent: creds.UserDeviceUserAgent,
	}

	span = ctx.span.StartChild("save-user-login-attempt")
	span.Status = sentry.SpanStatusOK
	if err := a.repo.SaveNewUserLogin(repository.NewDBOperationContext(ctx, span), loginRecord); nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SAVE_NEW_LOGIN_INFO")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed saving user new login attempt information")
		return nil, err
	}
	span.Finish()

	return &LoginResult{
		Token: LoginResultToken{
			ID:                 token.ID,
			Value:              token.Value,
			NotBeforeDateTime:  token.NotBeforeDateTime,
			ExpirationDateTime: token.ExpirationDateTime,
		},
	}, nil
}
