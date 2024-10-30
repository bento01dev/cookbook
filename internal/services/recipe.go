package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bento01dev/cookbook/internal/domain"
	"github.com/bento01dev/cookbook/internal/domain/recipe"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type recipeRepository interface {
	Get(context.Context, uuid.UUID) (recipe.Recipe, error)
	Add(context.Context, recipe.Recipe) error
	Update(context.Context, recipe.Recipe) (recipe.Recipe, error)
	Delete(context.Context, uuid.UUID) error
}

type RecipeService struct {
	recipes recipeRepository
}

type RecipeConfiguration func(rs *RecipeService) error

func NewRecipeService(cfgs ...RecipeConfiguration) (RecipeService, error) {
	rs := RecipeService{}

	for _, cfg := range cfgs {
		err := cfg(&rs)
		if err != nil {
			return rs, err
		}
	}

	return rs, nil
}

func WithMemoryRepository() RecipeConfiguration {
	return func(rs *RecipeService) error {
		mr := recipe.NewMemoryRepository()
		rs.recipes = mr
		return nil
	}
}

func WithMongoRepository(client *mongo.Client, getEnv func(string) string) RecipeConfiguration {
	return func(rs *RecipeService) error {
		databaseName := getEnv("MONGO_DB")
		if databaseName == "" {
			return errors.New("DB not set. Set env MONGO_DB")
		}

		collectionName := getEnv("RECIPE_COLLECTION")
		if collectionName == "" {
			return errors.New("recipe collection not set. Set env RECIPE_COLLECTION")
		}

		mr := recipe.NewMongoRepository(client, databaseName, collectionName)
		rs.recipes = mr
		return nil
	}
}

func (rs RecipeService) CreateRecipe(ctx context.Context, name string, description string, cuisine domain.CuisineType) (recipe.Recipe, error) {
	r, err := recipe.NewRecipe(name, description, cuisine)
	if err != nil {
		return r, err
	}

	err = rs.recipes.Add(ctx, r)
	if err != nil {
		return r, err
	}
	slog.InfoContext(ctx, "recipe successfully added", "recipe_id", r.ID().String())
	return r, nil
}

func (rs RecipeService) GetRecipe(ctx context.Context, uuidStr string) (recipe.Recipe, error) {
	slog.InfoContext(ctx, "retrieving recipe..", "recipe_id", uuidStr)
	recipeUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return recipe.Recipe{}, recipe.ErrInvalidID
	}

	r, err := rs.recipes.Get(ctx, recipeUuid)
	if err != nil {
		return recipe.Recipe{}, err
	}

	return r, nil
}
