package config

import "os"

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Server struct {
		Port int
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5435")
	cfg.Database.User = getEnv("DB_USER", "user")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.Name = getEnv("DB_NAME", "payments_db")
	cfg.Server.Port = 8004

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
