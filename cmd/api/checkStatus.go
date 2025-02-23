package main

import (
	"encoding/json"
	"net/http"
)

type requestPayload struct {
	Code string `json:"code"`
}

func (app *application) checkStatus(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	if payload.Code == "" {
		http.Error(w, `{"error": "Missing code"}`, http.StatusBadRequest)
		return
	}

	exists, err := app.store.CheckStatus.Exists(payload.Code)
	if err != nil || !exists {
		http.Error(w, `{"status": "failure", "message": "Processing ID not found"}`, http.StatusNotFound)
		return
	}

	result, err := app.store.CheckStatus.GetProcessedImages(payload.Code)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch processing data"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   result,
	})
}
