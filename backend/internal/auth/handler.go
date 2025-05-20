package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

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
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUserOAuth(ctx context.Context, fullName, email string) (*User, error)
	GetUserByID(ctx context.Context, id int32) (*User, error)
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

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := GetGoogleAuthURL("state-token") // Replace with dynamic state in future
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.URL.Query().Get("code")
	if code == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing code", nil)
		return
	}

	token, err := ExchangeCode(ctx, code)
	if err != nil {
		h.Logger.Error("OAuth token exchange failed", "err", err)
		response.WriteJSON(w, http.StatusInternalServerError, "OAuth exchange failed", nil)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		h.Logger.Error("Failed to get user info", "err", err)
		response.WriteJSON(w, http.StatusInternalServerError, "Google user info failed", nil)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &googleUser); err != nil {
		h.Logger.Error("User info unmarshal failed", "err", err)
		response.WriteJSON(w, http.StatusInternalServerError, "Invalid user info", nil)
		return
	}

	user, err := h.Service.GetUserByEmail(ctx, googleUser.Email)
	if err != nil {
		user, err = h.Service.CreateUserOAuth(ctx, googleUser.Name, googleUser.Email)
		if err != nil {
			h.Logger.Error("Failed to create OAuth user", "err", err)
			response.WriteJSON(w, http.StatusInternalServerError, "User creation failed", nil)
			return
		}
	}

	if err := h.SessionManager.SetUserID(w, r, user.ID); err != nil {
		h.Logger.Error("Failed to store session", "err", err)
		response.WriteJSON(w, http.StatusInternalServerError, "Session error", nil)
		return
	}

	http.Redirect(w, r, os.Getenv("FRONTEND_SUCCESS_REDIRECT_URL"), http.StatusSeeOther)
}

func (h *AuthHandler) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.SessionManager.GetUserID(r)
	fmt.Println("🧪 incoming /user/me request, userID:", userID, "err:", err)
	if err != nil || userID == 0 {
		response.WriteJSON(w, http.StatusUnauthorized, "Not authenticated", nil)
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), userID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "User lookup failed", nil)
		return
	}

	response.WriteJSON(w, http.StatusOK, "User details", user)
}
