package store

import (
	"database/sql"
	"errors"
)

type CheckStatusStore struct {
	db *sql.DB
}

func (s *CheckStatusStore) Exists(uniqueCode string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM csv_uploads WHERE unique_code=$1)", uniqueCode).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *CheckStatusStore) GetProcessedImages(uniqueCode string) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`
        SELECT product_name, input_image_urls, output_image_urls 
        FROM processed_images 
        WHERE csv_upload_id = (SELECT id FROM csv_uploads WHERE unique_code=$1)`, uniqueCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var productName string
		var inputURLs, outputURLs string

		err := rows.Scan(&productName, &inputURLs, &outputURLs)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"product_name":      productName,
			"input_image_urls":  inputURLs,
			"output_image_urls": outputURLs,
		})
	}

	if len(result) == 0 {
		return nil, errors.New("no data found for the given code")
	}

	return result, nil
}
