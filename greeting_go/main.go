package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	greeting := os.Getenv("GREETING")
	if greeting == "" {
		greeting = "Hello Docker!"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, greeting)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on: %s with greeting: %s", port, greeting)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
