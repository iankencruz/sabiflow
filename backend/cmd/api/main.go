package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/sabiflow/internal/application"
	"github.com/iankencruz/sabiflow/internal/routes"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	cfg := app.Config

	r := chi.NewRouter()
	routes.Register(r, app)

	log.Println("Starting server on port", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}
