package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "root",
		DB:       0,
	})
	//checking connection
	ctx := context.Background()
	if result, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Error connecting to redis: %v", result)
	}
	return rdb
}
