package auth

import (
	"github.com/iankencruz/sabiflow/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        int32              `db:"id" json:"id"`
	FirstName string             `db:"first_name" json:"first_name"`
	LastName  string             `db:"last_name" json:"last_name"`
	Email     string             `db:"email" json:"email"`
	Password  string             `db:"password" json:"password"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
}

type CreateUserParams struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Password  string `db:"password"`
}
type Queries struct {
	DB database.DBTX
}

func New(db database.DBTX) *Queries {
	return &Queries{DB: db}
}
