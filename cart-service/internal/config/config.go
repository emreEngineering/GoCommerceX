package config

import (
	"fmt"
	"os"
)

type Config struct {
	GRPCPort      string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	CartKeyPrefix string
}

func Load() *Config {
	return &Config{
		GRPCPort:      getEnv("GRPC_PORT", "50055"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		CartKeyPrefix: getEnv("CART_KEY_PREFIX", "cart:"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		var parsed int
		_, err := fmt.Sscanf(value, "%d", &parsed)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
