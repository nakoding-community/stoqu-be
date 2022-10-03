package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"time"
)

type (
	RateLimit struct {
		Allowed    bool
		RetryAfter time.Duration
	}
	RateLimiter interface {
		GetRateLimit(ctx context.Context, key string, limit int, limit_unit time.Duration) (rl RateLimit, err error)
	}
	rateLimiter struct {
		limiter *redis_rate.Limiter
	}
)

func StringToLimitUnit(s string) time.Duration {
	switch s {
	case "second":
		return time.Second
	case "minute":
		return time.Minute
	case "hour":
		return time.Hour
	default:
		return time.Second
	}
}

func New(rc *redis.Client) RateLimiter {
	limiter := &rateLimiter{
		limiter: redis_rate.NewLimiter(rc),
	}

	return limiter
}

func (r rateLimiter) GetRateLimit(ctx context.Context, key string, limit int, limit_unit time.Duration) (rl RateLimit, err error) {
	var req redis_rate.Limit
	switch limit_unit {
	case time.Hour:
		req = redis_rate.PerHour(limit)
	case time.Minute:
		req = redis_rate.PerMinute(limit)
	default:
		req = redis_rate.PerSecond(limit)
	}
	res, err := r.limiter.Allow(ctx, key, req)
	if err != nil {
		return
	}
	rl.Allowed = res.Remaining > 0 || (res.Remaining == 0 && res.Allowed > 0)
	rl.RetryAfter = res.RetryAfter
	return
}
