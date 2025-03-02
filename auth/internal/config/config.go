package config

type Config struct {
    GRPCPort    string
    UserService string
    JWTSecret   string
}

func New() *Config {
    return &Config{
        GRPCPort:    ":50051",
        UserService: "user:50053",
        JWTSecret:   "your-secret-key",
    }
} 