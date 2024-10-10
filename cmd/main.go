package main

import (
	"log"
	"net/http"
	"triple-s/config"
	"triple-s/internal"
)

func main() {
	log.Printf("http://localhost:%s/\n", config.PortNumber)
	http.HandleFunc("/put/", internal.PutHandler)
	http.HandleFunc("/get/", internal.GetHandler)
	http.HandleFunc("/delete/", internal.DeleteHandler)
	log.Fatal(http.ListenAndServe(":"+config.PortNumber, nil))
}
