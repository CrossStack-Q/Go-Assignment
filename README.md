# Go-Assignment
Assignment Line .

## Overview
This project is a **Golang-based image processing system** that reads image URLs from a CSV file, downloads them, compresses them to 50% quality, and saves them in an organized directory structure while maintaining the original filenames.

## Features
- **CSV File Processing**: Reads product image URLs from a CSV file.
- **Asynchronous Image Processing**: Uses Goroutines to process images concurrently.
- **Dynamic Naming**: Saves processed images with the same name as the input image.
- **Organized Output**: Stores images in a structured directory format:  
  `imagesOut/{uid}/{Product Name}/{original_filename}`
- **Supports PNG & JPEG Formats**: Converts images while maintaining format compatibility.
- **Error Handling**: Logs errors if an image fails to download or process.

## Project Structure
```
project-root/
│── cmd/api/         # Main application logic
│── cmd/migrate/     # Database migrations (if applicable)
│── internals/       # Core logic for image processing
│── scripts/         # Utility scripts
│── csv/input/       # CSV input folder
│── imagesOut/       # Processed image storage
│── main.go          # Entry point for application
│── README.md        # Project documentation
```

## Input CSV Format
The input CSV should be structured as follows:
```csv
SNo,Product Name,Input Image URLS
1,Product 1,"http://localhost:8080/v1/images/o1.png,http://localhost:8080/v1/images/o2.png"
2,Product 2,"http://localhost:8080/v1/images/t1.png,http://localhost:8080/v1/images/t2.png"
```

## Output Directory Structure
```
imagesOut/
│── {uid}/
    │── Product 1/
        │── o1.png
        │── o2.png
    │── Product 2/
        │── t1.png
        │── t2.png
```

## Installation & Setup
### Prerequisites
- **Golang 1.18+**
- **PostgreSQL**

## Functionality Breakdown
### `processAndSaveImage(imageURL, productName, uid)`
- Downloads the image.
- Extracts the filename from the URL.
- Saves the compressed image in `imagesOut/{uid}/{productName}/` with the same filename.

### Key Code Snippet:
```go
imageName := filepath.Base(imageURL) // Extracts original filename
outputDir := fmt.Sprintf("imagesOut/%s/%s/", uid, productName)
os.MkdirAll(outputDir, os.ModePerm)
outputFile := filepath.Join(outputDir, imageName)
```

## Future Improvements
- Add a database table to track processed images.
- Implement a REST API for querying processed images.
- Add support for more image formats.

For any issues or contributions, feel free to open a pull request or raise an issue! 🚀

