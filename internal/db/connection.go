package db

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/game-sales-analytics/users-service/internal/config"
	"github.com/game-sales-analytics/users-service/internal/db/repository"
)

func Connect(ctx context.Context, logger *logrus.Entry, cfg *config.DatabaseConfig) (*DB, error) {
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

	clientOptions.SetMonitor(apmmongo.CommandMonitor())

	logger.Trace("connecting database")
	client, err := mongo.Connect(ctx, clientOptions)
	if nil != err {
		return nil, err
	}

	logger.Trace("checking database connection healthiness")
	if err := client.Ping(ctx, nil); nil != err {
		return nil, err
	}

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
