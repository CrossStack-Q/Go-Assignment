package store

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

type CSVStore struct {
	db *sql.DB
}

func (s *CSVStore) CSVUpload(uid, inputFilePath, outputFilePath string) error {
	query := `INSERT INTO csv_uploads (unique_code, input_file_path, output_file_path) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, uid, inputFilePath, outputFilePath)
	if err != nil {
		return err
	}
	return nil
}

func (s *CSVStore) ProcessImages(uid, productName string, inputImageUrls, outputImageUrls []string) error {
	_, err := s.db.Exec(
		`INSERT INTO processed_images (csv_upload_id, product_name, input_image_urls, output_image_urls) 
		VALUES ((SELECT id FROM csv_uploads WHERE unique_code=$1), $2, $3, $4)`,
		uid, productName, pq.Array(inputImageUrls), pq.Array(outputImageUrls),
	)
	if err != nil {
		log.Printf("Failed to insert processed image data: %v", err)
		return err
	}
	return nil
}
