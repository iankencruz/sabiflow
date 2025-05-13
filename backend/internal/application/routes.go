package application

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/iankencruz/sabiflow/internal/application"
	"github.com/iankencruz/sabiflow/internal/handlers"
	"github.com/iankencruz/sabiflow/internal/response"
	"github.com/iankencruz/sabiflow/internal/sessions"
)

func Routes(app *application.Application) *chi.Mux {
	r := chi.NewRouter()

	// CORS (allow frontend app to make requests)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Update for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ✅ API Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			response.WriteJSON(w, http.StatusOK, "Connected to Sabiflow backend ✅", "Data will go here")
		})

		r.Route("/auth", func(r chi.Router) {
			// Only include handlers that exist
			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("TODO: Login handler"))
			})

			r.Post("/register", handlers.RegisterUser(app))
		})
	})

	// ✅ Final fallback
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If it's an API route, return a JSON 404
		if len(r.URL.Path) >= 5 && r.URL.Path[:5] == "/api/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Not found",
				"message": "No API route matches this path",
			})
			return
		}

		// ❗ Protect all fallback routes with session check
		userID, err := sessions.GetUserID(r)
		if err != nil || userID == 0 {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Serve real static file if it exists
		fs := http.Dir("./static")
		fileServer := http.FileServer(fs)
		if f, err := fs.Open(r.URL.Path); err == nil {
			defer f.Close()
			if stat, _ := f.Stat(); stat != nil && !stat.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		// Otherwise fallback to index.html (SvelteKit SPA)
		http.ServeFile(w, r, "./static/index.html")
	}))

	return r
}
