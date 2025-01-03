package recipe

import (
	"context"
	"errors"
	"time"

	"github.com/bento01dev/cookbook/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
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

type recipe struct {
	ID          uuid.UUID      `bson:"id"`
	Name        string         `bson:"name"`
	Description string         `bson:"description"`
	CreatedAt   bson.Timestamp `bson:"created_at"`
}

func (r recipe) ToRecipe() Recipe {
	return Recipe{
		item: &domain.Item{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		},
		createdAt: time.Unix(int64(r.CreatedAt.T), 0),
	}
}

func recipeFromRecipe(r Recipe) recipe {
	return recipe{
		ID:          r.item.ID,
		Name:        r.item.Name,
		Description: r.item.Description,
		CreatedAt:   bson.Timestamp{T: uint32(r.createdAt.Unix())},
	}
}

func (mr *MongoRepository) Get(ctx context.Context, id uuid.UUID) (Recipe, error) {
	collection := mr.client.Database(mr.databaseName).Collection(mr.collectionName)
	var result recipe
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Recipe{}, ErrRecipeNotFound
		}
		return Recipe{}, err
	}
	return result.ToRecipe(), nil
}

func (mr *MongoRepository) Add(ctx context.Context, recipe Recipe) error {
	collection := mr.client.Database(mr.databaseName).Collection(mr.collectionName)
	_, err := collection.InsertOne(ctx, recipeFromRecipe(recipe))
	return err
}

func (mr *MongoRepository) Update(ctx context.Context, recipe Recipe) (Recipe, error) {
	return Recipe{}, nil
}

func (mr *MongoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
