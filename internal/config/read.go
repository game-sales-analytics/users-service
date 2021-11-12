package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func readEnvironmetVariablesOrUseDefaults(logger *logrus.Entry) (Config, error) {
	logger.Trace("loading default configuration")
	conf := getDefaults()

	if value, exists := os.LookupEnv("SERVER_HOST"); exists {
		logger.WithField("variable", "SERVER_HOST").WithField("value", value).Debug("using provided environment variable")
		conf.Server.Host = value
	}

	if value, exists := os.LookupEnv("SERVER_PORT"); exists {
		value, err := strconv.ParseUint(value, 10, 32)
		if nil != err {
			return Config{}, err
		}

		logger.WithField("variable", "SERVER_PORT").WithField("value", value).Debug("using provided environment variable")
		conf.Server.Port = uint(value)
	}

	if value, exists := os.LookupEnv("DATABASE_HOST"); exists {
		logger.WithField("variable", "DATABASE_HOST").WithField("value", value).Debug("using provided environment variable")
		conf.Database.Host = value
	}

	if value, exists := os.LookupEnv("DATABASE_PORT"); exists {
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

	if value, exists := os.LookupEnv("DATABASE_USERNAME"); exists {
		logger.WithField("variable", "DATABASE_USERNAME").WithField("value", strings.Repeat("*", len(value))).Debug("using provided environment variable")
		conf.Database.Username = value
	}

	if value, exists := os.LookupEnv("DATABASE_PASSWORD"); exists {
		logger.WithField("variable", "DATABASE_PASSWORD").WithField("value", strings.Repeat("*", len(value))).Debug("using provided environment variable")
		conf.Database.Password = value
	}

	if value, exists := os.LookupEnv("DATABASE_NAME"); exists {
		logger.WithField("variable", "DATABASE_NAME").WithField("value", value).Debug("using provided environment variable")
		conf.Database.Name = value
	}

	if value, exists := os.LookupEnv("JWT_SECRET"); exists {
		logger.WithField("variable", "JWT_SECRET").WithField("value", strings.Repeat("*", len(value))).Debug("using provided environment variable")
		conf.Jwt.Secret = value
	} else {
		return Config{}, errors.New("'JWT_SECRET' environment variable is required")
	}

	return conf, nil
}
