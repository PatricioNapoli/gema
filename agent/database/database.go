package database

import "github.com/go-redis/redis"

func New() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		DB:       0,
	})

	return client
}