package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CrossStack-Q/Go-Assignment/internals/store"
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

	err = app.store.CompressImages.CSVUpload(uid, inputFilePath, outputFilePath)
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error": "Failed to save upload metadata"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("CSV uploaded successfully: %s", inputFilePath)

	go ProcessCSV(&app.store, inputFilePath, outputFilePath, uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "File uploaded successfully. Processing started.", "input_file": "%s", "output_file": "%s"}`, inputFilePath, outputFilePath)))
}

func ProcessCSV(store *store.Storage, inputFilePath, outputFilePath, uid string) {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Printf("Failed to open CSV file: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV file: %v", err)
		return
	}

	outFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Printf("Failed to create output CSV: %v", err)
		return
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	header := []string{"S.No", "Product Name", "Input Image URLs", "Output Image URLs"}
	writer.Write(header)

	for i, row := range rows[1:] {
		if len(row) < 3 {
			log.Printf("Skipping invalid row %d: %v", i+1, row)
			continue
		}

		sNo := row[0]
		productName := row[1]
		inputImageUrls := strings.Split(row[2], ",")
		outputImageUrls := []string{}

		for _, imageUrl := range inputImageUrls {
			imageUrl = strings.TrimSpace(imageUrl)
			outputUrl := processAndSaveImage(imageUrl, productName, uid)
			if outputUrl != "" {
				outputImageUrls = append(outputImageUrls, outputUrl)
			}
		}

		// db ..

		err := store.CompressImages.ProcessImages(uid, productName, inputImageUrls, outputImageUrls)

		if err != nil {
			log.Printf("Failed to process images for %s: %v", productName, err)
		}
		writer.Write([]string{
			sNo,
			productName,
			strings.Join(inputImageUrls, ", "),
			strings.Join(outputImageUrls, ", "),
		})
	}

	log.Printf("CSV processing completed. Output file: %s", outputFilePath)
}

func downloadImage(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image fetch failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %v", err)
	}

	return data, nil
}

func processAndSaveImage(imageURL, productName, uid string) string {

	fmt.Println("Processing Image:", imageURL)
	imgData, err := downloadImage(imageURL)
	if err != nil {
		log.Printf("Failed to download image: %v", err)
		return ""
	}

	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Printf("Failed to decode image: %s, format: %s, error: %v", imageURL, format, err)
		return ""
	}

	imageName := filepath.Base(imageURL)

	outputDir := fmt.Sprintf("imagesOut/%s/%s/", uid, productName)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		return ""
	}

	outputFile := filepath.Join(outputDir, imageName)
	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Printf("Failed to create output file: %v", err)
		return ""
	}
	defer outFile.Close()

	if format == "png" {
		err = png.Encode(outFile, img)
	} else if format == "jpeg" || format == "jpg" {
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 50})
	} else {
		log.Printf("Unsupported image format: %s", format)
		return ""
	}

	if err != nil {
		log.Printf("Failed to save compressed image: %v", err)
		return ""
	}

	log.Printf("Image processed and saved: %s", outputFile)
	return outputFile
}
