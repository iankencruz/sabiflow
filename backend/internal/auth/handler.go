package auth

import (
	"net/http"
	"strings"

	"github.com/iankencruz/sabiflow/internal/application"
	"github.com/iankencruz/sabiflow/internal/auth"
	"github.com/iankencruz/sabiflow/internal/logger"
	"github.com/iankencruz/sabiflow/internal/response"
	"github.com/iankencruz/sabiflow/internal/sessions"
	"github.com/iankencruz/sabiflow/internal/validators"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUser(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

		// 1. Decode JSON
		if err := response.DecodeJSON(w, r, &req); err != nil {
			logger.WriteJSONError(w, app.Logger, http.StatusBadRequest, "Invalid JSON input", err)
			return
		}
		// Ensure all input is trimmed and lowercased
		req.FirstName = strings.ToLower(strings.TrimSpace(req.FirstName))
		req.LastName = strings.ToLower(strings.TrimSpace(req.LastName))
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		// 2. Validate input
		v := validators.New()

		v.Require("firstname", req.FirstName)
		v.Require("lastname", req.LastName)
		v.Require("email", req.Email)
		v.Require("password", req.Password)

		v.MatchPattern("email", req.Email, validators.EmailRX, "Invalid Email address")

		v.Check("password", len(req.Password) >= 8, "Password must be at least 8 characters long")
		v.Check("password", validators.UppercaseRX.MatchString(req.Password), "Password must contain at least one uppercase letter")
		v.Check("password", validators.LowercaseRX.MatchString(req.Password), "Password must contain at least one lowercase letter")
		v.Check("password", validators.NumberRX.MatchString(req.Password), "Password must contain at least one number")
		v.Check("password", validators.SpecialRX.MatchString(req.Password), "Password must contain at least one special character")

		if !v.Valid() {
			response.WriteJSON(w, http.StatusUnprocessableEntity, "Validation error", map[string]interface{}{
				"errors": v.Errors,
			})
			return
		}

		// 3. Hash password
		hashed, err := auth.HashPassword(req.Password)
		if err != nil {
			logger.WriteJSONError(w, app.Logger, http.StatusInternalServerError, "Password hashing failed", err)
			return
		}

		// 4. Create user in DB
		user, err := app.Auth.CreateUser(r.Context(), auth.CreateUserParams{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  hashed,
		})
		if err != nil {
			logger.WriteJSONError(w, app.Logger, http.StatusInternalServerError, "Failed to create user", err)
			return
		}

		// 5. Set session
		if err := sessions.SetUserID(w, r, user.ID); err != nil {
			logger.WriteJSONError(w, app.Logger, http.StatusInternalServerError, "Failed to set session", err)
			return
		}

		// 6. Return success response
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
		}
		response.WriteJSON(w, http.StatusOK, "Registered successfully", data)
	}
}

// LoginUser handles user login. It verifies the user's credentials and sets a session cookie.
func LoginUser(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := response.DecodeJSON(w, r, &req); err != nil {
			logger.WriteJSONError(w, app.Logger, http.StatusBadRequest, "Invalid input", err)
			return
		}

		user, err := app.Auth.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			logger.WriteJSONError(w, app.Logger, http.StatusUnauthorized, "User not found", err)
			return
		}

		if err := auth.ComparePassword(user.Password, req.Password); err != nil {
			logger.WriteJSONError(w, app.Logger, http.StatusUnauthorized, "Invalid password", err)
			return
		}

		if err := sessions.SetUserID(w, r, user.ID); err != nil {
			logger.WriteJSONError(w, app.Logger, http.StatusInternalServerError, "Failed to set session", err)
			return
		}

		data := map[string]interface{}{
			"user": map[string]interface{}{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
		}

		response.WriteJSON(w, http.StatusOK, "Logged in successfully", data)
	}
}
