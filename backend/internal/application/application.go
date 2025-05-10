package application

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a new Application struct
type Application struct {
	Config *Config
	DB     *pgxpool.Pool
}

func NewApplication() (*Application, error) {

	cfg := LoadConfig()

	db, err := pgxpool.New(context.Background(), cfg.DB_DSN)
	if err != nil {
		return nil, err
	}

	return &Application{
		Config: cfg,
		DB:     db,
	}, nil
}
