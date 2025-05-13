package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	const sql = `
		INSERT INTO auth.users (
			first_name, last_name, email, password
		) VALUES (
			@first_name, @last_name, lower(@email), @password
		)
		RETURNING id, first_name, last_name, email, password, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"first_name": arg.FirstName,
		"last_name":  arg.LastName,
		"email":      arg.Email,
		"password":   arg.Password,
	}

	var u User
	err := q.DB.QueryRow(ctx, sql, args).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	return u, err
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	const sql = `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM auth.users
		WHERE email = lower(@email)
	`

	args := pgx.NamedArgs{
		"email": email,
	}

	var u User
	err := q.DB.QueryRow(ctx, sql, args).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)

	return u, err
}

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	const sql = `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM auth.users
		ORDER BY created_at DESC
	`

	rows, err := q.DB.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}
