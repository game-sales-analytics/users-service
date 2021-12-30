package db

import (
	"context"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/game-sales-analytics/users-service/internal/config"
	"github.com/game-sales-analytics/users-service/internal/db/repository"
)

type ConnectContext struct {
	context.Context
	span *sentry.Span
}

func NewConnectContext(ctx context.Context, span *sentry.Span) ConnectContext {
	return ConnectContext{
		ctx,
		span,
	}
}

func Connect(ctx ConnectContext, logger *logrus.Entry, cfg *config.DatabaseConfig) (*DB, error) {
	addr := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	logger.WithField("address", addr).Debug("connecting database using configured address")

	clientOptions := options.
		Client().
		ApplyURI(addr).
		SetConnectTimeout(time.Second * 10).
		SetServerSelectionTimeout(time.Second * 10)

	if cfg.UseAuth {
		logger.Debug("using provided credentials for database connection authentication")
		clientOptions = clientOptions.
			SetAuth(options.Credential{
				PasswordSet:   true,
				AuthMechanism: "SCRAM-SHA-256",
				AuthSource:    cfg.Name,
				Password:      cfg.Password,
				Username:      cfg.Username,
			})
	}

	logger.Trace("connecting database")
	child := ctx.span.StartChild("create-database-connection")
	child.Status = sentry.SpanStatusOK
	client, err := mongo.Connect(ctx, clientOptions)
	if nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		return nil, err
	}
	child.Finish()

	logger.Trace("checking database connection healthiness")
	child = ctx.span.StartChild("check-database-connection-healthiness")
	child.Status = sentry.SpanStatusOK
	if err := client.Ping(ctx, nil); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusUnavailable
		return nil, err
	}
	child.Finish()

	logger.WithField("database", cfg.Name).Debug("using configured database name")
	db := client.Database(cfg.Name)
	return &DB{
		client: client,
		logger: logger,
		Repo: repository.New(
			logger.WithField("srv", "repository"),
			repository.Collections{
				Users:      db.Collection(UsersCollectionName),
				UserLogins: db.Collection(UserLoginsCollectionName),
			},
		),
	}, nil
}

func (db *DB) Disconnect() error {
	db.logger.Trace("closing database connection")
	return db.client.Disconnect(context.Background())
}
