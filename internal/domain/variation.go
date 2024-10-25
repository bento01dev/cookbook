package domain

import "github.com/google/uuid"

type Variation struct {
	item      uuid.UUID
	variation string
}

func (v Variation) Variation() string {
	return v.variation
}
