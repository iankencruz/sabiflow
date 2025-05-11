package auth

import (
	"encoding/json"
	"net/http"

	"github.com/iankencruz/sabiflow/internal/sessions"
)

// LoginRequest is the expected body from frontend
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser validates credentials and starts a session
func (q *Queries) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := q.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := ComparePassword(user.Password, req.Password); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := sessions.SetUserSession(w, r, user.ID); err != nil {
		http.Error(w, "Failed to start session", http.StatusInternalServerError)
		return
	}

	// Return success JSON
	resp := map[string]any{
		"status":  "success",
		"message": "Logged in successfully",
		"data": map[string]any{
			"user": map[string]any{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
