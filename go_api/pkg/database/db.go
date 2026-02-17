package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"go_api/internal/config"
)

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", cfg.PostgresDSN())
		if err != nil {
			log.Printf("Attempt %d: Failed to open connection: %v", i, err)
			time.Sleep(time.Second * time.Duration(i))
			continue
		}

		err = db.Ping()
		if err == nil {
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(5 * time.Minute)

			if err := createTable(db); err != nil {
				return nil, err
			}

			log.Println("Successfully connected to database")
			return db, nil
		}

		log.Printf("Attempt %d: Failed to ping: %v", i, err)
		db.Close()
		time.Sleep(time.Second * time.Duration(i))
	}

	return nil, fmt.Errorf("failed to connect after 5 attempts: %v", err)
}

func createTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := db.Exec(query)
	return err
}
