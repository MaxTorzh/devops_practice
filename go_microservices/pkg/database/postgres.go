package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"go_microservices/internal/config"
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

		db.SetMaxOpenConns(cfg.DBMaxConns)
		db.SetMaxIdleConns(cfg.DBMaxIdleConns)
		db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to PostgreSQL")

			return db, nil
		}

		log.Printf("Attempt %d: Failed to ping: %v", i, err)
		db.Close()
		time.Sleep(time.Second * time.Duration(i))
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after 5 attempts: %v", err)
}
