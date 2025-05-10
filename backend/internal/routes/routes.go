package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/sabiflow/internal/application"
)

func Routes(app *application.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sabiflow CRM API root"))
	})

	// Future: r.Route("/leads", func(r chi.Router) { ... })

	return r
}
