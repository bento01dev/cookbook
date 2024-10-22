package config

import "strconv"

type Config struct {
	Host               string
	Port               string
	InMemory           bool
	GetRecipeTimeoutMs int
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

	var getRecipeTimeoutMs int = 1000
	if v := getEnv("GET_RECIPE_TIMEOUT_MS"); v != "" {
		getRecipeTimeoutMs, err = strconv.Atoi(v)
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
