package config

import (
	"os"
	"strconv"
)

type Config struct {
	DataDir       string
	Port          string
	MaxUploadSize int64
}

func Load() *Config {
	return &Config{
		DataDir:       getEnv("DATA_DIR", "/data"),
		Port:          getEnv("PORT", "8080"),
		MaxUploadSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 10*1024*1024),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err != nil {
			return intVal
		}
	}
	return defaultValue
}
