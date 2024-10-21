package domain

import "github.com/google/uuid"

type Variation struct {
	item      uuid.UUID
	variation string
}
