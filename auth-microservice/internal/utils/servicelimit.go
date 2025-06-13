package utils

import (
	"auth-microservice/internal/redis"
	"context"
	"time"
)

func RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	rdb := redis.GetRedisClient()
	luaScript := `
    local current
    current = redis.call("INCR", KEYS[1])
    if tonumber(current) == 1 then
        redis.call("EXPIRE", KEYS[1], ARGV[1])
    end
    if tonumber(current) > tonumber(ARGV[2]) then
        return 0
    end
    return 1
    `
	res, err := rdb.Eval(ctx, luaScript, []string{key}, int(window.Seconds()), limit).Result()
	if err != nil {
		return false, err
	}
	allowed := res.(int64) == 1
	return allowed, nil
}
