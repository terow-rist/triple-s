package internal

import (
	"encoding/xml"
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"triple-s/config"
)

type Bucket struct {
	Name             string `xml:"Name"`
	DateOfCreation   string `xml:"DateOfCreation"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

type Buckets struct {
	SliceBucket []Bucket `xml:"Bucket"`
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets Buckets  `xml:"Buckets"`
}

type ErrorResponse struct {
	XMLName    xml.Name `xml:"Error"`
	StatusCode string   `xml:"StatusCode"`
	Message    string   `xml:"Message"`
}

type Response struct {
	StatusCode string `xml:"StatusCode"`
	Message    string `xml:"Message"`
}

func listAllMyBucketsResult() ([]byte, error) {
	records, err := readCSV("/buckets.csv")
	if err != nil {
		return nil, err
	}
	buckets := Buckets{}
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			return nil, errors.New("csv record is not correct")
		}
		soloBucket := Bucket{
			Name:             record[0],
			DateOfCreation:   record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		}
		buckets.SliceBucket = append(buckets.SliceBucket, soloBucket)
	}

	result := ListAllMyBucketsResult{
		Buckets: buckets,
	}

	xmlData, err := xml.MarshalIndent(result, "", "   ")
	if err != nil {
		return nil, err
	}
	xmlData = append(xmlData, '\n')
	return xmlData, nil
}

func listObjectData(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) error {
	// Construct the file path
	filePath := filepath.Join(config.Directory, bucketName, objectKey)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file info for modification time
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Detect the content type based on the file extension
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		// Fallback content type if not determined by extension
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)

	// Serve the content using http.ServeContent
	http.ServeContent(w, r, objectKey, fileInfo.ModTime(), file)
	return nil
}

func writeXMLError(w http.ResponseWriter, statusCode, message string, code int) {
	errorResponse := ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
	}
	xmlData, err := xml.MarshalIndent(errorResponse, "", "   ")
	if err != nil {
		errorResponse.StatusCode = "Internal Server Error"
		errorResponse.Message = "Error: error in xml.MarshalIndent"
	}
	xmlData = append(xmlData, '\n')
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	w.Write(xmlData)
}

func writeXMLResponse(w http.ResponseWriter, statusCode, message string, code int) {
	response := Response{
		Message:    message,
		StatusCode: statusCode,
	}

	xmlData, err := xml.MarshalIndent(response, "", "   ")
	if err != nil {
		writeXMLError(w, "Internal Server Error", "Error: error in xml.MarshalIndent", code)
		return
	}
	xmlData = append(xmlData, '\n')
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	w.Write(xmlData)
}
