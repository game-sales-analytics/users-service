package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserLoginInfo struct {
	ID       string
	Password string
}

var (
	ErrUserNotExists = errors.New("user does not exist")
)

type NewUserToSave struct {
	ID           string
	Email        string
	Password     string
	RegisteredAt time.Time
}

func (r *Repo) SaveNewUser(ctx context.Context, user NewUserToSave) error {
	userDoc := bson.D{
		{Key: "id", Value: user.ID},
		{Key: "registered_at", Value: user.RegisteredAt.Format(time.RFC3339)},
		{Key: "email", Value: user.Email},
		{Key: "password", Value: user.Password},
	}
	_, err := r.collections.Users.InsertOne(ctx, userDoc)
	return err
}

func (r *Repo) GetUserLoginInfo(ctx context.Context, email string) (*UserLoginInfo, error) {
	filter := bson.M{
		"email": email,
	}
	projection := bson.D{
		bson.E{Key: "id", Value: 1},
		bson.E{Key: "password", Value: 1},
	}
	opts := options.FindOne().SetProjection(projection)
	user := bson.M{}
	result := r.collections.Users.FindOne(ctx, filter, opts)
	if err := result.Err(); nil != err {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotExists
		}
		return nil, err
	}
	if err := result.Decode(&user); nil != err {
		return nil, err
	}

	passwd, ok := user["password"].(string)
	if !ok {
		return nil, errors.New("could not parse user password")
	}
	id, ok := user["id"].(string)
	if !ok {
		return nil, errors.New("could not parse user id")
	}

	return &UserLoginInfo{
		Password: passwd,
		ID:       id,
	}, nil
}
