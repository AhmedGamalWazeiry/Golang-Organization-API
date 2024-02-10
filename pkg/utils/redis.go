package utils

import (
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis(address string,password string , db int) {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:      db,
	})
}