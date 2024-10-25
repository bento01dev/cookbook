package recipe

import (
	"errors"

	"github.com/bento01dev/cookbook/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrInvalidItemName    = errors.New("invalid name for item")
	ErrRecipeNotFound     = errors.New("recipe not found for given id")
	ErrRecipeUpdateFailed = errors.New("recipe could not be updated")
	ErrRecipeExists       = errors.New("recipe already exists for given id")
	ErrInvalidID          = errors.New("invalid id format")
)

type Recipe struct {
	item        *domain.Item
	ingredients []*domain.Ingredient
	variations  []domain.Variation
	prepSteps   []domain.Prep
	steps       []domain.Step
	pairings    []domain.Pairing
}

func NewRecipe(name string, description string, cuisine domain.CuisineType) (Recipe, error) {
	if name == "" {
		return Recipe{}, ErrInvalidItemName
	}

	item := &domain.Item{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Cuisine:     cuisine,
	}

	return Recipe{
		item:        item,
		ingredients: make([]*domain.Ingredient, 0),
		variations:  make([]domain.Variation, 0),
		prepSteps:   make([]domain.Prep, 0),
		steps:       make([]domain.Step, 0),
		pairings:    make([]domain.Pairing, 0),
	}, nil
}

func (r Recipe) ID() uuid.UUID {
	return r.item.ID
}

func (r Recipe) Name() string {
	return r.item.Name
}
