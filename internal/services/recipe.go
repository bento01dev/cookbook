package services

import (
	"fmt"
	"log/slog"

	"github.com/bento01dev/cookbook/internal/domain"
	"github.com/bento01dev/cookbook/internal/domain/recipe"
	"github.com/google/uuid"
)

type recipeRepository interface {
	Get(uuid.UUID) (recipe.Recipe, error)
	Add(recipe.Recipe) error
	Update(recipe.Recipe) (recipe.Recipe, error)
	Delete(uuid.UUID) error
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

func (rs RecipeService) CreateRecipe(name string, description string, cuisine domain.CuisineType) (recipe.Recipe, error) {
	r, err := recipe.NewRecipe(name, description, cuisine)
	if err != nil {
		return r, err
	}

	err = rs.recipes.Add(r)
	if err != nil {
		return r, err
	}
	slog.Info("recipe successfully added", "recipe_id", r.ID().String())
	return r, nil
}

func (rs RecipeService) GetRecipe(uuidStr string) (recipe.Recipe, error) {
	slog.Info("retrieving recipe..", "recipe_id", uuidStr)
	recipeUuid, err := uuid.FromBytes([]byte(uuidStr))
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("invalid uuid format: %w", err)
	}

	r, err := rs.recipes.Get(recipeUuid)
	if err != nil {
		return recipe.Recipe{}, err
	}

	return r, nil
}
