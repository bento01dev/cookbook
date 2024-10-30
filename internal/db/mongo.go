package db

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	clientOnce sync.Once
	client     *mongo.Client
	clientErr  error
)

func MongoClient(getEnv func(string) string) (*mongo.Client, error) {
	clientOnce.Do(func() {
		client, clientErr = mongo.Connect(options.Client().ApplyURI(getEnv("MONGO_DB_URL")))
	})
	return client, clientErr
}

func Close(ctx context.Context) error {
	return client.Disconnect(ctx)
}
