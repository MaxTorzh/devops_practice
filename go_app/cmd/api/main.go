package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go_app/internal/version"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/version", handleVersion)
	http.HandleFunc("/health", handleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s (env: %s)", port, os.Getenv("APP_ENV"))
	log.Printf("Version: %s, Commit: %s", version.Version, version.Commit)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Hello from Docker!",
		"env":     os.Getenv("APP_ENV"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version.GetInfo())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
