package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bento01dev/cookbook/internal/domain"
	"github.com/bento01dev/cookbook/internal/domain/recipe"
)

type recipeService interface {
	CreateRecipe(string, string, domain.CuisineType) (recipe.Recipe, error)
	GetRecipe(context.Context, string) (recipe.Recipe, error)
}

func handleCreateRecipe(rs recipeService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("recipe created.."))
	})
}

func handleGetRecipe(rs recipeService) http.Handler {
	type errResponse struct {
		ErrCode int    `json:"err_code"`
		Msg     string `json:"msg"`
	}

	type ingredient struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}

	type prep struct {
		IngredientID string `json:"ingredient_id"`
		Action       string `json:"action"`
	}

	type step struct {
		IngredientID string  `json:"ingredient_id"`
		Action       string  `json:"action"`
		Temperature  float64 `json:"temperature"`
	}

	type recipeResponse struct {
		Item struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Cuisine     string `json:"cuisine"`
		} `json:"item"`
		Ingredients []ingredient `json:"ingredients"`
		Variations  []string     `json:"variations"`
		Prep        []prep       `json:"prep"`
		Steps       []step       `json:"steps"`
	}

	convertResponse := func(r recipe.Recipe) recipeResponse {
		var res recipeResponse
		res.Item.ID = r.ID().String()
		res.Item.Name = r.Name()
		res.Item.Description = r.Description()
		for _, v := range r.Ingredients() {
			res.Ingredients = append(res.Ingredients, ingredient{ID: v.ID.String(), Name: v.Name, Type: string(v.Type)})
		}
		res.Variations = r.Variations()
		return res
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ctx := r.Context()
		recipeRes, err := rs.GetRecipe(ctx, id)

		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			var errRes errResponse
			if errors.Is(err, context.DeadlineExceeded) {
				w.WriteHeader(http.StatusGatewayTimeout)
				errRes = errResponse{ErrCode: 50001, Msg: "service time out"}
			}

			if errors.Is(err, recipe.ErrRecipeNotFound) {
				w.WriteHeader(http.StatusNotFound)
				errRes = errResponse{ErrCode: 40401, Msg: fmt.Sprintf("recipe not found for id: %s", id)}
			}

			if errors.Is(err, recipe.ErrInvalidID) {
				w.WriteHeader(http.StatusBadRequest)
				errRes = errResponse{ErrCode: 40001, Msg: fmt.Sprintf("invalid format for id: %s", id)}
			}

			json.NewEncoder(w).Encode(errRes)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(convertResponse(recipeRes))
	})
}
