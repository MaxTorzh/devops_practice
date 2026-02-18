package config

import (
    "os"
    "strconv"
    "strings"
    "time"
)

type Config struct {
    Sources map[string]string

    AppName    string
    AppEnv     string
    Hostname   string
    Port       int
    LogLevel   string

    DatabaseURL string
    RedisURL    string

    MaxWorkers int
    Timeout    time.Duration

    FeatureFlags map[string]bool
}

func Load() *Config {
    sources := make(map[string]string)
    
    cfg := &Config{
        Sources: sources,
        
        AppName:    getEnvWithSource("APP_NAME", "GoEnvApp", sources),
        AppEnv:     getEnvWithSource("APP_ENV", "development", sources),
        Hostname:   getHostname(sources),
        Port:       getEnvAsIntWithSource("PORT", 8080, sources),
        LogLevel:   getEnvWithSource("LOG_LEVEL", "info", sources),
        
        DatabaseURL: getEnvWithSource("DATABASE_URL", "postgres://localhost:5432/mydb?sslmode=disable", sources),
        RedisURL:    getEnvWithSource("REDIS_URL", "redis://localhost:6379", sources),
        
        MaxWorkers:  getEnvAsIntWithSource("MAX_WORKERS", 10, sources),
        Timeout:     getEnvAsDurationWithSource("TIMEOUT", 30*time.Second, sources),
        
        FeatureFlags: parseFeatureFlags(sources),
    }
    
    return cfg
}

func getEnvWithSource(key, defaultValue string, sources map[string]string) string {
    if value := os.Getenv(key); value != "" {
        sources[key] = "environment"
        return value
    }
    sources[key] = "default"
    return defaultValue
}

func getEnvAsIntWithSource(key string, defaultValue int, sources map[string]string) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            sources[key] = "environment"
            return intVal
        }
    }
    sources[key] = "default"
    return defaultValue
}

func getEnvAsDurationWithSource(key string, defaultValue time.Duration, sources map[string]string) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            sources[key] = "environment"
            return duration
        }
    }
    sources[key] = "default"
    return defaultValue
}

func getHostname(sources map[string]string) string {
    hostname, _ := os.Hostname()
    sources["HOSTNAME"] = "system"
    if hostname == "" {
        hostname = "unknown"
    }
    return hostname
}

func parseFeatureFlags(sources map[string]string) map[string]bool {
    flags := make(map[string]bool)
    
    flagNames := []string{"FEATURE_A", "FEATURE_B", "FEATURE_C"}
    
    for _, name := range flagNames {
        value := os.Getenv(name)
        if value == "" {
            flags[name] = false
            sources[name] = "default"
        } else {
            flags[name] = strings.ToLower(value) == "true" || value == "1"
            sources[name] = "environment"
        }
    }
    
    return flags
}

func (c *Config) SafeConfig() map[string]interface{} {
    return map[string]interface{}{
        "app_name":    c.AppName,
        "app_env":     c.AppEnv,
        "hostname":    c.Hostname,
        "port":        c.Port,
        "log_level":   c.LogLevel,
        "max_workers": c.MaxWorkers,
        "timeout_sec": int(c.Timeout.Seconds()),
        "feature_flags": c.FeatureFlags,
        "config_sources": c.Sources,
    }
}

func (c *Config) DatabaseDSN() string {
    if c.AppEnv == "development" {
        return c.DatabaseURL
    }
    return "postgres://****:****@****/****"
}