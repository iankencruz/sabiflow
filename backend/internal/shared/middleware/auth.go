package middleware

import (
	"net/http"

	"github.com/iankencruz/sabiflow/internal/shared/sessions"
)

func RequireAuth(sm *sessions.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := sm.GetUserID(r)
			if err != nil || userID == 0 {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Optionally add userID to context
			// ctx := context.WithValue(r.Context(), "user_id", userID)
			// r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
