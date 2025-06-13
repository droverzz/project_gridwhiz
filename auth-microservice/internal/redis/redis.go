package redis

import (
	"context"
	"os"
	"sync"

	"github.com/go-redis/redis/v8" // ถ้าเปลี่ยนมาใช้ redis v8 ตามที่แนะนำก่อนหน้านี้
)

var (
	rdb  *redis.Client
	once sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	})
	return rdb
}

func Ping(ctx context.Context) error {
	return GetRedisClient().Ping(ctx).Err()
}
