package internal

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"triple-s/config"
)

func PutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Error: only PUT command for /put/ url.", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Path[len("/put/"):]
	if bucketName == "" {
		http.Error(w, "Error: bucket name cannot be empty.", http.StatusBadRequest)
		return
	}
	err := validateBucketName(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	if readCSV(bucketName) {
		http.Error(w, "Error: The bucket name is already in use.", http.StatusConflict)
		return
	}

	err = os.Mkdir(filepath.Join("buckets", filepath.Join(config.Directory, bucketName)), 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeCSV(bucketName)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Succesfull creation of bucket!")
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error: only GET command in /get/ url.", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "NOTHING")
}

// NOT FINISHED
func validateBucketName(name string) error {
	err := errors.New("Error: does not meet Amazon S3 naming requirements.")

	if len(name) < 3 || len(name) > 63 {
		return err
	}

	for _, char := range name {
		if !(char >= 'a' && char <= 'z' || char >= '0' && char <= '9' || char == '-' || char == '.') {
			return err
		}
	}

	return nil
}
