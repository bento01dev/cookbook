package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

var mongoURL string
var mongoContainer *mongodb.MongoDBContainer

func getEnv(mongoURL string) func(string) string {
	return func(s string) string {
		switch s {
		case "DB_TYPE":
			return "mongo"
		case "LOG_LEVEL":
			return "info"
		case "SERVICE_NAME":
			return "test_cookbook"
		case "ENV":
			return "test"
		case "HOST_IP":
			return "127.0.0.1"
		case "HTTP_HOST":
			return "127.0.0.1"
		case "HTTP_PORT":
			return "8080"
		case "MONGO_DB_URL":
			return mongoURL
		case "MONGO_DB":
			return "cookbook"
		case "RECIPE_COLLECTION":
			return "recipe"
		default:
			return ""
		}
	}
}

// waitForReady calls the specified endpoint until it gets a 200
// response or until the context is cancelled or the timeout is
// reached.
func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		fmt.Println("checkign if server is up..")
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			// wait a little while between checks
			time.Sleep(250 * time.Millisecond)
		}
	}
}

type RecipeTestSuite struct {
	suite.Suite
	mongoContainer *mongodb.MongoDBContainer
	mongoURL       string
	ctx            context.Context
	cancel         context.CancelFunc
}

func (suite *RecipeTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	suite.ctx = ctx
	suite.cancel = cancel
	container, err := mongodb.Run(ctx, "mongo:8")
	if err != nil {
		fmt.Println("mongodb setup failed:", err.Error())
		os.Exit(1)
	}
	suite.mongoContainer = container
	url, err := container.ConnectionString(ctx)
	if err != nil {
		fmt.Println("mongodb setup failed:", err.Error())
		os.Exit(1)
	}
	suite.mongoURL = url
	go Run(ctx, io.Discard, io.Discard, []string{}, getEnv(suite.mongoURL))
	if err := waitForReady(ctx, 10*time.Second, "http://localhost:8080/healthz"); err != nil {
		os.Exit(1)
	}
}

func (suite *RecipeTestSuite) TearDownSuite() {
	if err := suite.mongoContainer.Terminate(suite.ctx); err != nil {
		fmt.Println("issue in stopping mongo:", err.Error())
		os.Exit(1)
	}
	suite.cancel()
}

func TestRecipe(t *testing.T) {
	suite.Run(t, new(RecipeTestSuite))
}

func (suite *RecipeTestSuite) TestCreateRecipe() {
	t := suite.T()

	client := http.Client{}
	body := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Cuisine     string `json:"cuisine"`
	}{
		Name:        "test_recipe",
		Description: "testing creating recipe",
		Cuisine:     "spanish",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(
		suite.ctx,
		http.MethodPost,
		"http://localhost:8080/recipe",
		&buf,
	)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
}

func (suite *RecipeTestSuite) TestGetRecipe() {
	t := suite.T()

	client := http.Client{}
	req, err := http.NewRequestWithContext(
		suite.ctx,
		http.MethodGet,
		"http://localhost:8080/recipe/abc",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
}

func (suite *RecipeTestSuite) TestNotFoundRecipe() {
	t := suite.T()

	client := http.Client{}
	req, err := http.NewRequestWithContext(
		suite.ctx,
		http.MethodGet,
		"http://localhost:8080/recipe/3f2a4244-d10b-464f-9985-b63fda452fec",
		nil,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("error in creating request: %w", err))
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(fmt.Errorf("error in sending request: %w", err))
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
}
