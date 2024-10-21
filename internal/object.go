package internal

import (
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

	path := strings.TrimPrefix(r.URL.Path, "/put/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		writeXMLError(w, "BadRequest", "Error: Invalid URL format.", http.StatusBadRequest)
		return
	}
	bucketName, objectKey := parts[0], parts[1]
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
	updateBucketCSV(bucketName)
	w.WriteHeader(http.StatusOK)
}

func RetrieveObject(w http.ResponseWriter, r *http.Request) {
	// http errors checking
	if r.Method != http.MethodGet {
		writeXMLError(w, "MethodNotAllowed", "Error: only GET command in /get/ url.", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/get/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		writeXMLError(w, "BadRequest", "Error: Invalid URL format.", http.StatusBadRequest)
		return
	}
	bucketName, objectKey := parts[0], parts[1]
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
