package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CrossStack-Q/Go-Assignment/internals"
	"github.com/theMitocondria/slimuuid"
)

func dirExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, os.ModePerm)
	}
	return nil
}

func (app *application) uploadCsvHandler(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	today := time.Now().Format("2006-01-02")
	inputDir := filepath.Join("csv/input", today)
	outputDir := filepath.Join("csv/output", today)

	if err := dirExists(inputDir); err != nil {
		http.Error(w, `{"error": "Failed to create input directory"}`, http.StatusInternalServerError)
		return
	}
	if err := dirExists(outputDir); err != nil {
		http.Error(w, `{"error": "Failed to create output directory"}`, http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("csv")
	if err != nil {
		http.Error(w, `{"error": "Invalid file upload"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".csv") {
		http.Error(w, `{"error": "Only CSV files are allowed"}`, http.StatusBadRequest)
		return
	}

	uid, err := slimuuid.Generate()
	if err != nil {
		http.Error(w, `{"error": "Failed to generate unique ID"}`, http.StatusInternalServerError)
		return
	}

	inputFilePath := filepath.Join(inputDir, fmt.Sprintf("%s_input.csv", uid))
	outputFilePath := filepath.Join(outputDir, fmt.Sprintf("%s_output.csv", uid))

	outFile, err := os.Create(inputFilePath)
	if err != nil {
		http.Error(w, `{"error": "Failed to save file"}`, http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, file); err != nil {
		http.Error(w, `{"error": "Failed to write file"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("CSV uploaded successfully: %s", inputFilePath)

	go internals.ProcessCSV(inputFilePath, outputFilePath)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "File uploaded successfully. Processing started.", "input_file": "%s", "output_file": "%s"}`, inputFilePath, outputFilePath)))
}
