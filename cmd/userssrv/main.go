package main

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
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

	logger.Trace("initializing Sentry sdk")
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              conf.APM.DSN,
		Environment:      conf.APM.Env,
		Release:          conf.APM.Release,
		Debug:            true,
		TracesSampleRate: 1.0,
		SampleRate:       1.0,
		AttachStacktrace: true,
	})
	if nil != err {
		logger.WithError(err).Fatal("unable to initialize Sentry sdk")
	}

	defer sentry.Flush(4 * time.Second)
	logger.Trace("Sentry sdk initialized")

	span := sentry.StartSpan(ctx, "startup", sentry.TransactionName("service-startup"))

	logger.Trace("initializing database connection")
	child := span.StartChild("create-database-connection")
	database, err := db.Connect(db.NewConnectContext(ctx, child), logger.WithField("srv", "db"), &conf.Database)
	if nil != err {
		defer child.Finish()

		logger.WithError(err).Fatal("unable to connect to database")
	}
	child.Finish()

	defer func() {
		logger.Debug("closing database connection before exit")
		if err := database.Disconnect(); nil != err {
			logger.WithError(err).Debug("unable to close database connection")
		}
	}()
	logger.Trace("connected to database")

	validator := validate.New(logger.WithField("srv", "validate"), &database.Repo)
	authSrv := auth.New(&database.Repo, logger.WithField("srv", "auth"), &conf.Jwt)

	span.Finish()

	server := grpcsrv.New(logger.WithField("srv", "grpc"), &database.Repo, validator, authSrv)
	logger.WithError(server.Listen(conf.Server.Host, conf.Server.Port)).Fatal("unable to start GRPC server")
}
