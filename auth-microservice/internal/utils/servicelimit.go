package utils

import (
	"context"
	"time"

	"auth-microservice/internal/redis"

	"github.com/go-redis/redis_rate/v9"
)

var limiter = redis_rate.NewLimiter(redis.GetRedisClient())

func RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	res, err := limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   limit,
		Period: window,
		Burst:  limit,
	})
	if err != nil {
		return false, err
	}
	return res.Allowed > 0, nil
}
