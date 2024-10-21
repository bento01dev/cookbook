package domain

import "github.com/google/uuid"

type Step struct {
	ingredient  uuid.UUID
	action      string
	temperature float64
}
