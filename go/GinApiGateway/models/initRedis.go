package models

import "github.com/redis/go-redis/v9"

var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}