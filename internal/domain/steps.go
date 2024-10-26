package domain

import "github.com/google/uuid"

type Step struct {
	ingredient  uuid.UUID
	action      string
	temperature float64
}

func (s Step) Action() string {
	return s.action
}

func (s Step) Ingredient() string {
	return s.ingredient.String()
}

func (s Step) Temperature() float64 {
	return s.temperature
}
