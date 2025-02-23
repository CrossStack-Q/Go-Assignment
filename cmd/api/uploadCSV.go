package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CrossStack-Q/Go-Assignment/internals"
	"github.com/theMitocondria/slimuuid"
)

func dirExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, os.ModePerm)
	}
	return nil
}

func validateCSV(file io.Reader) error {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header")
	}

	expectedHeader := []string{"SNo", "Product Name", "Input Image URLS"}
	for i, h := range expectedHeader {
		if i >= len(header) || strings.TrimSpace(header[i]) != h {
			return fmt.Errorf("invalid CSV format: incorrect header %v", header)
		}
	}

	lineNumber := 2
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV at line %d", lineNumber)
		}

		if len(record) != 3 {
			return fmt.Errorf("invalid CSV format: line %d should have 3 columns", lineNumber)
		}

		if _, err := strconv.Atoi(strings.TrimSpace(record[0])); err != nil {
			return fmt.Errorf("invalid SNo at line %d: must be an integer", lineNumber)
		}

		images := strings.Split(strings.TrimSpace(record[2]), ",")
		if len(images) == 0 || images[0] == "" {
			return fmt.Errorf("invalid image URLs at line %d: must contain at least one URL", lineNumber)
		}

		lineNumber++
	}
	return nil
}

func (app *application) uploadCsvHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	// today := time.Now().Format("2006-01-02")
	inputDir := "csv/input"
	outputDir := "csv/output"

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

	if err := validateCSV(file); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	file.Seek(0, io.SeekStart)

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

	go internals.ProcessCSV(inputFilePath, outputFilePath, uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "File uploaded successfully. Processing started.", "input_file": "%s", "output_file": "%s"}`, inputFilePath, outputFilePath)))
}
