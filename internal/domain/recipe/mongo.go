package recipe

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository struct {
	client         *mongo.Client
	databaseName   string
	collectionName string
}

func NewMongoRepository(client *mongo.Client, databaseName, collectionName string) *MongoRepository {
	return &MongoRepository{
		client:         client,
		databaseName:   databaseName,
		collectionName: collectionName,
	}
}

func (mr *MongoRepository) Get(ctx context.Context, id uuid.UUID) (Recipe, error) {
	return Recipe{}, nil
}

func (mr *MongoRepository) Add(ctx context.Context, recipe Recipe) error {
	return nil
}

func (mr *MongoRepository) Update(ctx context.Context, recipe Recipe) (Recipe, error) {
	return Recipe{}, nil
}

func (mr *MongoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
