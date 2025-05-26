package middleware

import (
	"net/http"
	"slices"

	"github.com/iankencruz/sabiflow/internal/auth"
	"github.com/iankencruz/sabiflow/internal/shared/response"
	"github.com/iankencruz/sabiflow/internal/shared/sessions"
)

var (
	sessionManager *sessions.Manager
	userRepo       auth.UserRepository
)

// InitPermissionMiddleware wires the shared session manager and repository
// exactly once — do this in routes.go or wherever you build the router.
func InitPermissionMiddleware(sm *sessions.Manager, r auth.UserRepository) {
	sessionManager = sm
	userRepo = r
}

// -----------------------------------------------------------------------------
// Core middle-wares
// -----------------------------------------------------------------------------

// RequireAuth ensures the request carries a valid session.
func RequireAuth(sm *sessions.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := sm.GetUserID(r)
			if err != nil || userID == 0 {
				_ = response.WriteJSON(w, http.StatusUnauthorized, "unauthorised", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PermissionRequired checks a user’s group permissions.
// If requireAll == true, the user must have *every* permission in requiredPerms.
// Otherwise, having *any* single permission suffices.
func PermissionRequired(
	sm *sessions.Manager,
	repo auth.UserRepository,
	requiredPerms []string,
	requireAll bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID, err := sm.GetUserID(r)
			if err != nil || userID == 0 {
				_ = response.WriteJSON(w, http.StatusUnauthorized, "unauthorised", nil)
				return
			}

			userPerms, err := repo.GetGroupPermissions(r.Context(), userID)
			if err != nil {
				_ = response.WriteJSON(w, http.StatusInternalServerError, "failed to fetch permissions", nil)
				return
			}

			if requireAll {
				// Must possess *all* required permissions
				for _, perm := range requiredPerms {
					if !slices.Contains(userPerms, perm) {
						_ = response.WriteJSON(w, http.StatusForbidden, "forbidden", nil)
						return
					}
				}
			} else {
				// Must possess *at least one*
				allowed := false
				for _, perm := range requiredPerms {
					if slices.Contains(userPerms, perm) {
						allowed = true
						break
					}
				}
				if !allowed {
					_ = response.WriteJSON(w, http.StatusForbidden, "forbidden", nil)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// -----------------------------------------------------------------------------
// Friendly wrappers — import middleware and use Can / CanAny / CanAll inline
// -----------------------------------------------------------------------------

func Can(permission string) func(http.Handler) http.Handler {
	return PermissionRequired(sessionManager, userRepo, []string{permission}, false)
}

func CanAny(perms ...string) func(http.Handler) http.Handler {
	return PermissionRequired(sessionManager, userRepo, perms, false)
}

func CanAll(perms ...string) func(http.Handler) http.Handler {
	return PermissionRequired(sessionManager, userRepo, perms, true)
}
