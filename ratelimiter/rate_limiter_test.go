package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
	"user_news_api/ratelimiter/mocks"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReached(t *testing.T) {
	tests := []struct {
		name           string
		mockApplier    func(mockRedis *mocks.RedisCounter)
		key            string
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Counter below max",
			mockApplier: func(mockRedis *mocks.RedisCounter) {
				mockRedis.On("Incr", mock.Anything, "testKey-suffix").Return(redis.NewIntResult(5, nil)).Once()
				mockRedis.On("TTL", mock.Anything, "testKey-suffix").Return(redis.NewDurationResult(10*time.Second, nil)).Once()
			},
			key:            "testKey",
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Counter above max",
			mockApplier: func(mockRedis *mocks.RedisCounter) {
				mockRedis.On("Incr", mock.Anything, "testKey-suffix").Return(redis.NewIntResult(15, nil)).Once()
				mockRedis.On("TTL", mock.Anything, "testKey-suffix").Return(redis.NewDurationResult(10*time.Second, nil)).Once()
			},
			key:            "testKey",
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Error increasing counter",
			mockApplier: func(mockRedis *mocks.RedisCounter) {
				mockRedis.On("Incr", mock.Anything, "testKey-suffix").Return(redis.NewIntResult(0, errors.New("error"))).Once()
			},
			key:            "testKey",
			expectedResult: false,
			expectedError:  fmt.Errorf("error increasing user counter due to: %w", errors.New("error")),
		},
		{
			name: "TTL not set, then set",
			mockApplier: func(mockRedis *mocks.RedisCounter) {
				mockRedis.On("Incr", mock.Anything, "testKey-suffix").Return(redis.NewIntResult(5, nil)).Once()
				mockRedis.On("TTL", mock.Anything, "testKey-suffix").Return(redis.NewDurationResult(-2, nil)).Once()
				mockRedis.On("Expire", mock.Anything, "testKey-suffix", mock.Anything).Return(redis.NewBoolResult(true, nil)).Once()
			},
			key:            "testKey",
			expectedResult: false,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := mocks.NewRedisCounter(t)

			tt.mockApplier(mockRedis)

			rl := rateLimiter{
				db:        mockRedis,
				max:       10,
				suffixKey: "suffix",
				ttl:       30 * time.Second,
			}

			result, err := rl.Reached(context.Background(), tt.key)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
