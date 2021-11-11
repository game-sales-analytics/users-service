package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Users      *mongo.Collection
	UserLogins *mongo.Collection
}

type Repo struct {
	collections Collections
}

func New(collections Collections) Repo {
	return Repo{
		collections,
	}
}
