package contract

import "time"

type RedisClient interface {
	Get(key string, dest any) bool
	Set(key string, value any, ttl time.Duration) bool
	Delete(key string) bool
}
