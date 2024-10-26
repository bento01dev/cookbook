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
	CreateRecipe(context.Context, string, string, domain.CuisineType) (recipe.Recipe, error)
	GetRecipe(context.Context, string) (recipe.Recipe, error)
}

type errResponse struct {
	ErrCode int    `json:"err_code"`
	Msg     string `json:"msg"`
}

func handleCreateRecipe(rs recipeService) http.Handler {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Cuisine     string `json:"cuisine"`
	}

	type response struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Add("Content-Type", "application/json")

		var reqObj request
		err := json.NewDecoder(r.Body).Decode(&reqObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResponse{ErrCode: 40002, Msg: "Issue in parsing request body"})
			return
		}

		var cuisine domain.CuisineType = domain.Indian
		// switch reqObj.Cuisine {
		// case string(domain.African):
		// 	cuisine = domain.African
		// case string(domain.Indian):
		// 	cuisine = domain.Indian
		// case string(domain.Japanese):
		// 	cuisine = domain.Japanese
		// case string(domain.French):
		// 	cuisine = domain.French
		// case string(domain.Spanish):
		// 	cuisine = domain.Spanish
		// case string(domain.Chinese):
		// 	cuisine = domain.Chinese
		// case string(domain.Western):
		// 	cuisine = domain.Western
		// }
		// if cuisine == domain.UnknownCuisine {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode(errResponse{ErrCode: 40003, Msg: "Unknown cuisine"})
		// 	return
		// }

		recipe, err := rs.CreateRecipe(ctx, reqObj.Name, reqObj.Description, cuisine)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				w.WriteHeader(http.StatusGatewayTimeout)
				json.NewEncoder(w).Encode(errResponse{ErrCode: 50001, Msg: "service time out"})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errResponse{ErrCode: 50002, Msg: "Uncaught exception"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{ID: recipe.ID().String(), Name: recipe.Name()})
	})
}

func handleGetRecipe(rs recipeService) http.Handler {

	type ingredient struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
	}

	type prep struct {
		IngredientID string `json:"ingredient_id,omitempty"`
		Action       string `json:"action,omitempty"`
	}

	type step struct {
		IngredientID string  `json:"ingredient_id,omitempty"`
		Action       string  `json:"action,omitempty"`
		Temperature  float64 `json:"temperature,omitempty"`
	}

	type recipeResponse struct {
		Item struct {
			ID          string `json:"id,omitempty"`
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
			Cuisine     string `json:"cuisine,omitempty"`
		} `json:"item"`
		Ingredients []ingredient `json:"ingredients,omitempty"`
		Variations  []string     `json:"variations,omitempty"`
		Prep        []prep       `json:"prep,omitempty"`
		Steps       []step       `json:"steps,omitempty"`
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
		for _, p := range r.Prep() {
			res.Prep = append(res.Prep, prep{IngredientID: p.Ingredient(), Action: p.Action()})
		}
		for _, s := range r.Steps() {
			res.Steps = append(res.Steps, step{IngredientID: s.Ingredient(), Action: s.Action(), Temperature: s.Temperature()})
		}
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
