package auth

import (
	"context"
	"errors"
	"math/rand"
	"time"

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

func (a authsrv) LoginWithEmail(ctx context.Context, creds LoginWithEmailCreds) (*LoginResult, error) {
	user, err := a.repo.GetUserLoginInfo(ctx, creds.Email)
	if nil != err {
		if errors.Is(err, repository.ErrUserNotExists) {
			randSpan := rand.Int63n(251)
			time.Sleep(time.Millisecond * time.Duration((6136 + randSpan)))
			return nil, ErrUnauthenticated
		}

		a.logger.WithError(err).WithField("err_code", "E_RETRIEVE_USER_LOGIN_INFO").Error("failed retrieving user login information")
		return nil, ErrInternal
	}

	matched, err := passhash.Verify(creds.Password, user.Password)
	if nil != err {
		a.logger.WithError(err).WithField("err_code", "E_VERIFY_PASSWORD").Error("failed verifying user password for login")
		return nil, ErrInternal
	}
	if !matched {
		return nil, ErrUnauthenticated
	}

	token, err := a.generateToken(user.ID, a.cfg.Secret)
	if nil != err {
		a.logger.WithError(err).WithField("err_code", "E_GENERATE_AUTH_TOKEN").Error("failed generating authentication token for user")
		return nil, ErrInternal
	}

	loginRecordID, err := id.GenerateUserLoginID()
	if nil != err {
		a.logger.WithError(err).WithField("err_code", "E_GENERATE_NEW_LOGIN_RECORD_ID").Error("failed generating id for new user login record")
		return nil, ErrInternal
	}

	loginRecord := repository.NewUserLoginToSave{
		ID:                  loginRecordID,
		UserID:              user.ID,
		LoggedInAt:          time.Now(),
		UserIPAddress:       creds.UserIPAddress,
		UserDeviceUserAgent: creds.UserDeviceUserAgent,
	}
	if err := a.repo.SaveNewUserLogin(ctx, loginRecord); nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SAVE_NEW_LOGIN_INFO").Error("failed saving user new login attempt information")
		return nil, err
	}

	return &LoginResult{
		Token: LoginResultToken{
			ID:                 token.ID,
			Value:              token.Value,
			NotBeforeDateTime:  token.NotBeforeDateTime,
			ExpirationDateTime: token.ExpirationDateTime,
		},
	}, nil
}
