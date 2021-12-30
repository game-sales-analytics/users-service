package repository

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/game-sales-analytics/users-service/internal/apm"
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

func (r *Repo) SaveNewUser(ctx DBOperationContext, user NewUserToSave) error {
	userDoc := bson.D{
		{Key: "id", Value: user.ID},
		{Key: "registered_at", Value: user.RegisteredAt},
		{Key: "email", Value: user.Email},
		{Key: "normalized_email", Value: user.NormalizedEmail},
		{Key: "password", Value: user.Password},
		{Key: "first_name", Value: user.FirstName},
		{Key: "last_name", Value: user.LastName},
	}

	span := ctx.span.StartChild("insert-user")
	span.Status = sentry.SpanStatusOK
	_, err := r.collections.Users.InsertOne(ctx, userDoc)
	if nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithError(err).WithField("err_code", "E_SAVE_USER")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("unable to save user to database")

		return err
	}
	span.Finish()

	return nil
}

func (r *Repo) GetUserLoginInfo(ctx DBOperationContext, email string) (*UserLoginInfo, error) {
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

	span := ctx.span.StartChild("query-user-login-info")
	span.Status = sentry.SpanStatusOK
	result := r.collections.Users.FindOne(ctx, filter, opts)
	if err := result.Err(); nil != err {
		defer span.Finish()

		if errors.Is(err, mongo.ErrNoDocuments) {
			span.Status = sentry.SpanStatusNotFound
			return nil, ErrUserNotExists
		}

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithError(err).WithField("err_code", "E_FIND_USER_LOGIN_INFO")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("unable to retrieve user login information")
		return nil, errors.New("unable to retrieve user login information")
	}
	span.Finish()

	span = ctx.span.StartChild("decode-queried-user-login-info")
	span.Status = sentry.SpanStatusOK
	if err := result.Decode(&user); nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithError(err).WithField("err_code", "E_DECODE_DOCUMENT")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("unable to decode user login information document")
		return nil, errors.New("unable to decode retrieved user login information")
	}
	span.Finish()

	span = ctx.span.StartChild("parse-decoded-user-login-info")
	span.Status = sentry.SpanStatusOK
	passwd, ok := user["password"].(string)
	if !ok {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithField("err_code", "E_CAST_USER_PASSWORD_FIELD")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("could not parse user document password field")
		return nil, errors.New("could not parse user password")
	}
	id, ok := user["id"].(string)
	if !ok {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithField("err_code", "E_CAST_USER_ID_FIELD")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("could not parse user document id field")
		return nil, errors.New("could not parse user id")
	}
	span.Finish()

	return &UserLoginInfo{
		Password: passwd,
		ID:       id,
	}, nil
}

func (r *Repo) NormalizedEmailExists(ctx DBOperationContext, normalizedEmail string) (bool, error) {
	filter := bson.M{
		"normalized_email": normalizedEmail,
	}

	span := ctx.span.StartChild("check-user-with-email-existence")
	span.Status = sentry.SpanStatusOK
	exists, err := r.userWithFilterExists(NewDBOperationContext(ctx, span), filter)
	if nil != err {
		defer span.Finish()

		return false, err
	}
	span.Finish()

	return exists, nil
}

func (r *Repo) UserWithIDExists(ctx DBOperationContext, userID string) (bool, error) {
	filter := bson.M{
		"id": userID,
	}

	span := ctx.span.StartChild("check-user-with-email-existence")
	span.Status = sentry.SpanStatusOK
	exists, err := r.userWithFilterExists(NewDBOperationContext(ctx, span), filter)
	if nil != err {
		defer span.Finish()

		return false, err
	}
	span.Finish()

	return exists, nil
}

func (r *Repo) userWithFilterExists(ctx DBOperationContext, filter bson.M) (bool, error) {
	projection := bson.D{}
	opts := options.FindOne().SetProjection(projection)
	doc := bson.M{}

	span := ctx.span.StartChild("query-user-with-filter")
	span.Status = sentry.SpanStatusOK
	result := r.collections.Users.FindOne(ctx, filter, opts)
	if err := result.Err(); nil != err {
		defer span.Finish()

		if errors.Is(err, mongo.ErrNoDocuments) {
			span.Status = sentry.SpanStatusNotFound
			return false, nil
		}

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithError(err).WithField("err_code", "E_RETRIEVE_USER_DOCUMENT")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed retrieving user document")
		return false, err
	}

	span = ctx.span.StartChild("decode-queried-user")
	span.Status = sentry.SpanStatusOK
	if err := result.Decode(&doc); nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithError(err).WithField("err_code", "E_DECODE_DOCUMENT")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("unable to decode user document")
		return false, err
	}
	span.Finish()

	if userID, ok := doc["_id"].(string); !ok || len(userID) > 0 {
		defer span.Finish()

		span.Status = sentry.SpanStatusOK
		return true, nil
	}

	log := r.logger.WithField("err_code", "E_UNEXPECTED_LOGICAL_SITUATION")
	apm.SetSpanTagsFromLogEntry(span, log)
	log.Error("expected to not to reach this branch")
	return false, errors.New("unexpected situation")
}

type UserAuthenticationInfo struct {
	FirstName string
	LastName  string
}

func (r *Repo) GetUserAuthenticationInfo(ctx DBOperationContext, userID string) (*UserAuthenticationInfo, error) {
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

	span := ctx.span.StartChild("query-user-authentication-info")
	span.Status = sentry.SpanStatusOK
	result := r.collections.Users.FindOne(ctx, filter, opts)
	if err := result.Err(); nil != err {
		defer span.Finish()

		if errors.Is(err, mongo.ErrNoDocuments) {
			span.Status = sentry.SpanStatusNotFound
			return nil, ErrUserNotExists
		}

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithError(err).WithField("err_code", "E_RETRIEVE_USER_DOCUMENT")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed retrieving user authentication information document")
		return nil, err
	}
	span.Finish()

	span = ctx.span.StartChild("decode-queried-authentication-info")
	span.Status = sentry.SpanStatusOK
	if err := result.Decode(&user); nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithError(err).WithField("err_code", "E_DECODE_DOCUMENT")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("unable to decode user authentication information document")
		return nil, err
	}
	span.Finish()

	span = ctx.span.StartChild("parse-decoded-authentication-info")
	span.Status = sentry.SpanStatusOK
	firstName, ok := user["first_name"].(string)
	if !ok {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithField("err_code", "E_CAST_USER_FIRST_NAME_FIELD")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("could not parse user document first_name field")
		return nil, errors.New("could not parse user first_name")
	}

	lastName, ok := user["last_name"].(string)
	if !ok {
		defer span.Finish()

		span.Status = sentry.SpanStatusInternalError
		log := r.logger.WithField("err_code", "E_CAST_USER_LAST_NAME_FIELD")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("could not parse user document last_name field")
		return nil, errors.New("could not parse user last_name")
	}
	span.Finish()

	return &UserAuthenticationInfo{
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
