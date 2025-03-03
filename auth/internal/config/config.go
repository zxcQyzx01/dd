package config

import (
	"os"
)

type Config struct {
	GRPCPort    string
	UserService string
	JWTSecret   string
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func New() *Config {
	return &Config{
		GRPCPort:    getEnvOrDefault("GRPC_PORT", ":50051"),
		UserService: getEnvOrDefault("USER_SERVICE", "user1:50053"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "your-secret-key"),
	}
} 