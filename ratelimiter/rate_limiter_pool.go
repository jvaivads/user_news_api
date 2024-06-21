package ratelimiter

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrMessageTypeNotValid = errors.New("message type not valid")
)

func NewLimiterPool(db redis.Client, configs map[string]Config) LimiterPool {
	limiterPool := LimiterPool{
		limiters: make(map[string]rateLimiter),
	}

	for msgType, config := range configs {
		limiterPool.limiters[msgType] = newRateLimiter(db, config.Max, msgType, config.TTL)
	}

	return limiterPool
}

type LimiterPool struct {
	limiters map[string]rateLimiter
}

func (lp LimiterPool) Reached(ctx context.Context, user string, msgType string) (bool, error) {
	limiter, ok := lp.limiters[msgType]
	if !ok {
		return false, ErrMessageTypeNotValid
	}

	return limiter.Reached(ctx, user)
}
