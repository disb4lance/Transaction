package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}
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
