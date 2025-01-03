package config

import (
	"context"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func InitRedis() {
	once.Do(func() {
		redisAddr := os.Getenv("REDIS_ADDR")
		if redisAddr == "" {
			redisAddr = "gau-redis:6379"
		}

		client := redis.NewClient(&redis.Options{
			Addr: redisAddr,
		})

		ctx := context.Background()
		_, err := client.Ping(ctx).Result()
		if err != nil {
			panic("Failed to connect to Redis: " + err.Error())
		}

		redisClient = client
	})
}

func GetRedisClient() *redis.Client {
	return redisClient
}
