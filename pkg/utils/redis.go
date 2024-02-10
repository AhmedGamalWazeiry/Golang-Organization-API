package utils

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func InitRedis() {
	viper.SetConfigType("yaml")


    viper.SetConfigFile("./config/app-config.yaml") 
	
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		panic("Failed to read the configuration file")
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
}