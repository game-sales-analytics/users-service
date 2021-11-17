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
	ID              string
	Email           string
	NormalizedEmail string
	Password        string
	FirstName       string
	LastName        string
	RegisteredAt    time.Time
}

func (r *Repo) SaveNewUser(ctx context.Context, user NewUserToSave) error {
	userDoc := bson.D{
		{Key: "id", Value: user.ID},
		{Key: "registered_at", Value: user.RegisteredAt},
		{Key: "email", Value: user.Email},
		{Key: "normalized_email", Value: user.NormalizedEmail},
		{Key: "password", Value: user.Password},
		{Key: "first_name", Value: user.FirstName},
		{Key: "last_name", Value: user.LastName},
	}
	_, err := r.collections.Users.InsertOne(ctx, userDoc)
	return err
}

func (r *Repo) GetUserLoginInfo(ctx context.Context, email string) (*UserLoginInfo, error) {
	filter := bson.M{
		"email": email,
	}
	projection := bson.D{
		bson.E{Key: "_id", Value: 0},
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
		r.logger.WithError(err).WithField("err_code", "E_FIND_USER_LOGIN_INFO").Error("unable to retrieve user login information")
		return nil, errors.New("unable to retrieve user login information")
	}

	if err := result.Decode(&user); nil != err {
		r.logger.WithError(err).WithField("err_code", "E_DECODE_DOCUMENT").Error("unable to decode user login information document")
		return nil, errors.New("unable to decode retrieved user login information")
	}

	passwd, ok := user["password"].(string)
	if !ok {
		r.logger.WithField("err_code", "E_CAST_USER_PASSWORD_FIELD").Error("could not parse user document password field")
		return nil, errors.New("could not parse user password")
	}
	id, ok := user["id"].(string)
	if !ok {
		r.logger.WithField("err_code", "E_CAST_USER_ID_FIELD").Error("could not parse user document id field")
		return nil, errors.New("could not parse user id")
	}

	return &UserLoginInfo{
		Password: passwd,
		ID:       id,
	}, nil
}

func (r *Repo) NormalizedEmailExists(ctx context.Context, normalizedEmail string) (bool, error) {
	filter := bson.M{
		"normalized_email": normalizedEmail,
	}
	return r.userWithFilterExists(ctx, filter)
}

func (r *Repo) UserWithIDExists(ctx context.Context, userID string) (bool, error) {
	filter := bson.M{
		"id": userID,
	}
	return r.userWithFilterExists(ctx, filter)
}

func (r *Repo) userWithFilterExists(ctx context.Context, filter bson.M) (bool, error) {
	projection := bson.D{}
	opts := options.FindOne().SetProjection(projection)
	doc := bson.M{}
	result := r.collections.Users.FindOne(ctx, filter, opts)
	if err := result.Err(); nil != err {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		r.logger.WithError(err).WithField("err_code", "E_RETRIEVE_USER_DOCUMENT").Error("failed retrieving user document")
		return false, err
	}

	if err := result.Decode(&doc); nil != err {
		r.logger.WithError(err).WithField("err_code", "E_DECODE_DOCUMENT").Error("unable to decode user document")
		return false, err
	}

	if userID, ok := doc["_id"].(string); !ok || len(userID) > 0 {
		return true, nil
	}

	r.logger.WithField("err_code", "E_UNEXPECTED_LOGICAL_SITUATION").Error("expected to not to reach this branch")
	return false, errors.New("unexpected situation")
}

type UserAuthenticationInfo struct {
	FirstName string
	LastName  string
}

func (r *Repo) GetUserAuthenticationInfo(ctx context.Context, userID string) (*UserAuthenticationInfo, error) {
	filter := bson.M{
		"id": userID,
	}
	projection := bson.D{
		bson.E{Key: "_id", Value: 0},
		bson.E{Key: "first_name", Value: 1},
		bson.E{Key: "last_name", Value: 1},
	}
	opts := options.FindOne().SetProjection(projection)
	user := bson.M{}
	result := r.collections.Users.FindOne(ctx, filter, opts)
	if err := result.Err(); nil != err {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotExists
		}

		r.logger.WithError(err).WithField("err_code", "E_RETRIEVE_USER_DOCUMENT").Error("failed retrieving user authentication information document")
		return nil, err
	}

	if err := result.Decode(&user); nil != err {
		r.logger.WithError(err).WithField("err_code", "E_DECODE_DOCUMENT").Error("unable to decode user authentication information document")
		return nil, err
	}

	firstName, ok := user["first_name"].(string)
	if !ok {
		r.logger.WithField("err_code", "E_CAST_USER_FIRST_NAME_FIELD").Error("could not parse user document first_name field")
		return nil, errors.New("could not parse user first_name")
	}

	lastName, ok := user["last_name"].(string)
	if !ok {
		r.logger.WithField("err_code", "E_CAST_USER_LAST_NAME_FIELD").Error("could not parse user document last_name field")
		return nil, errors.New("could not parse user last_name")
	}

	return &UserAuthenticationInfo{
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
