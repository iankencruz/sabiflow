package application

import (
	"net/http"

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

	r.Use(mw.RequireAdminAuth(app.AuthHandler.SessionManager))
	// ✅ API Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			response.WriteJSON(w, http.StatusOK, "Connected to Sabiflow backend ✅", "Data will go here")
		})

		r.Route("/auth", func(r chi.Router) {
			// Only include handlers that exist
			r.Post("/login", app.AuthHandler.LoginHandler)
			r.Post("/register", app.AuthHandler.RegisterHandler)

		})

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

	// ✅ Final fallback
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ✅ Allow login & register to pass through — don't block or 404 them
		if r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/register" {
			http.ServeFile(w, r, "./static/index.html")
			return
		}

		// ✅ Protect only non-API routes
		if len(r.URL.Path) < 5 || r.URL.Path[:5] != "/api/" {
			userID, err := app.AuthHandler.SessionManager.GetUserID(r)
			if err != nil || userID == 0 {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
		}

		// 🪄 Serve static asset or fallback to SvelteKit SPA
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
