package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


type application struct {
	config config
}
type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Get("/v1/health", app.healthcheckHandler)

	return mux
}

func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 15 * time.Second,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server has started at %v", app.config.addr)

	return server.ListenAndServe()
}