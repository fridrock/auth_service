package stores

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheStore struct {
	client *redis.Client
}

func CreateCacheStore(client *redis.Client) *CacheStore {
	return &CacheStore{
		client: client,
	}
}

func (cs CacheStore) PutUserId(to string, key string, userId int64) error {
	ctx := context.Background()
	status := cs.client.Set(ctx, fmt.Sprintf("%v:%v", to, key), userId, time.Second*600)
	return status.Err()
}

func (cs CacheStore) GetUserId(from, code string) (int64, error) {
	ctx := context.Background()
	status := cs.client.Get(ctx, fmt.Sprintf("%v:%v", from, code))
	if status.Err() != nil {
		return 0, status.Err()
	}
	userId, err := strconv.Atoi(status.Val())
	if err != nil {
		return 0, err
	}
	return int64(userId), nil
}
