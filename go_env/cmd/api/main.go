package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
	"strings"
    "time"

    "go_env/internal/config"
    "go_env/internal/models"
)

var (
    version   = "1.0.0"
    startTime = time.Now()
)

func main() {
    cfg := config.Load()

    log.Printf("=== Configuration (%s) ===", cfg.AppEnv)
    log.Printf("App: %s v%s", cfg.AppName, version)
    log.Printf("Port: %d", cfg.Port)
    log.Printf("Log Level: %s", cfg.LogLevel)
    log.Printf("Max Workers: %d", cfg.MaxWorkers)
    log.Printf("Timeout: %v", cfg.Timeout)
    log.Printf("Feature Flags: %v", cfg.FeatureFlags)
    log.Printf("Sources: %v", cfg.Sources)
    log.Printf("==========================")

    http.HandleFunc("/", rootHandler(cfg))
    http.HandleFunc("/config", configHandler(cfg))
    http.HandleFunc("/health", healthHandler(cfg))
    http.HandleFunc("/env", envHandler)

    addr := fmt.Sprintf(":%d", cfg.Port)
    log.Printf("Server starting on %s", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}

func rootHandler(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        response := map[string]interface{}{
            "app":         cfg.AppName,
            "version":     version,
            "environment": cfg.AppEnv,
            "hostname":    cfg.Hostname,
            "uptime":      time.Since(startTime).String(),
            "timestamp":   time.Now().Format(time.RFC3339),
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

func configHandler(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cfg.SafeConfig())
    }
}

func healthHandler(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        status := models.HealthStatus{
            Status:    "healthy",
            Version:   version,
            Uptime:    time.Since(startTime).String(),
            Timestamp: time.Now(),
            Env:       cfg.AppEnv,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(status)
    }
}

func envHandler(w http.ResponseWriter, r *http.Request) {
    envVars := make([]string, 0)
    for _, env := range os.Environ() {
        if strings.HasPrefix(env, "APP_") {
            parts := strings.SplitN(env, "=", 2)
            envVars = append(envVars, parts[0])
        }
    }
    
    response := map[string]interface{}{
        "message":          fmt.Sprintf("%d environment variables available", len(os.Environ())),
        "app_variables":    envVars,
        "config_endpoint":  "/config for full configuration",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}