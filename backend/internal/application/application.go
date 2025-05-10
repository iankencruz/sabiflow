package application

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a new Application struct
type Application struct {
	DB *pgxpool.Pool
}
