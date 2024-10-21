package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
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
	// validate for len of bucket or object
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

	contentType := r.Header.Get("Content-Type")
	contentLength := r.Header.Get("Content-Length")

	file, err := os.Create(filepath.Join(config.Directory+"/"+bucketName, objectKey))
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if contentLength == "" {
		fileInfo, err := file.Stat()
		if err != nil {
			writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		contentLength = strconv.FormatInt(fileInfo.Size(), 10)
	}

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
	updateBucketCSV("Acitve", bucketName)
	w.WriteHeader(http.StatusOK)
}

func RetrieveObject(w http.ResponseWriter, r *http.Request) {
	// http errors checking
	if r.Method != http.MethodGet {
		writeXMLError(w, "MethodNotAllowed", "Error: only GET command in /get/ url.", http.StatusMethodNotAllowed)
		return
	}
	bucketName, objectKey, err := checkPathURL("/get/", r)
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
	is, err = elementExists("/"+bucketName+"/objects.csv", objectKey)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: object key does not exists$.", http.StatusBadRequest)
		return
	}
	xmlData, err := listObjectData(bucketName, objectKey)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(xmlData)
}

func DeleteAnObject(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(len(records))
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
