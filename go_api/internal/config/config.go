package config

import (
    "os"
    "strconv"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    Port       string
    AppEnv     string
}

func Load() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "postgres"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        Port:       getEnv("PORT", "8080"),
        AppEnv:     getEnv("APP_ENV", "development"),
    }
}

func (c *Config) PostgresDSN() string {
    return "host=" + c.DBHost +
        " port=" + c.DBPort +
        " user=" + c.DBUser +
        " password=" + c.DBPassword +
        " dbname=" + c.DBName +
        " sslmode=" + c.DBSSLMode
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}