package recipe

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	recipes map[uuid.UUID]Recipe
	mu      sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		recipes: make(map[uuid.UUID]Recipe),
	}
}
func (mr *MemoryRepository) Get(ctx context.Context, id uuid.UUID) (Recipe, error) {
	if recipe, ok := mr.recipes[id]; ok {
		return recipe, nil
	}
	return Recipe{}, ErrRecipeNotFound
}

func (mr *MemoryRepository) Add(recipe Recipe) error {
	if _, ok := mr.recipes[recipe.ID()]; ok {
		return ErrRecipeExists
	}
	mr.mu.Lock()
	mr.recipes[recipe.ID()] = recipe
	mr.mu.Unlock()
	return nil
}

func (mr *MemoryRepository) Update(recipe Recipe) (Recipe, error) {
	if _, ok := mr.recipes[recipe.ID()]; ok {
		return recipe, ErrRecipeExists
	}
	mr.mu.Lock()
	mr.recipes[recipe.ID()] = recipe
	mr.mu.Unlock()
	return recipe, nil
}

func (mr *MemoryRepository) Delete(id uuid.UUID) error {
	if _, ok := mr.recipes[id]; !ok {
		return ErrRecipeNotFound
	}
	delete(mr.recipes, id)
	return nil
}
