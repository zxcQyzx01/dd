package config

type Config struct {
    GRPCPort string
    DB       DBConfig
}

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

func New() *Config {
    return &Config{
        GRPCPort: ":50053",
        DB: DBConfig{
            Host:     "postgres",
            Port:     "5432",
            User:     "postgres",
            Password: "postgres",
            DBName:   "userdb",
        },
    }
} 