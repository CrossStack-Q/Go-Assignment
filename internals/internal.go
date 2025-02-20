package internals

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "image/png"
)

func ProcessCSV(inputFilePath, outputFilePath string) {

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
			outputUrl := processAndSaveImage(imageUrl, productName)
			if outputUrl != "" {
				outputImageUrls = append(outputImageUrls, outputUrl)
			}
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

func processAndSaveImage(imageURL, productName string) string {

	resp, err := http.Get(imageURL)
	if err != nil {
		log.Printf("Failed to download image: %s, error: %v", imageURL, err)
		return ""
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("Failed to decode image: %s, error: %v", imageURL, err)
		return ""
	}

	today := time.Now().Format("2006-01-02")
	outputDir := filepath.Join("imagesOut", today, productName)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		return ""
	}

	imageID := filepath.Base(imageURL)
	outputFilePath := filepath.Join(outputDir, imageID)

	outFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Printf("Failed to create output image file: %s, error: %v", outputFilePath, err)
		return ""
	}
	defer outFile.Close()

	var opt jpeg.Options
	opt.Quality = 50
	err = jpeg.Encode(outFile, img, &opt)
	if err != nil {
		log.Printf("Failed to save compressed image: %s, error: %v", outputFilePath, err)
		return ""
	}

	log.Printf("Image processed and saved: %s", outputFilePath)

	return fmt.Sprintf("http://localhost:8080/%s", outputFilePath)
}
