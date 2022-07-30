package db

import (
	"github.com/goava/di"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionNames = []string{
	"people",
	"signals",
}
var Collections = di.Options(providers(collectionNames...)...)

func providers(collections ...string) []di.Option {
	opts := make([]di.Option, 0)
	for _, collectionName := range collections {
		opt := di.Provide(func(database *mongo.Database) *mongo.Collection {
			coll := database.Collection(collectionName)
			return coll
		}, di.Tags{"collection": collectionName})
		opts = append(opts, opt)
	}
	return opts
}
