package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthServiceImpl implements AuthService.
type AuthServiceImpl struct {
	Repo UserRepository
}

// Register creates a new user with hashed password.
func (s *AuthServiceImpl) Register(ctx context.Context, firstName, lastName, email, password string) (*User, error) {
	existing, _ := s.Repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashed),
	}

	if err := s.Repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login checks the email and password match.
func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := s.Repo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// Logout is currently a no-op.
func (s *AuthServiceImpl) Logout(ctx context.Context) error {
	return nil
}

func (s *AuthServiceImpl) CreateUserOAuth(ctx context.Context, firstName, lastName, email string) (*User, error) {
	return s.Repo.CreateUserOAuth(ctx, firstName, lastName, email)
}

func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id int32) (*User, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *AuthServiceImpl) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.Repo.GetByEmail(ctx, email)
}
