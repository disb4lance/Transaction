package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transaction-service/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClientWithError(cfg config.RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}, nil
}

func (r *RedisClient) Get(key string, dest interface{}) bool {
	if r == nil || r.client == nil {
		return false
	}

	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false
	}

	return true
}

func (r *RedisClient) Set(key string, value interface{}, ttl time.Duration) bool {
	if r == nil || r.client == nil {
		return false
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return false
	}

	err = r.client.Set(r.ctx, key, jsonData, ttl).Err()
	return err == nil
}

func (r *RedisClient) Delete(key string) bool {
	if r == nil || r.client == nil {
		return false
	}

	err := r.client.Del(r.ctx, key).Err()
	return err == nil
}
