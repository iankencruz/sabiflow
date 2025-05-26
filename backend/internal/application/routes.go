// internal/routes/routes.go
package application

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	mw "github.com/iankencruz/sabiflow/internal/shared/middleware" // RequireAuth, Can, …
	"github.com/iankencruz/sabiflow/internal/shared/response"      // WriteJSON helper
)

// Router returns the fully-wired chi router.
// Call this from cmd/api/main.go *after* you construct Application.
func Routes(app *Application) *chi.Mux {
	//--------------------------------------------------------------------
	// One-time injection so the simple Can/CanAny/CanAll wrappers work
	//--------------------------------------------------------------------
	mw.InitPermissionMiddleware(app.SessionManager, app.UserRepo)

	//--------------------------------------------------------------------
	// Root router + global middleware
	//--------------------------------------------------------------------
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//--------------------------------------------------------------------
	// API v1
	//--------------------------------------------------------------------
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			// ---------------- Health check ----------------
			r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
				response.WriteJSON(w, http.StatusOK,
					"Connected to Sabiflow backend ✅", nil)
			})

			// ---------------- Auth ------------------------
			r.Route("/auth", func(r chi.Router) {
				r.Post("/login", app.AuthHandler.LoginHandler)
				r.Post("/register", app.AuthHandler.RegisterHandler)
				r.Post("/logout", app.AuthHandler.LogoutHandler)

				// Google OAuth2
				r.HandleFunc("/google/login", app.AuthHandler.GoogleLogin)
				r.HandleFunc("/google/callback", app.AuthHandler.GoogleCallback)

				// Front-end uses this once per tab
				r.Get("/me", app.AuthHandler.GetAuthenticatedUser)
			})

			// ----------- Protected API group --------------
			r.Group(func(r chi.Router) {
				// Session must be valid
				r.Use(mw.RequireAuth(app.SessionManager))

				// Example of fine-grained authorisation:
				//
				// r.With(mw.Can("projects:read")).Get(
				//	   "/projects", app.ProjectHandler.List)
				//
				// r.With(mw.CanAny("contacts:read", "contacts:write")).Post(
				//     "/contacts", app.ContactHandler.Create)
			})
		})
	})

	//--------------------------------------------------------------------
	// SPA fallback & JSON 404 for unknown /api/* paths
	//--------------------------------------------------------------------
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		// If it looks like an API call, return structured JSON
		if strings.HasPrefix(req.URL.Path, "/api/") {
			response.WriteJSON(w, http.StatusNotFound,
				"No API route matches this path",
				map[string]string{"path": req.URL.Path})
			return
		}

		// Try to serve a real static asset first
		fs := http.Dir("./static")
		if f, err := fs.Open(req.URL.Path); err == nil {
			defer f.Close()
			if stat, _ := f.Stat(); stat != nil && !stat.IsDir() {
				http.FileServer(fs).ServeHTTP(w, req)
				return
			}
		}

		// Otherwise fall back to the SvelteKit index.html
		http.ServeFile(w, req, "./static/index.html")
	})

	return r
}
