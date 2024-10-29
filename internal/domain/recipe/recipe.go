package recipe

import (
	"errors"
	"time"

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
	createdAt   time.Time
	updatedAt   time.Time
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
		createdAt:   time.Now().UTC(),
	}, nil
}

func (r Recipe) ID() uuid.UUID {
	return r.item.ID
}

func (r Recipe) Name() string {
	return r.item.Name
}

func (r Recipe) Description() string {
	return r.item.Description
}

func (r Recipe) Cuisine() domain.CuisineType {
	return r.item.Cuisine
}

func (r Recipe) Ingredients() []*domain.Ingredient {
	return r.ingredients
}

func (r Recipe) Variations() []string {
	var res []string
	for _, v := range r.variations {
		res = append(res, v.Variation())
	}
	return res
}

func (r Recipe) Prep() []domain.Prep {
	return r.prepSteps
}

func (r Recipe) Steps() []domain.Step {
	return r.steps
}

func (r Recipe) CreatedAt() time.Time {
	return r.createdAt
}
