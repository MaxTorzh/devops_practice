package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_api/internal/config"
	"go_api/internal/handler"
	"go_api/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	h := handler.NewUserHandler(db, cfg)

	router := h.SetupRoutes()

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Printf("Server started on port: %s", cfg.Port)
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server stopped")
}
