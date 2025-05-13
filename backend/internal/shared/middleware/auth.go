package middleware

import (
	"net/http"

	"github.com/iankencruz/sabiflow/internal/sessions"
)

// RequireAdminAuth ensures the user is authenticated before accessing /api/admin/* routes.
// It checks for a valid user session and redirects to /login if unauthenticated.
func RequireAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := sessions.GetUserID(r)
		if err != nil || userID == 0 {
			// Redirect to login if not authenticated
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Optionally: attach user ID to context for later use
		// r = r.WithContext(context.WithValue(r.Context(), "user_id", userID))

		next.ServeHTTP(w, r)
	})
}
