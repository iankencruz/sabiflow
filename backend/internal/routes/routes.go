package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/iankencruz/sabiflow/internal/application"
)

func Routes(app *application.Application) *chi.Mux {
	r := chi.NewRouter()

	// Add this first
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // or "*" for dev
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Logger)

	r.Get("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Connected to Sabiflow backend ✅",
		})
	})

	// Future routes:
	r.Route("/api/auth", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{"route": "auth base"})
		})
	})

	r.Route("/api/crm", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{"route": "crm base"})
		})
	})

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs := http.Dir("./static")
		fileServer := http.FileServer(fs)

		// Try to open the requested file
		path := r.URL.Path
		f, err := fs.Open(path)
		if err == nil {
			defer f.Close()

			// Check if it's a directory (optional)
			stat, _ := f.Stat()
			if stat != nil && !stat.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		// Fallback to index.html
		http.ServeFile(w, r, "./static/index.html")
	}))

	return r
}
