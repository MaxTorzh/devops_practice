package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_microservices/internal/config"
	"go_microservices/internal/handler"
	"go_microservices/internal/repository/postgres"
	"go_microservices/internal/repository/redis"
	"go_microservices/internal/service"
	"go_microservices/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer db.Close()

	rdb, err := database.NewRedis(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer rdb.Close()

	userRepo := postgres.NewUserRepository(db)
	productRepo := postgres.NewProductRepository(db)
	cacheRepo := redis.NewCacheRepository(rdb, cfg.CacheTTL)

	userService := service.NewUserService(userRepo, cacheRepo)
	productService := service.NewProductService(productRepo, cacheRepo)

	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/health", handleHealth)

	userHandler.RegisterRoutes(mux)
	productHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		log.Printf("Environment: %s", cfg.AppEnv)
		log.Printf("PostgreSQL: %s:%s", cfg.DBHost, cfg.DBPort)
		log.Printf("Redis: %s:%s", cfg.RedisHost, cfg.RedisPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	gracefulShutdown(srv)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"name":    "Go Microservices API",
		"version": "1.0.0",
		"endpoints": []string{
			"GET    /health",
			"GET    /users",
			"POST   /users",
			"GET    /users/{id}",
			"PUT    /users/{id}",
			"DELETE /users/{id}",
			"GET    /products",
			"POST   /products",
			"GET    /products/{id}",
			"PUT    /products/{id}",
			"DELETE /products/{id}",
			"PATCH  /products/{id}/stock",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped")
}
