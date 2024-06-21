package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func newRateLimiter(db redis.Client, max int64, suffixKey string, ttl time.Duration) rateLimiter {
	return rateLimiter{
		db:        db,
		max:       max,
		suffixKey: suffixKey,
		ttl:       ttl,
	}
}

// RedisCounter is an abstraction for redis.Client making it mockeable
type RedisCounter interface {
	Incr(ctx context.Context, key string) *redis.IntCmd
	Decr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd
}

type rateLimiter struct {
	db        RedisCounter // db can be shared between different rateLimiter, however, suffixKey must be different
	max       int64        // max is the maximum hits
	suffixKey string       // suffixKey is used for avoiding collisions between different rateLimiter
	ttl       time.Duration
}

func (rl rateLimiter) Reached(ctx context.Context, key string) (bool, error) {
	key = fmt.Sprintf("%s-%s", key, rl.suffixKey)

	// It increases by one the counter associated with the key.
	// If the key does not exist, then is created with one value.
	counter, err := rl.db.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error increasing user counter due to: %w", err)
	}

	// It checks if the key resource has a TTL set. If it is not, then is created.
	// Both the TTL and Expire methods are not validated in case of fail, but they are retried in future requests.
	// This approach avoids to use methods that lock the resource.
	// Depending on the context this could be a good or bad strategy.
	if ttl := rl.db.TTL(ctx, key).Val(); ttl < 0 {
		_ = rl.db.Expire(ctx, key, rl.ttl)
	}

	return counter > rl.max, nil
}
