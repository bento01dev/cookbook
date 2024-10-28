package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bento01dev/cookbook/internal/domain"
	"github.com/bento01dev/cookbook/internal/domain/recipe"
	"github.com/bento01dev/cookbook/internal/stats"
)

type recipeService interface {
	CreateRecipe(context.Context, string, string, domain.CuisineType) (recipe.Recipe, error)
	GetRecipe(context.Context, string) (recipe.Recipe, error)
}

type errResponse struct {
	ErrCode int    `json:"err_code"`
	Msg     string `json:"msg"`
}

type cuisine string

func (c cuisine) MarshalText() ([]byte, error) {
	switch c {
	case japanese, french, spanish, indian, chinese, western:
		return []byte(c), nil
	default:
		return nil, fmt.Errorf("unknown type: %v", c)
	}
}

func (c *cuisine) UnmarshalText(data []byte) error {
	s := string(data)
	switch strings.ToLower(s) {
	case string(japanese):
		*c = japanese
		return nil
	case string(french):
		*c = french
		return nil
	case string(spanish):
		*c = spanish
		return nil
	case string(indian):
		*c = indian
		return nil
	case string(chinese):
		*c = chinese
		return nil
	case string(western):
		*c = western
		return nil
	default:
		return fmt.Errorf("unknown type: %s", s)
	}
}

func (c cuisine) ToDomain() domain.CuisineType {
	var dc domain.CuisineType
	switch c {
	case japanese:
		dc = domain.Japanese
	case french:
		dc = domain.French
	case spanish:
		dc = domain.Spanish
	case indian:
		dc = domain.Indian
	case chinese:
		dc = domain.Chinese
	case western:
		dc = domain.Western
	default:
		dc = domain.UnknownCuisine
	}
	return dc
}

func (c *cuisine) FromDomain(dc domain.CuisineType) {
	switch dc {
	case domain.Japanese:
		*c = japanese
	case domain.French:
		*c = french
	case domain.Spanish:
		*c = spanish
	case domain.Indian:
		*c = indian
	case domain.Chinese:
		*c = chinese
	case domain.Western:
		*c = western
	}
}

const (
	unknown  cuisine = "unknown"
	japanese cuisine = "japanese"
	french   cuisine = "french"
	spanish  cuisine = "spanish"
	indian   cuisine = "indian"
	chinese  cuisine = "chinese"
	western  cuisine = "western"
)

func handleCreateRecipe(rs recipeService, statsCollection *stats.StatsCollection) http.Handler {
	type request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Cuisine     cuisine `json:"cuisine"`
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
			slog.ErrorContext(ctx, "parsing request object failed")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResponse{ErrCode: 40002, Msg: "Issue in parsing request body"})
			return
		}

		cuisine := reqObj.Cuisine.ToDomain()
		if cuisine == domain.UnknownCuisine {
			slog.ErrorContext(ctx, "unknown cuisine in request", "cuisine", string(reqObj.Cuisine))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResponse{ErrCode: 40003, Msg: "Unknown cuisine"})
			return
		}

		slog.InfoContext(
			ctx,
			"new recipe request",
			slog.Group("payload",
				slog.String("name", reqObj.Name),
				slog.String("description", reqObj.Description),
				slog.String("cuisine", string(reqObj.Cuisine)),
			),
		)

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

		statsCollection.StatusOkInc("create_recipe")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{ID: recipe.ID().String(), Name: recipe.Name()})
	})
}

func handleGetRecipe(rs recipeService, statsCollection *stats.StatsCollection) http.Handler {

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
			ID          string  `json:"id,omitempty"`
			Name        string  `json:"name,omitempty"`
			Description string  `json:"description,omitempty"`
			Cuisine     cuisine `json:"cuisine,omitempty"`
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
		var c cuisine
		c.FromDomain(r.Cuisine())
		res.Item.Cuisine = c
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
				slog.ErrorContext(ctx, "get recipe exceeded timeout", "recipe_id", id)
				w.WriteHeader(http.StatusGatewayTimeout)
				errRes = errResponse{ErrCode: 50001, Msg: "service time out"}
			}

			if errors.Is(err, recipe.ErrRecipeNotFound) {
				slog.ErrorContext(ctx, "recipe not found for given id", "recipe_id", id)
				w.WriteHeader(http.StatusNotFound)
				errRes = errResponse{ErrCode: 40401, Msg: fmt.Sprintf("recipe not found for id: %s", id)}
			}

			if errors.Is(err, recipe.ErrInvalidID) {
				slog.ErrorContext(ctx, "invalid id format", "recipe_id", id)
				statsCollection.BadRequestInc("get_recipe")
				w.WriteHeader(http.StatusBadRequest)
				errRes = errResponse{ErrCode: 40001, Msg: fmt.Sprintf("invalid format for id: %s", id)}
			}

			json.NewEncoder(w).Encode(errRes)
			return
		}

		statsCollection.StatusOkInc("get_recipe")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(convertResponse(recipeRes))
	})
}
