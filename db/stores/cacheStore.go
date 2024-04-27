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
func (cs CacheStore) PutEmailConfirmation(confirmationCode string, userId int64) error {
	ctx := context.Background()
	status := cs.client.Set(ctx, fmt.Sprintf("email_confirmation:%v", confirmationCode), userId, time.Second*600)
	return status.Err()
}

func (cs CacheStore) GetUserId(confirmationCode string) (int64, error) {
	ctx := context.Background()
	status := cs.client.Get(ctx, fmt.Sprintf("email_confirmation:%v", confirmationCode))
	if status.Err() != nil {
		return 0, status.Err()
	}
	userId, err := strconv.Atoi(status.Val())
	if err != nil {
		return 0, err
	}
	return int64(userId), nil
}
