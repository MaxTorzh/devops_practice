package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	//app
	AppEnv string
	Port   string

	//psql
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxConns        int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	//redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	CacheTTL      time.Duration
}

func Load() *Config {
	return &Config{
		//app
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("PORT", "8080"),

		//psql
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "postgres"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBMaxConns:        getEnvAsInt("DB_MAX_CONNS", 20),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),

		//redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
		CacheTTL:      getEnvAsDuration("CACHE_TTL", 5*time.Minute),
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

func (c *Config) RedisAddr() string {
	return c.RedisHost + ":" + c.RedisPort
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

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
