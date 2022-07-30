package db

import (
	"context"
	middlewares "github.com/umangraval/Go-Mongodb-REST-boilerplate/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDatabase() (*mongo.Database, error) {

	mongoDbUrl := middlewares.DotEnvVariable("MONGO_URL")

	clientOptions := options.Client().ApplyURI(mongoDbUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database("golang"), err
}
