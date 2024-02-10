package utils

import (
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// InitRedis initializes the Redis client.
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Update with your Redis server address
		Password: "",                // No password for local Redis server
		DB:       0,
	})
}