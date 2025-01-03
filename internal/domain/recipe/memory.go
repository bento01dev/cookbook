package recipe

import (
	"context"
	"sync"
	"time"

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

type wrapper struct {
	recipe Recipe
	err    error
}

func (mr *MemoryRepository) Get(ctx context.Context, id uuid.UUID) (Recipe, error) {
	select {
	case <-ctx.Done():
		return Recipe{}, ctx.Err()
	case res := <-mr.get(id):
		return res.recipe, res.err
	}
}

// this really isnt needed. just for fun
func (mr *MemoryRepository) get(id uuid.UUID) <-chan wrapper {
	ch := make(chan wrapper)
	go func() {
		time.Sleep(200 * time.Millisecond)
		mr.mu.Lock()
		defer mr.mu.Unlock()
		if recipe, ok := mr.recipes[id]; ok {
			ch <- wrapper{recipe: recipe}
			return
		}
		ch <- wrapper{err: ErrRecipeNotFound}
	}()
	return ch
}

func (mr *MemoryRepository) Add(ctx context.Context, recipe Recipe) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-mr.add(recipe):
		return res
	}
}

func (mr *MemoryRepository) add(recipe Recipe) <-chan error {
	ch := make(chan error)
	go func() {
		time.Sleep(200 * time.Millisecond)
		mr.mu.Lock()
		defer mr.mu.Unlock()
		mr.recipes[recipe.ID()] = recipe
		ch <- nil
	}()
	return ch
}

func (mr *MemoryRepository) Update(ctx context.Context, recipe Recipe) (Recipe, error) {
	if _, ok := mr.recipes[recipe.ID()]; ok {
		return recipe, ErrRecipeExists
	}
	mr.mu.Lock()
	mr.recipes[recipe.ID()] = recipe
	mr.mu.Unlock()
	return recipe, nil
}

func (mr *MemoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := mr.recipes[id]; !ok {
		return ErrRecipeNotFound
	}
	delete(mr.recipes, id)
	return nil
}
