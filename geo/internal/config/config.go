package config

type Config struct {
    GRPCPort string
}

func New() *Config {
    return &Config{
        GRPCPort: ":50052",
    }
} 