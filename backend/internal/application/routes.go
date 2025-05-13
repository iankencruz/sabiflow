package application

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/iankencruz/sabiflow/internal/auth"
	mw "github.com/iankencruz/sabiflow/internal/shared/middleware"
	"github.com/iankencruz/sabiflow/internal/shared/response"
)

func Routes(app *Application) *chi.Mux {
	r := chi.NewRouter()

	// CORS (allow frontend app to make requests)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // <- frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			response.WriteJSON(w, http.StatusOK, "Connected to Sabiflow backend ✅", "Data will go here")
		})

		// Authentication Routes
		r.Route("/auth", func(r chi.Router) {
			// Only include handlers that exist
			r.Post("/login", app.AuthHandler.LoginHandler)
			r.Post("/register", app.AuthHandler.RegisterHandler)

			// r.With(mw.RequireAuth(app.AuthHandler.SessionManager)).Post("/logout", app.AuthHandler.LogoutHandler)
			r.Post("/logout", app.AuthHandler.LogoutHandler)

		})

		// ✅ API Routes
		r.Group(func(r chi.Router) {
			r.Use(mw.RequireAuth(app.AuthHandler.SessionManager))
			r.Get("/user/me", func(w http.ResponseWriter, r *http.Request) {
				userID, err := app.AuthHandler.SessionManager.GetUserID(r)
				if err != nil || userID == 0 {
					response.WriteJSON(w, http.StatusUnauthorized, "Not authenticated", nil)
					return
				}

				user, err := app.AuthHandler.Service.(*auth.AuthServiceImpl).Repo.GetByID(r.Context(), userID)
				if err != nil {
					response.WriteJSON(w, http.StatusInternalServerError, "Could not fetch user", nil)
					return
				}

				response.WriteJSON(w, http.StatusOK, "Authenticated", map[string]any{"user": user})
			})
		})

	})

	// ✅ Final fallback

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Not found",
				"message": "No API route matches this path",
			})
			return
		}

		// ✅ Allow /login to work even after session cleared
		if r.URL.Path == "/login" {
			http.ServeFile(w, r, "./static/index.html")
			return
		}

		userID, err := app.AuthHandler.SessionManager.GetUserID(r)
		if err != nil || userID == 0 {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Serve static file or fallback
		fs := http.Dir("./static")
		fileServer := http.FileServer(fs)
		if f, err := fs.Open(r.URL.Path); err == nil {
			defer f.Close()
			if stat, _ := f.Stat(); stat != nil && !stat.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		http.ServeFile(w, r, "./static/index.html")
	}))

	return r
}
