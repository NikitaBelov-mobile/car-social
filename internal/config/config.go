package config

import (
	"os"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string
	Mode string // gin mode (debug/release)
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
