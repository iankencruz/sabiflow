package application

import (
	"context"
	"log/slog"

	"github.com/iankencruz/sabiflow/internal/auth"
	"github.com/iankencruz/sabiflow/internal/shared/logger"
	"github.com/iankencruz/sabiflow/internal/shared/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a new Application struct
type Application struct {
	Config         *Config
	DB             *pgxpool.Pool
	Logger         *slog.Logger
	AuthHandler    *auth.AuthHandler
	SessionManager *sessions.Manager
	UserRepo       auth.UserRepository
}

func NewApplication() (*Application, error) {

	cfg := LoadConfig()

	db, err := pgxpool.New(context.Background(), cfg.DB_DSN)
	if err != nil {
		return nil, err
	}

	// Initialize the Logger
	log := logger.New(cfg.Env)
	// Initialize the Session Manager
	sessionManager := sessions.NewManager(db)

	userRepo := auth.NewUserRepository(db)
	authService := &auth.AuthServiceImpl{Repo: userRepo}
	authHandler := &auth.AuthHandler{
		Service:        authService,
		SessionManager: sessionManager,
		Logger:         log,
	}

	return &Application{
		Config:         cfg,
		DB:             db,
		Logger:         log,
		AuthHandler:    authHandler,
		SessionManager: sessionManager,
		UserRepo:       userRepo,
	}, nil
}
