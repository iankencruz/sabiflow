package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/iankencruz/sabiflow/internal/application"
	"github.com/iankencruz/sabiflow/internal/middleware"
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

	// ✅ Main API group
	r.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Connected to Sabiflow backend ✅",
			})
		})

		// Auth endpoints
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", app.Auth.LoginUser)
			// r.Post("/register", ...)
			// r.Post("/logout", ...)
		})

		// ✅ Admin-only endpoints (protected)
		// r.Route("/admin", func(r chi.Router) {
		// 	r.Use(authmw.RequireAdminAuth)
		//
		// 	r.Get("/clients", app.Admin.GetClients)
		// 	r.Post("/invoices", app.Admin.CreateInvoice)
		// 	// Add more secured admin endpoints here
		// })

		// Other modules like CRM can go here
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
