package domain

import "github.com/google/uuid"

type IngredientType int

const (
	UnknownIngredient = iota
	Vegetable
	Fruit
	Poultry
	Fish
	Condiments
)

type Ingredient struct {
	ID          uuid.UUID
	Name        string
	Description string
	Type        IngredientType
}
