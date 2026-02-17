package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go_microservices/internal/config"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr(),
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 1; i <= 3; i++ {
		err := client.Ping(ctx).Err()
		if err == nil {
			log.Println("Successfully connected to Redis")
			return client, nil
		}

		log.Printf("Attempt %d: Failed to connect to Redis: %v", i, err)
		time.Sleep(time.Second * 2)
	}

	return nil, fmt.Errorf("failed to connect to Redis after 3 attempts")
}
