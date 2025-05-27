package auth

import (
	"context"

	"github.com/iankencruz/sabiflow/internal/platform/database"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int32) (*User, error)
	CreateUserOAuth(ctx context.Context, firstName, lastName, email string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetGroupPermissions(ctx context.Context, userID int32) ([]string, error)
}

type PgxUserRepository struct {
	DB database.DBTX
}

func NewUserRepository(db database.DBTX) *PgxUserRepository {
	return &PgxUserRepository{DB: db}
}

func (r *PgxUserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO auth.users (first_name, last_name, email, password)
		VALUES (@first_name, @last_name, @email, @password)
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.Password,
	}

	return r.DB.QueryRow(ctx, query, args).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *PgxUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM auth.users
		WHERE email = @email
	`

	args := pgx.NamedArgs{"email": email}

	var user User
	err := r.DB.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PgxUserRepository) GetByID(ctx context.Context, id int32) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM auth.users
		WHERE id = @id
	`

	args := pgx.NamedArgs{"id": id}

	var user User
	err := r.DB.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PgxUserRepository) CreateUserOAuth(ctx context.Context, firstName, lastName, email string) (*User, error) {
	query := `
		INSERT INTO auth.users (first_name, last_name, email, password, provider)
		VALUES (@first_name, @last_name, @email, '', 'google')
		RETURNING id, first_name, last_name, email, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
	}

	var user User
	err := r.DB.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PgxUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email
		FROM auth.users
		WHERE email = @email
	`
	args := pgx.NamedArgs{"email": email}
	var u User
	err := r.DB.QueryRow(ctx, query, args).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PgxUserRepository) GetGroupPermissions(ctx context.Context, userID int32) ([]string, error) {
	query := `
		SELECT p.code
		FROM auth.users u
		JOIN auth.permission_groups pg ON u.group_id = pg.id
		JOIN auth.group_permissions gp ON gp.group_id = pg.id
		JOIN auth.permissions p ON p.id = gp.permission_id
		WHERE u.id = @user_id
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	rows, err := r.DB.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		perms = append(perms, code)
	}

	return perms, nil
}
