package application

import (
	"context"
	"log/slog"
	"os"

	"github.com/iankencruz/sabiflow/internal/auth"
	"github.com/iankencruz/sabiflow/internal/core/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a new Application struct
type Application struct {
	Config *Config
	DB     *pgxpool.Pool
	Logger *slog.Logger
	Auth   *auth.Queries
}

func NewApplication() (*Application, error) {

	cfg := LoadConfig()

	db, err := pgxpool.New(context.Background(), cfg.DB_DSN)
	if err != nil {
		return nil, err
	}

	log := logger.New(os.Getenv(cfg.Env))

	return &Application{
		Config: cfg,
		DB:     db,
		Logger: log,
		Auth:   auth.New(db),
	}, nil
}
