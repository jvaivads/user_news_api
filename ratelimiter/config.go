package ratelimiter

import "time"

// DefaultConfigs set the business rules needed for the rate limiter.
// Depending on the context, they could be migrated to a database for getting dynamism
var DefaultConfigs = map[string]Config{
	"Status": {
		Max: 2,
		TTL: time.Minute,
	},
	"News": {
		Max: 1,
		TTL: 24 * time.Hour,
	},
	"Marketing": {
		Max: 3,
		TTL: time.Hour,
	},
}

type Config struct {
	Max int64
	TTL time.Duration
}
