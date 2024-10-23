package config

import (
	"strconv"
	"time"
)

type Config struct {
	Host               string
	Port               string
	InMemory           bool
	GetRecipeTimeoutMs time.Duration
}

func NewConfig(getEnv func(string) string) (Config, error) {
	var err error

	host := "127.0.0.1"
	if getEnv("HTTP_HOST") != "" {
		host = getEnv("HTTP_HOST")
	}

	port := "8080"
	if getEnv("HTTP_PORT") != "" {
		port = getEnv("HTTP_PORT")
	}

	var inMemory bool
	if v := getEnv("MEMORY_REPO"); v != "" {
		inMemory, err = strconv.ParseBool(v)
		if err != nil {
			return Config{}, err
		}
	}

	var getRecipeTimeoutMs = 1000 * time.Millisecond
	if v := getEnv("GET_RECIPE_TIMEOUT"); v != "" {
		getRecipeTimeoutMs, err = time.ParseDuration(v)
		if err != nil {
			return Config{}, err
		}
	}

	return Config{
		Host:               host,
		Port:               port,
		InMemory:           inMemory,
		GetRecipeTimeoutMs: getRecipeTimeoutMs,
	}, err
}
