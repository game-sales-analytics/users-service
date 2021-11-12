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
		{Key: "registered_at", Value: user.RegisteredAt.Format(time.RFC3339)},
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
		return false, err
	}
	if err := result.Decode(&doc); nil != err {
		return false, err
	}

	if userID, ok := doc["_id"].(string); !ok || len(userID) > 0 {
		return true, nil
	}

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
		return nil, err
	}
	if err := result.Decode(&user); nil != err {
		return nil, err
	}

	firstName, ok := user["first_name"].(string)
	if !ok {
		return nil, errors.New("could not parse user first_name")
	}
	lastName, ok := user["last_name"].(string)
	if !ok {
		return nil, errors.New("could not parse user last_name")
	}

	return &UserAuthenticationInfo{
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
