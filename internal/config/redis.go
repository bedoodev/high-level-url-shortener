package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		return err
	}

	zap.L().Info("Redis connected successfully")

	return nil
}
