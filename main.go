package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	portNumber string
	directory  string
)

func init() {
	flag.StringVar(&portNumber, "port", "8080", "Port number")
	flag.StringVar(&directory, "dir", "", "Path to the directory")

	helpMessage :=
		`Simple Storage Service.

**Usage:**
	triple-s [-port <N>] [-dir <S>]  
	triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, helpMessage)
	}
	flag.Parse()
}

func putHandler(w http.ResponseWriter, r *http.Request) {
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
	err = os.Mkdir(filepath.Join(directory, bucketName), 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Succesfull creation of bucket!")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	http.HandleFunc("/put/", putHandler)
	http.HandleFunc("/get/", getHandler)
	log.Fatal(http.ListenAndServe(":"+portNumber, nil))
}
