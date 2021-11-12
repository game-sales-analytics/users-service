package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Users      *mongo.Collection
	UserLogins *mongo.Collection
}

type Repo struct {
	collections Collections
	logger      *logrus.Entry
}

func New(logger *logrus.Entry, collections Collections) Repo {
	return Repo{
		collections,
		logger,
	}
}
