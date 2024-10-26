package domain

import "github.com/google/uuid"

type Prep struct {
	ingredient uuid.UUID
	action     string
	index      int
}

func (p Prep) Ingredient() string {
	return p.ingredient.String()
}

func (p Prep) Action() string {
	return p.action
}
