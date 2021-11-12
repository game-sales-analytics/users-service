package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/game-sales-analytics/users-service/internal/auth"
	"github.com/game-sales-analytics/users-service/internal/config"
	"github.com/game-sales-analytics/users-service/internal/db"
	"github.com/game-sales-analytics/users-service/internal/grpcsrv"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func main() {
	logger := logrus.New()

	// TODO: make level configurable
	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

	logger.Trace("initializing background context")
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Trace("canceling background context")
		cancel()
	}()

	logger.Trace("loading configuration")
	conf, err := config.Load(logger)
	if nil != err {
		logger.WithError(err).Fatal("unable to load configuration")
	}

	logger.Trace("initializing database connection")
	database, err := db.Connect(ctx, logger, &conf.Database)
	if nil != err {
		logger.WithError(err).Fatal("unable to connect to database")
	}
	defer func() {
		logger.Debug("closing database connection before exit")
		if err := database.Disconnect(); nil != err {
			logger.WithError(err).Debug("unable to close database connection")
		}
	}()
	logger.Trace("connected to database")

	validator := validate.New(logger, &database.Repo)
	authSrv := auth.New(&database.Repo, logger, &conf.Jwt)

	server := grpcsrv.New(logger, &database.Repo, validator, authSrv)
	logger.WithError(server.Listen(conf.Server.Host, conf.Server.Port)).Fatal("unable to start GRPC server")
}
