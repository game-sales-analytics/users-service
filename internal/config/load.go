package config

import (
	"github.com/sirupsen/logrus"
)

func Load(logger *logrus.Entry) (Config, error) {
	return readEnvironmetVariablesOrUseDefaults(logger)
}
