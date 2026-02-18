package models

import "time"

type ConfigInfo struct {
    AppName         string            `json:"app_name"`
    AppEnv          string            `json:"app_env"`
    Hostname        string            `json:"hostname"`
    Port            int               `json:"port"`
    LogLevel        string            `json:"log_level"`
    DatabaseURL     string            `json:"database_url,omitempty"`
    RedisURL        string            `json:"redis_url,omitempty"`
    MaxWorkers      int               `json:"max_workers"`
    TimeoutSeconds  int               `json:"timeout_seconds"`
    FeatureFlags    map[string]bool   `json:"feature_flags"`
    ConfigSources   map[string]string `json:"config_sources"`
}

type HealthStatus struct {
    Status    string    `json:"status"`
    Version   string    `json:"version"`
    Uptime    string    `json:"uptime"`
    Timestamp time.Time `json:"timestamp"`
    Env       string    `json:"environment"`
}