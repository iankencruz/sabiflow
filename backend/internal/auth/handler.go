package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/iankencruz/sabiflow/internal/shared/errors"
	"github.com/iankencruz/sabiflow/internal/shared/response"
	"github.com/iankencruz/sabiflow/internal/shared/sessions"
	"github.com/iankencruz/sabiflow/internal/shared/validators"
)

// AuthService defines the business logic interface for authentication.
type AuthService interface {
	Register(ctx context.Context, firstName, lastName, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	Logout(ctx context.Context) error
}

// AuthHandler handles HTTP requests for authentication-related operations.
type AuthHandler struct {
	Service        AuthService
	SessionManager *sessions.Manager
	Logger         *slog.Logger
}

// RegisterHandler handles user registration.
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Firstname string `json:"firstName"`
		Lastname  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errResp := errors.BadRequest("Invalid request payload")
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	v := validators.New()
	v.Require("firstName", input.Firstname)
	v.Require("lastName", input.Lastname)
	v.Require("email", input.Email)
	v.MatchPattern("email", input.Email, validators.EmailRX, "Must be a valid email address")
	v.Require("password", input.Password)
	v.MatchPattern("password", input.Password, validators.UppercaseRX, "Must include at least one uppercase letter")
	v.MatchPattern("password", input.Password, validators.NumberRX, "Must include at least one number")

	if !v.Valid() {
		errResp := errors.BadRequest("Validation failed")
		response.WriteJSON(w, errResp.Code, errResp.Message, v.Errors)
		return
	}

	user, err := h.Service.Register(r.Context(), input.Firstname, input.Lastname, input.Email, input.Password)
	if err != nil {
		errResp := errors.Internal(err.Error())
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	if err := h.SessionManager.SetUserID(w, r, user.ID); err != nil {
		errResp := errors.Internal("Failed to set session")
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	response.WriteJSON(w, http.StatusCreated, "User registered", map[string]any{"user": user})
}

// LoginHandler handles user login.
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errResp := errors.BadRequest("Invalid login payload")
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	v := validators.New()
	v.Require("email", input.Email)
	v.MatchPattern("email", input.Email, validators.EmailRX, "Must be a valid email address")
	v.Require("password", input.Password)

	if !v.Valid() {
		errResp := errors.BadRequest("Validation failed")
		response.WriteJSON(w, errResp.Code, errResp.Message, v.Errors)
		return
	}

	// Use the service to Check the credentials
	user, err := h.Service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		errResp := errors.Internal(err.Error())
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	if err := h.SessionManager.SetUserID(w, r, user.ID); err != nil {
		errResp := errors.Internal("Failed to set session")
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	response.WriteJSON(w, http.StatusOK, "Logged in", map[string]any{"user": user})
}

// LogoutHandler clears the user session.
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LogoutHandler called")

	if h.Logger != nil {
		h.Logger.Info("LogoutHandler called")
	}

	if err := h.SessionManager.Clear(w, r); err != nil {
		errResp := errors.Internal("Failed to clear session")
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	response.WriteJSON(w, http.StatusOK, "Logged out successfully", nil)
}
