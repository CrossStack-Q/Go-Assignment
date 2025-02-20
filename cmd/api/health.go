package main

import "net/http"

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health"))
}
