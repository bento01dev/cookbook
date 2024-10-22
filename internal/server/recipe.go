package server

import (
	"net/http"

	"github.com/bento01dev/cookbook/internal/domain"
	"github.com/bento01dev/cookbook/internal/domain/recipe"
)

type recipeService interface {
	CreateRecipe(string, string, domain.CuisineType) (recipe.Recipe, error)
	GetRecipe(uuidStr string) (recipe.Recipe, error)
}

func handleCreateRecipe(rs recipeService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("recipe created.."))
	})
}

func handleGetRecipe(rs recipeService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		_, _ = rs.GetRecipe(id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("recipe get.."))
	})
}
