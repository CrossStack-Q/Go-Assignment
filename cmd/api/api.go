package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	r.Use(corsHandler.Handler)

	fileServer := http.FileServer(http.Dir("./images"))
	r.Handle("/v1/images/*", http.StripPrefix("/v1/images/", fileServer))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheck)
		r.Post("/upload", app.uploadCsvHandler)
		r.Post("/checkStatus", app.checkStatus)
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 2,
		ReadTimeout:  time.Second * 2,
		IdleTimeout:  time.Second * 2,
	}

	log.Println("Server started at ", app.config.addr)

	return srv.ListenAndServe()

}
