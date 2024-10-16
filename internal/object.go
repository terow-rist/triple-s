package internal

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

	bucketName := ""
	for _, v := range r.URL.Path[len("/put/"):] {
		if v == '/' {
			break
		}
		bucketName += string(v)
	}
	if bucketName == "" {
		writeXMLError(w, "BadRequest", "Error: bucket name cannot be empty.", http.StatusBadRequest)
		return
	}
	is, err := elementExists("/buckets.csv", bucketName)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: bucket name does not exists$.", http.StatusBadRequest)
		return
	}
	// checking that --dir=path exists
	if _, err := os.Stat(config.Directory); os.IsNotExist(err) {
		err = os.Mkdir(config.Directory, 0o755)
		if err != nil {
			writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// checking that '--dir=' is standard or not
	if isStandardPackage(config.Directory) {
		writeXMLError(w, "BadRequest", "Error: directory(--dir=) cannot be one of the used ones.", http.StatusBadRequest)
		return
	}

	objectKey := r.URL.Path[len("/put/"+bucketName)+1:]
	if objectKey == "" {
		writeXMLError(w, "BadRequest", "Error: object key cannot be empty.", http.StatusBadRequest)
		return
	}

	file, err := os.Create(filepath.Join(config.Directory+"/"+bucketName, objectKey))
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")
	contentLength := r.Header.Get("Content-Length")

	if contentType == "" {
		writeXMLError(w, "BadRequest", "Error: Content-Type cannot be empty.", http.StatusBadRequest)
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

	w.WriteHeader(http.StatusOK)
}

// NOT FINISHED XML RESPONSE [434234 bytes of object data]!!!
func RetrieveObject(w http.ResponseWriter, r *http.Request) {
	// http errors checking
	if r.Method != http.MethodGet {
		writeXMLError(w, "MethodNotAllowed", "Error: only GET command in /get/ url.", http.StatusMethodNotAllowed)
		return
	}
	bucketName := ""
	for _, v := range r.URL.Path[len("/put/"):] {
		if v == '/' {
			break
		}
		bucketName += string(v)
	}
	if bucketName == "" {
		writeXMLError(w, "BadRequest", "Error: bucket name cannot be empty.", http.StatusBadRequest)
		return
	}
	is, err := elementExists("/buckets.csv", bucketName)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: bucket name does not exists$.", http.StatusBadRequest)
		return
	}

	// checking that --dir=path exists
	if _, err := os.Stat(config.Directory); os.IsNotExist(err) {
		err = os.Mkdir(config.Directory, 0o755)
		if err != nil {
			writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// checking that '--dir=' is standard or not
	if isStandardPackage(config.Directory) {
		writeXMLError(w, "BadRequest", "Error: directory(--dir=) cannot be one of the used ones.", http.StatusBadRequest)
		return
	}
	objectKey := r.URL.Path[len("/put/"+bucketName)+1:]
	if objectKey == "" {
		writeXMLError(w, "BadRequest", "Error: object key cannot be empty.", http.StatusBadRequest)
		return
	}

	is, err = elementExists("/"+bucketName+"objects.csv", objectKey)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		writeXMLError(w, "BadRequest", "Error: object key does not exists$.", http.StatusBadRequest)
		return
	}
}
