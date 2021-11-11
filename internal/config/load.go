package config

import (
	"github.com/sirupsen/logrus"
)

func Load(logger *logrus.Logger) (Config, error) {
	return readEnvironmetVariablesOrUseDefaults(logger)
}
