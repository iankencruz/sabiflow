package main

import (
	"log"
	"net/http"

	"github.com/iankencruz/sabiflow/internal/application"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	cfg := app.Config

	log.Println("Starting server on port", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, application.Routes(app))
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}
