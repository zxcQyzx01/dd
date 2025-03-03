package config

import (
	"os"
)

type Config struct {
	GRPCPort   string
	RedisAddr  string
	JWTSecret  string
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func New() *Config {
	return &Config{
		GRPCPort:   getEnvOrDefault("GRPC_PORT", ":50052"),
		RedisAddr:  getEnvOrDefault("REDIS_ADDR", "redis:6379"),
		JWTSecret:  getEnvOrDefault("JWT_SECRET", "your-secret-key"),
	}
} 