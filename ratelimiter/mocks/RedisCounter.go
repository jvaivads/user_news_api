// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	redis "github.com/redis/go-redis/v9"

	time "time"
)

// RedisCounter is an autogenerated mock type for the RedisCounter type
type RedisCounter struct {
	mock.Mock
}

// Decr provides a mock function with given fields: ctx, key
func (_m *RedisCounter) Decr(ctx context.Context, key string) *redis.IntCmd {
	ret := _m.Called(ctx, key)

	var r0 *redis.IntCmd
	if rf, ok := ret.Get(0).(func(context.Context, string) *redis.IntCmd); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.IntCmd)
		}
	}

	return r0
}

// Expire provides a mock function with given fields: ctx, key, expiration
func (_m *RedisCounter) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	ret := _m.Called(ctx, key, expiration)

	var r0 *redis.BoolCmd
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Duration) *redis.BoolCmd); ok {
		r0 = rf(ctx, key, expiration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.BoolCmd)
		}
	}

	return r0
}

// Incr provides a mock function with given fields: ctx, key
func (_m *RedisCounter) Incr(ctx context.Context, key string) *redis.IntCmd {
	ret := _m.Called(ctx, key)

	var r0 *redis.IntCmd
	if rf, ok := ret.Get(0).(func(context.Context, string) *redis.IntCmd); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.IntCmd)
		}
	}

	return r0
}

// TTL provides a mock function with given fields: ctx, key
func (_m *RedisCounter) TTL(ctx context.Context, key string) *redis.DurationCmd {
	ret := _m.Called(ctx, key)

	var r0 *redis.DurationCmd
	if rf, ok := ret.Get(0).(func(context.Context, string) *redis.DurationCmd); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.DurationCmd)
		}
	}

	return r0
}

// NewRedisCounter creates a new instance of RedisCounter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRedisCounter(t interface {
	mock.TestingT
	Cleanup(func())
}) *RedisCounter {
	mock := &RedisCounter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}