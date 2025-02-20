package main

import "net/http"

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
}
