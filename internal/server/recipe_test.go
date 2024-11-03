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
)

func getEnv(key string) string {
	switch key {
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
	default:
		return ""
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

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go Run(ctx, io.Discard, io.Discard, []string{}, getEnv)
	if err := waitForReady(ctx, 10*time.Second, "http://localhost:8080/healthz"); err != nil {
		os.Exit(1)
	}
	fmt.Println("service started successfully..")
	os.Exit(m.Run())
}

func TestCreateRecipe(t *testing.T) {
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
		context.Background(),
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
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetRecipe(t *testing.T) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(
		context.Background(),
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
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
