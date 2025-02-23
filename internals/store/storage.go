package store

import "database/sql"

type Storage struct {
	CompressImages interface {
		CSVUpload(uid, inputFilePath, outputFilePath string) error
		ProcessImages(uid, productName string, inputImageUrls, outputImageUrls []string) error
	}

	CheckStatus interface {
		Exists(uniqueCode string) (bool, error)
		GetProcessedImages(uniqueCode string) ([]map[string]interface{}, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		CompressImages: &CSVStore{db},
		CheckStatus:    &CheckStatusStore{db},
	}
}
