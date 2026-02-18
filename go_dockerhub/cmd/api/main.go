package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    "go_dockerhub/internal/version"
)

var startTime = time.Now()

func main() {
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/version", handleVersion)
    http.HandleFunc("/health", handleHealth)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server...")
    log.Printf("Version: %s", version.Version)
    log.Printf("Commit: %s", version.Commit)
    log.Printf("Build Time: %s", version.BuildTime)
    log.Printf("Listening on port %s", port)

    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "message":    "Hello from Docker Hub!",
        "version":    version.Version,
        "uptime":     time.Since(startTime).String(),
        "timestamp":  time.Now().Format(time.RFC3339),
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-App-Version", version.Version)
    json.NewEncoder(w).Encode(response)
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(version.Get())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}