package application

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration...
type Config struct {
	Env    string
	Port   string
	DB_DSN string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Env:    getEnv("APP_ENV", "development"),
		Port:   getEnv("PORT", "8080"),
		DB_DSN: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/sabiflow?sslmode=disable"),
	}

	if cfg.DB_DSN == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
