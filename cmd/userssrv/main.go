package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"

	"github.com/game-sales-analytics/users-service/internal/auth"
	"github.com/game-sales-analytics/users-service/internal/config"
	"github.com/game-sales-analytics/users-service/internal/db"
	"github.com/game-sales-analytics/users-service/internal/grpcsrv"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func main() {
	logger := logrus.New()

	logger.AddHook(&apmlogrus.Hook{})
	// TODO: make level configurable
	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		PrettyPrint:       true,
		TimestampFormat:   time.RFC3339Nano,
	})

	logger.Trace("initializing background context")
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Trace("canceling background context")
		cancel()
	}()

	logger.Trace("loading configuration")
	conf, err := config.Load(logger.WithField("srv", "config"))
	if nil != err {
		logger.WithError(err).Fatal("unable to load configuration")
	}

	logger.Trace("initializing database connection")
	database, err := db.Connect(ctx, logger.WithField("srv", "db"), &conf.Database)
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

	validator := validate.New(logger.WithField("srv", "validate"), &database.Repo)
	authSrv := auth.New(&database.Repo, logger.WithField("srv", "auth"), &conf.Jwt)

	server := grpcsrv.New(logger.WithField("srv", "grpc"), &database.Repo, validator, authSrv)
	logger.WithError(server.Listen(conf.Server.Host, conf.Server.Port)).Fatal("unable to start GRPC server")
}
