package database

import (
	"log"

	"github.com/go-redis/redis"
)

func New() *redis.Client {
	log.Printf("Connecting to database.")

	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	return client
}
