package recipe

import "github.com/bento01dev/cookbook/internal/domain"

type Recipe struct {
	item        *domain.Item
	ingredients []*domain.Ingredient
	variations  []domain.Variation
	prepSteps   []domain.Prep
	steps       []domain.Step
	pairings    []domain.Pairing
}
