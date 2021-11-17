package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type NewUserLoginToSave struct {
	ID                  string
	UserID              string
	LoggedInAt          time.Time
	UserIPAddress       string
	UserDeviceUserAgent string
}

func (r *Repo) SaveNewUserLogin(ctx context.Context, userLogin NewUserLoginToSave) error {
	doc := bson.D{
		{Key: "id", Value: userLogin.ID},
		{Key: "logged_in_at", Value: userLogin.LoggedInAt},
		{Key: "user", Value: bson.D{
			{Key: "id", Value: userLogin.UserID},
			{Key: "ip", Value: userLogin.UserIPAddress},
			{Key: "device_agent", Value: userLogin.UserDeviceUserAgent},
		}},
	}
	_, err := r.collections.UserLogins.InsertOne(ctx, doc)
	return err
}
