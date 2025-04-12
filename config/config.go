package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	ServerPort  int
	Environment string
}

func LoadConfig() *Config {
	config := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "root:gggggggg@tcp(localhost:3306)/twitter_monitor"),
		ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
