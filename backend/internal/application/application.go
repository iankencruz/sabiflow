package application

import (
	"context"
	"log/slog"
	"os"

	"github.com/iankencruz/sabiflow/internal/auth"
	"github.com/iankencruz/sabiflow/internal/shared/logger"
	"github.com/iankencruz/sabiflow/internal/shared/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a new Application struct
type Application struct {
	Config      *Config
	DB          *pgxpool.Pool
	Logger      *slog.Logger
	AuthHandler *auth.AuthHandler
}

func NewApplication() (*Application, error) {

	cfg := LoadConfig()

	db, err := pgxpool.New(context.Background(), cfg.DB_DSN)
	if err != nil {
		return nil, err
	}

	// Initialize the Logger
	log := logger.New(os.Getenv(cfg.Env))

	// Initialize the Session Manager
	sessionManager := sessions.NewManager()

	userRepo := auth.NewUserRepository(db)
	authService := &auth.AuthServiceImpl{Repo: userRepo}
	authHandler := &auth.AuthHandler{
		Service:        authService,
		SessionManager: sessionManager,
	}

	return &Application{
		Config:      cfg,
		DB:          db,
		Logger:      log,
		AuthHandler: authHandler,
	}, nil
}
