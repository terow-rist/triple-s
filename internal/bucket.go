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
	// http errors
	if r.Method != http.MethodPut {
		http.Error(w, "Error: only PUT command for /put/ url.", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Path[len("/put/"):]
	if bucketName == "" {
		http.Error(w, "Error: bucket name cannot be empty.", http.StatusBadRequest)
		return
	}
	// checking the correctness of bucket name
	err := validateBucketName(bucketName)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusNotAcceptable)
		return
	}
	// checking that --dir=path exists
	if _, err = os.Stat(config.Directory); os.IsNotExist(err) {
		err = os.Mkdir(config.Directory, 0755)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// checking that '--dir=' is empty
	is, err := isBucketEmpty(config.Directory)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !is {
		http.Error(w, "Error: directory from '--dir=' is not empty.", http.StatusConflict)
		return
	}
	// checking the uniqueness of bucket name
	elementIn, err := elementExists(bucketName)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if elementIn {
		http.Error(w, "Error: The bucket name is already in use.", http.StatusConflict)
		return
	}
	// the creation of bucket
	err = os.Mkdir(filepath.Join(config.Directory, bucketName), 0o755)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// fullfilling the bucket metadata
	err = writeCSV(bucketName)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// the finish line
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Succesfull creation of bucket!") // NOT IN XML < REFACTOR!!!
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
	elementIn, err := elementExists(target)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if elementIn {
		if is, err := isBucketEmpty(path); err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		} else if is {
			err = os.Remove(path)
			if err != nil {
				http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			err = deleteRecord(target)
			if err != nil {
				http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return

		} else {
			http.Error(w, "Error: bucket is not empty.", http.StatusConflict)
			return
		}
	} else {
		http.Error(w, "Error: bucket does not exist.", http.StatusNotFound)
		return
	}
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
