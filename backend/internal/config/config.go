package config

import (
	"os"
)

type Config struct {
	Environment string
	Port        string
	WorkerURL   string
	WorkerKey   string
}

func Load() *Config {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	workerURL := os.Getenv("WORKER_API_URL")
	if workerURL == "" {
		workerURL = "https://api.nawthtech.com"
	}

	workerKey := os.Getenv("WORKER_API_KEY")

	return &Config{
		Environment: env,
		Port:        port,
		WorkerURL:   workerURL,
		WorkerKey:   workerKey,
	}
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}