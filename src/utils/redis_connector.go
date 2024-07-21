package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
)

func GetRedisPool(ctx context.Context, log *logrus.Logger) *redis.Client {
	godotenv.Load()
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password:     "",               // No password set
		DB:           0,                // Use default DB
		PoolSize:     10,               // Default is 10
		MinIdleConns: 3,                // Minimum number of idle connections
		PoolTimeout:  30 * time.Second, // Amount of time client waits for a connection if all connections are busy
		IdleTimeout:  5 * time.Minute,  // Amount of time after which idle connections are closed
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return rdb
}
