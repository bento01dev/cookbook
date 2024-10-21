package domain

import "github.com/google/uuid"

type Prep struct {
	ingredient uuid.UUID
	action     string
	index      int
}
