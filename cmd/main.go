package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"triple-s/config"
)

func elementExists(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func writeCSV(bucketName string) {
	file, err := os.OpenFile("buckets/buckets.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if fileInfo, _ := os.Stat(file.Name()); fileInfo.Size() == 0 {
		writer.Write([]string{"NameOfBucket", "DateOfCreation"})
	}

	err = writer.Write([]string{bucketName, time.Now().Format("2006/01/02 15:04:05")})
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func readCSV(bucketName string) bool {
	file, err := os.OpenFile("buckets/buckets.csv", os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	for _, v := range records {
		if elementExists(v, bucketName) {
			return true
		}
	}
	return false
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

	err = os.Mkdir(filepath.Join("buckets", filepath.Join(config.Directory, bucketName)), 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if readCSV(bucketName) {
		http.Error(w, "Error: The bucket name is already in use.", http.StatusConflict)
		return
	}

	writeCSV(bucketName)

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
	log.Printf("http://localhost:%s/\n", config.PortNumber)
	http.HandleFunc("/put/", putHandler)
	http.HandleFunc("/get/", getHandler)
	log.Fatal(http.ListenAndServe(":"+config.PortNumber, nil))
}
