package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	// mux.HandleFunc("/songs", handleSongs)

	log.Println("server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
