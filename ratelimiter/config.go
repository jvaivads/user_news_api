package ratelimiter

import "time"

const (
	StatusType    = "Status"
	NewsType      = "News"
	MarketingType = "Marketing"
)

// DefaultConfigs set the business rules needed for the rate limiter.
// Depending on the context, they could be migrated to a database for getting dynamism
var DefaultConfigs = map[string]Config{
	StatusType: {
		Max: 2,
		TTL: time.Minute,
	},
	NewsType: {
		Max: 1,
		TTL: 24 * time.Hour,
	},
	MarketingType: {
		Max: 3,
		TTL: time.Hour,
	},
}

type Config struct {
	Max int64
	TTL time.Duration
}
