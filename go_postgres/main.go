package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Connected to PostgreSQL! Users count: %d", count)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectDB() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "postgres")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Failed to open connection: %v", i+1, err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to database")
			return db, nil
		}

		log.Printf("Attempt %d: Failed to ping: %v", i+1, err)
		db.Close()
		time.Sleep(time.Second * time.Duration(i+1))
	}

	return nil, fmt.Errorf("failed to connect after 10 attempts: %v", err)
}

func createTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	_, err := db.Exec(query)
	return err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
