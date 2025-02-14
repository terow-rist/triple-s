package internal

import (
	"net/http"
	"os"
	"path/filepath"
	"triple-s/config"
)

func PutHandler(w http.ResponseWriter, r *http.Request) {
	// http errors handling
	if r.Method != http.MethodPut {
		writeXMLError(w, "MethodNotAllowed", "Error: only PUT command for /put/ url.", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Path[len("/put/"):]
	if bucketName == "" {
		writeXMLError(w, "BadRequest", "Error: bucket name cannot be empty.", http.StatusBadRequest)
		return
	}
	// checking the correctness of bucket name
	err := validateBucketName(bucketName)
	if err != nil {
		writeXMLError(w, "BadRequest", "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// checking the uniqueness of bucket name
	elementIn, err := elementExists("/buckets.csv", bucketName)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if elementIn {
		writeXMLError(w, "Conflict", "Error: The bucket name is already in use.", http.StatusConflict)
		return
	}
	// the creation of bucket
	err = os.Mkdir(filepath.Join(config.Directory, bucketName), 0o755)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// fullfilling the bucket metadata
	err = writeBucketCSV(bucketName)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// the finish line
	writeXMLResponse(w, "OK", "Successful creation of bucket!", http.StatusOK)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// http errors checking
	if r.Method != http.MethodGet {
		writeXMLError(w, "MethodNotAllowed", "Error: only GET command in /get/ url.", http.StatusMethodNotAllowed)
		return
	}
	if len(r.URL.Path[len("/get/"):]) > 0 {
		writeXMLError(w, "Conflict", "Error: too much data after /get/", http.StatusConflict)
		return
	}
	xmlData, err := listAllMyBucketsResult()
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(xmlData)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// http errors handling
	if r.Method != http.MethodDelete {
		writeXMLError(w, "MethodNotAllowed", "Error: only DELETE command in /delete/ url.", http.StatusMethodNotAllowed)
		return
	}
	target := r.URL.Path[len("/delete/"):]
	if target == "" {
		writeXMLError(w, "BadRequest", "Error: bucket name cannot be empty.", http.StatusBadRequest)
		return
	}

	// handle if bucket does not exists
	is, err := isBucketEmpty(config.Directory)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if is {
		writeXMLError(w, "NotFound", "Error: bucket does not exist.", http.StatusNotFound)
		return
	}

	// checking for existing bucket
	path := config.Directory + "/" + target
	elementIn, err := elementExists("/buckets.csv", target)
	if err != nil {
		writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if elementIn {
		if is, err = bucketForDeletion("/buckets.csv", target); err != nil {
			writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
			return
		} else if is {
			err = os.RemoveAll(path)
			if err != nil {
				writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			err = deleteRecord("/buckets.csv", target)
			if err != nil {
				writeXMLError(w, "InternalServerError", "Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return

		} else {
			writeXMLError(w, "Confilct", "Error: bucket is not empty.", http.StatusConflict)
			return
		}
	} else {
		writeXMLError(w, "NotFound", "Error: bucket does not exist.", http.StatusNotFound)
		return
	}
}
