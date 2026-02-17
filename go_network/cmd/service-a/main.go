package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "go_network/internal/client"
)

type ServiceInfo struct {
    Name    string `json:"name"`
    Message string `json:"message"`
    Peer    string `json:"peer,omitempty"`
}

func main() {
    serviceName := "service-a"
    peerName := os.Getenv("PEER_NAME")
    if peerName == "" {
        peerName = "service-b"
    }

    // Клиент для обращения к peer сервису
    peerClient := client.NewHTTPClient(fmt.Sprintf("http://%s:8080", peerName))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        info := ServiceInfo{
            Name:    serviceName,
            Message: fmt.Sprintf("Hello from %s", serviceName),
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(info)
    })

    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })

    http.HandleFunc("/peer", func(w http.ResponseWriter, r *http.Request) {
        // Обращаемся к peer сервису
        response, err := peerClient.Get("/")
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to reach peer: %v", err), http.StatusInternalServerError)
            return
        }

        info := ServiceInfo{
            Name:    serviceName,
            Message: fmt.Sprintf("Response from %s", peerName),
            Peer:    response,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(info)
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("%s starting on port %s", serviceName, port)
    log.Printf("Peer: %s", peerName)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}