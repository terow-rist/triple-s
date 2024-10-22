package internal

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"triple-s/config"
)

type ObjectMD struct {
	ObjectKey    string
	Size         string
	ContentType  string
	LastModified string
}

func UploadNewObject(w http.ResponseWriter, r *http.Request) {
	// http errors handling
	if r.Method != http.MethodPut {
		writeXMLError(w, "MethodNotAllowed", "Error: only PUT command for /put/ url.", http.StatusMethodNotAllowed)
		return
	}

	bucketName, objectKey, err := checkPathURL("/put/", r)
	if err != nil {
		writeXMLError(w, "BadRequest", "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if objectKey == "" {
		writeXMLError(w, "BadRequest", "Error: object key cannot be empty.", http.StatusBadRequest)
		return
	}

	// validate bucket existence
	is, err := elementExists("/buckets.csv", bucketName)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: bucket name does not exist.", http.StatusBadRequest)
		return
	}

	// Open file for writing
	file, err := os.Create(filepath.Join(config.Directory+"/"+bucketName, objectKey))
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Write the file content from the request body to the file
	_, err = io.Copy(file, r.Body)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error writing file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the content length from the file after writing
	fileInfo, err := file.Stat()
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	contentLength := strconv.FormatInt(fileInfo.Size(), 10)
	contentType := r.Header.Get("Content-Type")

	// Save object metadata in CSV
	o := ObjectMD{
		ObjectKey:    objectKey,
		Size:         contentLength,
		ContentType:  contentType,
		LastModified: time.Now().Format("2006/01/02 15:04:05"),
	}
	err = writeObjectCSV(bucketName, o)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	updateBucketCSV("Active", bucketName)
	writeXMLResponse(w, "OK", "Successful creation of object!", http.StatusOK)
}

func RetrieveObject(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		writeXMLError(w, "MethodNotAllowed", "Error: only GET command in /get/ url.", http.StatusMethodNotAllowed)
		return
	}

	// Parse bucket name and object key
	bucketName, objectKey, err := checkPathURL("/get/", r)
	if err != nil {
		writeXMLError(w, "BadRequest", "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate object key
	if objectKey == "" {
		writeXMLError(w, "BadRequest", "Error: object key cannot be empty.", http.StatusBadRequest)
		return
	}

	// Validate bucket existence
	is, err := elementExists("/buckets.csv", bucketName)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: bucket name does not exist.", http.StatusBadRequest)
		return
	}

	// Validate object existence
	is, err = elementExists("/"+bucketName+"/objects.csv", objectKey)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: object key does not exist.", http.StatusBadRequest)
		return
	}

	// Get the object metadata (e.g., Content-Type) from the CSV file
	records, err := readCSV("/" + bucketName + "/objects.csv")
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var contentType string
	for _, v := range records {
		if v[0] == objectKey {
			contentType = v[2] // Assuming column 2 has Content-Type
			break
		}
	}

	// Open and serve the object
	filePath := filepath.Join(config.Directory, bucketName, objectKey)
	file, err := os.Open(filePath)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the correct Content-Type for the object
	w.Header().Set("Content-Type", contentType)

	// Serve the file content
	http.ServeContent(w, r, objectKey, time.Now(), file)
}

func DeleteAnObject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeXMLError(w, "MethodNotAllowed", "Error: only DELETE command in /delete/ url.", http.StatusMethodNotAllowed)
		return
	}
	bucketName, objectKey, err := checkPathURL("/delete/", r)
	if err != nil {
		writeXMLError(w, "BadRequest", "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	// validate for len of object
	if objectKey == "" {
		writeXMLError(w, "BadRequest", "Error: object key cannot be empty.", http.StatusBadRequest)
		return
	}
	// validate bucket existence
	is, err := elementExists("/buckets.csv", bucketName)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: bucket name does not exists$.", http.StatusBadRequest)
		return
	}

	// validate object existence
	pathToCSV := "/" + bucketName + "/objects.csv"
	is, err = elementExists(pathToCSV, objectKey)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: object key does not exists$.", http.StatusBadRequest)
		return
	}
	err = os.Remove(config.Directory + "/" + bucketName + "/" + objectKey)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = deleteRecord(pathToCSV, objectKey)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.OpenFile(config.Directory+pathToCSV, os.O_RDONLY, 0o644)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(records) == 1 {
		updateBucketCSV("MarkedForDeletion", bucketName)
	}
	w.WriteHeader(http.StatusNoContent)
}

func checkPathURL(what string, r *http.Request) (string, string, error) {
	path := strings.TrimPrefix(r.URL.Path, what)
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return "", "", errors.New("Invalid URL format.")
	}
	return parts[0], parts[1], nil
}
