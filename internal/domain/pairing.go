package domain

import "github.com/google/uuid"

type Pairing struct {
	base        uuid.UUID
	with        uuid.UUID
	description string
}
