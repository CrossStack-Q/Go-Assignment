package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	outputFilePath := filepath.Join("csv/output", fmt.Sprintf("%s_output.csv", payload.Code))

	file, err := os.Open(outputFilePath)
	if err != nil {
		http.Error(w, `{"status": "failure", "message": "Output file not found"}`, http.StatusNotFound)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		http.Error(w, `{"status": "failure", "message": "Failed to read CSV file"}`, http.StatusInternalServerError)
		return
	}

	expectedHeaders := []string{"S.No", "Product Name", "Input Image URLs", "Output Image URLs"}
	if len(header) < len(expectedHeaders) {
		http.Error(w, `{"status": "failure", "message": "Invalid CSV format: Missing required columns"}`, http.StatusBadRequest)
		return
	}

	for i, col := range expectedHeaders {
		if header[i] != col {
			http.Error(w, fmt.Sprintf(`{"status": "failure", "message": "Invalid column name: expected '%s' but found '%s'"}`, col, header[i]), http.StatusBadRequest)
			return
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, `{"status": "failure", "message": "Error reading CSV rows"}`, http.StatusInternalServerError)
			return
		}

		if len(record) < len(expectedHeaders) {
			http.Error(w, `{"status": "failure", "message": "Invalid CSV format: Missing data in some columns"}`, http.StatusBadRequest)
			return
		}

		inputUrls := strings.Split(strings.TrimSpace(record[2]), ",")
		outputUrls := strings.Split(strings.TrimSpace(record[3]), ",")

		if len(inputUrls) != len(outputUrls) {
			http.Error(w, `{"status": "failure", "message": "Mismatch in Input and Output Image URLs count"}`, http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Output file is valid"}`))
}
