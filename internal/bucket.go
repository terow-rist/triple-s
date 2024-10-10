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
		http.Error(w, "Error: "+err.Error(), http.StatusNotAcceptable)
		return
	}
	if elementExists(bucketName) {
		http.Error(w, "Error: The bucket name is already in use.", http.StatusConflict)
		return
	}

	err = os.Mkdir(filepath.Join("buckets", filepath.Join(config.Directory, bucketName)), 0755)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeCSV(bucketName)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Succesfull creation of bucket!")
}

// NOT FINISHED
func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error: only GET command in /get/ url.", http.StatusMethodNotAllowed)
		return
	}
	if len(r.URL.Path[len("/get/"):]) > 0 {
		http.Error(w, "Error: too much data after '/get/'", http.StatusConflict)
		return
	}

	fmt.Fprintln(w, len(r.URL.Path[len("/get/"):]))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Error: only DELETE command in /delete/ url.", http.StatusMethodNotAllowed)
		return
	}
	target := r.URL.Path[len("/delete/"):]
	if target == "" {
		http.Error(w, "Error: bucket name cannot be empty.", http.StatusBadRequest)
		return
	}

	path := "buckets/" + target
	if elementExists(target) {
		if is, err := isBucketEmpty(path); err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		} else if is {
			err = os.Remove(path)
			if err != nil {
				http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			http.Error(w, "Error: bucket is not empty.", http.StatusConflict)
			return
		}
	} else {
		http.Error(w, "Error: bucket does not exist.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Deleted %s bucket\n", target)
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

func isBucketEmpty(path string) (bool, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(dirEntries) == 0, nil
}
