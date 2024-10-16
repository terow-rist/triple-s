package main

import (
	"log"
	"net/http"
	"strings"
	"triple-s/config"
	"triple-s/internal"
)

func main() {
	mux := http.NewServeMux()
	log.Printf("http://localhost:%s/\n", config.PortNumber)
	// mux.HandleFunc("/put/", internal.PutHandler)

	mux.HandleFunc("/put/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Count(strings.TrimPrefix(r.URL.Path, "/put/"), "/") == 0 {
			internal.PutHandler(w, r)
		} else {
			internal.UploadNewObject(w, r)
		}
	})

	mux.HandleFunc("/get/", internal.GetHandler)
	mux.HandleFunc("/delete/", internal.DeleteHandler)
	log.Fatal(http.ListenAndServe(":"+config.PortNumber, mux))
}
