package repository

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/game-sales-analytics/users-service/internal/apm"
)

type NewUserLoginToSave struct {
	ID                  string
	UserID              string
	LoggedInAt          time.Time
	UserIPAddress       string
	UserDeviceUserAgent string
}

func (r *Repo) SaveNewUserLogin(ctx DBOperationContext, userLogin NewUserLoginToSave) error {
	doc := bson.D{
		{Key: "id", Value: userLogin.ID},
		{Key: "logged_in_at", Value: userLogin.LoggedInAt},
		{Key: "user", Value: bson.D{
			{Key: "id", Value: userLogin.UserID},
			{Key: "ip", Value: userLogin.UserIPAddress},
			{Key: "device_agent", Value: userLogin.UserDeviceUserAgent},
		}},
	}

	span := ctx.span.StartChild("insert-user-login-info")
	span.Status = sentry.SpanStatusOK
	_, err := r.collections.UserLogins.InsertOne(ctx, doc)
	if nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithError(err).WithField("err_code", "E_SAVE_USER_LOGIN_ATTEMPT")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed saving user login attempt record")
		return err
	}

	return nil
}
