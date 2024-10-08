package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

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

	err := os.Mkdir(bucketName, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Succesfull creation of bucket!")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error: only GET command in /get/ url.", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "NOTHING")
}

func main() {
	http.HandleFunc("/put/", putHandler)
	http.HandleFunc("/get/", getHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
