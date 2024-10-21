package domain

import "github.com/google/uuid"

type CuisineType int

const (
	UnknownCuisine = iota
	Japanese
	French
	Spanish
	Indian
	Chinese
	Western
	African
)

type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
	Cuisine     CuisineType
}
