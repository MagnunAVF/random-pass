package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{
		client: client,
		ctx:    context.Background(),
	}
}

func generateKey(userID string) string {
	return "passwords:" + userID
}

func (r *RedisRepo) Save(userID string, password string) error {
	key := generateKey(userID)
	r.client.LPush(r.ctx, key, password)
	r.client.LTrim(r.ctx, key, 0, 4)

	return nil
}

func (r *RedisRepo) GetLastFive(userID string) ([]string, error) {
	return r.client.LRange(r.ctx, generateKey(userID), 0, 4).Result()
}
