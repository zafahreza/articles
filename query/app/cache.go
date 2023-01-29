package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
)

func InitCache() *redis.Client {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	host := os.Getenv("REDIS_HOST")
	addr := fmt.Sprintf("%s:6379", host)

	cache := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: "article",
		Password: "password",
		DB:       0,
	})

	return cache
}
