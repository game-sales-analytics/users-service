package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"

	"github.com/game-sales-analytics/users-service/internal/constants"
)

func readEnvironmetVariablesOrUseDefaults(logger *logrus.Entry) (Config, error) {
	logger.Trace("loading default configuration")
	conf := getDefaults()

	if value, exists := os.LookupEnv("SERVER_HOST"); exists && len(value) != 0 {
		logger.WithField("variable", "SERVER_HOST").WithField("value", value).Debug("using provided environment variable")
		conf.Server.Host = value
	}

	if value, exists := os.LookupEnv("SERVER_PORT"); exists && len(value) != 0 {
		value, err := strconv.ParseUint(value, 10, 32)
		if nil != err {
			return Config{}, err
		}

		logger.WithField("variable", "SERVER_PORT").WithField("value", value).Debug("using provided environment variable")
		conf.Server.Port = uint(value)
	}

	if value, exists := os.LookupEnv("DATABASE_HOST"); exists && len(value) != 0 {
		logger.WithField("variable", "DATABASE_HOST").WithField("value", value).Debug("using provided environment variable")
		conf.Database.Host = value
	}

	if value, exists := os.LookupEnv("DATABASE_PORT"); exists && len(value) != 0 {
		value, err := strconv.ParseUint(value, 10, 32)
		if nil != err {
			return Config{}, err
		}

		logger.WithField("variable", "DATABASE_PORT").WithField("value", value).Debug("using provided environment variable")
		conf.Database.Port = uint(value)
	}

	if _, exists := os.LookupEnv("DATABASE_USE_AUTH"); exists {
		logger.WithField("variable", "DATABASE_USE_AUTH").Debug("enabling authentication in database connection due to existence of environment variable")
		conf.Database.UseAuth = true
	}

	if value, exists := os.LookupEnv("DATABASE_USERNAME"); exists && len(value) != 0 {
		logger.WithField("variable", "DATABASE_USERNAME").WithField("value", strings.Repeat("*", len(value))).Debug("using provided environment variable")
		conf.Database.Username = value
	}

	if value, exists := os.LookupEnv("DATABASE_PASSWORD"); exists && len(value) != 0 {
		logger.WithField("variable", "DATABASE_PASSWORD").WithField("value", strings.Repeat("*", len(value))).Debug("using provided environment variable")
		conf.Database.Password = value
	}

	if value, exists := os.LookupEnv("DATABASE_NAME"); exists && len(value) != 0 {
		logger.WithField("variable", "DATABASE_NAME").WithField("value", value).Debug("using provided environment variable")
		conf.Database.Name = value
	}

	if value, exists := os.LookupEnv("JWT_SECRET"); exists && len(value) != 0 {
		logger.WithField("variable", "JWT_SECRET").WithField("value", strings.Repeat("*", len(value))).Debug("using provided environment variable")
		conf.Jwt.Secret = value
	} else {
		return Config{}, errors.New("'JWT_SECRET' environment variable is required")
	}

	if value, exists := os.LookupEnv("SENTRY_DSN"); exists && len(value) != 0 {
		dsn, err := sentry.NewDsn(value)
		if nil != err {
			return Config{}, fmt.Errorf("invalid 'SENTRY_DSN' environment variable is provided: %s", err)
		}

		maskedDsn := dsn.EnvelopeAPIURL()
		var userInfo *url.Userinfo
		if passwd, hasPasswd := dsn.EnvelopeAPIURL().User.Password(); hasPasswd {
			userInfo = url.UserPassword(strings.Repeat("*", len(dsn.EnvelopeAPIURL().User.Username())), strings.Repeat("*", len(passwd)))
		} else {
			userInfo = url.User(strings.Repeat("*", len(dsn.EnvelopeAPIURL().User.Username())))
		}
		maskedDsn.User = userInfo
		logger.WithField("variable", "SENTRY_DSN").WithField("value", maskedDsn.String()).Debug("using provided Sentry DSN environment variable")

		conf.APM.DSN = value
	} else {
		return Config{}, errors.New("'SENTRY_DSN' environment variable is required")
	}

	if value, exists := os.LookupEnv("SENTRY_RELEASE"); exists && len(value) != 0 {
		logger.WithField("variable", "SENTRY_RELEASE").WithField("value", value).Debug("using provided Sentry Release environment variable")
		conf.APM.Release = value
	} else {
		conf.APM.Release = constants.VERSION
	}

	if value, exists := os.LookupEnv("SENTRY_ENVIRONMENT"); exists && len(value) != 0 {
		logger.WithField("variable", "SENTRY_ENVIRONMENT").WithField("value", value).Debug("using provided Sentry Environment environment variable")
		conf.APM.Env = value
	} else {
		conf.APM.Env = "prod"
	}

	return conf, nil
}
